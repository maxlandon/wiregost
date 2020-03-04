#!/usr/bin/env python3.7.2

# API
import requests
import json


__author__ = 'Maxime Landon'
__copyright__ = 'Copyright 2019, EffectiveCouscous Project'
__credits__ = []

__license__ = 'GPL'
__version__ = '0.2'
__maintainer__ = 'Maxime Landon'
__email__ = 'maximelandon@gmail.com'
__status__ = 'Development'



# GET METHODS
#------------------------------------------------
def get_raw_json(url, config, params=None, headers=None):
    headers = {
        'accept': 'application/json',
        'Authorization': 'Bearer {0}'.format(config['EffectiveCouscous.local.userapitoken']),
    }

    if headers is None:
        headers = self.headers

    response = requests.get(url, headers, params=None, verify=False)
    return response

def get_json_dict(url, config, params=None, headers=None):
    headers = {
        'accept': 'application/json',
        'Authorization': 'Bearer {0}'.format(config['EffectiveCouscous.local.userapitoken']),
    }

    if headers is None:
        headers = headers

    response = requests.get(url, headers=headers, params=params, verify=False)

    json_dict = response.json()['data']
    return json_dict 
    

# POST METHODS
#------------------------------------------------
def post_json(url, data, config, headers=None):
    headers = {
            'accept': 'application/json',
            'Authorization': 'Bearer {0}'.format(config['EffectiveCouscous.local.userapitoken']),
            'Content-Type': 'application/json',
        }

    response = requests.post(url, headers=headers, data=data, verify=False)
    return response


# PUT METHODS
#------------------------------------------------
def put_json(url, data, config, headers=None):
    headers = {
            'accept': 'application/json',
            'Authorization': 'Bearer {0}'.format(config['EffectiveCouscous.local.userapitoken']),
            'Content-Type': 'application/json',
        }

    response = requests.put(url, headers=headers, data=data, verify=False)
    return response


# DELETE METHODS
#------------------------------------------------
def delete_json(url, data, config, headers=None):
    headers = {
            'accept': 'application/json',
            'Authorization': 'Bearer {0}'.format(config['EffectiveCouscous.local.userapitoken']),
            'Content-Type': 'application/json',
        }

    response = requests.delete(url, headers=headers, data=data, verify=False)
    return response




