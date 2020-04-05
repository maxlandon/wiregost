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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/evilsocket/islazy/fs"
	"github.com/lmorg/readline"
	"github.com/maxlandon/wiregost/client/util"
)

func completeLocalPath(line []rune, pos int) (string, []string, map[string]string, readline.TabDisplayType) {

	// Completions
	var suggestions []string
	listSuggestions := map[string]string{}

	// Get last path
	splitLine := strings.Split(string(line), " ")
	last := splitLine[len(splitLine)-1]

	// Any parsing error is silently ignored, for not messing the prompt
	processedPath, _ := util.ParseEnvVariables([]string{last})

	// Check if processed input is empty
	var inputPath string
	if len(processedPath) == 1 {
		inputPath = processedPath[0]
	}

	// Add a slash if the raw input has one but not the processed input
	if line[len(line)-1] == '/' {
		inputPath += "/"
	}

	// linePath is the curated version of the inputPath
	var linePath string
	// absPath is the absolute path (excluding suffix) of the inputPath
	var absPath string
	// lastPath is the last directory in the input path
	var lastPath string

	if strings.HasSuffix(string(inputPath), "/") {
		// Trim the non needed slash
		linePath = strings.TrimSuffix(string(inputPath), "/")
		linePath = filepath.Dir(string(inputPath))
		// Get absolute path
		absPath, _ = fs.Expand(string(linePath))

	} else if string(inputPath) == "" {
		linePath = "."
		absPath, _ = fs.Expand(string(linePath))
	} else {
		linePath = string(inputPath)
		linePath = filepath.Dir(string(inputPath))
		// Get absolute path
		absPath, _ = fs.Expand(string(linePath))
		// Save filter
		lastPath = filepath.Base(string(inputPath))
	}

	// 2) We take the absolute path we found, and get all dirs in it.
	var dirs []string
	files, _ := ioutil.ReadDir(absPath)
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		}
	}

	switch lastPath {
	case "":
		for _, dir := range dirs {
			if strings.HasPrefix(dir, lastPath) || lastPath == dir {
				suggestions = append(suggestions, dir[len(lastPath):]+"/")
			}
		}
	default:
		filtered := []string{}
		for _, dir := range dirs {
			if strings.HasPrefix(dir, lastPath) {
				filtered = append(filtered, dir)
			}
		}

		for _, dir := range filtered {
			if !hasPrefix([]rune(lastPath), []rune(dir)) || lastPath == dir {
				suggestions = append(suggestions, dir[len(lastPath):]+"/")
			}
		}

	}

	return string(lastPath), suggestions, listSuggestions, readline.TabDisplayGrid
}

func completeLocalPathAndFiles(line []rune, pos int) (string, []string, map[string]string, readline.TabDisplayType) {

	// Completions
	var suggestions []string
	listSuggestions := map[string]string{}

	// Get last path
	splitLine := strings.Split(string(line), " ")
	last := splitLine[len(splitLine)-1]

	// Any parsing error is silently ignored, for not messing the prompt
	processedPath, _ := util.ParseEnvVariables([]string{last})

	// Check if processed input is empty
	var inputPath string
	if len(processedPath) == 1 {
		inputPath = processedPath[0]
	}

	// Add a slash if the raw input has one but not the processed input
	if line[len(line)-1] == '/' {
		inputPath += "/"
	}

	// linePath is the curated version of the inputPath
	var linePath string
	// absPath is the absolute path (excluding suffix) of the inputPath
	var absPath string
	// lastPath is the last directory in the input path
	var lastPath string

	if strings.HasSuffix(string(inputPath), "/") {
		// Trim the non needed slash
		linePath = strings.TrimSuffix(string(inputPath), "/")
		linePath = filepath.Dir(string(inputPath))
		// Get absolute path
		absPath, _ = fs.Expand(string(linePath))

	} else if string(inputPath) == "" {
		linePath = "."
		absPath, _ = fs.Expand(string(linePath))
	} else {
		linePath = string(inputPath)
		linePath = filepath.Dir(string(inputPath))
		// Get absolute path
		absPath, _ = fs.Expand(string(linePath))
		// Save filter
		lastPath = filepath.Base(string(inputPath))
	}

	// 2) We take the absolute path we found, and get all dirs in it.
	var dirs []string
	files, _ := ioutil.ReadDir(absPath)
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		}
	}

	switch lastPath {
	case "":
		for _, file := range files {
			if strings.HasPrefix(file.Name(), lastPath) || lastPath == file.Name() {
				if file.IsDir() {
					suggestions = append(suggestions, file.Name()[len(lastPath):]+"/")
				} else {
					suggestions = append(suggestions, file.Name()[len(lastPath):]+" ")
				}
			}
		}
	default:
		filtered := []os.FileInfo{}
		for _, file := range files {
			if strings.HasPrefix(file.Name(), lastPath) {
				filtered = append(filtered, file)
			}
		}

		for _, file := range filtered {
			if !hasPrefix([]rune(lastPath), []rune(file.Name())) || lastPath == file.Name() {
				if file.IsDir() {
					suggestions = append(suggestions, file.Name()[len(lastPath):]+"/")
				} else {
					suggestions = append(suggestions, file.Name()[len(lastPath):]+" ")
				}
			}
		}

	}

	return string(lastPath), suggestions, listSuggestions, readline.TabDisplayGrid
}
