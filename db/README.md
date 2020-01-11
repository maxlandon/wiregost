### Database

The `database/` directory contains all the code necessary
to Wiregost for interacting with the PostgreSQL database.
It all also contains all entities used in Wiregost that will 
need some form of persistence.
This package should *not* import any other package from Wiregost,
maybe except from some configuration.

* `db/`       - Functions for querying and writing to DB
* `models/`   - All entities that need persistence in Wiregost
* `.config`   - Configuration file for database access
