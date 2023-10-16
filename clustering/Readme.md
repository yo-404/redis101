## Redis replication

[Redis Replication blog](https://redis.io/docs/management/replication/)

How Redis supports high availability and failover with replication

At the base of Redis replication (excluding the high availability features provided as an additional layer by Redis Cluster or Redis Sentinel) there is a leader follower (master-replica) replication that is simple to use and configure. It allows replica Redis instances to be exact copies of master instances. The replica will automatically reconnect to the master every time the link breaks, and will attempt to be an exact copy of it regardless of what happens to the master.

## configuration

change the following configurations in the redis.conf file

```
#security
masterauth "verystrongpassword"
requirepass verystrongpassword

#persistence
dir "/data"
dbfilename dump.rdb
appendonly yes
appendfilename "appendonly.aof"

```

### Redis-0 Configuration

```
protected-mode no
port 6379

#authentication
masterauth verystrongpassword
requirepass verystrongpassword
```

### Redis-1 Configuration

```
protected-mode no
port 6379
slaveof redis-0 6379

#authentication
masterauth verystrongpassword
requirepass verystrongpassword

```

### Redis-2 Configuration

```
protected-mode no
port 6379
slaveof redis-0 6379

#authentication
masterauth verystrongpassword
requirepass verystrongpassword

```

We have to run all the containers on the same network , hence first we will create a network 

```
docker network create redis
```

### To run

```

cd .\storage\redis\clustering\

#redis-0
docker run -d --rm --name redis-0 `
    --net redis `
    -v ${PWD}/redis-0:/etc/redis/ `
    redis:6.0-alpine redis-server /etc/redis/redis.conf

#redis-1
docker run -d --rm --name redis-1 `
    --net redis `
    -v ${PWD}/redis-1:/etc/redis/ `
    redis:6.0-alpine redis-server /etc/redis/redis.conf


#redis-2
docker run -d --rm --name redis-2 `
    --net redis `
    -v ${PWD}/redis-2:/etc/redis/ `
    redis:6.0-alpine redis-server /etc/redis/redis.conf

```

## To run our counter application

```
cd .\storage\redis\applications\client\
docker build . -t aimvector/redis-client:v1.0.0

docker run -it --net redis `
-e REDIS_HOST=redis-0 `
-e REDIS_PORT=6379 `
-e REDIS_PASSWORD="verystrongpassword" `
-p 80:80 `
aimvector/redis-client:v1.0.0

```

## To test replication

technically writted data should be now on the replicas

```
# go to one of the clients
docker exec -it redis-2 sh
redis-cli
auth "a-very-complex-password-here"
keys *

```

## Sentinel in redis

[sentinel blog](https://redis.io/docs/management/sentinel/)

Redis Sentinel provides high availability for Redis when not using Redis Cluster.
This is the full list of Sentinel capabilities at a macroscopic level (i.e. the big picture):

- Monitoring. Sentinel constantly checks if your master and replica instances are working as expected.
- Notification. Sentinel can notify the system administrator, or other computer programs, via an API, that something is wrong with one of the monitored Redis instances.
- Automatic failover. If a master is not working as expected, Sentinel can start a failover process where a replica is promoted to master, the other additional replicas are reconfigured to use the new master, and the applications using the Redis server are informed about the new address to use when connecting.
- Configuration provider. Sentinel acts as a source of authority for clients service discovery: clients connect to Sentinels in order to ask for the address of the current Redis master responsible for a given service. If a failover occurs, Sentinels will report the new address.


## Starting Redis in sentinel mode

```
cd .\storage\redis\clustering\

docker run -d --rm --name sentinel-0 --net redis `
    -v ${PWD}/sentinel-0:/etc/redis/ `
    redis:6.0-alpine `
    redis-sentinel /etc/redis/sentinel.conf

docker run -d --rm --name sentinel-1 --net redis `
    -v ${PWD}/sentinel-1:/etc/redis/ `
    redis:6.0-alpine `
    redis-sentinel /etc/redis/sentinel.conf

docker run -d --rm --name sentinel-2 --net redis `
    -v ${PWD}/sentinel-2:/etc/redis/ `
    redis:6.0-alpine `
    redis-sentinel /etc/redis/sentinel.conf


docker logs sentinel-0
docker exec -it sentinel-0 sh
redis-cli -p 5000
info
sentinel master mymaster

# clean up 

docker rm -f redis-0 redis-1 redis-2
docker rm -f sentinel-0 sentinel-1 sentinel-2


```

## Note -

replicas are read only and masters are write only

## Cleanup 

```
# clean up 

docker rm -f redis-0 redis-1 redis-2
docker rm -f sentinel-0 sentinel-1 sentinel-2
```