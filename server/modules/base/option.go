package base

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	modulepb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
)

// type OptionType int
//
// const (
//         TypeDefault = iota
//         TypeInt
//         TypeUint
//         TypeString
//         TypeAddress
//         TypeURL
//         TypeBool
// )

// option - A module option, with many methods for processing its content
// and returning it in an optimal, such as a range of net.IPs if we want it.
type Option struct {
	proto modulepb.Option // This is the base information.
}

// NewOption - Instantiates an option for this module, populating it with given parameters.
func (m *Module) NewOption(name, category, value, description string, required bool) (opt *Option, err error) {

	// Invalid cases: Name is nil, and option is required but has not been provided a default value
	if name == "" || (value == "" && required) {
		return nil, errors.New("parsed an option with no name, or required but with no default value")
	}

	// Populate
	opt = &Option{
		proto: modulepb.Option{
			Name:        name,
			Category:    category,
			Value:       value,
			Description: description,
			Required:    required,
		},
	}

	return
}

// Value - Returns the value of the option as a string
func (o *Option) Value() string {
	return o.proto.Value
}

// String - Implements the string interface
func (o *Option) String() string {
	return o.proto.Value
}

// Bool - Returns the option value as a boolean, with an optional error, just in case
func (o *Option) Bool() (bool, error) {
	// Allow for both expressions, and it must be required because its bool
	if (o.proto.Value == "true" || o.proto.Value == "yes") && o.proto.Required {
		return true, nil
	}
	return false, fmt.Errorf("option %s is not a boolean, or uncorrectly set (should not happen)", o.proto.Name)
}

// ToInt - Returns the option value as an int
func (o *Option) ToInt() (int, error) {
	port, err := strconv.Atoi(o.proto.Value)
	if err != nil {
		return -1, err
	}
	return port, nil
}

// ToIP - Parses a string into a net.IP object, supporting both IPv4 and IPv6.
// This function is not redundant with option.ToAddressRange():
// both provide different objects for different uses, and we assume the
// module author will know which of these parsing function to use.
func (o *Option) ToIP() (net.IP, error) {
	ip := net.ParseIP(o.proto.Value)
	if ip.String() == "<nil>" {
		return nil, fmt.Errorf("failed to parse IP address for option %s ", o.proto.Name)
	}
	return ip, nil
}

// ToAddressRange - Parses an option
func (o *Option) ToAddressRange() {

}

func (o *Option) Path() *os.File {

}
