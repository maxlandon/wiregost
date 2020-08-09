package module

import (
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"

	pb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
)

// Option - A module option, with methods for processing its content and returning it in an optimal format,
// such as a range of net.IPs if we want it. Note that all option values are strings: you can process
// or transform them as you wish with Go's std/community libraries.
// However, the option type provides some utility methods for processing its content and
// returning it in an optimal format, such as an int, a file object, a range of net.IPs, etc.
type Option interface {
	Name() string
	Category() string
	Value() string
	Description() string
	Required() bool
	Set(string) error
	// toProtobuf - Unexported, so that No types in any other
	// package can implement the Option interface. Included because
	// it is always at the module.Module level that we pack protobuf
	// information for sending to sharing with other components.
	toProtobuf() *pb.Option
}

// option - The concrete type if unexported because we do not want
// any module to have such control on what an option may do.
type option struct {
	i *pb.Option // This is the base information.
}

// Name - Implements the Option interface
func (o *option) Name() string {
	return o.i.Name
}

// Category - Implements the Option interface
func (o *option) Category() string {
	return o.i.Category
}

// Value - Implements the Option interface
func (o *option) Value() string {
	return o.i.Value
}

// Description - Implements the Option interface
func (o *option) Description() string {
	return o.i.Description
}

// Required - Implements the Option interface
func (o *option) Required() bool {
	return o.i.Required
}

// Required - Implements the Option interface
func (o *option) Set(value string) (err error) {
	if value == "" && o.i.Required {
		return fmt.Errorf("Cannot set %s to empty: option is required.", o.i.Name)
	}
	o.i.Value = value
	return
}

// AddOption - Instantiates an option for this module, populating it with given parameters.
// name           - The option name (cannot be empty)
// category       - The option category name (can be empty), used for pretty printing on consoles.
// value          - If non-empty, this will serve a default value. This field is used for the option value.
// description    - Option description
// required       - Is this option required ? If true, the Value field cannot be empty.
//
// This function is called twice: on server and on stack binary. The server should call it
// for every entry in the result of ToProtobuf(), called by the stack binary.
func (m *Module) AddOption(name, category, value, description string, required bool) (err error) {

	// Invalid cases: Name is nil, and option is required but has not been provided a default value
	if name == "" || (value == "" && required) {
		return errors.New("parsed an option with no name, or required but with no default value")
	}

	// Populate
	opt := &option{
		&pb.Option{
			Name:        name,
			Category:    category,
			Value:       value,
			Description: description,
			Required:    required,
		},
	}

	// Add option to list
	m.Opts[name] = opt

	return
}

// Option - Returns one of the module's options, by name, with an error. This function is a safer but less handy
// alternative to the Options["Name"] way of querying options, because you can check immediately for key presence
func (m *Module) Option(name string) (opt Option, err error) {
	if opt, found := m.Opts[name]; found {
		return opt, nil
	}

	return nil, fmt.Errorf("module asked for an unknown option: %s", name)
}

func (o *option) toProtobuf() (opt *pb.Option) {

	return &pb.Option{
		Name:        o.i.Name,
		Category:    o.i.Category,
		Value:       o.i.Value,
		Description: o.i.Description,
		Required:    o.i.Required,
	}

}

// String - Implements the string interface, and returns the option value as a string.
func (o *option) String() string {
	return o.i.Value
}

// Bool - Returns the option value as a boolean, with an optional error, just in case
func (o *option) Bool() (bool, error) {
	// Allow for both expressions, and it must be required because its bool
	if (o.i.Value == "true" || o.i.Value == "yes") && o.i.Required {
		return true, nil
	}
	return false, fmt.Errorf("option %s is not a boolean, or uncorrectly set (should not happen)", o.Name)
}

// Int - Returns the option value as an int
func (o *option) Int() (int, error) {
	port, err := strconv.Atoi(o.i.Value)
	if err != nil {
		return -1, err
	}
	return port, nil
}

// Regexp - Returns a Regexp object, enabling to deal with various regular expression problems
func (o *option) Regexp() (reg *regexp.Regexp, err error) {
	reg, err = regexp.Compile(o.i.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to parse option %s as regexp: %s", o.Name, err)
	}
	return
}

// File - Returns a file object. We can then get/set its properties and read/write it.
func (o *option) File() (file *os.File) {
	return
}

// IP - Parses a string into a net.IP object, supporting both IPv4 and IPv6. This function is not redundant with option.ToAddressRange():
// both provide different objects for different uses, and we assume the module author will know which of these parsing function to use.
func (o *option) IP() (net.IP, error) {
	ip := net.ParseIP(o.i.Value)
	if ip.String() == "<nil>" {
		return nil, fmt.Errorf("failed to parse IP address for option %s ", o.Name)
	}
	return ip, nil
}

// ToAddressRange - Parses an option
func (o *option) ToAddressRange() {

}
