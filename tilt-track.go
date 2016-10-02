package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo"
	"gopkg.in/redis.v4"
	"log"
	"net/http"
	"os"
	"strconv"
)

var redisClient *redis.Client

func check(function string, e error) {
	if e != nil {
		log.Fatal(function, e)
	}
}
func responseHandler(w http.ResponseWriter, r *http.Request) {
	response := redisClient.Keys("{lat:*")
	fmt.Printf("%T|%v\n", response, response)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func main() {

	var port string
	if os.Getenv("VCAP_APP_PORT") != "" {
		port = os.Getenv("VCAP_APP_PORT")
	} else if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	} else {
		port = "8080"
	}

	servicesJSON := os.Getenv("VCAP_SERVICES")
	fmt.Println(servicesJSON)
	parser, err := gojq.NewStringQuery(servicesJSON)
	if err != nil {
		log.Fatal(err)
	}

	serviceName := "p-redis"
	redisHost, err := parser.Query(serviceName + ".[0].credentials.host")
	redisPassword, err := parser.Query(serviceName + ".[0].credentials.password")
	redisPort, err := parser.Query(serviceName + ".[0].credentials.port")

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost.(string) + ":" + strconv.FormatFloat(redisPort.(float64), 'f', -1, 64),
		Password: redisPassword.(string),
		DB:       0,
	})

	pong, err := redisClient.Ping().Result()
	fmt.Println(pong, err)
	redisClient.Del("ghiorzi1")
	redisClient.Del("ghiorzi2")
	err = redisClient.Set("{lat: 49, lng: -75}", "20161002150600", 0).Err()
	if err != nil {
		log.Fatal(err)
	}
	err = redisClient.Set("{lat: 80, lng: -95}", "20161001150600", 0).Err()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/data", responseHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
