## using Docker Image 

### Running Redis 

```
docker network create redis
docker run -it --rm --name redis --net redis -p 6379:6379 redis:7.2.1-alpine
```