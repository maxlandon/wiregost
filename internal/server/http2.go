package server

import (
	// Standard
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	// 3rd Party
	"github.com/cretz/gopaque/gopaque"
	"github.com/evilsocket/islazy/tui"
	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/h2quic"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"go.dedis.ch/kyber"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	// Wiregost
	"github.com/maxlandon/wiregost/internal/agents"
	"github.com/maxlandon/wiregost/internal/core"
	"github.com/maxlandon/wiregost/internal/logging"
	"github.com/maxlandon/wiregost/internal/messages"
	"github.com/maxlandon/wiregost/internal/util"
)

// Server is a structure for creating and instantiating new server objects
type Server struct {
	ID          uuid.UUID      // Unique identifier for the Server object
	Interface   string         // The network adapter interface the server will listen on
	Port        int            // The port the server will listen on
	Protocol    string         // The protocol (i.e. HTTP/2 or HTTP/3) the server will use
	Key         string         // The x.509 private key used for TLS encryption
	Certificate string         // The x.509 public key used for TLS encryption
	Server      interface{}    // A Golang server object (i.e http.Server or h3quic.Server)
	Mux         *http.ServeMux // The message handler/multiplexer
	jwtKey      []byte         // The password used by the server to create JWTs
	Psk         string         // The pre-shared key password used prior to Password Authenticated Key Exchange (PAKE)
	opaqueKey   kyber.Scalar   // OPAQUE server's keys

	// Added
	Workspace   string
	WorkspaceID int
	Running     bool
	log         *logging.WorkspaceLogger
}

