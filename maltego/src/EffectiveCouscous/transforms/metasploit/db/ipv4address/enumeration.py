#!/usr/bin/env python3.7.2

# -------------------- Imports ----------------------- #

# Custom Entities
from EffectiveCouscous.entities.infrastructure import IPv4Address

# Maltego Message
from canari.maltego.message import Field, MaltegoException, LinkStyle, LinkColor

# Maltego Transforms
from canari.maltego.transform import Transform

# API
import json
from EffectiveCouscous.tools.msftools import apitools

# OS Factory
from EffectiveCouscous.tools.entitytools.os_factory import getOsEntity

# Service Factory
from EffectiveCouscous.tools.entitytools.service_factory import getServiceEntity

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



# IPv4Addresses ------------------------------------------------------------------------------------------------------------------------------------------#


class EnumerateSubnetIP(Transform):
    """ Enumerates all other IPv4 Addresses in Subnet"""

    display_name = "To Subnet IPs"
    transform_set = 'Msf__DB                    | IPv4Address         | Enumerate'
    input_type = IPv4Address 

    def do_transform(self, request, response, config):
        ip = request.entity
        
        # Test for properties -------------------------------------------//
        try:
            ws_url = config['EffectiveCouscous.local.baseurl'] + 'workspaces/{0}'.format(ip['workspace_id'])
            workspace = apitools.get_json_dict(ws_url, config)['name']
            url = config['EffectiveCouscous.local.baseurl'] + 'hosts' 
            params = (('workspace', '{0}'.format(workspace)),)
            hosts = apitools.get_json_dict(url, config, params=params)
        except KeyError:
            raise MaltegoException("This IPv4Address is not tied to a Metasploit Host. \
                                    Please associate it with a Host before running this transform")
            return response

        # Enumerate IP addresses ----------------------------------------//
        for host in hosts:
            if host['id'] == ip.host_id:
                pass
            else:
                ip_entity = IPv4Address()
                ip_entity['ipv4-address'] = host['address']
                ip_entity.host_id = host['id'] 
                ip_entity.workspace_id = host['workspace_id']
                ip_entity.icon_url = network_interface
                ip_entity.origin_tool = 'Metasploit'
                # Link Style
                ip_entity.link_color = LinkColor.LightGray
                ip_entity.link_thickness = 4
                response += ip_entity 

        return response




# Hosts --------------------------------------------------------------------------------------------------------------------------------------------------#


class EnumerateHost(Transform):
    """ Enumerates a Host with this IP address"""

    display_name = 'To Host'
    transform_set = 'Msf__DB                    | IPv4Address         | Enumerate'
    input_type = IPv4Address

    def do_transform(self, request, response, config):
        ip_address = request.entity 

        # Test for properties
        try:
            url = config['EffectiveCouscous.local.baseurl'] + 'hosts/{0}'.format(ip_address['id'])
            host = apitools.get_json_dict(url, config)
        except KeyError:
            raise MaltegoException("This IPv4Address is not tied to a Metasploit Host. \
                                    Please associate it with a Host before running this transform")
            return response

        # Spawn Host (NEW)
        h = getOsEntity(host['os_name'], host['name'])

        h.ipv4address = ip_address['ipv4-address']
        h.id = '-' if host['id'] is None else host['id']
        h.mac = '-' if host['mac'] is None else host['mac']
        h.comm = '-' if host['comm'] == '' else host['comm']
        h.name = '-' if host['name'] is None else host['name']
        h.state = '-' if host['state'] is None else host['state']
        h.os_family = '-' if host['os_family'] is None else host['os_family']
        h.os_name = '-' if host['os_name'] is None else host['os_name']
        h.os_flavor = '-' if host['os_flavor'] is None else host['os_flavor']
        h.os_sp = '-' if host['os_sp'] is None else host['os_sp']
        h.os_lang = '-' if host['os_lang'] is None else host['os_lang']
        h.arch = '-' if host['arch'] is None else host['arch']
        h.workspace_id = '-' if host['workspace_id'] is None else host['workspace_id']
        h.purpose = '-' if host['purpose'] is None else host['purpose']
        h.info = '-' if host['info'] is None else host['info']
        h.comments = '-' if host['comments'] is None else host['comments']
        h.scope = '-' if host['scope'] is None else host['scope']
        h.virtual_host = '-' if host['virtual_host'] is None else host['virtual_host']
        h.note_count = '-' if host['note_count'] is None else host['note_count']
        h.vuln_count = '-' if host['vuln_count'] is None else host['vuln_count']
        h.service_count = '-' if host['service_count'] is None else host['service_count']
        h.host_detail_count = '-' if host['host_detail_count'] is None else host['host_detail_count']
        h.exploit_attempt_count = '-' if host['exploit_attempt_count'] is None else host['exploit_attempt_count']
        h.cred_count = '-' if host['cred_count'] is None else host['cred_count']
        h.detected_arch = '-' if host['detected_arch'] is None else host['detected_arch']
        h.created_at = host['created_at']
        h.updated_at = host['updated_at']
        # Origin Tool
        h.origin_tool = 'Metasploit'
        # Link Style
        h.link_color = LinkColor.Black
        h.link_thickness = 3

        response += h 
        return response



# Services ----------------------------------------------------------------------------------------------------------------------------------------------#


class EnumerateServices(Transform):
    """ Enumerates Services bound to this IP address"""

    display_name = 'To Services'
    transform_set = 'Msf__DB                    | IPv4Address         | Enumerate'
    input_type = IPv4Address

    def do_transform(self, request, response, config):
        ip = request.entity

        # Test for properties ------------------------------------------------//
        try:
            url = config['EffectiveCouscous.local.baseurl'] + 'workspaces/{0}'.format(ip['workspace_id'])
            workspace = apitools.get_json_dict(url, config)['name']
            s_url = config['EffectiveCouscous.local.baseurl'] + 'services'
            params = (('workspace', '{0}'.format(workspace)),)
            services = apitools.get_json_dict(s_url, config, params=params)
        except KeyError:
            raise MaltegoException("This IPv4Address is not tied to a Metasploit Host. \
                                    Please associate it with a Host before running this transform")
            return response

        # Fill properties ----------------------------------------------------//
        for service in services:
            if int(service['host']['id']) == int(ip['id']):
                # Get Service Entity
                msf_service = getServiceEntity(service['name'], service['info'])
                
                msf_service.display = "{port}:{proto}/{name}".format(port=service['port'],
                                                                proto=service['proto'],
                                                                name=service['name'])
                msf_service.info = service['info']
                msf_service.name = service['name']
                msf_service.proto = service['proto']
                msf_service.port = service['port']
                msf_service.host_id = service['host_id']
                msf_service.id = service['id']
                msf_service.workspace_id = ip['workspace_id']
                msf_service.created_at = service['created_at']
                msf_service.updated_at = service['updated_at']
                msf_service.state = service['state']

                # Icons
                msf_service.origin_tool = 'Metasploit'
                # Link Style
                msf_service.link_color = LinkColor.LightGray
                msf_service.link_thickness = 3
                
                response += msf_service 

        return response

