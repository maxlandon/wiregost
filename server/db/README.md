## DB 

The `db` package is a key:value database implemented using badger-db.
This package differs from the `data_service` package, in that this one is used merely for
storing binary objects such as Certificate Key Pairs, which are in binary form and for which
we want some storage safety without too much hassle.
