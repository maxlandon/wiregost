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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/maxlandon/wiregost/server/assets"
	"github.com/maxlandon/wiregost/server/db"
	"github.com/maxlandon/wiregost/server/log"
)

const (
	// ghostBucketName - Name of the bucket that stores data related to ghosts
	ghostBucketName = "ghosts"

	// ghostConfigNamespace - Namespace that contains ghosts configs
	ghostConfigNamespace   = "config"
	ghostFileNamespace     = "file"
	ghostDatetimeNamespace = "datetime"
)

var (
	storageLog = log.ServerLogger("generate", "storage")

	// ErrGhostNotFound - More descriptive 'key not found' error
	ErrGhostNotFound = errors.New("Ghost not found")
)

// GhostConfigByName - Get a ghost's config by it's codename
func GhostConfigByName(name string) (*GhostConfig, error) {
	bucket, err := db.GetBucket(ghostBucketName)
	if err != nil {
		return nil, err
	}
	rawConfig, err := bucket.Get(fmt.Sprintf("%s.%s", ghostConfigNamespace, name))
	if err != nil {
		return nil, err
	}
	config := &GhostConfig{}
	err = json.Unmarshal(rawConfig, config)
	return config, err
}

// GhostConfigMap - Get a ghost's config by it's codename
func GhostConfigMap() (map[string]*GhostConfig, error) {
	bucket, err := db.GetBucket(ghostBucketName)
	if err != nil {
		return nil, err
	}
	ls, err := bucket.List(ghostConfigNamespace)
	configs := map[string]*GhostConfig{}
	for _, config := range ls {
		ghostName := config[len(ghostConfigNamespace)+1:]
		config, err := GhostConfigByName(ghostName)
		if err != nil {
			continue
		}
		configs[ghostName] = config
	}
	return configs, nil
}

// GhostConfigSave - Save a configuration to the database
func GhostConfigSave(config *GhostConfig) error {
	bucket, err := db.GetBucket(ghostBucketName)
	if err != nil {
		return err
	}
	rawConfig, err := json.Marshal(config)
	if err != nil {
		return err
	}
	storageLog.Infof("Saved config for '%s'", config.Name)
	return bucket.Set(fmt.Sprintf("%s.%s", ghostConfigNamespace, config.Name), rawConfig)
}

// GhostFileSave - Saves a binary file into the database
func GhostFileSave(name, fpath string) error {
	bucket, err := db.GetBucket(ghostBucketName)
	if err != nil {
		return err
	}

	rootAppDir, _ := filepath.Abs(assets.GetRootAppDir())
	fpath, _ = filepath.Abs(fpath)
	if !strings.HasPrefix(fpath, rootAppDir) {
		return fmt.Errorf("Invalid path '%s' is not a subdirectory of '%s'", fpath, rootAppDir)
	}

	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}
	storageLog.Infof("Saved '%s' file to database %d byte(s)", name, len(data))
	bucket.Set(fmt.Sprintf("%s.%s", ghostDatetimeNamespace, name), []byte(time.Now().Format(time.RFC1123)))
	return bucket.Set(fmt.Sprintf("%s.%s", ghostFileNamespace, name), data)
}

// GhostFileByName - Saves a binary file into the database
func GhostFileByName(name string) ([]byte, error) {
	bucket, err := db.GetBucket(ghostBucketName)
	if err != nil {
		return nil, err
	}
	ghost, err := bucket.Get(fmt.Sprintf("%s.%s", ghostFileNamespace, name))
	if err != nil {
		return nil, ErrGhostNotFound
	}
	return ghost, nil
}

// GhostFiles - List all ghost files
func GhostFiles() ([]string, error) {
	bucket, err := db.GetBucket(ghostBucketName)
	if err != nil {
		return nil, err
	}
	keys, err := bucket.List(ghostFileNamespace)
	if err != nil {
		return nil, err
	}

	// Remove namespace prefix
	names := []string{}
	for _, key := range keys {
		names = append(names, key[len(ghostFileNamespace):])
	}
	return names, nil
}