// New instantiates a new server object and returns it
func New(iface string, port int, protocol string, key string, certificate string, psk string, workspace string, workspaceID int, logger *logging.WorkspaceLogger) (Server, error) {
	s := Server{
		ID:          uuid.NewV4(),
		Protocol:    protocol,
		Interface:   iface,
		Port:        port,
		Mux:         http.NewServeMux(),
		jwtKey:      []byte(core.RandStringBytesMaskImprSrc(32)), // Used to sign and encrypt JWT
		Psk:         psk,
		Workspace:   workspace,
		WorkspaceID: workspaceID,
		Running:     false,
		log:         logger,
	}
	// Set context logger
	log := s.log.WithFields(logrus.Fields{"component": "server", "workspaceId": s.WorkspaceID})

	// OPAQUE Server Public/Private keys; Can be used with every agent
	s.opaqueKey = gopaque.CryptoDefault.NewKey(nil)

	var cer tls.Certificate
	var err error
	// Check if certificate exists on disk
	_, errCrt := os.Stat(certificate)
	if os.IsNotExist(errCrt) {
		// generate a new ephemeral certificate
		log.Warnf("No certificate found at %s", certificate)
		t := "Creating in-memory x.509 certificate used for this session only."
		log.Infof(t)
		cerp, err := util.GenerateTLSCert(nil, nil, nil, nil, nil, nil, true) //ec certs not supported (yet) :(
		if err != nil {
			m := fmt.Sprintf("There was an error generating the SSL/TLS certificate:\r\n%s", err.Error())
			log.Warnf(m)
			return s, err
		}
		cer = *cerp
	} else {
		if errCrt != nil {
			m := fmt.Sprintf("There was an error importing the SSL/TLS x509 certificate:\r\n%s", errCrt.Error())
			log.Warnf(m)
			return s, errCrt
		}
		s.Certificate = certificate

		_, errKey := os.Stat(key)
		if errKey != nil {
			m := fmt.Sprintf("There was an error importing the SSL/TLS x509 key:\r\n%s", errKey.Error())
			log.Warnf(m)
			return s, errKey
		}
		s.Key = key

		cer, err = tls.LoadX509KeyPair(certificate, key)
		if err != nil {
			m := fmt.Sprintf("There was an error importing the SSL/TLS x509 key pair\r\n%s", err.Error())
			log.Warnf(m)
			log.Warnf("Ensure a keypair is located in the data/x509 directory")
			return s, err
		}
	}

	if len(cer.Certificate) < 1 || cer.PrivateKey == nil {
		m := "Unable to import certificate for use in Merlin: empty certificate structure."
		log.Warnf(m)
		return s, errors.New("empty certificate structure")
	}

	// Parse into X.509 format
	x, errX509 := x509.ParseCertificate(cer.Certificate[0])
	if errX509 != nil {
		m := fmt.Sprintf("There was an error parsing the tls.Certificate structure into a x509.Certificate"+
			" structure:\r\n%s", errX509.Error())
		log.Warnf(m)
		return s, errX509
	}
	// Create fingerprint
	S256 := sha256.Sum256(x.Raw)
	sha256Fingerprint := hex.EncodeToString(S256[:])

	// merlinCRT is the string representation of the SHA1 fingerprint for the public x.509 certificate distributed with Merlin
	merlinCRT := "4af9224c77821bc8a46503cfc2764b94b1fc8aa2521afc627e835f0b3c449f50"

	// Check to see if the Public Key SHA1 finger print matches the certificate distributed with Merlin for testing
	if merlinCRT == sha256Fingerprint {
		log.Warnf("Insecure publicly distributed Merlin x.509 testing certificate in use")
		log.Infof("Additional details: https://github.com/Ne0nd0g/merlin/wiki/TLS-Certificates")
	}

	// Log certificate information
	log.Debugf(fmt.Sprintf("Starting Merlin Server using an X.509 certificate with a %s signature of %s",
		x.SignatureAlgorithm.String(), hex.EncodeToString(x.Signature)))
	log.Debugf(fmt.Sprintf("Starting Merlin Server using an X.509 certificate with a public key of %v", x.PublicKey))
	log.Debugf(fmt.Sprintf("Starting Merlin Server using an X.509 certificate with a serial number of %d", x.SerialNumber))
	log.Debugf(fmt.Sprintf("Starting Merlin Server using an X.509 certifcate with a subject of %s", x.Subject.String()))
	log.Debugf(fmt.Sprintf("Starting Merlin Server using an X.509 certificate with a SHA256 hash, "+
		"calculated by Merlin, of %s", sha256Fingerprint))

	// Configure TLS
	TLSConfig := &tls.Config{
		Certificates:             []tls.Certificate{cer},
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		},
		//NextProtos: []string{protocol}, //Dont need to specify because server will pick
	}

	s.Mux.HandleFunc("/", s.agentHandler)

	srv := &http.Server{
		Addr:           s.Interface + ":" + strconv.Itoa(s.Port),
		Handler:        s.Mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		TLSConfig:      TLSConfig,
		//TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0), // <- Disables HTTP/2
	}

	if s.Protocol == "h2" {
		s.Server = srv
	} else if s.Protocol == "hq" {
		s.Server = &h2quic.Server{
			Server: srv,
			QuicConfig: &quic.Config{
				KeepAlive:                   false,
				IdleTimeout:                 168 * time.Hour,
				RequestConnectionIDOmission: false,
			},
		}

	} else {
		return s, fmt.Errorf("%s is an invalid server protocol", s.Protocol)
	}
	return s, nil
}

