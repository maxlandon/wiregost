#!/usr/bin/env python3.7.2

# -------------------- Imports ----------------------- #

from pkg_resources import resource_filename
from os import path
# ---------------------------------------------------- #


__author__ = 'Maxime Landon'
__copyright__ = 'Copyright 2019, EffectiveCouscous Project'
__credits__ = []

__license__ = 'GPL'
__version__ = '0.2'
__maintainer__ = 'Maxime Landon'
__email__ = 'maximelandon@gmail.com'
__status__ = 'Development'

images = 'EffectiveCouscous.resources.images'
etc = 'EffectiveCouscous.resources.etc'



# Image Directory Search Functions   ------------------------- #

def imageicon(cat, name):
    return 'file://%s' % resource_filename('.'.join([ images, cat ]), name)

def imagepath(cat, name):
    return '%s' % resource_filename('.'.join([ images, cat ]), name)
#------------------------------------------------------------- #



#####################################################################################################
#                                              IMAGES                                               #
#####################################################################################################


# Operating Systems --------------------------------------------------------------------------#

systems = dict(
    # Generic ------------------------------- /
    generic = imageicon('os', 'generic_host.png' ),
    # Linux --------------------------------- /
    archlinux = imageicon('os', 'archlinux.png'),
    debian = imageicon('os', 'debian.png'),
    gentoo = imageicon('os', 'gentoo.png'),
    linux = imageicon('os', 'linux.png'),
    ubuntu = imageicon('os', 'ubuntu.png'),

    # Microsoft ----------------------------- /
    windows = imageicon('os', 'windows.png'),
    windows2000 = imageicon('os', 'windows2000.jpeg'),
    windowsxp = imageicon('os', 'windowsxp.jpeg'),
    windows2003 = imageicon('os', 'windows2003.jpg'),
    windows2008 = imageicon('os', 'windows2008.jpeg'),
    #  windowsvista
    windows2012 = imageicon('os', 'windows-server-2012.jpg'),
    windows7 = imageicon('os', 'windows7.jpeg'),
    windows10 = imageicon('os', 'windows-10.jpg'),
    

    # HP ------------------------------------ /
    hp = imageicon('os', 'hp.png'),

    # Apple --------------------------------- /
    apple = imageicon('os', 'apple.gif'),

    # Cisco --------------------------------- /
    cisco = imageicon('os', 'cisco.gif'),

    # BSD ----------------------------------- /
    freebsd = imageicon('os', 'freebsd.png'),

    # Others -------------------------------- /

)


# Devices ------------------------------------------------------------------------------------#

devices = dict(
    # Apple --------------------------------- /
    ipad = imageicon('device', 'ipad.png'),
    macbook = imageicon('device', 'macbook.png')

    # Others -------------------------------- /
)


# Networking ---------------------------------------------------------------------------------#

    # Ports ----------------------- /
unavailableport = imageicon('networking', 'unavailableport.gif')
openport = imageicon('networking', 'openport.gif')
timedoutport = imageicon('networking', 'timedoutport.gif')
closedport = imageicon('networking', 'closedport.gif')

    # Interfaces ------------------ /
network_interface = imageicon('networking', 'networkinterface.png')



# Services -----------------------------------------------------------------------------------#

services = dict(
    # __________________________APPLICATION____________________________ #
    # Generic ----------------------------------------------------- / 
    generic_service = imageicon('services', 'metasploit-service.png'),

    # Web --------------------------------------------------------- / 
    generic_webservice = imageicon('services.application.web', 'web-service.png'),
    # Microsoft ......../
    iis_web_service = imageicon('services.application.web', 'microsoft-iis.png'),
    microsoft_http_api = imageicon('services.application.web', 'microsoft-http-api.jpg'),
    # RPC ............../
    # Database ........./
    oracle_xml_db = imageicon('services.application.web', 'oracle_xml-db.png'),
    # Apache .........../
    apache_server = imageicon('services.application.web', 'apache-server.jpg'),
    apache_tomcat = imageicon('services.application.web', 'apache-tomcat.png'),
    apache_php = imageicon('services.application.web', 'apache-php.jpeg'),
    # HTTP File Server../
    http_file_server = imageicon('services.application.web', 'http-file-server.jpeg'),
    # Ruby ............./
    webrick = imageicon('services.application.web', 'ruby-on-rails.png'),
    # Java ............./
    jetty_server = imageicon('services.application.web', 'jetty-server.jpeg'),
    # JavaScript ......./
    node_js = imageicon('services.application.web', 'node-js.png'),
    # Others .........../ 
    nginx_web_service = imageicon('services.application.web', 'nginx.png'),
    glassfish = imageicon('services.application.web.application', 'glassfish.jpg'),

    # RPC ----------------------------------------------------------/
    ms_rpc = imageicon('services.application.rpc', 'rpc-service.png'),
    generic_rpc = imageicon('services.application.rpc', 'rpc-service.png'),
    java_rmi = imageicon('services.application.rpc', 'java-rmi.png'),

    # VPN ----------------------------------------------------------/ 
    cisco_vpn = imageicon('services.application.vpn', 'cisco-vpn-client.png'),
    vpn_service = imageicon('services.application.vpn', 'vpn-service.png'),

    # Samba --------------------------------------------------------/ 
    smb_service = imageicon('services.application.smb', 'smb-service.png'),
    
    # RDP ----------------------------------------------------------/
    rdp_service = imageicon('services.application.rdp', 'rdp-service.png'),

    # SSH ----------------------------------------------------------/
    ssh_service = imageicon('services.application.ssh', 'ssh-service.png'),
    putty = imageicon('services.application.ssh', 'putty.png'),

    # DNS ----------------------------------------------------------/
    dns_service = imageicon('services.application.dns', 'DNS.png'),

    # SQL ----------------------------------------------------------/
    sql_service = imageicon('services.application.sql', 'sql-service.png'),
    sql_server = imageicon('services.application.sql', 'sql-server.jpg'),
    postgresql = imageicon('services.application.sql', 'postgresql.jpg'),
    mariadb = imageicon('services.application.sql', 'mariadb.jpg'),
    mysql = imageicon('services.application.sql', 'my-sql.svg'),

    # Virtualization -----------------------------------------------/
    virtualization_software = imageicon('services.application.virtualization', 'virtualization-software.jpg'),
    vwmare_workstation = imageicon('services.application.virtualization', 'vmware.png'),

    # Others -------------------------------------------------------/
    elastic_search_api = imageicon('services.application', 'elastic-search.png'),



    )



# Tools --------------------------------------------------------------------------------------#

tools = dict(
    metasploit = imageicon('logos', 'metasploit.png'),
    nmap = imageicon('logos', 'nmap.gif'),
    nessus = imageicon('logos', 'nessus.png')
    )





# -------------------------    MISCELLANEOUS   -------------------------- #

# flag
def flag(c):
    f = imageicon('flags', '%s.png' % c.lower())
    if path.exists(f[7:]):
        return f
    return None
