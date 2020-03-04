#!/usr/bin/env python3.7.2

# -------------------- Imports ----------------------- #

# Custom Entities
# Custom entities are directly imported in functions using them,
# so to avoid circular references during imports.

# Icons
from EffectiveCouscous.resource import systems, devices

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
#                                        OS FACTORY                                     #
#---------------------------------------------------------------------------------------#

# The OS Factory is responsible for determining the Type of Entity to be spawn, related 
# to Operating Systems, throughout EffectiveCouscous.


# I - General Functions
#__________________________________
# 1) getOsEntity(os_name, name):    Is responsible for determining the type of Entity
#                                   to be returned for transforms spawning Hosts.
#                                   It looks at different fields of a Metasploit Host, 
#                                   prioritizes them and finds the good Entity.



# ----------------------------   General Functions   --------------------------------- #

def getOsEntity(os_name, name):
    from EffectiveCouscous.entities.host.base import MetasploitHost, AppleHost, LinuxHost, WindowsHost

    os_entity = MetasploitHost()
    os_entity.icon_url = systems['generic']

    if os_name or name:
        if name:
            # 1) ________________________________________________________
            # Try with Name first, more reliable for Apple devices
            #....................................
            # Apple
            if re.search("macbook", name.lower(), re.I):
                os_entity = AppleHost()
                return os_entity
            elif re.search("ipad", name.lower(), re.I):
                os_entity = AppleHost()
                return os_entity
            #....................................
            # Linux
            elif ".home" in name.lower():
                os_entity = LinuxHost()
            # 2) ________________________________________________________
            # If the Name has not confirmed any device, use OS Name 
            elif os_name:
                #....................................
                # Windows
                if "windows" in os_name.lower():
                    if "windows 2003" in os_name.lower():
                        os_entity = WindowsHost()
                    elif "windows 2000" in os_name.lower():
                        os_entity = WindowsHost()
                    elif "windows 2008" in os_name.lower():
                        os_entity = WindowsHost()
                    elif "windows 2012" in os_name.lower():
                        os_entity = WindowsHost()
                    elif "windows xp" in os_name.lower():
                        os_entity = WindowsHost()
                    elif "windows 7" in os_name.lower():
                        os_entity = WindowsHost()
                    elif "windows vista" in os_name.lower():
                        os_entity = WindowsHost()
                    elif "windows 10" in os_name.lower():
                        os_entity = WindowsHost()
                    else:
                        os_entity = WindowsHost()
                    return os_entity
                #....................................
                # Apple
                elif re.search("ios", os_name.lower(), re.I):
                    os_entity = AppleHost()
                    return os_entity

                #....................................
                # Linux
                elif "linux" in os_name.lower():
                    if "arch" in os_name.lower():
                        os_entity = LinuxHost()
                    else:
                        os_entity = LinuxHost()
                    return os_entity
                #....................................
                # BSD
                #....................................
                # Embedded
                elif "embedded" in os_name.lower():
                    os_entity = EmbeddedOS()
        # 3) _______________________________________________________________
        # Both lookup methods failed if this point is reached. Spawn generic
        else:
            return os_entity
    # A Host needs to be returned anyway
    return os_entity



