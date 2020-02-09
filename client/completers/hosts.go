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

package completers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/data_service/remote"
)

// AutoCompleter is the autocompletion engine
type HostCompleter struct {
	Command *commands.Command
}

// Do is the completion function triggered at each line
func (hc *HostCompleter) Do(ctx *commands.ShellContext, line []rune, pos int) (options [][]rune, offset int) {

	// Complete command args
	splitLine := strings.Split(string(line), " ")
	line = trimSpaceLeft([]rune(splitLine[len(splitLine)-1]))

	switch splitLine[0] {
	case "delete":
		return hc.yieldHostValues(&ctx.Context, line, pos)
	case "add":
		return hc.yieldHostValues(&ctx.Context, line, pos)
	case "update":
		return hc.yieldHostValues(&ctx.Context, line, pos)
	case "search":
		return hc.yieldHostValues(&ctx.Context, line, pos)
	}

	return options, offset
}

func (hc *HostCompleter) yieldHostValues(ctx *context.Context, line []rune, pos int) (options [][]rune, offset int) {

	hosts, err := remote.Hosts(*ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, arg := range hc.Command.Args {
		search := arg.Name
		if !hasPrefix(line, []rune(search)) {
			sLine, sOffset := doInternal(line, pos, len(line), []rune(search+"="))
			options = append(options, sLine...)
			offset = sOffset
		} else {
			words := strings.Split(string(line), "=")
			argInput := lastString(words)

			// For some arguments, the split results in a last empty item.
			if words[len(words)-1] == "" {
				argInput = words[0]
			}

			// All boolean values
			if arg.Type == "boolean" {
				for _, search := range []string{"true ", "false "} {
					offset = 0
					if strings.HasPrefix(search, argInput) {
						options = append(options, []rune(search[len(argInput):]))
						offset = len(argInput)
					}
				}
				return
			}

			// Host ID
			if arg.Name == "host-id" {
				for _, h := range hosts {
					id := strconv.FormatUint(uint64(h.ID), 10)
					options = append(options, []rune(id))
				}
				return
			}

			// Host OS Properties
			if arg.Name == "os-name" {
				// osNames := []string{}
				for _, h := range hosts {
					options = append(options, []rune(h.OSName))
					// options = append(options, []rune(h.OSName))
					// osNames = append(osNames, h.OSName)
				}
				// for _, os := range osNames {
				//         offset = 0
				//         if strings.HasPrefix(os, argInput) {
				//                 options = append(options, []rune("\""+os[len(argInput):]+"\""+" "))
				//                 offset = len(argInput)
				//         }
				// }
				return
			}

			// Host Addresses
			if arg.Name == "addresses" {
				addrs := []string{}
				for _, h := range hosts {
					hAddrs := []string{}
					for _, a := range h.Addresses {
						hAddrs = append(hAddrs, a.String())
					}
					addrs = append(addrs, strings.Join(hAddrs, ","))
				}
				for _, addr := range addrs {
					offset = 0
					if strings.HasPrefix(addr, argInput) {
						options = append(options, []rune(addr[len(argInput):]+" "))
						offset = len(argInput)
					}
				}
				return
			}

			// Host names, users,
			if arg.Name == "hostnames" {
				names := []string{}
				for _, h := range hosts {
					hNames := []string{}
					for _, hn := range h.Hostnames {
						hNames = append(hNames, hn.Name)
					}
					names = append(names, strings.Join(hNames, ","))
				}
				for _, hn := range names {
					offset = 0
					if strings.HasPrefix(hn, argInput) {
						options = append(options, []rune(hn[len(argInput):]+" "))
						offset = len(argInput)
					}
				}
				return
			}
		}

	}
	return options, offset
}