// Run function starts the server on the preconfigured port for the preconfigured service
func (s *Server) Run() (status string, err error) {
	// Test context logger
	log := s.log.WithFields(logrus.Fields{"component": "server", "workspaceId": s.WorkspaceID})

	// Sleep to allow the shell to start up
	time.Sleep(45 * time.Millisecond)
	if s.Psk == "merlin" {
		log.Warnf("Listener was started using \"merlin\" as the Pre-Shared Key (PSK) allowing anyone" +
			" decrypt message traffic.")
		log.Infof("Consider changing the PSK by using the -Psk command line flag.")
	}
	m := fmt.Sprintf("%s[*]%s Starting %s listener on %s:%d %s(pre-shared key: %s%s)",
		tui.GREEN, tui.RESET, s.Protocol, s.Interface, s.Port, tui.DIM, s.Psk, tui.RESET)

	if s.Protocol == "h2" {
		server := s.Server.(*http.Server)

		defer func() {
			err := server.Close()
			if err != nil {
				m := fmt.Sprintf("There was an error starting the %s server:\r\n%s", s.Protocol, err.Error())
				log.Warnf(m)
				return
			}
		}()
		s.Running = true
		go log.Errorf(server.ListenAndServeTLS(s.Certificate, s.Key).Error())
		return m, nil
	} else if s.Protocol == "hq" {
		server := s.Server.(*h2quic.Server)

		defer func() {
			err := server.Close()
			if err != nil {
				m := fmt.Sprintf("There was an error starting the hq server:\r\n%s", err.Error())
				log.Warnf(m)
				return
			}
		}()
		go log.Errorf(server.ListenAndServeTLS(s.Certificate, s.Key).Error())
		// Server is now running
		s.Running = true
		return m, nil
	}
	return fmt.Errorf("%s is an invalid server protocol", s.Protocol).Error(), nil
}

