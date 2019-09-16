package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

const (
	DefaultHost = "http://localhost"
	DefaultPort = "8080"
)

func main() {
	// Retrieve proxy host
	proxyHost := os.Getenv("HOST")
	if proxyHost == "" {
		proxyHost = DefaultHost
	}

	// Retrieve proxy port
	proxyPort := os.Getenv("PORT")
	if proxyPort == "" {
		proxyPort = DefaultPort
	}

	// Proxy URL
	proxyUrl := proxyHost + ":" + proxyPort

	// The favicon.ico
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadFile("favicon.ico")
		w.Header().Set("Content-Type", "image/x-icon")
		_, _ = fmt.Fprintln(w, string(bytes))
	})

	// Proxy main route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Parse the query string
		queries, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			displayError(w, "Bad Request.")
			return
		}

		// Get the target url
		target, err := url.Parse(queries["url"][0])
		if err != nil || target.IsAbs() == false {
			displayError(w, "Bad Request.")
			return
		}

		// Make reverse proxy
		proxy := &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				// Pretend that the referrer is its homepage!
				req.Header.Set("Referer", target.Scheme+"://"+target.Host+"/")
				req.Header.Set("Origin", target.Host)
				req.Host = target.Host
				req.URL = target

				log.Println(req)
			},

			ModifyResponse: func(response *http.Response) error {
				response.Header.Set("Access-Control-Allow-Origin", "*")

				// Handle redirection responses
				if response.Header.Get("Location") != "" {
					response.Header.Set("Location", proxyUrl+"/?"+response.Header.Get("Location"))
				}

				return nil
			},
		}

		// Handle request with reverse proxy
		proxy.ServeHTTP(w, r)
	})

	// Log and listen to the web port
	log.Println("Started: " + proxyUrl)
	log.Fatal(http.ListenAndServe(":"+proxyPort, nil))
}

// Create an error response
func displayError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_, err := w.Write([]byte(message))
	if err != nil {
		panic("Cannot respond to the request.")
	}
}
