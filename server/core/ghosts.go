// Wiregost - Golang Exploitation Framework
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

package core

import (
	"errors"
	"sync"
	"time"

	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
	"github.com/sirupsen/logrus"
)

var (
	// Wire - Stores all ghost implant sessions
	Wire = &GhostWire{
		Ghosts: &map[uint32]*Ghost{},
		mutex:  &sync.RWMutex{},
	}

	wireID = new(uint32)
)

// Ghost Implant
type Ghost struct {
	ID            uint32
	Name          string
	Hostname      string
	Username      string
	UID           string
	GID           string
	OS            string
	Version       string
	Arch          string
	Transport     string
	RemoteAddress string
	PID           int32
	Filename      string
	LastCheckin   *time.Time
	Send          chan *ghostpb.Envelope
	Resp          map[uint64]chan *ghostpb.Envelope
	RespMutex     *sync.RWMutex
	ActiveC2      string

	// Added
	WorkspaceID uint
	HostID      uint

	// Logging
	Logger *logrus.Entry
}

// ToProtobuf - Returns the protobuf version of a Ghost object
func (g *Ghost) ToProtobuf() *clientpb.Ghost {
	var lastCheckin string
	if g.LastCheckin == nil {
		lastCheckin = time.Now().Format(time.RFC1123) // Stateful connections have a nil .LastCheckin
	} else {
		lastCheckin = g.LastCheckin.Format(time.RFC1123)
	}

	return &clientpb.Ghost{
		ID:            uint32(g.ID),
		Name:          g.Name,
		Hostname:      g.Hostname,
		Username:      g.Username,
		UID:           g.UID,
		GID:           g.GID,
		OS:            g.OS,
		Version:       g.Version,
		Arch:          g.Arch,
		Transport:     g.Transport,
		RemoteAddress: g.RemoteAddress,
		PID:           int32(g.PID),
		Filename:      g.Filename,
		LastCheckin:   lastCheckin,
		ActiveC2:      g.ActiveC2,

		// Added
		WorkspaceID: uint32(g.WorkspaceID),
		HostID:      uint32(g.HostID),
	}
}

// Config - Get the config the Ghost was generated with
func (g *Ghost) Config() error {

	return nil
}

// Request - Sends a protobuf request to an active Ghost and returns the response
func (g *Ghost) Request(msgType uint32, timeout time.Duration, data []byte) ([]byte, error) {

	resp := make(chan *ghostpb.Envelope)
	reqID := EnvelopeID()
	g.RespMutex.Lock()
	g.Resp[reqID] = resp
	g.RespMutex.Unlock()

	defer func() {
		g.RespMutex.Lock()
		defer g.RespMutex.Unlock()
		delete(g.Resp, reqID)
	}()

	g.Send <- &ghostpb.Envelope{
		ID:   reqID,
		Type: msgType,
		Data: data,
	}

	var respEnvelope *ghostpb.Envelope

	select {
	case respEnvelope = <-resp:
	case <-time.After(timeout):
		return nil, errors.New("timeout")
	}

	return respEnvelope.Data, nil
}

// GhostWire - Manages the Ghosts, provides atomic access
type GhostWire struct {
	mutex  *sync.RWMutex
	Ghosts *map[uint32]*Ghost
}

// Ghost - Get Ghost by ID
func (gw *GhostWire) Ghost(ghostID uint32) *Ghost {
	gw.mutex.Lock()
	defer gw.mutex.Unlock()

	return (*gw.Ghosts)[ghostID]
}

// AddGhost - Adds a ghost to the GhostWire, atomically
func (gw *GhostWire) AddGhost(ghost *Ghost) {
	gw.mutex.Lock()
	defer gw.mutex.Unlock()
	(*gw.Ghosts)[ghost.ID] = ghost
}

// RemoveGhost - Add a ghost to the hive (atomically)
func (gw *GhostWire) RemoveGhost(ghost *Ghost) {
	gw.mutex.Lock()
	defer gw.mutex.Unlock()
	delete((*gw.Ghosts), ghost.ID)
}

// GetGhostWireID - Returns an incremental nonce as an ID
func GetGhostWireID() uint32 {
	newID := (*wireID) + 1
	(*wireID)++
	return newID
}