// agentHandler function is responsible for all Merlin agent traffic
func (s *Server) agentHandler(w http.ResponseWriter, r *http.Request) {
	// Test context logger
	log := s.log.WithFields(logrus.Fields{"component": "server", "workspaceId": s.WorkspaceID})

	log.Infof(fmt.Sprintf("Received %s %s connection from %s", r.Proto, r.Method, r.RemoteAddr))

	log.Debugf(fmt.Sprintf("HTTP Connection Details:"))
	log.Debugf(fmt.Sprintf("Host: %s", r.Host))
	log.Debugf(fmt.Sprintf("URI: %s", r.RequestURI))
	log.Debugf(fmt.Sprintf("Method: %s", r.Method))
	log.Debugf(fmt.Sprintf("Protocol: %s", r.Proto))
	log.Debugf(fmt.Sprintf("Headers: %s", r.Header))
	log.Debugf(fmt.Sprintf("TLS Negotiated Protocol: %s", r.TLS.NegotiatedProtocol))
	log.Debugf(fmt.Sprintf("TLS Cipher Suite: %d", r.TLS.CipherSuite))
	log.Debugf(fmt.Sprintf("TLS Server Name: %s", r.TLS.ServerName))
	log.Debugf(fmt.Sprintf("Content Length: %d", r.ContentLength))

	// Check for Merlin PRISM activity
	if r.UserAgent() == "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/40.0.2214.85 Safari/537.36 " {
		log.Warnf(fmt.Sprintf("Someone from %s is attempting to fingerprint this Merlin server", r.RemoteAddr))
		//w.WriteHeader(404)
	}
	// Make sure the message has a JWT
	token := r.Header.Get("Authorization")
	if token == "" {
		log.Warnf("incoming request did not contain an Authorization header")
		w.WriteHeader(404)
		return
	}

	if r.Method == http.MethodPost {

		var returnMessage messages.Base
		var err error
		var key []byte

		//Read the request message until EOF
		requestBytes, errRequestBytes := ioutil.ReadAll(r.Body)
		if errRequestBytes != nil {
			log.Warnf(fmt.Sprintf("There was an error reading a POST message sent by an "+
				"agent:\r\n%s", errRequestBytes))
			return
		}

		// Decode gob to JWE string
		var jweString string
		errDecode := gob.NewDecoder(bytes.NewReader(requestBytes)).Decode(&jweString)
		if errDecode != nil {
			log.Warnf(fmt.Sprintf("there was an error decoding JWE payload message sent by an agent:\r\n%s", errDecode.Error()))
			return
		}

		// Validate JWT and get claims
		var agentID uuid.UUID
		var errValidate error

		// Set return headers
		//w.Header().Set("Content-Type", "application/octet-stream")

		// Validate JWT using HTTP interface JWT key; Given to authenticated agents by server
		agentID, errValidate = validateJWT(strings.Split(token, " ")[1], s.jwtKey, s)
		// If agentID was returned, then message contained a JWT encrypted with the HTTP interface key
		if (errValidate != nil) && (agentID == uuid.Nil) {
			log.Warnf(errValidate.Error())
			log.Infof("trying again with interface PSK")
			// Validate JWT using interface PSK; Used by unauthenticated agents
			hashedKey := sha256.Sum256([]byte(s.Psk))
			key = hashedKey[:]
			agentID, errValidate = validateJWT(strings.Split(token, " ")[1], key, s)
			if errValidate != nil {
				log.Warnf(errValidate.Error())
				w.WriteHeader(404)
				return
			}
			log.Infof("Unauthenticated JWT")

			// Decrypt the HTTP payload, a JWE, using interface PSK
			k, errDecryptPSK := decryptJWE(jweString, key, s)
			// Successfully decrypted JWE with interface PSK
			if errDecryptPSK == nil {
				log.Debugf(fmt.Sprintf("POST DATA: %v", k))
				log.Infof(fmt.Sprintf("Received %s message, decrypted JWE with interface PSK", k.Type))

				messagePayloadBytes := new(bytes.Buffer)

				// Allowed unauthenticated message types w/ PSK signed JWT and PSK encrypted JWT
				switch k.Type {
				case "AuthInit":
					serverAuthInit, err := agents.OPAQUEAuthenticateInit(k)
					if err != nil {
						log.Warnf(err.Error())
						w.WriteHeader(404)
						return
					}
					log.Infof(fmt.Sprintf("Received new agent OPAQUE authentication from %s", agentID))

					// Encode return message into a gob
					errAuthInit := gob.NewEncoder(messagePayloadBytes).Encode(serverAuthInit)
					if errAuthInit != nil {
						m := fmt.Sprintf("there was an error encoding the return message into a gob:\r\n%s", errAuthInit.Error())
						log.Warnf(m)
						w.WriteHeader(404)
						return
					}
				case "RegInit":
					// Added server ID as parameter for registering agent to this server.
					serverRegInit, err := agents.OPAQUERegistrationInit(k, s.opaqueKey, s.ID)
					if err != nil {
						log.Errorf(err.Error())
						w.WriteHeader(404)
						return
					}
					log.Infof(fmt.Sprintf("Received new agent OPAQUE user registration initialization from %s", agentID))

					// Encode return message into a gob
					errRegInit := gob.NewEncoder(messagePayloadBytes).Encode(serverRegInit)
					if errRegInit != nil {
						m := fmt.Sprintf("there was an error encoding the return message into a gob:\r\n%s", errRegInit.Error())
						log.Warnf(m)
						w.WriteHeader(404)
						return
					}
				case "RegComplete":
					serverRegComplete, err := agents.OPAQUERegistrationComplete(k)
					if err != nil {
						log.Warnf(err.Error())
						w.WriteHeader(404)
						return
					}
					log.Infof(fmt.Sprintf("Received new agent OPAQUE user registration complete from %s", agentID))

					// Encode return message into a gob
					errRegInit := gob.NewEncoder(messagePayloadBytes).Encode(serverRegComplete)
					if errRegInit != nil {
						m := fmt.Sprintf("there was an error encoding the return message into a gob:\r\n%s", errRegInit.Error())
						log.Warnf(m)
						w.WriteHeader(404)
						return
					}
				default:
					log.Warnf("invalid message type")
					w.WriteHeader(404)
					return
				}
				// Get JWE
				jwe, errJWE := core.GetJWESymetric(messagePayloadBytes.Bytes(), key)
				if errJWE != nil {
					log.Warnf(errJWE.Error())
					w.WriteHeader(404)
					return
				}

				// Set return headers
				w.Header().Set("Content-Type", "application/octet-stream")

				// Encode JWE into gob
				errJWEBuffer := gob.NewEncoder(w).Encode(jwe)
				if errJWEBuffer != nil {
					m := fmt.Errorf("there was an error writing the %s response message to the HTTP stream:\r\n%s", k.Type, errJWEBuffer.Error())
					log.Warnf(m.Error())
					w.WriteHeader(404)
					return
				}

				return
			}
			log.Infof("Unauthenticated JWT w/ Authenticated JWE agent session key")
			// Decrypt the HTTP payload, a JWE, using agent session key
			j, errDecrypt := decryptJWE(jweString, agents.GetEncryptionKey(agentID), s)
			if errDecrypt != nil {
				log.Warnf(errDecrypt.Error())
				w.WriteHeader(404)
				return
			}

			log.Debugf(fmt.Sprintf("POST DATA: %v", j))
			log.Infof(fmt.Sprintf("Received %s message from %s at %s", j.Type, j.ID, time.Now().UTC().Format(time.RFC3339)))

			// Allowed authenticated message with PSK JWT and JWE encrypted with derived secret
			switch j.Type {
			case "AuthComplete":
				returnMessage, err = agents.OPAQUEAuthenticateComplete(j)
				if err != nil {
					log.Infof(fmt.Sprintf("Received new agent OPAQUE authentication from %s", agentID))
				}
			default:
				log.Warnf(fmt.Sprintf("Invalid Activity: %s", j.Type))
				w.WriteHeader(404)
				return
			}
		} else {
			// If not using the PSK, the agent has previously authenticated
			log.Infof("Authenticated JWT")

			// Decrypt JWE
			key = agents.GetEncryptionKey(agentID)

			j, errDecrypt := decryptJWE(jweString, key, s)
			if errDecrypt != nil {
				log.Warnf(errDecrypt.Error())
				w.WriteHeader(404)
				return
			}

			log.Debugf(fmt.Sprintf("POST DATA: %v", j))
			log.Infof("Authenticated JWT w/ Authenticated JWE agent session key")
			log.Infof(fmt.Sprintf("Received %s message from %s at %s", j.Type, j.ID, time.Now().UTC().Format(time.RFC3339)))

			// If both an agentID and error were returned, then the claims were likely bad and the agent needs to re-authenticate
			if (errValidate != nil) && (agentID != uuid.Nil) {
				log.Warnf(fmt.Sprintf("Agent %s connected with expired JWT. Instructing agent to re-authenticate", agentID))
				j.Type = "ReAuthenticate"
			}

			// Authenticated and authorized message types
			switch j.Type {
			case "KeyExchange":
				returnMessage, err = agents.KeyExchange(j)
			case "StatusCheckIn":
				returnMessage, err = agents.StatusCheckIn(j)
			case "CmdResults":
				err = agents.JobResults(j)
			case "AgentInfo":
				err = agents.UpdateInfo(j)
			case "FileTransfer":
				err = agents.FileTransfer(j)
			case "ReAuthenticate":
				returnMessage, err = agents.OPAQUEReAuthenticate(agentID)
			default:
				err = fmt.Errorf("invalid message type: %s", j.Type)
			}
		}

		if err != nil {
			m := fmt.Sprintf("There was an error during while handling a message from agent %s:\r\n%s", agentID.String(), err.Error())
			log.Warnf(m)
			w.WriteHeader(404)
			return
		}

		if returnMessage.Type == "" {
			returnMessage.Type = "ServerOk"
			returnMessage.ID = agentID
		}
		log.Infof(fmt.Sprintf("Sending " + returnMessage.Type + " message type to agent"))

		// Get JWT to add to message.Base for all messages except re-authenticate messages
		if returnMessage.Type != "ReAuthenticate" {
			jsonWebToken, errJWT := getJWT(agentID, s.jwtKey, s) // Passing server for logging purposes
			if errJWT != nil {
				log.Warnf(errJWT.Error())
				w.WriteHeader(404)
				return
			}
			returnMessage.Token = jsonWebToken
		}

		// Encode messages.Base into a gob
		returnMessageBytes := new(bytes.Buffer)
		errReturnMessageBytes := gob.NewEncoder(returnMessageBytes).Encode(returnMessage)
		if errReturnMessageBytes != nil {
			m := fmt.Sprintf("there was an error encoding the %s return message for agent %s into a GOB:\r\n%s", returnMessage.Type, agentID.String(), errReturnMessageBytes.Error())
			log.Warnf(m)
			return
		}

		// Get JWE
		key = agents.GetEncryptionKey(agentID)
		jwe, errJWE := core.GetJWESymetric(returnMessageBytes.Bytes(), key)
		if errJWE != nil {
			log.Warnf(errJWE.Error())
		}

		// Set return headers
		w.Header().Set("Content-Type", "application/octet-stream")

		// Encode JWE to GOB and send it to the agent
		errEncode := gob.NewEncoder(w).Encode(jwe)
		if errEncode != nil {
			m := fmt.Sprintf("There was an error encoding the server AuthComplete GOB message:\r\n%s", errEncode.Error())
			log.Warnf(m)
			return
		}

		// Remove the agent from the server after successfully sending the kill message
		if returnMessage.Type == "AgentControl" {
			if returnMessage.Payload.(messages.AgentControl).Command == "kill" {
				err := agents.RemoveAgent(agentID)
				if err != nil {
					log.Warnf(err.Error())
					return
				}
				log.Infof(fmt.Sprintf("Agent %s was removed from the server", agentID.String()))
				return
			}
		}

	} else if r.Method == "GET" {
		w.WriteHeader(404)
	} else {
		w.WriteHeader(404)
	}
	log.Debugf("Leaving http2.agentHandler function without error")
}

