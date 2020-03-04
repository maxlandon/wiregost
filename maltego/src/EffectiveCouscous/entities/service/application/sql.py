#!/usr/bin/env python3.7.2

# -------------------- Imports --------------------- #

# Maltego Entities
from canari.maltego.entities import Entity

# Custom Entities
from EffectiveCouscous.entities.service.base import MetasploitService

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



# Base SQL Service ------------------------------------------------------------------------------#


class SQLService(MetasploitService):
    _category_ = 'Metasploit | Services | SQL'
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



# MySQL -----------------------------------------------------------------------------------------#


class MySQL(SQLService):
    _category_ = 'Metasploit | Services | SQL'
    _namespace_ = 'EffectiveCouscous.MetasploitService.SQLService'

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


# SQL Server ------------------------------------------------------------------------------------#


class SQLServer(SQLService):
    _category_ = 'Metasploit | Services | SQL'
    _namespace_ = 'EffectiveCouscous.MetasploitService.SQLService'

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


# PostgreSQL ------------------------------------------------------------------------------------#


class PostgreSQL(SQLService):
    _category_ = 'Metasploit | Services | SQL'
    _namespace_ = 'EffectiveCouscous.MetasploitService.SQLService'

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


# MariaDB ---------------------------------------------------------------------------------------#


class MariaDB(SQLService):
    _category_ = 'Metasploit | Services | SQL'
    _namespace_ = 'EffectiveCouscous.MetasploitService.SQLService'

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


