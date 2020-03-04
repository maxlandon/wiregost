#!/usr/bin/env python3.7.2

# -------------------- Imports --------------------- #

# Maltego Entities
from canari.maltego.entities import Entity

# Fields
from canari.maltego.message import *

# System-wide Icons
from EffectiveCouscous.tools.entitytools.icon_factory import getOriginTool, getStateIcon
# -------------------------------------------------- #


__author__ = 'Maxime Landon'
__copyright__ = 'Copyright 2019, EffectiveCouscous Project'
__credits__ = []

__license__ = 'GPLv3'
__version__ = '0.2'
__maintainer__ = 'Maxime Landon'
__email__ = 'maximelandon@gmail.com'
__status__ = 'Development'


#####################################################################################################
#                                           WEB SERVICES                                            #
#####################################################################################################


#  Base Web Service --------------------------------------------------------------------- {{{1 # 

class WebService(Entity):
    _category_ = 'Metasploit | Services | Web'
    _namespace_ = 'EffectiveCouscous.MetasploitService'

    # Entity properties
    display = StringEntityField('display', display_name='Display Name', is_value=True)
    proto = StringEntityField('proto', display_name='Protocol')
    name = StringEntityField('name', display_name='Name')
    info = StringEntityField('info', display_name='Info')
    port = StringEntityField('port', display_name='Port')
    state = StringEntityField('state', display_name='State', decorator=getStateIcon)
    host_id = StringEntityField('host_id', display_name='Host ID')
    id = StringEntityField('id', display_name='Service ID')
    workspace_id = StringEntityField('workspace_id', display_name='Workspace ID')
    created_at = StringEntityField('created_at', display_name='Created At')
    updated_at = StringEntityField('updated_at', display_name='Updated At')

    # Decorator properties
    origin_tool = StringEntityField('origin_tool', display_name='Origin Tool', decorator=getOriginTool)
    tool_icon = StringEntityField('tool_icon', display_name='Tool Icon')
    state_icon = StringEntityField('state_icon', display_name='State Icon')


#  Microsoft ---------------------------------------------------------------------------- {{{2 # 

class MicrosoftHTTPAPI(WebService):
    _category_ = 'Metasploit | Services | Web'
    _namespace_ = 'EffectiveCouscous.MetasploitService.WebService'

    # Entity properties
    display = StringEntityField('display', display_name='Display Name', is_value=True)
    proto = StringEntityField('proto', display_name='Protocol')
    name = StringEntityField('name', display_name='Name')
    info = StringEntityField('info', display_name='Info')
    port = StringEntityField('port', display_name='Port')
    state = StringEntityField('state', display_name='State', decorator=getStateIcon)
    host_id = StringEntityField('host_id', display_name='Host ID')
    id = StringEntityField('id', display_name='Service ID')
    workspace_id = StringEntityField('workspace_id', display_name='Workspace ID')
    created_at = StringEntityField('created_at', display_name='Created At')
    updated_at = StringEntityField('updated_at', display_name='Updated At')

    # Decorator properties
    origin_tool = StringEntityField('origin_tool', display_name='Origin Tool', decorator=getOriginTool)
    tool_icon = StringEntityField('tool_icon', display_name='Tool Icon')
    state_icon = StringEntityField('state_icon', display_name='State Icon')



#  RPC ---------------------------------------------------------------------------------- {{{3 # 

class RPCoverHttp(WebService):
    _category_ = 'Metasploit | Services | Web'
    _namespace_ = 'EffectiveCouscous.MetasploitService.WebService'

    # Entity properties
    display = StringEntityField('display', display_name='Display Name', is_value=True)
    proto = StringEntityField('proto', display_name='Protocol')
    name = StringEntityField('name', display_name='Name')
    info = StringEntityField('info', display_name='Info')
    port = StringEntityField('port', display_name='Port')
    state = StringEntityField('state', display_name='State', decorator=getStateIcon)
    host_id = StringEntityField('host_id', display_name='Host ID')
    id = StringEntityField('id', display_name='Service ID')
    workspace_id = StringEntityField('workspace_id', display_name='Workspace ID')
    created_at = StringEntityField('created_at', display_name='Created At')
    updated_at = StringEntityField('updated_at', display_name='Updated At')

    # Decorator properties
    origin_tool = StringEntityField('origin_tool', display_name='Origin Tool', decorator=getOriginTool)
    tool_icon = StringEntityField('tool_icon', display_name='Tool Icon')
    state_icon = StringEntityField('state_icon', display_name='State Icon')



