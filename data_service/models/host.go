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

package models

type Host struct {
	ID                  int `json:"host_id" sql:"host_id,pk"`
	Address             string
	MAC                 string
	Comm                string
	Name                string
	State               string
	OSName              string
	OSFlavor            string
	OSSp                string
	OSLang              string
	OSFamily            string
	Arch                string
	DetectedArch        string
	WorkspaceID         int `sql:"workspace_id,notnull,on_delete:CASCADE"`
	Workspace           *Workspace
	Purpose             string
	Info                string
	Comments            string
	Scope               string
	VirtualHost         string
	NoteCount           int
	VulnCount           int
	ServiceCount        int
	HostDetailCount     int
	ExploitAttemptCount int
	CredCount           int
	CreatedAt           string
	UpdatedAt           string
}
