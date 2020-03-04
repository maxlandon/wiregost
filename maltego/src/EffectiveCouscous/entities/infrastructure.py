#!/usr/bin/env python3.7.2

# -------------------- Imports --------------------- #

# Maltego Entities
from canari.maltego.entities import Entity

# Fields
from canari.maltego.message import (StringEntityField, 
                                    IntegerEntityField, 
                                    BooleanEntityField,
                                    LinkColor)

# System-wide Icons
from EffectiveCouscous.tools.entitytools.icon_factory import getOriginTool
# -------------------------------------------------- #


__author__ = 'Maxime Landon'
__copyright__ = 'Copyright 2019, EffectiveCouscous Project'
__credits__ = []

__license__ = 'GPLv3'
__version__ = '0.2'
__maintainer__ = 'Maxime Landon'
__email__ = 'maximelandon@gmail.com'
__status__ = 'Development'



# This file redefines Netblock and IPv4Addresses, by overloading them with new attributes.



class Netblock(Entity):
    _category_ = 'Infrastructure'
    _namespace_ = 'maltego'

    # Main properties ---------------------------//
    name = StringEntityField('name', display_name='Workspace Name',
                                    description='The name of the workspace. This is the unique identifier \
                                                for determining which workspace is being accessed.')
    workspace_id = IntegerEntityField('workspace_id', display_name='ID', 
                                    description='The primary key used to identify this object in the database.')
    boundary = StringEntityField('boundary', display_name='Boundary', 
                                    description='Comma separated list of IP ranges (in various formats) \
                                                and IP addresses that users of this workspace are allowed to interact \
                                                with if limit_to_network is true.')
    description = StringEntityField('description', display_name='Description', 
                                    description='Long description that explains the purpose of this workspace.')
    owner_id = StringEntityField('owner_id', display_name='Owner ID', 
                                    description='ID of the user who owns this workspace.')
    limit_to_network = BooleanEntityField('limit_to_network', display_name='Limit to Network', 
                                    description='true to restrict the hosts and services in this workspace \
                                                to the IP addresses listed in boundary')
    import_fingerprint = BooleanEntityField('import_fingerprint', display_name='Import fingerprint', 
                                    description='Identifier that indicates if and where this workspace was imported from.')
    created_at = StringEntityField('created_at', display_name='Created at')
    updated_at = StringEntityField('updated_at', display_name='Updated at')

    # Icon Properties ---------------------------//
    origin_tool = StringEntityField('origin_tool', display_name='Origin Tool', decorator=getOriginTool)
    tool_icon = StringEntityField('tool_icon', display_name='Tool Icon')



class IPv4Address(Entity):
    _category_ = 'Infrastructure'
    _namespace_ = 'maltego'

    # Main properties ---------------------------//
    host_id = IntegerEntityField('id', display_name='Host ID')
    workspace_id = IntegerEntityField('workspace_id', display_name='Workspace ID')

    # Icon Properties ---------------------------//
    origin_tool = StringEntityField('origin_tool', display_name='Origin Tool', decorator=getOriginTool)
    tool_icon = StringEntityField('tool_icon', display_name='Tool Icon')
