### Scanner

The `scanner/` directory contains all the code necessary to scanning functionality in Wiregost.
Ideally, it should not import anything other than the `database/` code, which will be used to populate entities in Wiregost, or use them for scans.
On the other hand, the `scanner/` code can be imported by `server` handlers or `ghost` agents, if 
agent scanning is implemented in conjunction with port forwarding.
