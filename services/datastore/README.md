# Datastore

![datastore test status badge](https://github.com/IamCathal/neo/actions/workflows/buildDatastore.yml/badge.svg)   ![datastore deploy status](https://github.com/IamCathal/neo/actions/workflows/deployDataStore.yml/badge.svg) 

Datastore is the thin client that sits infront of the database

## Configuration

Datastore expects the following variables to be set in .env

#### Datastore specific

| Variable     | Description |
| ----------- | ----------- |
| `MONGODB_USER`      |  MongoDB account username  |
| `MONGODB_PASSWORD`      |  MongoDB account password  |
| `MONGO_INSTANCE_IP`      |  IP for the MongoDB instance |
| `DB_NAME`      |  Database name for the stored data |
| `USER_COLLECTION`      |  Collection name for the user data |
| `CRAWLING_STATS_COLLECTION`      |  Collection name for the crawling stats |
| `SHORTEST_DISTANCE_COLLECTION`    | Collection name for shortest distance info |
| `POSTGRES_USER`      |  Username for postgres worker account |
| `POSTGRES_PASSWORD`      |  Password for postgres worker account |
| `POSTGRES_DB`      |  DB name for postgres saved graphs table |
| `POSTGRES_INSTANCE_IP`      |  IP for the postgres instance |



## Running 

`docker-compose up` to start with docker-compose (preferred)

`docker build -f Dockerfile -t iamcathal/datastore:0.0.1 .` and `docker run -it --rm -p PORT:PORT iamcathal/datastore:0.0.1` to start as a standalone container