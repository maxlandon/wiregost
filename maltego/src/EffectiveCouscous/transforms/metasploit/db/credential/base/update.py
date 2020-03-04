#!/usr/bin/env python3.7.2

# -------------------- Imports --------------------- #

# Custom Entities
from EffectiveCouscous.entities.credential.base import MetasploitCredential

# Maltego Messages
from canari.maltego.message import MaltegoException

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


########### TO BE FULLY REWRITTEN !! TO BE FULLY REWRITTEN !! TO BE FULLY REWRITTEN !! TO BE FULLY REWRITTEN !! ############

# Credentials -------------------------------------------------------------------------------------------------------------------#


class PutCredential(Transform):
    """PUT Maltego properties to Metasploit Credential."""

    display_name = 'PUT Credential'
    transform_set = 'Metasploit | Cred | Update'
    input_type = MetasploitCredential

    def do_transform(self, request, response, config):
        # Get Credential JSON dictionary
        credential = request.entity

        # Test for properties
        try:
            url = config['EffectiveCouscous.local.baseurl'] + 'credentials/{0}'.format(credential['id'])
        except KeyError:
            raise MaltegoException("This Credential is not tied to a Metasploit Credential. Please associate it with a Credential before running this transform")
            return response

        cred = apitools.get_json_dict(url, config)[0]

        # THIS IS NOT WORKING: THE DICTIONARY ASKED FOR THE PUT METHOD IS
        # WEIRD, THEREFORE IT NEEDS TO BE COPIED ANOTHER WAY THAN THIS ONE
        # Push values
        #  cred['id'] = int(credential.id)
        #  cred['logins_count'] = int(credential.logins_count)
        #  cred['public']['username'] = credential.pub_username
        #  cred['public']['type'] = credential.pub_type
        #  cred['private']['data'] = credential.priv_data
        #  cred['private']['type'] = credential.priv_type
        #  cred['private']['jtr_format'] = credential.priv_jtr_format
        #  cred['origin']['service_id'] = credential.origin_service_id
        #  cred['origin']['type'] = credential.origin_type
        #  cred['origin']['module_full_name'] = credential.origin_module

        #  data = json.dumps(cred)
        #  update = apitools.post_json(url, data) 
        return response


class GetCredential(Transform):
    """ GET Metasploit properties to Maltego Credential"""

    display_name = 'GET Credential'
    transform_set = 'Metasploit | Cred | Update'
    input_type = MetasploitCredential

    def do_transform(self, request, response, config):
        # Get Credential JSON dictionary
        credential = request.entity

        # Test for properties
        try:
            url = config['EffectiveCouscous.local.baseurl'] + 'credentials/{0}'.format(credential['id'])
        except KeyError:
            raise MaltegoException("This Credential is not tied to a Metasploit Credential. Please associate it with a Credential before running this transform")
            return response

        cred = apitools.get_json_dict(url, config) 

        # Fetch values
        credential.id = cred['id']
        credential.logins_count = cred['logins_count']
        credential.pub_username = cred['public']['username']
        credential.pub_type = cred['public']['type']
        credential.priv_data = cred['private']['data']
        credential.priv_type = cred['private']['type']
        credential.priv_jtr_format = cred['private']['jtr_format']
        credential.origin_service_id = cred['origin']['service_id']
        credential.origin_type = cred['origin']['type']
        credential.origin_module = cred['origin']['module_full_name']
        raise MaltegoException('Due to implementation limitations, it is not possible to refresh this Entity "as is". Please suppress it and spawn it again from its parent Entity to see it with its new properties.')

        return response


class DeleteCredential(Transform):
    """ Delete Metasploit Credential associated to this Credential Entity"""

    display_name = 'DELETE Credential'
    transform_set = 'Metasploit | Cred | Update'
    input_type = MetasploitCredential 

    def do_transform(self, request, response, config):
        credential = request.entity

        # Test for properties
        try:
            url = config['EffectiveCouscous.local.baseurl'] + 'credentials'
            data = '{ "ids": [%s] }' % credential['id']
        except KeyError:
            raise MaltegoException("This Credential is not tied to a Metasploit Credential. Please associate it with a Credential before running this transform")
            return response

        delete = apitools.delete_json(url, data, config)
        return response
