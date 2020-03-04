Effective Couscous - Metasploit & Maltego for Python 3.X
================================================================


![home-pic](./docpics/home-pic.png) 


  Overview
-----------------------------------

EffectiveCouscous is a set of Maltego Transforms that interact with various security-related tools, 
among which is the Metasploit Framework. Metasploit is used in this description to illustrate the intents
and aims of this project. Later on, other tools will be integrated in the same manner.
The transforms, at the higher level, either:

*   Act upon data stored in the Metasploit Database, mainly through its REST API.
*   Execute commands of the Metasploit Framework, mainly through its RPC server.


If well combined, Metasploit and Maltego can greatly benefit each other, and therefore the 
penetration tester using them. Maltego is here to be considered as a sort of "meta-framework" 
which would manage a complex, versatile set of data and expose it to various tools, such as the 
Metasploit Framework. 
This is conditioned by sound, balanced and planned integration of Metasploit toolsets in Maltego.
The major benefits would be:

*   **Exhaustive, flexible and versatile representation/visualisation of Metasploit Data 
    into Maltego's entity graphs**. Computer networks of various kinds, comprising various objects,
    can be viewed with accurate topology representation, multi-layer information, etc...
    Each host, netblock exposes various kinds of data, which can be further used as input to
    other transforms. 

*   **Focused, context-sensitive availability of Metasploit toolset into Maltego**. Netblocks
    can become workspaces, as well as hosts, so that they can selfishly profit from Metasploit
    workspace structure. Host, services, sessions, or consoles are entities upon which one
    can act, exactly like in the Metasploit console. Sessions, for instance, should offer the
    right subset of their tools in the Maltego graph. 
    
*   **Various layers of information can be integrated into one Maltego graph, or can be separated
    into multiple Maltego graphs**. For example, mapping Metasploit routes can be done in such a
    layer. It could be then used to further visualize Metasploit related traffic, and potential
    "points of application" can be identified for further scans, exploits, etc...
    Separate graphs can be used for a particular host, subnet, document set, website struture.

*   In addition to Maltego free transform set, various transform librairies can be paid for and
    used. A correct integration of Metasploit entities into Maltego entity structure can expose
    Metasploit-stored data to transforms unrelated to Metasploit. This might apply to Loot,
    Notes, files, passwords, mail addresses, etc...
    Another example is to imagine other tools integrated into Maltego, such as Nmap operations,
    on a host entity, a netblock entity, a service in a netblock, etc...

*   **Taking advantage of Maltego's inference capabilities**, so to find corresponding entities/attributes
    in complex graphs. This is particularly important because it may analyze the full dataset contained in
    the graph while never setting any connection of any sort in the target network.
    One example can be found with passwords: There might be cases where Maltego can infer the services and
    domains that are accessible for a given credential set, without ever actually logging into them.
    This serves the main purpose of this project: **offering visualization and discovery capabilities to
    a black-box security audit, which requires a level of discretion and attention to detail that is 
    significantly higher than for defensive security audits.**
    Maltego's inference capacibilities are also highly modular (thanks to MatchingRules).

<!-- *A demonstration of Maltego transforms interacting with Metasploit Database. (Much faster than what -->
<!-- the Gif demo implies.)* -->
<!-- ![full_demo](./docpics/full_demo.gif)  -->

______

  Documentation
-----------------------------------

The documentation for the EffectiveCouscous project is articulated around the following pages:

1. [Installation](https://github.com/maxlandon/EffectiveCouscous/wiki/Installation)
2. [Usage](https://github.com/maxlandon/EffectiveCouscous/wiki/Usage)
3. [Data Model Presentation](https://github.com/maxlandon/EffectiveCouscous/wiki/Data-Model-Presentation)
4. [Entity Code Structure (Developers)](https://github.com/maxlandon/EffectiveCouscous/wiki/Entity-Code-Structure-(Developers))
5. [Transform Code Structure (Developers)](https://github.com/maxlandon/EffectiveCouscous/wiki/Transform-Code-Structure-(Developers))
6. [Previous Approaches & Implementations](https://github.com/maxlandon/EffectiveCouscous/wiki/Previous-Approaches-&-Implementations)
7. [Transforms List](https://github.com/maxlandon/EffectiveCouscous/wiki/Transforms-List)
8. [Machines List](https://github.com/maxlandon/EffectiveCouscous/wiki/Machines-List)
9. [Issues, ToDo & DoNot](https://github.com/maxlandon/EffectiveCouscous/wiki/Issues,-ToDo-&-DoNot)

______

## Warmest Thanks
To @allfro, for the Canari Framework, an absolutely amazing tool.
