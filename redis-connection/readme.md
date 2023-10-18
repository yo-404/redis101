# Storing data in Redis Database

starting our container
```
cd redis101/redis-connection

docker build --target dev . -t go
docker run -it -p 80:80 -v ${PWD}:/work go sh

```

### Creating application

```
mkdir videos
cd videos
go mod init videos
```

### main.go

```
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var redisClient *redis.Client

func main() {

	var redis_sentinels = os.Getenv("REDIS_SENTINELS")
	var redis_master = os.Getenv("REDIS_MASTER_NAME")
	var redis_password = os.Getenv("REDIS_PASSWORD")

	sentinelAddrs := strings.Split(redis_sentinels, ",")

	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    redis_master,
		SentinelAddrs: sentinelAddrs,
		Password:      redis_password,
		DB:            0,
	})

	redisClient = rdb

	rdb.Ping(ctx)

	http.HandleFunc("/", HandleGetVideos)
	http.HandleFunc("/update", HandleUpdateVideos)

	http.ListenAndServe(":80", nil)
}

func HandleGetVideos(w http.ResponseWriter, r *http.Request) {

	videos := getVideos()
	videoBytes, err := json.Marshal(videos)

	if err != nil {
		panic(err)
	}

	w.Write(videoBytes)
}

func HandleUpdateVideos(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		var videos []video
		err = json.Unmarshal(body, &videos)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Bad request")
		}

		saveVideos(videos)

	} else {
		w.WriteHeader(405)
		fmt.Fprintf(w, "Method not Supported!")
	}
}

```

### videos.go

```
package main

import (
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

type video struct {
	Id          string
	Title       string
	Description string
	Imageurl    string
	Url         string
}

func getVideos() (videos []video) {

	keys, err := redisClient.Keys(ctx, "*").Result()

	if err != nil {
		panic(err)
	}

	for _, key := range keys {
		video := getVideo(key)
		videos = append(videos, video)
	}

	return videos
}

func saveVideo(video video) {

	videoBytes, err := json.Marshal(video)
	if err != nil {
		panic(err)
	}

	err = redisClient.Set(ctx, video.Id, videoBytes, 0).Err()
	if err != nil {
		panic(err)
	}
}

func saveVideos(videos []video) {
	for _, video := range videos {
		saveVideo(video)
	}
}

func getVideo(id string) (video video) {

	value, err := redisClient.Get(ctx, id).Result()

	if err != nil {
		panic(err)
	}

	if err != redis.Nil {
		err = json.Unmarshal([]byte(value), &video)
	}

	return video
}

```

### Redis Go package

inside docker go sh
```
go get github.com/go-redis/redis/v9
```

Dont forget to import it as well

```
import (
  "github.com/go-redis/redis/v9"
)
```

### Running application again with sentinels

```
docker run -it -p 80:80 `
  --net redis `
  -e REDIS_SENTINELS="sentinel-0:5000,sentinel-1:5000,sentinel-2:5000" `
  -e REDIS_MASTER_NAME="mymaster" `
  -e REDIS_PASSWORD="a-very-complex-password-here" `
  -v ${PWD}:/work go sh
```

### Building docker file

```
FROM golang:1.15-alpine as dev

WORKDIR /work

FROM golang:1.15-alpine as build

WORKDIR /videos
COPY ./videos/* /videos/
RUN go build -o videos

FROM alpine as runtime 
COPY --from=build /videos/videos /
CMD ./videos

```

Build:

```
cd redis101/redis-connection
docker build . -t videos
```

RUN:

```
docker run -it -p 80:80 `
  --net redis `
  -e REDIS_SENTINELS="sentinel-0:5000,sentinel-1:5000,sentinel-2:5000" `
  -e REDIS_MASTER_NAME="mymaster" `
  -e REDIS_PASSWORD="a-very-complex-password-here" `
  videos
```

### Cleanup

```
# clean up 

docker rm -f redis-0 redis-1 redis-2
docker rm -f sentinel-0 sentinel-1 sentinel-2
```
