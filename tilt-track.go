package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/redis.v4"
	"log"
	"net/http"
	"os"
	"strconv"
)

type pcfDevData struct {
	Serv []service `json:"p-redis"`
}

type pwsData struct {
	Serv []service `json:"rediscloud"`
}

type bluemixData struct {
	Serv []service `json:"redis-2.6"`
}

type service struct {
	Credentials creds `json:"credentials"`
}

type creds struct {
	Host     string `json:"host"`
	Password string `json:"password"`
	Port     int    `json:"port"`
}

func check(function string, e error) {
	if e != nil {
		log.Fatal(function, e)
	}
}
func responseHandler(w http.ResponseWriter, r *http.Request) {
	response := `{"Key1": "val1"}`
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
	var services pcfDevData
	err := json.Unmarshal([]byte(servicesJSON), &services)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%T: %v\n", services.Serv[0].Credentials.Host, services.Serv[0].Credentials.Host)
	credentials := services.Serv[0].Credentials

	client := redis.NewClient(&redis.Options{
		Addr:     credentials.Host + ":" + strconv.Itoa(credentials.Port),
		Password: credentials.Password,
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	http.HandleFunc("/data", responseHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
