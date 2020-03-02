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

	"github.com/maxlandon/wiregost/server/db"
)

const (
	// ProfilesBucketName - DB file for implant profiles
	ProfilesBucketName = "profiles"
)

// ProfileSave - Save a ghost profile to disk
func ProfileSave(name string, config *GhostConfig) error {
	bucket, err := db.GetBucket(ProfilesBucketName)
	if err != nil {
		return err
	}
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return bucket.Set(name, configJSON)
}

// ProfileByName - Fetch a single profile from the database
func ProfileByName(name string) (*GhostConfig, error) {
	bucket, err := db.GetBucket(ProfilesBucketName)
	if err != nil {
		return nil, err
	}
	rawProfile, err := bucket.Get(name)
	config := &GhostConfig{}
	err = json.Unmarshal(rawProfile, config)
	return config, err
}

// Profiles - Fetch a map of name<->profiles current in the database
func Profiles() map[string]*GhostConfig {
	bucket, err := db.GetBucket(ProfilesBucketName)
	if err != nil {
		return nil
	}
	rawProfiles, err := bucket.Map("")
	if err != nil {
		return nil
	}

	profiles := map[string]*GhostConfig{}
	for name, rawProfile := range rawProfiles {
		config := &GhostConfig{}
		err := json.Unmarshal(rawProfile, config)
		if err != nil {
			continue // We should probably log these failures ...
		}
		profiles[name] = config
	}
	return profiles
}
