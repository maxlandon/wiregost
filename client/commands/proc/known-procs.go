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

package proc

import "github.com/evilsocket/islazy/tui"

var (
	// Stylizes known processes in the `ps` command
	knownProcs = map[string]string{
		"ccSvcHst.exe": tui.RED, // SEP
		"cb.exe":       tui.RED, // Carbon Black

		// Antivirus Scanners
		"MCODS":        tui.RED, //McAfee On Demand Scanner
		"MCSHIELD":     tui.RED, //McAfee Realtime Scanner
		"msmpeng.exe":  tui.RED, //Microsoft Security Essentials
		"navapsvc.exe": tui.RED, //Norton Antivirus
		"avkwctl.exe":  tui.RED, //MicroWorld eScan
		"fsav32.exe":   tui.RED, //F-Secure Anti-Virus
		"mcshield.exe": tui.RED, //McAfee VirusScan
		"ntrtscan.exe": tui.RED, //TrendMicro OfficeScan
		"avguard.exe":  tui.RED, //AntiVir PersonalEdition
		"ashServ.exe":  tui.RED, //avast! antivirus
		"AVENGINE.EXE": tui.RED, //Panda AntiVirus Titanium
		"avgemc.exe":   tui.RED, //AVG Anti-Virus
		"tmntsrv.exe":  tui.RED, //TrendMicro Internet Security

		// Firewalls
		"BLACKD":       tui.RED, // BlackICE
		"efpeadm.exe":  tui.RED, // CA eTrust EZ Firewall
		"VPNGUI":       tui.RED, // Cisco Systems VPN Client
		"CVPND":        tui.RED, // Cisco Systems VPN Client
		"IPSECLOG":     tui.RED, //Cisco Systems VPN Client
		"cfp.exe":      tui.RED, //Comodo Firewall Pro - Defense+
		"fsdfwd.exe":   tui.RED, // F-Secure Distributed Firewall
		"fsguiexe.exe": tui.RED, //F-Secure Internet Security
		"blackd.exe":   tui.RED, //ISS BlackIce
		"kpf4gui.exe":  tui.RED, //Kerio Personal Firewall
		"MSSCLLI":      tui.RED, //McAfee
		"MCSHELL":      tui.RED, //McAfee
		"MPFSERVICE":   tui.RED, //McAfee Personal Firewall Service
		"MPFAGENT":     tui.RED, //McAfee Personal Firewall Service
		"nisum.exe":    tui.RED, //Norton Personal Firewall
		"smc.exe":      tui.RED, //Sygate Personal Firewall
		"persfw.exe":   tui.RED, //Tiny Firewall
		"pccpfw.exe":   tui.RED, //Trend Micro Internet Security
		"WINSS":        tui.RED, //WindowsLive One Care
		"ZLCLIENT":     tui.RED, //Zone Alarm

		// Others
		"ExDomusKNX_Service":   tui.RED, //ExDomus
		"McAfeeDataBackup.exe": tui.RED, //McAfee Backup
		"GoogleDesktop.exe":    tui.RED, //Google Desktop
		"CCEVTMGR":             tui.RED, //Symantec Event Manager
		"CCSETMGR":             tui.RED, //Symantec Settings Manager
	}
)
