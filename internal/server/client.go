package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"

	"github.com/maxlandon/wiregost/internal/db"
	"github.com/maxlandon/wiregost/internal/logging"
	"github.com/maxlandon/wiregost/internal/messages"
)

type Client struct {
	// Connection
	conn       net.Conn
	writer     *bufio.Writer
	reader     *bufio.Reader
	disconnect chan bool
	status     int
	// User-specific
	User               *db.User
	id                 int
	CurrentWorkspaceId int    // Will influence things like logging policy. CHANGE TO POINTER TO WORKSPACE
	Context            string // Will influence how commands are dispatched.
	// Message-specific
	requests  chan messages.ClientRequest
	responses chan messages.Message // Commands will always be sent as a list of strings

	// Client logger
	Logger *logging.ClientLogger
}

func CreateClient(conn net.Conn) *Client {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	client := &Client{
		conn:       conn,
		writer:     writer,
		responses:  make(chan messages.Message), // Generic message can pass anything needed
		reader:     reader,
		requests:   make(chan messages.ClientRequest),
		disconnect: make(chan bool),
		status:     1,
		// User-specific
		id: rand.Int(),
		// User: Add user
		Context: "main", // Default context is always main when a shell is spawned
		// Initilialize Logger
		Logger: new(logging.ClientLogger),
	}
	// Setup logger
	client.Logger.ClientId = client.id

	go client.Write()
	go client.Read()

	return client
}

func (client *Client) Write() {
	for {
		select {
		case <-client.disconnect:
			client.status = 0
		default:
			msg := <-client.responses
			enc := json.NewEncoder(client.writer)
			err := enc.Encode(msg)
			if err != nil {
				fmt.Println(err.Error())
			}
			err = client.writer.Flush()
		}
	}
}

// Reads messages from the client
func (client *Client) Read() {
	for {
		// Decode request
		var message messages.ClientRequest
		dec := json.NewDecoder(client.reader)
		err := dec.Decode(&message)
		if err != nil {
			fmt.Println(err)
			client.status = 0
			client.disconnect <- true
			client.conn.Close()
			break
		}
		// Fill message with client information
		message.ClientId = client.id

		if message.CurrentWorkspaceId != 0 {
			// Fill client with message information
			client.CurrentWorkspaceId = message.CurrentWorkspaceId
			// Fill client Logger with message information
			client.Logger.CurrentWorkspaceId = message.CurrentWorkspaceId
		}

		// Forward message
		client.requests <- message
	}
}
