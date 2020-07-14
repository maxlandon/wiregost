package route

import (
	"net"

	"github.com/maxlandon/wiregost/server/c2/resolver"
)

// This file is the second part of chain.go in the same package.
// While chain.go declares the Chain object and various methods to populate it with
// route nodes and groups, this file provides the code for the chain to connect to these nodes.

// Dial - Connects to the target address addr through the chain.
// If the chain is empty, it will use the net.Dial directly.
func (c *Chain) Dial(addr string, opts ...ChainOption) (conn net.Conn, err error) {
	return
}

// DiaDialWithOptions - Sets options needed by the Chain and connects to to the target address.
func (c *Chain) DialWithOptions(addr string, options *ChainOptions) (conn net.Conn, err error) {
	return
}

// Resolve - Finds the correct address/hostname for a name.
// SIGNATURE TO MODIFY
// func (c *Chain) Resolve(addr string, resolver resolver.Resolver, hosts *Hosts) (address string) {
// This signature should be modified in case Wiregost needs custom host/address resolution
func (c *Chain) Resolve(addr string, resolver resolver.Resolver) (address string) {
	return
}

// Conn - Obtains a handshaked connection to the last node of the chain.
func (c *Chain) Conn(opts ...ChainOption) (conn net.Conn, err error) {

	// Retries

	// Select Route

	// Check empty chain

	// Get first node we want to contact in the chain

	// Dial & Handshake

	// For each remaining node, connect and handshake in loop

	// Return the last connection

	return
}

// SelSelectRoute - Selects a route with bypass testing, Wiregost permissions, etc.
func (c *Chain) SelectRoute() (route *Chain, err error) {

	// Check route is empty, return if yes

	// Initialize a list of nodes, and an empty route

	// For each group in the Chain groups
	for {
		// Check bypass address for each node in group

		// Check multiplex mode and cutoff the chain if yes

		// Add the node to route and add to Node list
	}

	// Set the chain route with the node list

	return
}
