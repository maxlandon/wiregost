## Certs

X.509 certificate generation and management code. 
4 separate certificate chains are used (4 CAs):

* `GhostCA`     - Used to encrypt and authenticate client-side C2 channels between the Server and the Ghosts.
                  Uses both ECC and RSA certificates, depending on the use case.
* `OperatorCA`  - Used to sign certs that authenticate and encrypt the mutual TLS connection between an operator and the server.
* `ServerCA`    - Used to secure server-side C2, the ServerCA public key is embedded into the Ghost binaries.
* `HTTPSCA`     - Used to generate self-signed HTTPS certificates (that are not used to encrypt C2 data).

----
Certificates are all stored CA-specific Badger databases managed by the db package.
* Key:      Common Name of the certificate
* Value:    JSON object (CertificateKeyPair) containing key type (RSA or ECC), certificate and private key


#### ACME

The package can also interact with Let's Encrypt (ACME) services to generate certificates that are trusted in the browser
(alternative to HTTPSCA). These certificates are used with the HTTPS servers/listeners, but not used to encrypt any C2.

----

## Package Structure 

* `ca.go`           - Manage Certificate Authorities (generate/get/save)
* `certs.go`        - CertificateKeyPair objects, with its generic functions (save/get/remove/generate).
* `certs_test.go`   - Testing for generic CertificateKeyPair objects.
* `acme.go`         - Helper functions for ACME certs.
* `https.go`        - Generate HTTPS certs.
* `operators.go`    - Manage operator/client certificates (generate/get/save/list).
* `servers.go`      - Generate Server certs.
* `ghosts.go`       - Generate Ghost certs. 