// getJWT returns a JSON Web Token for the provided agent using the interface JWT Key
func getJWT(agentID uuid.UUID, key []byte, s *Server) (string, error) {
	// Test context logger
	log := s.log.WithFields(logrus.Fields{"component": "server", "workspaceId": s.WorkspaceID})

	log.Debugf("Entering into agents.GetJWT function")

	encrypter, encErr := jose.NewEncrypter(jose.A256GCM,
		jose.Recipient{
			Algorithm: jose.DIRECT,
			Key:       key},
		(&jose.EncrypterOptions{}).WithType("JWT").WithContentType("JWT"))
	if encErr != nil {
		return "", fmt.Errorf("there was an error creating the JWE encryptor:\r\n%s", encErr)
	}

	signer, errSigner := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.HS256,
		Key:       key},
		(&jose.SignerOptions{}).WithType("JWT"))
	if errSigner != nil {
		return "", fmt.Errorf("there was an error creating the JWT signer:\r\n%s", errSigner.Error())
	}

	lifetime, errLifetime := agents.GetLifetime(agentID)
	if errLifetime != nil && errLifetime.Error() != "agent WaitTime is equal to zero" {
		return "", errLifetime
	}

	// This is for when the server hasn't received an AgentInfo struct and doesn't know the agent's lifetime yet or sleep is set to zero
	if lifetime == 0 {
		lifetime = time.Second * 30
	}

	// TODO Add in the rest of the JWT claim info
	cl := jwt.Claims{
		ID:        agentID.String(),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Expiry:    jwt.NewNumericDate(time.Now().Add(lifetime)),
	}

	agentJWT, err := jwt.SignedAndEncrypted(signer, encrypter).Claims(cl).CompactSerialize()
	if err != nil {
		return "", fmt.Errorf("there was an error serializing the JWT:\r\n%s", err.Error())
	}

	// Parse it to check for errors
	_, errParse := jwt.ParseEncrypted(agentJWT)
	if errParse != nil {
		return "", fmt.Errorf("there was an error parsing the encrypted JWT:\r\n%s", errParse.Error())
	}
	log.Infof(fmt.Sprintf("Created authenticated JWT for %s", agentID))
	log.Debugf(fmt.Sprintf("Sending agent %s an authenticated JWT with a lifetime of %v:\r\n%v",
		agentID.String(), lifetime, agentJWT))

	return agentJWT, nil
}

