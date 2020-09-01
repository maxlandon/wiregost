package ghost

import transportpb "github.com/maxlandon/wiregost/proto/v1/gen/go/transport"

// AddTransport - Instructs the implant to add a transport to its list, if possible.
func (c *Client) AddTransport(t *transportpb.Transport) (err error) {

	return
}

// RemoveTransport - Instructs the implant to remove a transport from its list.
func (c *Client) RemoveTransport(t *transportpb.Transport) (err error) {
	return
}

// GetSupportedTransports - Get the list of all transport protocols supported by this implant
func (c *Client) GetSupportedTransports() (ts []transportpb.Transport) {
	return
}
