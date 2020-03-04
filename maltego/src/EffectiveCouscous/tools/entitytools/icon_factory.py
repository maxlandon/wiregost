#!/usr/bin/env python3.7.2

# -------------------- Imports ----------------------- #

# Icons
from EffectiveCouscous.resource import (systems,
                                        devices,
                                        tools,
                                        openport,
                                        closedport,
                                        unavailableport,
                                        timedoutport)

# String search
import re
# ---------------------------------------------------- #


__author__ = 'Maxime Landon'
__copyright__ = 'Copyright 2019, EffectiveCouscous Project'
__credits__ = []

__license__ = 'GPL'
__version__ = '0.2'
__maintainer__ = 'Maxime Landon'
__email__ = 'maximelandon@gmail.com'
__status__ = 'Development'



#---------------------------------------------------------------------------------------#
#                                   ICON FACTORY                                        #
#---------------------------------------------------------------------------------------#

# The Icon Factory is responsible for determining the right Icon for various Entities'
# attributes, which may be used for side-icons, for instance.

# Functions are regrouped by the type of Icons they fetch


# Origin Tool ------------------------------------------------------------------------- #

def getOriginTool(entity, origin_tool):
    tool = origin_tool.lower()

    # Test for tools
    if 'metasploit' in tool:
        entity.tool_icon = tools['metasploit']
 


# Service State ----------------------------------------------------------------------- #

def getStateIcon(service, state):
    from EffectiveCouscous.resource import (openport, closedport, timedoutport, unavailableport)
    if state == 'open':
        service.state_icon = openport
    if state == 'closed':
        service.state_icon = closedport
    if state == 'filtered':
        service.state_icon = timedoutport
    if state == 'unknown':
        service.state_icon = unavailableport


# Operating System / Device ----------------------------------------------------------- #

def getOsIcon(host, os):
    from EffectiveCouscous.entities.host.base import MetasploitHost, AppleHost, LinuxHost, WindowsHost

    if host.os_name or host.name:
        if host.name:
            # 1) ________________________________________________________
            # Try with Name first, more reliable for Apple devices
            #....................................
            # Apple
            if re.search("macbook", host.name.lower(), re.I):
                host.value = host.name
                host.icon_url = devices['macbook']
                return host
            elif re.search("ipad", host.name.lower(), re.I):
                host.value = host.name
                host.icon_url = devices['ipad']
                return host
            #....................................
            # Linux
            elif ".home" in host.name.lower():
                host.icon_url = systems['linux']
            # 2) ________________________________________________________
            # If the Name has not confirmed any device, use OS Name 
            elif host.os_name:
                    #....................................
                    # Windows
                    if "windows" in host.os_name.lower():
                        if "windows 2003" in host.os_name.lower():
                            host.icon_url = systems['windows2003']
                        elif "windows 2000" in host.os_name.lower():
                            host.icon_url = systems['windows2000']
                        elif "windows 2008" in host.os_name.lower():
                            host.icon_url = systems['windows2008']
                        elif "windows 2012" in host.os_name.lower():
                            host.icon_url = systems['windows2012']
                        elif "windows xp" in host.os_name.lower():
                            host.icon_url = systems['windowsxp']
                        elif "windows 7" in host.os_name.lower():
                            host.icon_url = systems['windows7']
                        elif "windows vista" in host.os_name.lower():
                            host.icon_url = systems['windowsvista']
                        elif "windows 10" in host.os_name.lower():
                            host.icon_url = systems['windows10']
                        else:
                            host.icon_url = systems['windows']
                    #....................................
                    # Apple
                    elif re.search("ios", host.os_name.lower(), re.I):
                        host.icon_url = systems['apple']

                    #....................................
                    # Linux
                    elif "linux" in host.os_name.lower():
                        if "arch" in host.os_name.lower():
                            host.icon_url = systems['archlinux']
                        #  if "debian"
                        #  if "ubuntu", etc...
                        else:
                            host.icon_url = systems['linux']
                    #....................................
                    # BSD
                    #....................................
                    # Embedded
                    elif "embedded" in host.os_name.lower():
                        host.icon_url = systems['generic']
        # 3) _______________________________________________________________
        # Both lookup methods failed if this point is reached. Spawn generic
        else:
            pass



