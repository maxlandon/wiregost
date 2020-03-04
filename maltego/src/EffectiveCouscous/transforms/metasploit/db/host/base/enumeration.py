#!/usr/bin/env python3.7.2

# -------------------- Imports ----------------------- #

# Custom Entities
from EffectiveCouscous.entities.host.base import MetasploitHost
from EffectiveCouscous.entities.infrastructure import IPv4Address

# Maltego Message
from canari.maltego.message import Field, MaltegoException, LinkStyle, LinkColor

# Maltego Transforms
from canari.maltego.transform import Transform

# API
import json
from EffectiveCouscous.tools.msftools import apitools

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



# IPv4 Addresses -------------------------------------------------------------------------------------------------------------------#


class EnumerateHostIP(Transform):
    """ Enumerates all IP addresses that match this Host in Metasploit."""

    display_name = "To Host IPv4 Addresses"
    transform_set = 'Msf__DB                    | Host                     | Enumerate'
    input_type = MetasploitHost

    def do_transform(self, request, response, config):
        host_input = request.entity
        
        # Get workspace & Test for Properties ---------------------------------//
        url = config['EffectiveCouscous.local.baseurl'] + 'workspaces' 
        workspaces = apitools.get_json_dict(url, config)
        workspace_name = ''
        try:
            for workspace in workspaces:
                if workspace['id'] == host_input.workspace_id:
                    workspace_name = workspace['name']
                    url = config['EffectiveCouscous.local.baseurl'] + 'hosts'
                    params = (('workspace', '{0}'.format(workspace_name)),)
                    hosts = apitools.get_json_dict(url, config, params=params)
        except KeyError:
            raise MaltegoException("This Host is not tied to a Metasploit Workspace. \
                                    Please associate it with a workspace before running this transform")
            return response
        
        # Select Hosts --------------------------------------------------------//
        for host in hosts:
            if (host['name'] == host_input.name) and (host['os_name'] == host_input.os_name):
                ip_entity = IPv4Address()
                ip_entity['ipv4-address'] = host_input['address']
                ip_entity.host_id = host_input['id']
                ip_entity.workspace_id = host_input['workspace_id']
                ip_entity.icon_url = network_interface
                ip_entity.origin_tool = 'Metasploit'
                # Link Style
                ip_entity.link_color = LinkColor.Black
                ip_entity.link_thickness = 3
                response += ip_entity

        return response
