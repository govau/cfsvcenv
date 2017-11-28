package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	hackernewsEndpointEnv = "HACKERNEWS_ENDPOINT"
	portEnv               = "PORT"
)

func main() {
	endpointStr := os.Getenv(hackernewsEndpointEnv)
	port := os.Getenv(portEnv)

	if endpointStr == "" {
		fmt.Printf("%s must be provided\n", hackernewsEndpointEnv)
		os.Exit(1)
	}

	endpoint, err := url.Parse(endpointStr)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		u, err := endpoint.Parse("items/192327.json")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("GET %s\n", u)

		res, err := http.Get(u.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(body)
	})

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
