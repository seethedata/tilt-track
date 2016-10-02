package main

import (
	"encoding/json"
	"fmt"
	"github.com/elgs/gojq"
	"github.com/garyburd/redigo/redis"
	"log"
	"net/http"
	"os"
	"strconv"
)

var redisClient redis.Conn

func check(function string, e error) {
	if e != nil {
		log.Fatal(function, e)
	}
}
func responseHandler(w http.ResponseWriter, r *http.Request) {
	response, err := redisClient.Do("KEYS", "{lat:*")
	m := response.(map[string]interface{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%T|%v\n", m, m)
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
	credentials, err := parser.Query(serviceName + ".[0].credentials")

	// Convert credentials to Go map so we can reference the values
	m := credentials.(map[string]interface{})
	hostKey := "host"
	redisHost := m[hostKey].(string)
	redisPassword := m["password"]
	redisPort := strconv.FormatFloat(m["port"].(float64), 'f', -1, 64)

	redisClient, err = redis.Dial("tcp", redisHost+":"+redisPort)

	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	_, err = redisClient.Do("AUTH", redisPassword.(string))
	if err != nil {
		log.Fatal(err)
	}
	pong, err := redisClient.Do("PING")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pong, err)

	http.HandleFunc("/data", responseHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
