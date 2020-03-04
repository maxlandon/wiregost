#!/usr/bin/env python3.7.2

# -------------------- Imports --------------------- #

# Maltego Entities
from canari.maltego.entities import Service

# Custom Entities
from EffectiveCouscous.entities.service.base import MetasploitService

# Maltego Transforms
from canari.maltego.transform import Transform

# API
import json
from EffectiveCouscous.tools.msftools import apitools

# GUI
import canari.easygui as gui
#--------------------------------------------------- #


__author__ = 'Maxime Landon'
__copyright__ = 'Copyright 2019, EffectiveCouscous Project'
__credits__ = []

__license__ = 'GPLv3'
__version__ = '0.2'
__maintainer__ = 'Maxime Landon'
__email__ = 'maximelandon@gmail.com'
__status__ = 'Development'




# Services ---------------------------------------------------------------------------------------------------------------------------------------------#


class ToMetasploitService(Transform):
    """Spawns a Metasploit Service"""

    display_name = 'To Metasploit Service'
    transform_set = 'Msf__DB                    | Service                 | Associate'
    input_type = Service 

    def do_transform(self, request, response, config):
        service_entity = request.entity

        # Select workspace & Service --------------------------------------//
        url = config['EffectiveCouscous.local.baseurl'] + 'workspaces' 
        workspaces = apitools.get_json_dict(url, config)
        title = "Workspace Choice"
        msg = """Choose a Metasploit Workspace for Service selection"""
        ws_names = [ workspace['name'] for workspace in workspaces ]
        ws_choice = gui.choicebox(title=title, msg=msg, choices=(ws_names))

        # Select Service -------------------------------------------------//
        service_url = config['EffectiveCouscous.local.baseurl'] + 'services'
        params = (('workspace', '{0}'.format(ws_choice)),)
        services = apitools.get_json_dict(service_url, config, params=params)
        title = "Service Choice"
        msg = """Choose a Metasploit Service to associate with this Service"""
        service_names = []
        service_infos = []
        for service in services:
            info = '%s      %s       %s' % (service['host']['address'], service['port'], service['info'])
            service_infos.append(info)
            service_names.append(service['info'])
        service_infos.append("Add Service")
        raw_choice = gui.choicebox(title=title, msg=msg, choices=(service_infos))
        service = {} 
        if "Add Service" in raw_choice:
            service['info'] = "Add Service" 
        else:
            for s in services:
                if (s['info'] in raw_choice) and (str(s['port']) in raw_choice):
                    service = s

        # If existing Service --------------------------------------// 
        if service['info'] != "Add Service":
            msf_service = getServiceEntity(service['name'], service['info'])

            if service['info'] == '': msf_service.info = '-'
            else: msf_service.info = service['info']
            if service['name'] == '': msf_service.name = '-'
            else: msf_service.name = service['name']
            if service['proto'] == '': msf_service.proto = '-'
            else: msf_service.proto = service['proto']
            if service['port'] == '': msf_service.port = '-'
            else: msf_service.port = service['port']
            if service['host']['id'] is None: msf_service.host_id = '-'
            else: msf_service.host_id = service['host']['id']
            if service['id'] == '': msf_service.id = '-'
            else: msf_service.service_id = service['id']
            if service['host']['workspace_id'] == '': msf_service.workspaceid = '-'
            else: msf_service.workspace_id = service['host']['workspace_id']
            msf_service.display = "{port}:{proto}/{name}".format(port=service['port'],
                                                                proto=service['proto'],
                                                                name=service['name'])
            msf_service.state = service['state']
            msf_service.created_at = service['created_at']
            msf_service.updated_at = service['updated_at']

            response += msf_service 

        # If new Service -------------------------------------------------//
        if service['info'] == "Add Service":
            title = "New Service"
            msg = "Add properties to create a Service in Metasploit"
            field_names = ['Workspace', 'Host IP', 'Port number', 'Protocol',
                            'Service Name', 'Text (Info)', 'State']  
            field_values = []
            field_values = gui.multenterbox(msg, fields=field_names)
            while 1:                                                 
                if field_values == None: break                        
                errmsg = ""                                          
                for i in range(len(field_names)):                     
                    if field_values[i].strip() == "":                 
                        errmsg += ('"%s" is a required field.\n\n' % field_names[i])                                       
                if errmsg == "":                                     
                    break 
                field_values = gui.multenterbox(errmsg, field_values, fields=field_names)                                         

            # Create Service in Metasploit
            dict = {}
            dict['workspace'] = field_values[0]
            dict['host'] = field_values[1]
            dict['port'] = field_values[2]
            dict['proto'] = field_values[3]
            dict['name'] = field_values[4]
            dict['info'] = field_values[5]
            dict['state'] = field_values[6]
            data = json.dumps(dict)
            post = apitools.post_json(service_url, data) 

            # Fetch new Service in Metasploit
            new = post.json()['data']

            msf_service = getServiceEntity(service['name'], service['info'])
            msf_service.info = new['info']
            msf_service.name = new['name']
            msf_service.proto = new['proto']
            msf_service.host_id = new['host']['id']
            msf_service.id = new['id']
            msf_service.workspace_id = new['host']['workspace_id']
            msf_service.display = "{port}:{proto}/{name}".format(port=new['port'],
                                                                proto=new['proto'],
                                                                name=new['name'])
            msf_service.state = new['state']
            msf_service.created_at = new['created_at']
            msf_service.updated_at = new['updated_at']

            response += msf_service
         
        return response


class PostMetasploitService(Transform):
    """ Spawns a Metasploit Service based on the properties of this Entity"""

    display_name = 'POST Metasploit Service'
    transform_set = 'Msf__DB                    | Service                 | Associate'
    input_type = MetasploitService 

    def do_transform(self, request, response, config):
        service_entity = request.entity

        # Find workspace with ID -------------------------------//
        url = config['EffectiveCouscous.local.baseurl'] + 'workspaces' 
        workspaces = apitools.get_json_dict(url, config)
        ws_name = ''
        for workspace in workspaces:
            if workspace['id'] == service_entity.workspace_id:
                ws_name = workspace['name']

        # Create Service in Metasploit -------------------------//
        dict = {}
        dict['workspace'] = ws_name 
        dict['host'] = service_entity.host_id 
        dict['port'] = service_entity.port 
        dict['proto'] = service_entity.proto 
        dict['name'] = service_entity.name 
        dict['info'] = service_entity.info 
        dict['state'] = service_entity.state 
        data = json.dumps(dict)
        post = apitools.post_json(service_url, data) 

        # Fetch new Service in Metasploit ----------------------//
        new = post.json()['data']

        msf_service = getServiceEntity(new['name'], new['info'])
        msf_service.info = new['info']
        msf_service.name = new['name']
        msf_service.proto = new['proto']
        msf_service.host_id = new['host']['id']
        msf_service.id = new['id']
        msf_service.workspace_id = new['host']['workspace_id']
        msf_service.display = "{port}:{proto}/{name}".format(port=new['port'],
                                                            proto=new['proto'],
                                                            name=new['name'])
        msf_service.state = new['state']
        msf_service.created_at = new['created_at']
        msf_service.updated_at = new['updated_at']

        response += msf_service
        return response

