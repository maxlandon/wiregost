#!/usr/bin/env python3.7.2

# -------------------- Imports --------------------- #

# Custom Entities
from EffectiveCouscous.entities.infrastructure import Netblock

# Maltego Messages
from canari.maltego.message import MaltegoException

# Maltego Transforms
from canari.maltego.transform import Transform

# API
import json
from EffectiveCouscous.tools.msftools import apitools

# GUI
import canari.easygui as gui
# ---------------------------------------------------- #


__author__ = 'Maxime Landon'
__copyright__ = 'Copyright 2019, EffectiveCouscous Project'
__credits__ = []

__license__ = 'GPLv3'
__version__ = '0.2'
__maintainer__ = 'Maxime Landon'
__email__ = 'maximelandon@gmail.com'
__status__ = 'Development'



# Workspaces ---------------------------------------------------------------------------------------------------------------------------------#


# PUT ------------------------------------------//
class NetblockPutWorkspace(Transform):
    """ Push Maltego properties to Metasploit workspace."""

    display_name = 'PUT Workspace'
    transform_set = 'Msf__DB                    | Workspace           | Update'
    input_type = Netblock 

    def do_transform(self, request, response, config):
        ws = request.entity

        # Test for properties
        try:
            test = ws['workspace_id']
        except KeyError:
            raise MaltegoException("This Netblock/Host is not tied to a Metasploit Workspace. \
                                    Please associate it with a workspace before running this transform")
            return response

        dict = {}
        dict['id'] = ws['workspace_id']
        dict['name'] = ws['name']
        dict['created_at'] = ws['created_at']
        dict['updated_at'] = ws['updated_at']
        dict['boundary'] = ws['boundary']
        dict['description'] = ws['description']
        dict['owner_id'] = ws['owner_id']
        dict['limit_to_network'] = ws['limit_to_network']
        dict['import_fingerprint'] = ws['import_fingerprint']
        data = json.dumps(dict)
        url = config['EffectiveCouscous.local.baseurl'] + 'workspaces/{0}'.format(ws['workspace_id'])
        update = apitools.put_json(url, data, config)
        return response


# GET ------------------------------------------//
class NetblockGetWorkspace(Transform):
    """ Pull Metasploit properties to Maltego workspace."""

    display_name = 'GET Workspace'
    transform_set = 'Msf__DB                    | Workspace           | Update'
    input_type = Netblock 

    def do_transform(self, request, response, config):
        ws = request.entity

        # Test for properties
        try:
            test = ws['workspace_id']
        except KeyError:
            raise MaltegoException("This Netblock/Host is not tied to a Metasploit Workspace. \
                                    Please associate it with a workspace before running this transform")
            return response

        url = config['EffectiveCouscous.local.baseurl'] + 'workspaces/{0}'.format(ws['workspace_id'])
        workspace = apitools.get_json_dict(url, config)
        ws['workspace_id'] = workspace['id']
        ws['name'] = workspace['name']
        ws['created_at'] = workspace['created_at']
        ws['updated_at'] = workspace['updated_at']
        ws['boundary'] = workspace['boundary']
        ws['description'] = workspace['description']
        ws['owner_id'] = workspace['owner_id']
        ws['limit_to_network'] = workspace['limit_to_network']
        ws['import_fingerprint'] = workspace['import_fingerprint']
        return response


# DELETE ---------------------------------------//
class NetblockDeleteWorkspace(Transform):
    """ Delete workspace associated to this Netblock."""

    display_name = 'DELETE Workspace'
    transform_set = 'Msf__DB                    | Workspace           | Update'
    input_type = Netblock 

    def do_transform(self, request, response, config):
        ws = request.entity

        # Test for properties
        try:
            test = ws['workspace_id']
        except KeyError:
            raise MaltegoException("This Netblock/Host is not tied to a Metasploit Workspace. \
                                    Please associate it with a workspace before running this transform")
            return response

        url = config['EffectiveCouscous.local.baseurl'] + 'workspaces'
        data = '{ "ids": [%s] }' % ws['workspace_id']
        delete = apitools.delete_json(url, data, config)
        return response