// validateJWT validates the provided JSON Web Token
func validateJWT(agentJWT string, key []byte, s *Server) (uuid.UUID, error) {
	// Test context logger
	log := s.log.WithFields(logrus.Fields{"component": "server", "workspaceId": s.WorkspaceID})

	var agentID uuid.UUID
	log.Debugf("Entering into http2.ValidateJWT")
	log.Debugf(fmt.Sprintf("Input JWT: %v", agentJWT))

	claims := jwt.Claims{}

	// Parse to make sure it is a valid JWT
	nestedToken, err := jwt.ParseSignedAndEncrypted(agentJWT)
	if err != nil {
		return agentID, fmt.Errorf("there was an error parsing the JWT:\r\n%s", err.Error())
	}

	// Decrypt JWT
	token, errToken := nestedToken.Decrypt(key)
	if errToken != nil {
		return agentID, fmt.Errorf("there was an error decrypting the JWT:\r\n%s", errToken.Error())
	}

	// Deserialize the claims and validate the signature
	errClaims := token.Claims(key, &claims)
	if errClaims != nil {
		return agentID, fmt.Errorf("there was an deserializing the JWT claims:\r\n%s", errClaims.Error())
	}

	agentID = uuid.FromStringOrNil(claims.ID)

	AgentWaitTime, errWait := agents.GetAgentFieldValue(agentID, "WaitTime")
	// An error will be returned during OPAQUE registration & authentication
	if errWait != nil {
		log.Debugf(fmt.Sprintf("there was an error getting the agent's wait time:\r\n%s", errWait.Error()))
	}
	if AgentWaitTime == "" {
		AgentWaitTime = "10s"
	}

	WaitTime, errParse := time.ParseDuration(AgentWaitTime)
	if errParse != nil {
		return agentID, fmt.Errorf("there was an error parsing the agent's wait time into a duration:\r\n%s", errParse.Error())
	}
	// Validate claims; Default Leeway is 1 minute; Set it to 1x the agent's WaitTime setting
	errValidate := claims.ValidateWithLeeway(jwt.Expected{
		Time: time.Now(),
	}, WaitTime)

	if errValidate != nil {
		log.Warnf(fmt.Sprintf("The JWT claims were not valid for %s", agentID))
		log.Infof(fmt.Sprintf("JWT Claim Expiry: %s", claims.Expiry.Time()))
		log.Infof(fmt.Sprintf("JWT Claim Issued: %s", claims.IssuedAt.Time()))
		return agentID, errValidate
	}
	log.Debugf(fmt.Sprintf("agentID: %s", agentID.String()))
	log.Debugf("Leaving http2.ValidateJWT without error")
	// TODO I need to validate other things like token age/expiry
	return agentID, nil
}

