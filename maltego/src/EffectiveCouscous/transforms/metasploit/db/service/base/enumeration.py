#!/usr/bin/env python3.7.2

# -------------------- Imports ----------------------- #

# Custom Entities
from EffectiveCouscous.entities.service.base import MetasploitService 
from EffectiveCouscous.entities.credential.base import MetasploitCredential

# Maltego Message
from canari.maltego.message import Field, MaltegoException, LinkStyle, LinkColor

# Maltego Transforms
from canari.maltego.transform import Transform

# API
import json
from EffectiveCouscous.tools.msftools import apitools
# ---------------------------------------------------- #


__author__ = 'Maxime Landon'
__copyright__ = 'Copyright 2019, EffectiveCouscous Project'
__credits__ = []

__license__ = 'GPLv3'
__version__ = '0.2'
__maintainer__ = 'Maxime Landon'
__email__ = 'maximelandon@gmail.com'
__status__ = 'Development'



# Credentials ---------------------------------------------------------------------------------------------------------------------------------------------#


class EnumerateServiceCredentials(Transform):
    """ Enumerates Credentials for this Service"""

    display_name = 'To Service Credentials'
    transform_set = 'Msf__DB                    | Service                 | Enumerate'
    input_type = MetasploitService

    def do_transform(self, request, response, config):
        service = request.entity

        # Test for properties & get Workspace Credentials --------------------------------//
        try:
            ws_url = config['EffectiveCouscous.local.baseurl'] + 'workspaces/{0}'.format(service.workspace_id)
            workspace = apitools.get_json_dict(ws_url, config)['name'] 
            url = config['EffectiveCouscous.local.baseurl'] + 'credentials'
            params = (('workspace', '{0}'.format(workspace)), ('svcs', '{0}'.format(service.name)),)
            creds = apitools.get_json_dict(url, config, params=params)
        except KeyError:
            raise MaltegoException("This IPv4Address is not tied to a Metasploit Host. \
                                    Please associate it with a Host before running this transform")
            return response

        # REWRITE FOR ALL CASES WHERE CREDENTIAL PROPERTIES ARE NOT EXISTING !!!!!!!!!!!!!!!!!
        # Filter for Service
        for cred in creds:
            try:
                if (int(cred['logins'][0]['service_id'])) == int(service.id):
                    cred_entity = MetasploitCredential()
                    cred_entity.id = cred['id']
                    cred_entity.logins_count = cred['logins_count']
                    cred_entity.pub_username = cred['public']['username']
                    cred_entity.pub_type = cred['public']['type']
                    cred_entity.priv_data = cred['private']['data']
                    cred_entity.priv_type = cred['private']['type']
                    cred_entity.priv_jtr_format = cred['private']['jtr_format']
                    cred_entity.origin_service_id = cred['origin']['service_id']
                    cred_entity.origin_type = cred['origin']['type']
                    cred_entity.origin_module = cred['origin']['module_full_name']
                    cred_entity.name = '{}/{}'.format(cred_entity.pub_username, cred_entity.priv_data)
                    # Link Style
                    cred_entity.link_color = LinkColor.DarkGreen
                    cred_entity.link_thickness = 3
                    response += cred_entity
            except KeyError as e:
                    continue 
            except IndexError as f:
                    continue 

        return response
