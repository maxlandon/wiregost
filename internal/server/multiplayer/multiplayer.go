package multiplayer

// Wiregost - Post-Exploitation & Implant Framework
// Copyright © 2020 Para
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"

	consts "github.com/maxlandon/wiregost/internal/client/constants"
	"github.com/maxlandon/wiregost/internal/server/certs"
	"github.com/maxlandon/wiregost/internal/server/configs"
	"github.com/maxlandon/wiregost/internal/server/core"
	"github.com/maxlandon/wiregost/internal/server/db"
	"github.com/maxlandon/wiregost/internal/server/db/models"
	"github.com/maxlandon/wiregost/internal/server/transport"
)

var namePattern = regexp.MustCompile("^[a-zA-Z0-9_-]*$") // Only allow alphanumeric chars

// ClientConfig - Client JSON config
type ClientConfig struct {
	Operator      string `json:"operator"`
	Token         string `json:"token"`
	LHost         string `json:"lhost"`
	LPort         int    `json:"lport"`
	CACertificate string `json:"ca_certificate"`
	PrivateKey    string `json:"private_key"`
	Certificate   string `json:"certificate"`
}

// NewOperatorConfig - Generate a new player/client/operator configuration
func NewOperatorConfig(operatorName string, lhost string, lport uint16) ([]byte, error) {
	if !namePattern.MatchString(operatorName) {
		return nil, errors.New("invalid operator name (alphanumerics only)")
	}
	if operatorName == "" {
		return nil, errors.New("operator name required")
	}
	if lhost == "" {
		return nil, errors.New("invalid lhost")
	}

	rawToken := models.GenerateOperatorToken()
	digest := sha256.Sum256([]byte(rawToken))
	dbOperator := &models.Operator{
		Name:  operatorName,
		Token: hex.EncodeToString(digest[:]),
	}
	err := db.Session().Save(dbOperator).Error
	if err != nil {
		return nil, err
	}

	publicKey, privateKey, err := certs.OperatorClientGenerateCertificate(operatorName)
	if err != nil {
		return nil, fmt.Errorf("failed to generate certificate %s", err)
	}
	caCertPEM, _, _ := certs.GetCertificateAuthorityPEM(certs.OperatorCA)
	config := ClientConfig{
		Operator:      operatorName,
		Token:         rawToken,
		LHost:         lhost,
		LPort:         int(lport),
		CACertificate: string(caCertPEM),
		PrivateKey:    string(privateKey),
		Certificate:   string(publicKey),
	}
	return json.Marshal(config)
}

// StartPersistentJobs starts all multiplayer listeners.
func StartPersistentJobs(cfg *configs.ServerConfig) error {
	if cfg.Jobs == nil {
		return nil
	}
	for _, j := range cfg.Jobs.Multiplayer {
		StartClientListener(j.Host, j.Port)
	}
	return nil
}

// StartClientListener starts serving the wiregost RPC.
func StartClientListener(host string, port uint16) (int, error) {
	_, ln, err := transport.StartClientListener(host, port)
	if err != nil {
		return -1, err // If we fail to bind don't setup the Job
	}

	job := &core.Job{
		ID:          core.NextJobID(),
		Name:        "grpc",
		Description: "client listener",
		Protocol:    "tcp",
		Port:        port,
		JobCtrl:     make(chan bool),
	}

	go func() {
		<-job.JobCtrl
		log.Printf("Stopping client listener (%d) ...\n", job.ID)
		ln.Close() // Kills listener GoRoutines in startMutualTLSListener() but NOT connections

		core.Jobs.Remove(job)
		core.EventBroker.Publish(core.Event{
			Job:       job,
			EventType: consts.JobStoppedEvent,
		})
	}()

	core.Jobs.Add(job)
	return job.ID, nil
}
