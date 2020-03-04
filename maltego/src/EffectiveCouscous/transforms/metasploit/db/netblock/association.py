#!/usr/bin/env python3.7.2

# -------------------- Imports --------------------- #

# Custom Entities
from EffectiveCouscous.entities.infrastructure import Netblock

# Maltego Entities
from canari.maltego.message import StringEntityField, Field, IntegerEntityField, BooleanEntityField

# Maltego Transforms
from canari.maltego.transform import Transform

# API
import json
from EffectiveCouscous.tools.msftools import apitools

# GUI
import canari.easygui as gui

# Icons
from EffectiveCouscous.resource import tools
#--------------------------------------------------- #


__author__ = 'Maxime Landon'
__copyright__ = 'Copyright 2019, EffectiveCouscous Project'
__credits__ = []

__license__ = 'GPLv3'
__version__ = '0.2'
__maintainer__ = 'Maxime Landon'
__email__ = 'maximelandon@gmail.com'
__status__ = 'Development'




# Workspaces ---------------------------------------------------------------------------------------------------------------------------------------------#


class NetblockToMetasploitWorkspace(Transform):
    """Adds properties associated to a Metasploit workspace"""

    display_name = "To Workspace"
    transform_set = 'Msf__DB                    | Netblock              | Associate'
    input_type = Netblock

    def do_transform(self, request, response, config):
        netblock = request.entity

        # Select Workspaces
        url = config['EffectiveCouscous.local.baseurl'] + 'workspaces' 
        workspaces = apitools.get_json_dict(url, config)
        title = "Workspace Choice"
        msg = "Choose a Metasploit workspace to associate with this Netblock"
        ws_names = [workspace['name'] for workspace in workspaces]
        ws_names.append('Add Workspace')
        workspace = {} 
        choice = gui.choicebox(msg=msg, title=title, choices=(ws_names))
        if choice == "Add Workspace":
            workspace['name'] = "Add Workspace" 
        else:
            for ws in workspaces:
                if ws['name'] == choice:
                    workspace = ws

        # If Existing Workspace --------------------------------------------------------------- //
        if workspace['name'] != "Add Workspace":
            # Set Values 
            netblock.name = '-' if workspace['name'] is None else workspace['name']
            netblock.workspace_id = workspace['id']
            netblock.boundary = '-' if workspace['boundary'] is None else workspace['boundary']
            netblock.description = '-' if workspace['description'] is None else workspace['description']
            netblock.owner_id = '-' if workspace['owner_id'] is None or '-' else workspace['owner_id']
            netblock.limit_to_network = workspace['limit_to_network']
            netblock.import_fingerprint = workspace['import_fingerprint']
            netblock.created_at = workspace['created_at']
            netblock.updated_at = workspace['updated_at']
            netblock.origin_tool = 'Metasploit'
            # IP Range and Boundary
            if workspace['boundary'] is not None: netblock['ipv4-range'] = netblock.boundary
            # Add to response
            response + netblock

        # If New Workspace ------------------------------------------------------------------- //
        if workspace['name'] == "Add Workspace":
            msg = "New Workspace"
            fieldNames = ["Name"]
            fieldValues = gui.multenterbox(msg, fields=fieldNames)
            while 1:                                                 
                if fieldValues == None: break                        
                errmsg = ""                                          
                for i in range(len(fieldNames)):                     
                    if fieldValues[i].strip() == "":                 
                        errmsg += ('"%s" is a required field.\n\n' % fieldNames[i])                                       
                if errmsg == "":                                     
                    break 
                fieldValues = gui.multenterbox(errmsg, fieldValues, fieldNames)                                         

            # Create and Fetch Workspace in Metasploit
            dict = {}
            dict['name'] = fieldValues[0]
            data = json.dumps(dict)
            post = apitools.post_json(url, data, config) 
            workspaces = apitools.get_json_dict(url, config)
            ws = []
            for workspace in workspaces:
                if workspace['name'] == dict['name']:
                    ws.append(workspace)      
            workspace = ws[0]

            # Set Values 
            netblock.name = '-' if workspace['name'] is None else workspace['name']
            netblock.workspace_id = workspace['id']
            netblock.boundary = '-' if workspace['boundary'] is None else workspace['boundary']
            netblock.description = '-' if workspace['description'] is None else workspace['description']
            netblock.owner_id = '-' if workspace['owner_id'] is None else workspace['owner_id']
            netblock.limit_to_network = workspace['limit_to_network']
            netblock.import_fingerprint = workspace['import_fingerprint']
            netblock.created_at = workspace['created_at']
            netblock.updated_at = workspace['updated_at']
            netblock.origin_tool = 'Metasploit'
            # IP Range and Boundary
            if workspace['boundary'] is not None: netblock['ipv4-range'] = netblock.boundary
            # Add to response
            response + netblock

        return response
