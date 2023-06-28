package cli

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
	"fmt"
	"sort"

	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/maxlandon/wiregost/internal/client/assets"
)

func selectConfig() *assets.ClientConfig {
	configs := assets.GetConfigs()

	if len(configs) == 0 {
		return nil
	}

	if len(configs) == 1 {
		for _, config := range configs {
			return config
		}
	}

	answer := struct{ Config string }{}
	qs := getPromptForConfigs(configs)
	err := survey.Ask(qs, &answer)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return configs[answer.Config]
}

func getPromptForConfigs(configs map[string]*assets.ClientConfig) []*survey.Question {
	keys := []string{}
	for k := range configs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return []*survey.Question{
		{
			Name: "config",
			Prompt: &survey.Select{
				Message: "Select a server:",
				Options: keys,
				Default: keys[0],
			},
		},
	}
}