#  Oracle ------------------------------------------------------------------------------- {{{4 # 

class OracleXMLDB(WebService):
    _category_ = 'Metasploit | Services | Web'
    _namespace_ = 'EffectiveCouscous.MetasploitService.WebService'

    # Entity properties
    display = StringEntityField('display', display_name='Display Name', is_value=True)
    proto = StringEntityField('proto', display_name='Protocol')
    name = StringEntityField('name', display_name='Name')
    info = StringEntityField('info', display_name='Info')
    port = StringEntityField('port', display_name='Port')
    state = StringEntityField('state', display_name='State', decorator=getStateIcon)
    host_id = StringEntityField('host_id', display_name='Host ID')
    id = StringEntityField('id', display_name='Service ID')
    workspace_id = StringEntityField('workspace_id', display_name='Workspace ID')
    created_at = StringEntityField('created_at', display_name='Created At')
    updated_at = StringEntityField('updated_at', display_name='Updated At')

    # Decorator properties
    origin_tool = StringEntityField('origin_tool', display_name='Origin Tool', decorator=getOriginTool)
    tool_icon = StringEntityField('tool_icon', display_name='Tool Icon')
    state_icon = StringEntityField('state_icon', display_name='State Icon')



#  Apache ------------------------------------------------------------------------------- {{{5 #

class ApachePHP(WebService):
    _category_ = 'Metasploit | Services | Web'
    _namespace_ = 'EffectiveCouscous.MetasploitService.WebService'

    # Entity properties
    display = StringEntityField('display', display_name='Display Name', is_value=True)
    proto = StringEntityField('proto', display_name='Protocol')
    name = StringEntityField('name', display_name='Name')
    info = StringEntityField('info', display_name='Info')
    port = StringEntityField('port', display_name='Port')
    state = StringEntityField('state', display_name='State', decorator=getStateIcon)
    host_id = StringEntityField('host_id', display_name='Host ID')
    id = StringEntityField('id', display_name='Service ID')
    workspace_id = StringEntityField('workspace_id', display_name='Workspace ID')
    created_at = StringEntityField('created_at', display_name='Created At')
    updated_at = StringEntityField('updated_at', display_name='Updated At')

    # Decorator properties
    origin_tool = StringEntityField('origin_tool', display_name='Origin Tool', decorator=getOriginTool)
    tool_icon = StringEntityField('tool_icon', display_name='Tool Icon')
    state_icon = StringEntityField('state_icon', display_name='State Icon')

class ApacheHttpd(WebService):
    _category_ = 'Metasploit | Services | Web'
    _namespace_ = 'EffectiveCouscous.MetasploitService.WebService'

    # Entity properties
    display = StringEntityField('display', display_name='Display Name', is_value=True)
    proto = StringEntityField('proto', display_name='Protocol')
    name = StringEntityField('name', display_name='Name')
    info = StringEntityField('info', display_name='Info')
    port = StringEntityField('port', display_name='Port')
    state = StringEntityField('state', display_name='State', decorator=getStateIcon)
    host_id = StringEntityField('host_id', display_name='Host ID')
    id = StringEntityField('id', display_name='Service ID')
    workspace_id = StringEntityField('workspace_id', display_name='Workspace ID')
    created_at = StringEntityField('created_at', display_name='Created At')
    updated_at = StringEntityField('updated_at', display_name='Updated At')

    # Decorator properties
    origin_tool = StringEntityField('origin_tool', display_name='Origin Tool', decorator=getOriginTool)
    tool_icon = StringEntityField('tool_icon', display_name='Tool Icon')
    state_icon = StringEntityField('state_icon', display_name='State Icon')



#  HTTP --------------------------------------------------------------------------------- {{{6 # 

