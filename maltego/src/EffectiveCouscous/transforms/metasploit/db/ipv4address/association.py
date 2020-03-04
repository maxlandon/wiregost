#!/usr/bin/env python3.7.2

# -------------------- Imports --------------------- #

# Custom Entities
from EffectiveCouscous.entities.infrastructure import IPv4Address
# Maltego Entities
from canari.maltego.message import Field

# Maltego Transforms
from canari.maltego.transform import Transform

# API
import json
from EffectiveCouscous.tools.msftools import apitools

# GUI
import canari.easygui as gui

# Icons
from EffectiveCouscous.resource import network_interface
#--------------------------------------------------- #


__author__ = 'Maxime Landon'
__copyright__ = 'Copyright 2019, EffectiveCouscous Project'
__credits__ = []

__license__ = 'GPLv3'
__version__ = '0.2'
__maintainer__ = 'Maxime Landon'
__email__ = 'maximelandon@gmail.com'
__status__ = 'Development'



# Host ---------------------------------------------------------------------------------------------------------------------------------------------#


class AppendHostProperties(Transform):
    """Adds Host properties to this IPv4Address"""

    display_name = 'Append Host properties'
    transform_set = 'Msf__DB                    | IPv4Address         | Associate'
    input_type = IPv4Address 

    def do_transform(self, request, response, config):
        ip_entity = request.entity

        # Test for properties -------------------------------------------//
        try: 
            test = ip_entity['id']
            if ip_entity['id'] is not None:
                title = "Confirmation"
                msg = """This IPv4Address is already bound to a Metasploit Host. \n
                    Do you really want to change the concerned properties ?"""
                confirm = gui.choicebox(title=title, msg=msg, choices=['Yes', 'No'])
            if confirm == 'No':
                return response
        except KeyError:
            pass

        # Select Workspaces & Hosts -------------------------------------//
        url = config['EffectiveCouscous.local.baseurl'] + 'workspaces' 
        workspaces = apitools.get_json_dict(url, config)
        title = "Workspace Choice"
        msg = "Please choose a workspace for Host selection"
        ws_names = [ workspace['name'] for workspace in workspaces ]
        choice = gui.choicebox(title=title, msg=msg, choices=(ws_names))

        # Select Hosts 
        url = config['EffectiveCouscous.local.baseurl'] + 'hosts' 
        params = (('workspace', '{0}'.format(choice)),)
        hosts = apitools.get_json_dict(url, config, params=params)
        title = "Host Choice"
        msg = "Choose a Metasploit Host to associate with this IPv4Address"
        host_infos = []
        host_names = []
        for host in hosts:
            info = '{0}      {1}'.format(host['address'], host['name'])
            host_infos.append(info)
            host_names.append(host['name'])
        host_infos.append("Add Host")
        raw_choice = gui.choicebox(title=title, msg=msg, choices=(host_infos))
        host = {}
        if "Add Host" in raw_choice:
            host['name'] = "Add Host" 
        else:
            for h in hosts:
                if h['address'] in raw_choice:
                    host = h

        # If existing host ---------------------------------------------//
        if host['name'] != "Add Host":
            ip_entity['ipv4-address'] = host['address']
            ip_entity.host_id = host['id']
            ip_entity.workspace_id = host['workspace_id']
            ip_entity.icon_url = network_interface
            ip_entity.origin_tool = 'Metasploit'
            response + ip_entity

        # If New Host --------------------------------------------------//
        if host['name'] == 'Add Host':
            url = config['EffectiveCouscous.local.baseurl'] + 'hosts'
            title = "New Host"
            msg = """Enter Host properties for creating a Host in Metasploit"""
            field_names = ["Address", "MAC", "Host Name", "OS Name", "OS Flavor",
                            'OS SP', 'OS Language', 'Purpose', 'Info', 'Comments', 'Scope',
                            'Virtual Host', 'Architecture', 'State']
            field_values = []
            field_values = gui.multenterbox(title=title, msg=msg, fields=field_names, values=field_values)
            while 1:                                                 
                if field_values == None: break                        
                errmsg = ""                                          
                for i in range(len(field_names)):                     
                    if field_values[i].strip() == "":                 
                        errmsg += ('"%s" is a required field.\n\n' % field_names[i]) 
                if errmsg == "":                                     
                    break 
                field_values = gui.multenterbox(errmsg, field_values, fields=field_names)                                         
            # Post Host
            dict = {}
            dict['workspace'] = choice 
            dict['host'] = field_values[0]
            dict['mac'] = field_values[1]
            dict['name'] = field_values[2]
            dict['os_name'] = field_values[3]
            dict['os_flavor'] = field_values[4]
            dict['os_sp'] = field_values[5]
            dict['os_lang'] = field_values[6]
            dict['purpose'] = field_values[7]
            dict['info'] = field_values[8]
            dict['comments'] = field_values[9]
            dict['scope'] = field_values[10]
            dict['virtual_host'] = field_values[11]
            dict['arch'] = field_values[12]
            dict['state'] = field_values[13]
            data = json.dumps(dict)
            post = apitools.post_json(url, data, config)

            # Fetch attributes of new Host
            host_dict = post.json()['data']
            ip_entity['ipv4-address'] = host_dict['address']
            ip_entity.host_id = host_dict['id']
            ip_entity.workspace_id = host_dict['workspace_id']
            ip_entity.icon_url = network_interface
            ip_entity.origin_tool = 'Metasploit'
            response + ip_entity

        return response
