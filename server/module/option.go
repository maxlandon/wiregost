package module

import (
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"

	modulepb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
)

// Option - A module option, with many methods for processing its content and
// returning it in an optimal format, such as a range of net.IPs if we want it.
// Note that all option values are strings: you can process/transform them as
// you wish with Go's std/community libraries.
// However, the Option type provides some utility methods for processing its content
// and returning it in an optimal format, such as an int, a file object, a range of net.IPs, etc...
type Option struct {
	info *modulepb.Option // This is the base information.
}

// AddOption - Instantiates an option for this module, populating it with given parameters.
// name           - The option name (cannot be empty)
// category       - The option category name (can be empty), used for pretty printing on consoles.
// value          - If non-empty, this will serve a default value. This field is used for the option value.
// description    - Option description
// required       - Is this option required ? If true, the Value field cannot be empty.
func (m *module) AddOption(name, category, value, description string, required bool) (err error) {

	// Invalid cases: Name is nil, and option is required but has not been provided a default value
	if name == "" || (value == "" && required) {
		return errors.New("parsed an option with no name, or required but with no default value")
	}

	// Populate
	opt := &Option{
		&modulepb.Option{
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
func (m *module) Option(name string) (opt *Option, err error) {
	if opt, found := m.Opts[name]; found {
		return opt, nil
	}

	return nil, fmt.Errorf("module asked for an unknown option: %s", name)
}

// String - Implements the string interface, and returns the option value as a string.
func (o *Option) String() string {
	return o.info.Value
}

// Bool - Returns the option value as a boolean, with an optional error, just in case
func (o *Option) Bool() (bool, error) {
	// Allow for both expressions, and it must be required because its bool
	if (o.info.Value == "true" || o.info.Value == "yes") && o.info.Required {
		return true, nil
	}
	return false, fmt.Errorf("option %s is not a boolean, or uncorrectly set (should not happen)", o.info.Name)
}

// Int - Returns the option value as an int
func (o *Option) Int() (int, error) {
	port, err := strconv.Atoi(o.info.Value)
	if err != nil {
		return -1, err
	}
	return port, nil
}

// Regexp - Returns a Regexp object, enabling to deal with various regular expression problems
func (o *Option) Regexp() (reg *regexp.Regexp, err error) {
	reg, err = regexp.Compile(o.info.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to parse option %s as regexp: %s", o.info.Name, err)
	}
	return
}

// File - Returns a file object. We can then get/set its properties and read/write it.
func (o *Option) File() (file *os.File) {
	return
}

// IP - Parses a string into a net.IP object, supporting both IPv4 and IPv6. This function is not redundant with option.ToAddressRange():
// both provide different objects for different uses, and we assume the module author will know which of these parsing function to use.
func (o *Option) IP() (net.IP, error) {
	ip := net.ParseIP(o.info.Value)
	if ip.String() == "<nil>" {
		return nil, fmt.Errorf("failed to parse IP address for option %s ", o.info.Name)
	}
	return ip, nil
}

// ToAddressRange - Parses an option
func (o *Option) ToAddressRange() {

}
