#!/usr/bin/env python3.7.2

# PyMetasploit Modules
from metasploit.msfrpc import MsfRpcClient, MsfRpcError

# Sploitego Modules
#  import socket
#  from urllib.parse import parse_qsl
#  from urllib.parse import urlencode
#  from os import path, unlink
#  
#  from canari.easygui import multpasswordbox
#  from canari.utils.fs import CookieFile, FileSemaphore
#  import canari.config as config

__author__ = 'Maxime Landon'
__copyright__ = ''
__credits__ = []

__license__ = 'GPL'
__version__ = '0.1'
__maintainer__ = 'Maxime Landon'
__email__ = 'maximelandon@gmail.com'
__status__ = 'Development'

# Module locator, in case this file would only contain functions,
# Like in Sploitego project
#  __all__ = [
#      'login'
#  ]


class MetasploitRPC(object):

    def __init__(self):
        #  self.user = user
        #  self._password = password
        #  self.db = db
        self._cur = self._login()

        
    def _login(self):
        conn = MsfRpcClient('testpassword')
        return conn
    
    # Sploitego approach to Login, should be used for further
    # readability on DB/RPC/API connections from within Maltego
    #  def _login(self, **kwargs):
    #      s = None
    #      host = kwargs.get('host', config['msfrpcd/server'])
    #      port = kwargs.get('port', config['msfrpcd/port'])
    #      uri = kwargs.get('uri', config['msfrpcd/uri'])
    #      fn = cookie('%s.%s.%s.msfrpcd' % (host, port, uri.replace('/', '.')))
    #      if not path.exists(fn):
    #          f = fsemaphore(fn, 'wb')
    #          f.lockex()
    #          fv = [ host, port, uri, 'msf' ]
    #          errmsg = ''
    #          while True:
    #              fv = multpasswordbox(errmsg, 'Metasploit Login', ['Server:', 'Port:', 'URI', 'Username:', 'Password:'], fv)
    #              if not fv:
    #                  return
    #              try:
    #                  s = MsfRpcClient(fv[4], server=fv[0], port=fv[1], uri=fv[2], username=fv[3])
    #              except MsfRpcError as e:
    #                  errmsg = str(e)
    #                  continue
    #              except socket.error as e:
    #                  errmsg = str(e)
    #                  continue
    #              break
    #          f.write(urlencode({'host' : fv[0], 'port' : fv[1], 'uri': fv[2], 'token': s.sessionid}))
    #          f.unlock()
    #  
    #          if 'db' not in s.db.status:
    #              s.db.connect(
    #                  config['metasploit/dbusername'],
    #                  database=config['metasploit/dbname'],
    #                  driver=config['metasploit/dbdriver'],
    #                  host=config['metasploit/dbhost'],
    #                  port=config['metasploit/dbport'],
    #                  password=config['metasploit/dbpassword']
    #              )
    #      else:
    #          f = fsemaphore(fn)
    #          f.locksh()
    #          try:
    #              d = dict(parse_qsl(f.read()))
    #              s = MsfRpcClient('', **d)
    #          except MsfRpcError:
    #              unlink(fn)
    #              return login()
    #          except socket.error:
    #              unlink(fn)
    #              return login()
    #      return s