class HTTPFileServer(WebService):
    _category_ = 'Metasploit | Services | Web'
    _namespace_ = 'EffectiveCouscous.MetasploitService.WebService'

    # Entity properties
    display = StringEntityField('display', display_name='Display Name', is_value=True)
    proto = StringEntityField('proto', display_name='Protocol')
    name = StringEntityField('name', display_name='Name')
    info = StringEntityField('info', display_name='Info')
    port = StringEntityField('port', display_name='Port')
    state = StringEntityField('state', display_name='State', decorator=getStateIcon)
    host_id = StringEntityField('host_id', display_name='Host ID')
    id = StringEntityField('id', display_name='Service ID')
    workspace_id = StringEntityField('workspace_id', display_name='Workspace ID')
    created_at = StringEntityField('created_at', display_name='Created At')
    updated_at = StringEntityField('updated_at', display_name='Updated At')

    # Decorator properties
    origin_tool = StringEntityField('origin_tool', display_name='Origin Tool', decorator=getOriginTool)
    tool_icon = StringEntityField('tool_icon', display_name='Tool Icon')
    state_icon = StringEntityField('state_icon', display_name='State Icon')



#  Ruby --------------------------------------------------------------------------------- {{{7 # 

class WEBrick(WebService):
    _category_ = 'Metasploit | Services | Web'
    _namespace_ = 'EffectiveCouscous.MetasploitService.WebService'

    # Entity properties
    display = StringEntityField('display', display_name='Display Name', is_value=True)
    proto = StringEntityField('proto', display_name='Protocol')
    name = StringEntityField('name', display_name='Name')
    info = StringEntityField('info', display_name='Info')
    port = StringEntityField('port', display_name='Port')
    state = StringEntityField('state', display_name='State', decorator=getStateIcon)
    host_id = StringEntityField('host_id', display_name='Host ID')
    id = StringEntityField('id', display_name='Service ID')
    workspace_id = StringEntityField('workspace_id', display_name='Workspace ID')
    created_at = StringEntityField('created_at', display_name='Created At')
    updated_at = StringEntityField('updated_at', display_name='Updated At')

    # Decorator properties
    origin_tool = StringEntityField('origin_tool', display_name='Origin Tool', decorator=getOriginTool)
    tool_icon = StringEntityField('tool_icon', display_name='Tool Icon')
    state_icon = StringEntityField('state_icon', display_name='State Icon')


# Other Web Servers --------------------------------------------------------------------------#

class Lighttpd(WebService):
    _category_ = 'Metasploit | Services | Web'
    _namespace_ = 'EffectiveCouscous.MetasploitService.WebService'

    # Entity properties
    display = StringEntityField('display', display_name='Display Name', is_value=True)
    proto = StringEntityField('proto', display_name='Protocol')
    name = StringEntityField('name', display_name='Name')
    info = StringEntityField('info', display_name='Info')
    port = StringEntityField('port', display_name='Port')
    state = StringEntityField('state', display_name='State', decorator=getStateIcon)
    host_id = StringEntityField('host_id', display_name='Host ID')
    id = StringEntityField('id', display_name='Service ID')
    workspace_id = StringEntityField('workspace_id', display_name='Workspace ID')
    created_at = StringEntityField('created_at', display_name='Created At')
    updated_at = StringEntityField('updated_at', display_name='Updated At')

    # Decorator properties
    origin_tool = StringEntityField('origin_tool', display_name='Origin Tool', decorator=getOriginTool)
    tool_icon = StringEntityField('tool_icon', display_name='Tool Icon')
    state_icon = StringEntityField('state_icon', display_name='State Icon')

class Nginx(WebService):
    _category_ = 'Metasploit | Services | Web'
    _namespace_ = 'EffectiveCouscous.MetasploitService.WebService'

    # Entity properties
    display = StringEntityField('display', display_name='Display Name', is_value=True)
    proto = StringEntityField('proto', display_name='Protocol')
    name = StringEntityField('name', display_name='Name')
    info = StringEntityField('info', display_name='Info')
    port = StringEntityField('port', display_name='Port')
    state = StringEntityField('state', display_name='State', decorator=getStateIcon)
    host_id = StringEntityField('host_id', display_name='Host ID')
    id = StringEntityField('id', display_name='Service ID')
    workspace_id = StringEntityField('workspace_id', display_name='Workspace ID')
    created_at = StringEntityField('created_at', display_name='Created At')
    updated_at = StringEntityField('updated_at', display_name='Updated At')

    # Decorator properties
    origin_tool = StringEntityField('origin_tool', display_name='Origin Tool', decorator=getOriginTool)
    tool_icon = StringEntityField('tool_icon', display_name='Tool Icon')
    state_icon = StringEntityField('state_icon', display_name='State Icon')


