#!/usr/bin/env python3.7.2

# -------------------- Imports --------------------- #

# Web 
from EffectiveCouscous.entities.service.application.web import *
# VPN 
from EffectiveCouscous.entities.service.application.vpn import *
# RPC
from EffectiveCouscous.entities.service.application.rpc import *
# SMB
from EffectiveCouscous.entities.service.application.smb import *
# RDP
from EffectiveCouscous.entities.service.application.rdp import *
# SSH
from EffectiveCouscous.entities.service.application.ssh import *
# DNS
from EffectiveCouscous.entities.service.application.dns import *
# SQL
from EffectiveCouscous.entities.service.application.sql import *
# Virtualization
from EffectiveCouscous.entities.service.application.virtualization import *

# Icons
from EffectiveCouscous.resource import services
# -------------------------------------------------- #


__author__ = 'Maxime Landon'
__copyright__ = 'Copyright 2019, EffectiveCouscous Project'
__credits__ = []

__license__ = 'GPLv3'
__version__ = '0.2'
__maintainer__ = 'Maxime Landon'
__email__ = 'maximelandon@gmail.com'
__status__ = 'Development'




# Application Layer Protocol Identification Strings -------------------------------------------------------------------------------------------- #

web_protocols = ['http', 
                'https', 
                'possible_wls', 
                "www", 
                "ncacn_http", 
                "ccproxy-http", 
                "ssl/http", 
                "http-proxy", 
                'wsdapi'        # Windows Http API
                ]
samba_protocols = ['samba', 
                'netbios-ssn', 
                'smb', 
                'microsoft-ds', 
                'netbios-ns', 
                'netbios-dgm', 
                'netbios',
                ]
vpn_protocols = ['vpn', ]
rpc_protocols = ['rpc', 
                'msrpc',
                ]
rdp_protocols = ['ms-wbt-server',
                'rdp',
                ]
ssh_protocols = ['ssh', ]
dns_protocols = ['domain', ]
sql_protocols = ['mysql', ]


# Service Info & Banner Identification Strings ------------------------------------------------------------------------------------------------- #
web_application_infos = ['tomcat',
                        'nodejs',
                        'jetty',
                        'iis',
                        'glassfish',]
virtualization_infos = ['vmware',]
others_infos = ['elasticsearch',]



# Entity Generation --------------------------------------------------------------------------------------------------------------------------- #

def getServiceEntity(service_name, service_info):
    from EffectiveCouscous.entities.service.base import MetasploitService

    service_entity = MetasploitService()
    name = service_name.lower()
    info = service_info.lower()

    if service_name or service_info:
         #  Test directly for service names & service info: unlikely to have info without name

        #  Web ------------------------------------------------------------------------------------- //
        if any(x in name for x in web_protocols):
            service_entity = WebService()
            if service_info:
                #  Web Application ------------------- //
                if any(x in info for x in web_application_infos):
                    if 'glassfish' in info:
                        service_entity = OracleGlassFish()
                    elif 'iis' in info:
                        service_entity = IISWebService()
                    elif 'jetty' in info:
                        service_entity = Jetty()
                    elif 'tomcat' in info:
                        service_entity = ApacheTomcat()
                    elif 'nodejs' in info:
                        service_entity = NodeJS()
                    return service_entity 

                # Microsoft ............................/
                elif 'httpapi' in info:
                    service_entity = MicrosoftHTTPAPI()
                # RPC ................................../
                elif 'rpc over http' in info:
                    service_entity = RPCoverHttp()
                # Databases ............................/
                elif 'oracle xml db' in info:
                    service_entity = OracleXMLDB()
                # Apache .............................../
                elif 'apache' in info:
                    if 'apache php' in info:
                        service_entity = ApachePHP()
                    if 'apache httpd' in info:
                        service_entity = ApacheHttpd()
                # HTTP File Servers ..................../
                elif 'httpfileserver' in info:
                    service_entity = HTTPFileServer()
                # Ruby on Rails Servers ................/
                elif 'webrick' in info:
                    service_entity = WEBrick()
                # Java ................................./
                # Other Web Servers ..................../
                elif 'lighttpd' in info:
                    service_entity = Lighttpd()
                elif 'nginx' in info:
                    service_entity = Nginx()
                elif 'waf' in info:
                    service_entity = WAF()
                elif 'oracle http server' in info:
                    service_entity = OracleHTTPServer()
                elif 'goahead' in info:
                    service_entity = GoAheadWebServer()
                elif 'webmin' in info:
                    service_entity = Webmin()
                elif 'rocket' in info:
                    service_entity = RocketWebServer()
                elif 'squid' in info:
                    service_entity = SquidProxyServer()
                elif 'communigate' in info:
                    service_entity = CommuniGatePro()

                """ Return Web Service if one is found """                 
                return service_entity

        #  VPN ------------------------------------------------------------------------------------- //
        # TO BE CHANGED TO ACCOUNT FOR SPECIFIC VENDORS/SOFTWARE
        elif any(x in name for x in vpn_protocols):
            # Cisco ................................/
            if 'vpn' in info:
                service_entity = CiscoVPN()
                
        #  RPC ------------------------------------------------------------------------------------- //
        elif any(x in name for x in rpc_protocols):
            # Microsoft ............................/
            if ('msprc' in name) or ('windows rpc' in info):
                service_entity = MicrosoftWindowsRPC()
                return service_entity
            # Generic ............................../
            elif 'rpc' in name:
                service_entity = RPCService()
     

        #  SMB --------------------------------------------------------------------------------------//
        elif any(x in name for x in samba_protocols):
            service_entity = SMBService()


        #  RDP ------------------------------------------------------------------------------------- //
        elif any(x in name for x in rdp_protocols):
            if 'ms-wbt-server' in name:
                service_entity = MicrosoftWindowsTerminal()
            elif 'rdp' in name:
                service_info = RDPService()


        #  SSH ------------------------------------------------------------------------------------- //
        elif any(x in name for x in ssh_protocols):
            if 'openssh' in info:
                service_entity = OpenSSH()
            elif 'weonlydo ssh' in info:
                service_entity = WeOnlyDoSSH()
            # Add filter for puTTY SSH Service
            else:
                service_entity = SSHService()


        #  DNS ------------------------------------------------------------------------------------- //
        elif any(x in name for x in dns_protocols):
            if 'dnsmasq' in info:
                service_entity = DNSService()
            else:
                service_entity = DNSService()


        #  SQL ------------------------------------------------------------------------------------- //
        elif any(x in name for x in sql_protocols):
            if 'mysql' in name:
                service_entity = MySQL()
                return service_entity
            elif 'sql' in name:
                service_entity = SQLService()


        #  Virtualization -------------------------------------------------------------------------- //
        elif any(x in info for x in virtualization_infos):
            if 'vmware' in info:
                service_entity = VMWareWorkstation()
            else:
                service_entity = VirtualizationSoftware()


        #  Others ---------------------------------------------------------------------------------- //
    elif any(x in info for x in others_infos):
        if 'elasticsearch' in info:
            service_entity = 'test' 



    # Return Service Entity anyway
    return service_entity




