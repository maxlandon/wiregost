#!/usr/bin/env python3.7.2

# -------------------- Imports --------------------- #

# Base Transforms
from EffectiveCouscous.transforms.metasploit.db.netblock.association import NetblockToMetasploitWorkspace

# Custom Host Entities
from EffectiveCouscous.entities.host.base import MetasploitHost
#--------------------------------------------------- #


__author__ = 'Maxime Landon'
__copyright__ = 'Copyright 2019, EffectiveCouscous Project'
__credits__ = []

__license__ = 'GPLv3'
__version__ = '0.2'
__maintainer__ = 'Maxime Landon'
__email__ = 'maximelandon@gmail.com'
__status__ = 'Development'



# Workspaces -------------------------------------------------------------------------------------------------------------------#

class HostToMetasploitWorkspace(NetblockToMetasploitWorkspace):
    """Adds properties associated to a Metasploit workspace"""

    display_name = "To Host Workspace"
    transform_set = 'Msf__DB                    | Host                     | Associate'
    input_type = MetasploitHost