class WAF(WebService):
    _category_ = 'Metasploit | Services | Web'
    _namespace_ = 'EffectiveCouscous.MetasploitService.WebService'

    # Entity properties
    display = StringEntityField('display', display_name='Display Name', is_value=True)
    proto = StringEntityField('proto', display_name='Protocol')
    name = StringEntityField('name', display_name='Name')
    info = StringEntityField('info', display_name='Info')
    port = StringEntityField('port', display_name='Port')
    state = StringEntityField('state', display_name='State', decorator=getStateIcon)
    host_id = StringEntityField('host_id', display_name='Host ID')
    id = StringEntityField('id', display_name='Service ID')
    workspace_id = StringEntityField('workspace_id', display_name='Workspace ID')
    created_at = StringEntityField('created_at', display_name='Created At')
    updated_at = StringEntityField('updated_at', display_name='Updated At')

    # Decorator properties
    origin_tool = StringEntityField('origin_tool', display_name='Origin Tool', decorator=getOriginTool)
    tool_icon = StringEntityField('tool_icon', display_name='Tool Icon')
    state_icon = StringEntityField('state_icon', display_name='State Icon')


class OracleHTTPServer(WebService):
    _category_ = 'Metasploit | Services | Web'
    _namespace_ = 'EffectiveCouscous.MetasploitService.WebService'

    # Entity properties
    display = StringEntityField('display', display_name='Display Name', is_value=True)
    proto = StringEntityField('proto', display_name='Protocol')
    name = StringEntityField('name', display_name='Name')
    info = StringEntityField('info', display_name='Info')
    port = StringEntityField('port', display_name='Port')
    state = StringEntityField('state', display_name='State', decorator=getStateIcon)
    host_id = StringEntityField('host_id', display_name='Host ID')
    id = StringEntityField('id', display_name='Service ID')
    workspace_id = StringEntityField('workspace_id', display_name='Workspace ID')
    created_at = StringEntityField('created_at', display_name='Created At')
    updated_at = StringEntityField('updated_at', display_name='Updated At')

    # Decorator properties
    origin_tool = StringEntityField('origin_tool', display_name='Origin Tool', decorator=getOriginTool)
    tool_icon = StringEntityField('tool_icon', display_name='Tool Icon')
    state_icon = StringEntityField('state_icon', display_name='State Icon')


class Webmin(WebService):
    _category_ = 'Metasploit | Services | Web'
    _namespace_ = 'EffectiveCouscous.MetasploitService.WebService'

    # Entity properties
    display = StringEntityField('display', display_name='Display Name', is_value=True)
    proto = StringEntityField('proto', display_name='Protocol')
    name = StringEntityField('name', display_name='Name')
    info = StringEntityField('info', display_name='Info')
    port = StringEntityField('port', display_name='Port')
    state = StringEntityField('state', display_name='State', decorator=getStateIcon)
    host_id = StringEntityField('host_id', display_name='Host ID')
    id = StringEntityField('id', display_name='Service ID')
    workspace_id = StringEntityField('workspace_id', display_name='Workspace ID')
    created_at = StringEntityField('created_at', display_name='Created At')
    updated_at = StringEntityField('updated_at', display_name='Updated At')

    # Decorator properties
    origin_tool = StringEntityField('origin_tool', display_name='Origin Tool', decorator=getOriginTool)
    tool_icon = StringEntityField('tool_icon', display_name='Tool Icon')
    state_icon = StringEntityField('state_icon', display_name='State Icon')


class RocketWebServer(WebService):
    _category_ = 'Metasploit | Services | Web'
    _namespace_ = 'EffectiveCouscous.MetasploitService.WebService'

    # Entity properties
    display = StringEntityField('display', display_name='Display Name', is_value=True)
    proto = StringEntityField('proto', display_name='Protocol')
    name = StringEntityField('name', display_name='Name')
    info = StringEntityField('info', display_name='Info')
    port = StringEntityField('port', display_name='Port')
    state = StringEntityField('state', display_name='State', decorator=getStateIcon)
    host_id = StringEntityField('host_id', display_name='Host ID')
    id = StringEntityField('id', display_name='Service ID')
    workspace_id = StringEntityField('workspace_id', display_name='Workspace ID')
    created_at = StringEntityField('created_at', display_name='Created At')
    updated_at = StringEntityField('updated_at', display_name='Updated At')

    # Decorator properties
    origin_tool = StringEntityField('origin_tool', display_name='Origin Tool', decorator=getOriginTool)
    tool_icon = StringEntityField('tool_icon', display_name='Tool Icon')
    state_icon = StringEntityField('state_icon', display_name='State Icon')


