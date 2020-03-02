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

package generate

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	// clientpb "github.com/maxlandon/wiregost/protobuf/client"
)

// ParseMTLSc2 - MTLS to protobuf
// func ParseMTLSc2(args string) []*clientpb.GhostC2 {
//         c2s := []*clientpb.GhostC2{}
//         if args == "" {
//                 return c2s
//         }
//         for index, arg := range strings.Split(args, ",") {
//                 uri := url.URL{Scheme: "mtls"}
//                 uri.Host = arg
//                 if uri.Port() == "" {
//                         uri.Host = fmt.Sprintf("%s:%d", uri.Host, DefaultMTLSLPort)
//                 }
//                 c2s = append(c2s, &clientpb.GhostC2{
//                         Priority: uint32(index),
//                         URL:      uri.String(),
//                 })
//         }
//         return c2s
// }

// ParseMTLSc2 - MTLS connection strings
func ParseMTLSc2(args string) []GhostC2 {
	c2s := []GhostC2{}
	if args == "" {
		return c2s
	}
	for index, arg := range strings.Split(args, ",") {
		uri := url.URL{Scheme: "mtls"}
		uri.Host = arg
		if uri.Port() == "" {
			uri.Host = fmt.Sprintf("%s:%d", uri.Host, DefaultMTLSLPort)
		}
		c2s = append(c2s, GhostC2{
			Priority: uint32(index),
			URL:      uri.String(),
		})
	}
	return c2s
}

// func ParseHTTPc2(args string) []*clientpb.GhostC2 {
//         c2s := []*clientpb.GhostC2{}
//         if args == "" {
//                 return c2s
//         }
//         for index, arg := range strings.Split(args, ",") {
//                 arg = strings.ToLower(arg)
//                 var uri *url.URL
//                 var err error
//                 if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
//                         uri, err = url.Parse(arg)
//                         if err != nil {
//                                 log.Printf("Failed to parse c2 URL %v", err)
//                                 continue
//                         }
//                 } else {
//                         uri = &url.URL{Scheme: "https"} // HTTPS is the default, will fallback to HTTP
//                         uri.Host = arg
//                 }
//                 c2s = append(c2s, &clientpb.GhostC2{
//                         Priority: uint32(index),
//                         URL:      uri.String(),
//                 })
//         }
//         return c2s
// }

// ParseHTTPc2 - Parse HTTP connection strings
func ParseHTTPc2(args string) []GhostC2 {
	c2s := []GhostC2{}
	if args == "" {
		return c2s
	}
	for index, arg := range strings.Split(args, ",") {
		arg = strings.ToLower(arg)
		var uri *url.URL
		var err error
		if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
			uri, err = url.Parse(arg)
			if err != nil {
				log.Printf("Failed to parse c2 URL %v", err)
				continue
			}
		} else {
			uri = &url.URL{Scheme: "https"} // HTTPS is the default, will fallback to HTTP
			uri.Host = arg
		}
		c2s = append(c2s, GhostC2{
			Priority: uint32(index),
			URL:      uri.String(),
		})
	}
	return c2s
}

// func ParseDNSc2(args string) []*clientpb.GhostC2 {
//         c2s := []*clientpb.GhostC2{}
//         if args == "" {
//                 return c2s
//         }
//         for index, arg := range strings.Split(args, ",") {
//                 uri := url.URL{Scheme: "dns"}
//                 if len(arg) < 1 {
//                         continue
//                 }
//                 // Make sure we have the FQDN
//                 if !strings.HasSuffix(arg, ".") {
//                         arg += "."
//                 }
//                 if strings.HasPrefix(arg, ".") {
//                         arg = arg[1:]
//                 }
//
//                 uri.Host = arg
//                 c2s = append(c2s, &clientpb.GhostC2{
//                         Priority: uint32(index),
//                         URL:      uri.String(),
//                 })
//         }
//         return c2s
// }

// ParseDNSc2 - Parse DNS domains from string
func ParseDNSc2(args string) []GhostC2 {
	c2s := []GhostC2{}
	if args == "" {
		return c2s
	}
	for index, arg := range strings.Split(args, ",") {
		uri := url.URL{Scheme: "dns"}
		if len(arg) < 1 {
			continue
		}
		// Make sure we have the FQDN
		if !strings.HasSuffix(arg, ".") {
			arg += "."
		}
		if strings.HasPrefix(arg, ".") {
			arg = arg[1:]
		}

		uri.Host = arg
		c2s = append(c2s, GhostC2{
			Priority: uint32(index),
			URL:      uri.String(),
		})
	}
	return c2s
}
