package generate

const (
	// OS ----------------------------------------------

	// Windows OS
	Windows = "windows"
	// Linux OS
	Linux = "linux"
	// Darwin / MacOS
	Darwin = "darwin"

	// Compiler ----------------------------------------

	// CC64EnvVar - Environment variable for specifying MinGW 64-bit path
	CC64EnvVar = "WIREGOST_CC_64"
	// CC32EnvVar - Environment variable for specifying MinGW 32-bit path
	CC32EnvVar = "WIREGOST_CC_32"

	// Transports --------------------------------------

	// DefaultReconnectInterval - In seconds
	DefaultReconnectInterval = 60
	// DefaultMTLSLPort - Default MTLS listener port
	DefaultMTLSLPort = 8888
	// DefaultHTTPLPort - Default listener port
	DefaultHTTPLPort = 443 // Assume SSL, it will fallback
)

var (
	// DefaultMinGWPath -
	DefaultMinGWPath = map[string]string{
		"386":   "/usr/bin/i686-w64-mingw32-gcc",
		"amd64": "/usr/bin/x86_64-w64-mingw32-gcc",
	}
)
