#!/usr/bin/env python3.7.2

# -------------------- Imports --------------------- #

# Base Transforms
from EffectiveCouscous.transforms.metasploit.db.netblock.update import (NetblockPutWorkspace,
                                                                        NetblockGetWorkspace,
                                                                        NetblockDeleteWorkspace)
# Custom Entities
from EffectiveCouscous.entities.host.base import MetasploitHost

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



# Workspaces -------------------------------------------------------------------------------------------------------------------#


class HostPutWorkspace(NetblockPutWorkspace):
    """ Push Maltego properties to Metasploit workspace."""

    display_name = 'PUT Host Workspace'
    transform_set = 'Msf__DB                    | Workspace           | Update'
    input_type = MetasploitHost


class HostGetWorkspace(NetblockGetWorkspace):
    """ Pull Metasploit properties to Maltego workspace."""

    display_name = 'GET Host Workspace'
    transform_set = 'Msf__DB                    | Workspace           | Update'
    input_type = MetasploitHost 


class HostDeleteWorkspace(NetblockDeleteWorkspace):
    """ Delete Workspace associated to this Netblock."""

    display_name = 'DELETE Host Workspace'
    transform_set = 'Msf__DB                    | Workspace           | Update'
    input_type = MetasploitHost 



# Host -------------------------------------------------------------------------------------------------------------------------#


class PutHost(Transform):
    """ Push Maltego properties to Metasploit Host"""
    
    display_name = 'PUT Host'
    transform_set = 'Msf__DB                    | Host                     | Update'
    input_type = MetasploitHost

    def do_transform(self, request, response, config):
        msf_host = request.entity

        # Test for properties
        try:
            test = msf_host['id']
        except KeyError:
            raise MaltegoException("This Host is not tied to a Metasploit Host. \
                                    Please associate it with a Host before running this transform")
            return response

        dict = {}
        dict['id'] = msf_host['id']
        dict['address'] = msf_host['address']
        dict['mac'] = msf_host['mac']
        dict['comm'] = msf_host.comm
        dict['name'] = msf_host.name
        dict['state'] = msf_host['state']
        dict['os_name'] = msf_host['os_name']
        dict['os_flavor'] = msf_host['os_flavor']
        dict['os_sp'] = msf_host['os_sp']
        dict['os_lang'] = msf_host['os_lang']
        dict['os_family'] = msf_host['os_family']
        dict['arch'] = msf_host['arch']
        dict['detected_arch'] = msf_host['detected_arch']
        dict['workspace_id'] = msf_host['workspace_id']
        dict['purpose'] = msf_host['purpose']
        dict['info'] = msf_host['info']
        dict['comments'] = msf_host['comments']
        dict['scope'] = msf_host['scope']
        dict['virtual_msf_host'] = msf_host['virtual_host']
        dict['note_count'] = msf_host['note_count']
        dict['vuln_count'] = msf_host.vuln_count
        dict['service_count'] = msf_host.service_count
        dict['host_detail_count'] = msf_host.msf_host_detail_count
        dict['exploit_attempt_count'] = msf_host.exploit_attempt_count
        dict['cred_count'] = msf_host.cred_count
        dict['created_at'] = msf_host['created_at']
        dict['updated_at'] = msf_host['updated_at']

        data = json.dumps(dict)
        url = config['EffectiveCouscous.local.baseurl'] + 'hosts/{0}'.format(msf_host['id'])
        update = apitools.put_json(url, data, config)
        return response


class GetHost(Transform):
    """ Pull Metasploit properties to Maltego Host"""

    display_name = 'GET Host'
    transform_set = 'Msf__DB                    | Host                     | Update'
    input_type = MetasploitHost

    def do_transform(self, request, response, config):
        msf_host = request.entity

        # Test for properties
        try:
            url = config['EffectiveCouscous.local.baseurl'] + 'hosts/{0}'.format(msf_host['id'])
            host = apitools.get_json_dict(url, config)
        except KeyError:
            raise MaltegoException("This Host is not tied to a Metasploit Host. Please associate it with a Host before running this transform")
            return response

        msf_host['id'] = host['id']
        msf_host['ipv4-address'] = host['address']
        msf_host['mac'] = host['mac']
        msf_host['comm'] = host['comm']
        msf_host['name'] = host['name']
        msf_host['state'] = host['state']
        msf_host['os_name'] = host['os_name']
        msf_host['os_flavor'] = host['os_flavor']
        msf_host['os_sp'] = host['os_sp']
        msf_host['os_lang'] = host['os_lang']
        msf_host['os_family'] = host['os_family']
        msf_host['arch'] = host['arch']
        msf_host['detected_arch'] = host['detected_arch']
        msf_host['workspace_id'] = host['workspace_id']
        msf_host['purpose'] = host['purpose']
        msf_host['info'] = host['info']
        msf_host['comments'] = host['comments']
        msf_host['scope'] = host['scope']
        msf_host['virtual_host'] = host['virtual_host']
        msf_host['note_count'] = host['note_count']
        msf_host['vuln_count'] = host['vuln_count']
        msf_host['service_count'] = host['service_count']
        msf_host['host_detail_count'] = host['host_detail_count']
        msf_host['exploit_attempt_count'] = host['exploit_attempt_count']
        msf_host['cred_count'] = host['cred_count']
        msf_host['created_at'] = host['created_at']
        msf_host['updated_at'] = host['updated_at']
        return response


class DeleteHost(Transform):
    """ Delete Metasploit Host associated to this Host Entity"""

    display_name = 'DELETE Host'
    transform_set = 'Msf__DB                    | Host                     | Update'
    input_type = MetasploitHost 

    def do_transform(self, request, response, config):
        msf_host = request.entity

        # Test for properties
        try:
            test = msf_host['id']
        except KeyError:
            raise MaltegoException("This Host is not tied to a Metasploit Host. \
                                    Please associate it with a Host before running this transform")
            return response

        url = config['EffectiveCouscous.local.baseurl'] + 'hosts'
        data = '{ "ids": [%s] }' % msf_host['id']
        delete = apitools.delete_json(url, data, config)
        return response
