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
