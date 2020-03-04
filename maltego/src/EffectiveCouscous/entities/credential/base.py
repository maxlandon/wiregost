#!/usr/bin/env python3.7.2

# -------------------- Imports --------------------- #

# Fields
from canari.maltego.message import *
# -------------------------------------------------- #


__author__ = 'Maxime Landon'
__copyright__ = 'Copyright 2019, EffectiveCouscous Project'
__credits__ = []

__license__ = 'GPLv3'
__version__ = '0.2'
__maintainer__ = 'Maxime Landon'
__email__ = 'maximelandon@gmail.com'
__status__ = 'Development'


#_____________________________________________________________________________________________#

#                                           CREDENTIALS 
#_____________________________________________________________________________________________#


class MetasploitCredential(Entity):
    # The MetasploitCredential class does not contain all the properties stored in the 
    # Metasploit Credential JSON object (which has deeply nested components). 
    # It is not necessary to retrieve them all, as only having the cred ID 
    # and its service ID is sufficient for retrieving it and processing it in other ways.

    # Static properties
    _category_ = 'Penetration Testing'
    _namespace_ = 'EffectiveCouscous.metasploit.credential'
    _alias_ = 'Metasploit Credential'

    
    # Entity properties
    name = Field('name', display_name='Name', is_value=True)
    id = Field('id', display_name='Credential ID' ) 
    logins_count = Field('logins_count', display_name='Logins Count')
    pub_username = Field('pub_username', display_name='Public Username')
    pub_type = Field('pub_type', display_name='Public Type')
    priv_data = Field('priv_data', display_name='Private Data')
    priv_type = Field('priv_type', display_name='Private Type')
    priv_jtr_format = Field('priv_jtr_format', display_name='Private JTR Format')
    origin_service_id = Field('origin_service_id', display_name='Origin Service ID')
    origin_type = Field('origin_type', display_name='Origin Type')
    origin_module = Field('origin_module', display_name='Origin Module')
