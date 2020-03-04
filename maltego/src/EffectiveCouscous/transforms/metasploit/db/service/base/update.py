#!/usr/bin/env python3.7.2

# -------------------- Imports --------------------- #

# Custom Entities
from EffectiveCouscous.entities.service.base import MetasploitService

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



# Services ---------------------------------------------------------------------------------------------------------------------------------------------#


class PutService(Transform):
    """ Push Maltego properties to Metasploit Service"""

    display_name = 'PUT Service'
    transform_set = 'Msf__DB                    | Service                 | Update'
    input_type = MetasploitService 

    def do_transform(self, request, response, config):
        msf_service = request.entity

       # Test for properties ------------------------------------------------//
        try:
            test = msf_service['id']
        except KeyError:
            raise MaltegoException("This Service is not tied to a Metasploit Service. \
                                    Please associate it with a Service before running this transform")
            return response

        dict = {}
        dict['id'] = msf_service.workspace_id
        dict['name'] = msf_service.name
        dict['created_at'] = msf_service.created_at
        dict['updated_at'] = msf_service.updated_at
        dict['info'] = msf_service.info
        dict['proto'] = msf_service.proto
        dict['port'] = msf_service.port
        dict['host_id'] = msf_service.host_id
        dict['state'] = msf_service.state

        data = json.dumps(dict)
        url = config['EffectiveCouscous.local.baseurl'] + 'services/{0}'.format(msf_service['id'])
        update = apitools.put_json(url, data, config)
        return response


class GetService(Transform):
    """ Pull Metasploit properties to Maltego Service"""

    display_name = 'GET Service'
    transform_set = 'Msf__DB                    | Service                 | Update'
    input_type = MetasploitService 

    def do_transform(self, request, response, config):
        msf_service = request.entity

        # Test for properties ----------------------------------------------//
        try:
            test = msf_service['id']
        except KeyError:
            raise MaltegoException("This Service is not tied to a Metasploit Service. \
                                    Please associate it with a Service before running this transform")
            return response

        url = config['EffectiveCouscous.local.baseurl'] + 'services/{0}'.format(msf_service['id'])
        service = apitools.get_json_dict(url, config)
        msf_service.info = service['info']
        msf_service.name = service['name']
        msf_service.proto = service['proto']
        msf_service.port = service['port']
        msf_service.host_id = service['host']['id']
        msf_service.service_id = service['id']
        msf_service.workspace_id = service['host']['workspace_id']
        msf_service.created_at = service['created_at']
        msf_service.updated_at = service['updated_at']
        msf_service.state = service['state']
        raise MaltegoException('Due to implementation limitations, it is not possible to refresh this Entity "as is". \
                                Please suppress it and spawn it again from its parent Entity to see it with its new properties.')
        return response


class DeleteService(Transform):
    """ Delete Metasploit Service associated to this Service Entity"""

    display_name = 'DELETE Service'
    transform_set = 'Msf__DB                    | Service                 | Update'
    input_type = MetasploitService 

    def do_transform(self, request, response, config):
        msf_service = request.entity

        # Test for properties ---------------------------------------------//
        try:
            url = config['EffectiveCouscous.local.baseurl'] + 'services'
            data = '{ "ids": [%s] }' % msf_service['id']
        except KeyError:
            raise MaltegoException("This Service is not tied to a Metasploit Service. \
                                    Please associate it with a Service before running this transform")
            return response

        delete = apitools.delete_json(url, data, config)
        return response
