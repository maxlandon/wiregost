#!/usr/bin/env python3.7.2

# -------------------- Imports ----------------------- #

# Custom Entities
from EffectiveCouscous.entities.infrastructure import Netblock, IPv4Address

# Maltego Message
from canari.maltego.message import Field, MaltegoException, LinkColor

# Maltego Transforms
from canari.maltego.transform import Transform

# API
import json
from EffectiveCouscous.tools.msftools import apitools

# GUI
import canari.easygui as gui

# Icons
from EffectiveCouscous.resource import network_interface
# ---------------------------------------------------- #


__author__ = 'Maxime Landon'
__copyright__ = 'Copyright 2019, EffectiveCouscous Project'
__credits__ = []

__license__ = 'GPLv3'
__version__ = '0.2'
__maintainer__ = 'Maxime Landon'
__email__ = 'maximelandon@gmail.com'
__status__ = 'Development'




# IPv4 Addresses ---------------------------------------------------------------------------------------------------------------------------------#


class EnumerateWorkspaceIP(Transform):
    """ Enumerates one or more Hosts in Workspace"""

    display_name = "To Workspace IPv4 addresses "
    transform_set = 'Msf__DB                    | Netblock              | Enumerate'
    input_type = Netblock

    def do_transform(self, request, response, config):
        netblock = request.entity
        
        # Test for properties
        try:
            url = config['EffectiveCouscous.local.baseurl'] + 'hosts'
            params = (('workspace', '{0}'.format(netblock['name'])),)
            hosts = apitools.get_json_dict(url, config, params=params)
            title = "Host Choice"
            msg = "Choose one or more Hosts for IPv4Address enumeration"
        except KeyError:
            raise MaltegoException("This Netblock is not tied to a Metasploit Workspace. \
                                    Please associate it with a workspace before running this transform")
            return response

        # Select Hosts
        host_infos = []
        host_map = {}
        for host in hosts:
            info = '{0}      {1}'.format(host['address'], host['name'])
            host_infos.append(info)
            host_map[host['address']] = info
        raw_choices = gui.multchoicebox(title=title, msg=msg, choices=(host_infos))

        for choice in raw_choices:
            for (address, mapped) in host_map.items():
                if choice == mapped:
                    for host in hosts:
                        if address == str(host['address']):
                            ip_entity = IPv4Address()
                            ip_entity['ipv4-address'] = host['address']
                            ip_entity.host_id = host['id']
                            ip_entity.workspace_id = host['workspace_id']
                            ip_entity.icon_url = network_interface
                            ip_entity.origin_tool = 'Metasploit'
                            response += ip_entity
        return response