// decryptJWE takes provided JWE string and decrypts it using the per-agent key
func decryptJWE(jweString string, key []byte, s *Server) (messages.Base, error) {
	// Test context logger
	log := s.log.WithFields(logrus.Fields{"component": "server", "workspaceId": s.WorkspaceID})

	log.Debugf("Entering into http2.DecryptJWE function")
	log.Debugf(fmt.Sprintf("Input JWE String: %s", jweString))

	var m messages.Base

	// Parse JWE string back into JSONWebEncryption
	jwe, errObject := jose.ParseEncrypted(jweString)
	if errObject != nil {
		return m, fmt.Errorf("there was an error parseing the JWE string into a JSONWebEncryption object:\r\n%s", errObject)
	}

	log.Debugf(fmt.Sprintf("Parsed JWE:\r\n%+v", jwe))

	// Decrypt the JWE
	jweMessage, errDecrypt := jwe.Decrypt(key)
	if errDecrypt != nil {
		return m, fmt.Errorf("there was an error decrypting the JWE:\r\n%s", errDecrypt.Error())
	}

	// Decode the JWE payload into a messages.Base struct
	errDecode := gob.NewDecoder(bytes.NewReader(jweMessage)).Decode(&m)
	if errDecode != nil {
		return m, fmt.Errorf("there was an error decoding JWE payload message sent by an agent:\r\n%s", errDecode.Error())
	}

	log.Debugf("Leaving http2.DecryptJWE function without error")
	log.Debugf(fmt.Sprintf("Returning message base: %+v", m))
	return m, nil
}
