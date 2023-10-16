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