class SquidProxyServer(WebService):
    _category_ = 'Metasploit | Services | Web'
    _namespace_ = 'EffectiveCouscous.MetasploitService.WebService'

    # Entity properties
    display = StringEntityField('display', display_name='Display Name', is_value=True)
    proto = StringEntityField('proto', display_name='Protocol')
    name = StringEntityField('name', display_name='Name')
    info = StringEntityField('info', display_name='Info')
    port = StringEntityField('port', display_name='Port')
    state = StringEntityField('state', display_name='State', decorator=getStateIcon)
    host_id = StringEntityField('host_id', display_name='Host ID')
    id = StringEntityField('id', display_name='Service ID')
    workspace_id = StringEntityField('workspace_id', display_name='Workspace ID')
    created_at = StringEntityField('created_at', display_name='Created At')
    updated_at = StringEntityField('updated_at', display_name='Updated At')

    # Decorator properties
    origin_tool = StringEntityField('origin_tool', display_name='Origin Tool', decorator=getOriginTool)
    tool_icon = StringEntityField('tool_icon', display_name='Tool Icon')
    state_icon = StringEntityField('state_icon', display_name='State Icon')


class CommuniGatePro(Entity):
    _category_ = 'Metasploit | Services | Web'
    _namespace_ = 'EffectiveCouscous.MetasploitService.WebService'

    # Entity properties
    display = StringEntityField('display', display_name='Display Name', is_value=True)
    proto = StringEntityField('proto', display_name='Protocol')
    name = StringEntityField('name', display_name='Name')
    info = StringEntityField('info', display_name='Info')
    port = StringEntityField('port', display_name='Port')
    state = StringEntityField('state', display_name='State', decorator=getStateIcon)
    host_id = StringEntityField('host_id', display_name='Host ID')
    id = StringEntityField('id', display_name='Service ID')
    workspace_id = StringEntityField('workspace_id', display_name='Workspace ID')
    created_at = StringEntityField('created_at', display_name='Created At')
    updated_at = StringEntityField('updated_at', display_name='Updated At')

    # Decorator properties
    origin_tool = StringEntityField('origin_tool', display_name='Origin Tool', decorator=getOriginTool)
    tool_icon = StringEntityField('tool_icon', display_name='Tool Icon')
    state_icon = StringEntityField('state_icon', display_name='State Icon')






#  APPLICATION SERVERS ______________________________________________________________________________ {{{1 # 

#  Java -------------------------------------------------------------------------- {{{1 # 
class Jetty(WebService):
    _category_ = 'Metasploit | Services | Web | Application'
    _namespace_ = 'EffectiveCouscous.MetasploitService.WebService'

    # Entity properties
    display = StringEntityField('display', display_name='Display Name', is_value=True)
    proto = StringEntityField('proto', display_name='Protocol')
    name = StringEntityField('name', display_name='Name')
    info = StringEntityField('info', display_name='Info')
    port = StringEntityField('port', display_name='Port')
    state = StringEntityField('state', display_name='State', decorator=getStateIcon)
    host_id = StringEntityField('host_id', display_name='Host ID')
    id = StringEntityField('id', display_name='Service ID')
    workspace_id = StringEntityField('workspace_id', display_name='Workspace ID')
    created_at = StringEntityField('created_at', display_name='Created At')
    updated_at = StringEntityField('updated_at', display_name='Updated At')

    # Decorator properties
    origin_tool = StringEntityField('origin_tool', display_name='Origin Tool', decorator=getOriginTool)
    tool_icon = StringEntityField('tool_icon', display_name='Tool Icon')
    state_icon = StringEntityField('state_icon', display_name='State Icon')


#  Apache Tomcat ----------------------------------------------------------------- {{{2 # 

