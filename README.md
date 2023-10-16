## using Docker Image 

### Running Redis 

```
docker network create redis
docker run -it --rm --name redis --net redis -p 6379:6379 redis:7.0-alpine
```

## Configuration

To start redis with a custom config

```
cd .\storage\redis\
docker run -it --rm --name redis --net redis -v ${PWD}/config:/etc/redis/ redis:7.0-alpine redis-server /etc/redis/redis.conf

```

## How to secure Redis

Since redis instance does not have a password by default and hence it should not be exposed to public traffic .

For setting password , use a strong password in `redis.conf`

```
requirepass VeryStrongpassword
```

## Redis Persistence

Documentation here [Redis persistence](https://redis.io/docs/management/persistence/)

RDB Snapshotting:
- Redis can periodically save a snapshot of the dataset to disk. This is done by forking a child process, which writes the data to a binary dump file.
- The RDB snapshots are point-in-time backups and are useful for recovery after a system restart.
- You can configure the frequency of snapshots using the save configuration in the redis.conf file or using the SAVE or BGSAVE commands.

AOF (Append-Only File):
- AOF persistence mode logs every write operation to a log file as a series of append-only commands.
- When Redis restarts, it can reconstruct the dataset by replaying these commands.
- AOF provides more granular recovery and is considered more robust than RDB in some cases.
- You can configure AOF options in the redis.conf file.

Hybrid Persistence:
- Redis allows you to use both RDB snapshots and AOF logs together, providing both point-in-time backups and a log of write operations.
- In this setup, Redis can be configured to use RDB snapshots as a background process and AOF for detailed changes.

Snapshots and AOF Rewrite:
- AOF log files can become large over time, so Redis supports a process called AOF rewrite, which can reduce the AOF file size by creating a new log file containing only the necessary operations.
- This process helps manage the size of AOF files while still retaining the benefits of AOF persistence.

You can configure these persistence options in the redis.conf file or through runtime configurations using the CONFIG command. The choice of persistence method depends on your specific use case and requirements, including data durability and recovery speed.

To turn on RDB mode -

find `dbfilename dump.rdb` 
change the filename /location to your desired location/filename . It will start writing at various intervals on that file. The intervals of writing can also be changed

To turn on append Mode

Search for `appendonly no` in the redis.conf . Change it to `appendonly yes`

filename can also be specified in the same section under the redis.conf file

## Docker Volume

```
docker volume create redis
cd .\storage\redis\
docker run -it --rm --name redis --net redis -v ${PWD}/config:/etc/redis/ -v redis:/data/  redis:7.0-alpine redis-server /etc/redis/redis.conf

```

## TO run the client application

```
cd .\storage\redis\applications\client\
# start go dev environment
docker run -it -v ${PWD}:/go/src -w /go/src --net redis -p 80:80 golang:1.21-alpine

# to add dependencies from git
 apk add --no-cache git

#to download dependencies
go mod tidy

go build client.go
# start the app
./client

# build the container
docker build . -t aimvector/redis-client:v1.0.0
```

## To run our application

```
cd .\storage\redis\applications\client\
docker build . -t aimvector/redis-client:v1.0.0

docker run -it --net redis `
-e REDIS_HOST=redis `
-e REDIS_PORT=6379 `
-e REDIS_PASSWORD="VeryStrongpassword" `
-p 80:80 `
aimvector/redis-client:v1.0.0

```

