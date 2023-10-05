# Setup

0. Setup ImmuDB https://docs.immudb.io/0.9.2/quickstart.html

1. Delete "data" directory to reset the database

2. Start two ImmuDB instances (with different data directories, and running on different ports). Parameters: https://docs.immudb.io/1.2.1/reference/configuration.html
```
./immudb
./immudb --dir=./data2 --port=3323 --web-server-port=8081 --pgsql-server-port=5433
```
For some reason, the web server still points to the first instance...

3. Delete ".identity-XXXXX" and ".state-XXXXX" files (generated by the clients) if reseted the databases

4. Start two ImmuDB clients pointing to the respectives databases
````
./balances-manager-sample
./balances-manager-sample -port 3323
````