class ApacheTomcat(WebService):
    _category_ = 'Metasploit | Services | Web | Application'
    _namespace_ = 'EffectiveCouscous.MetasploitService.WebService'

    # Entity properties
    display = StringEntityField('display', display_name='Display Name', is_value=True)
    proto = StringEntityField('proto', display_name='Protocol')
    name = StringEntityField('name', display_name='Name')
    info = StringEntityField('info', display_name='Info')
    port = StringEntityField('port', display_name='Port')
    state = StringEntityField('state', display_name='State', decorator=getStateIcon)
    host_id = StringEntityField('host_id', display_name='Host ID')
    id = StringEntityField('id', display_name='Service ID')
    workspace_id = StringEntityField('workspace_id', display_name='Workspace ID')
    created_at = StringEntityField('created_at', display_name='Created At')
    updated_at = StringEntityField('updated_at', display_name='Updated At')

    # Decorator properties
    origin_tool = StringEntityField('origin_tool', display_name='Origin Tool', decorator=getOriginTool)
    tool_icon = StringEntityField('tool_icon', display_name='Tool Icon')
    state_icon = StringEntityField('state_icon', display_name='State Icon')


#  Micosoft----------------------------------------------------------------------- {{{3 # 

class IISWebService(Entity):
    _category_ = 'Metasploit | Services | Web | Application'
    _namespace_ = 'EffectiveCouscous.MetasploitService.WebService'

    # Entity properties
    display = StringEntityField('display', display_name='Display Name', is_value=True)
    proto = StringEntityField('proto', display_name='Protocol')
    name = StringEntityField('name', display_name='Name')
    info = StringEntityField('info', display_name='Info')
    port = StringEntityField('port', display_name='Port')
    state = StringEntityField('state', display_name='State', decorator=getStateIcon)
    host_id = StringEntityField('host_id', display_name='Host ID')
    id = StringEntityField('id', display_name='Service ID')
    workspace_id = StringEntityField('workspace_id', display_name='Workspace ID')
    created_at = StringEntityField('created_at', display_name='Created At')
    updated_at = StringEntityField('updated_at', display_name='Updated At')

    # Decorator properties
    origin_tool = StringEntityField('origin_tool', display_name='Origin Tool', decorator=getOriginTool)
    tool_icon = StringEntityField('tool_icon', display_name='Tool Icon')
    state_icon = StringEntityField('state_icon', display_name='State Icon')

    
#  JavaScript -------------------------------------------------------------------- {{{4 # 

class NodeJS(WebService):
    _category_ = 'Metasploit | Services | Web | Application'
    _namespace_ = 'EffectiveCouscous.MetasploitService.WebService'

    # Entity properties
    display = StringEntityField('display', display_name='Display Name', is_value=True)
    proto = StringEntityField('proto', display_name='Protocol')
    name = StringEntityField('name', display_name='Name')
    info = StringEntityField('info', display_name='Info')
    port = StringEntityField('port', display_name='Port')
    state = StringEntityField('state', display_name='State', decorator=getStateIcon)
    host_id = StringEntityField('host_id', display_name='Host ID')
    id = StringEntityField('id', display_name='Service ID')
    workspace_id = StringEntityField('workspace_id', display_name='Workspace ID')
    created_at = StringEntityField('created_at', display_name='Created At')
    updated_at = StringEntityField('updated_at', display_name='Updated At')

    # Decorator properties
    origin_tool = StringEntityField('origin_tool', display_name='Origin Tool', decorator=getOriginTool)
    tool_icon = StringEntityField('tool_icon', display_name='Tool Icon')
    state_icon = StringEntityField('state_icon', display_name='State Icon')


#  Glassfish --------------------------------------------------------------------- {{{5 # 

class OracleGlassFish(WebService):
    _category_ = 'Metasploit | Services | Web | Application'
    _namespace_ = 'EffectiveCouscous.MetasploitService.WebService'

    # Entity properties
    display = StringEntityField('display', display_name='Display Name', is_value=True)
    proto = StringEntityField('proto', display_name='Protocol')
    name = StringEntityField('name', display_name='Name')
    info = StringEntityField('info', display_name='Info')
    port = StringEntityField('port', display_name='Port')
    state = StringEntityField('state', display_name='State', decorator=getStateIcon)
    host_id = StringEntityField('host_id', display_name='Host ID')
    id = StringEntityField('id', display_name='Service ID')
    workspace_id = StringEntityField('workspace_id', display_name='Workspace ID')
    created_at = StringEntityField('created_at', display_name='Created At')
    updated_at = StringEntityField('updated_at', display_name='Updated At')

    # Decorator properties
    origin_tool = StringEntityField('origin_tool', display_name='Origin Tool', decorator=getOriginTool)
    tool_icon = StringEntityField('tool_icon', display_name='Tool Icon')
    state_icon = StringEntityField('state_icon', display_name='State Icon')






