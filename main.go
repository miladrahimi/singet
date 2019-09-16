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

	proxyUrl := proxyHost + ":" + proxyPort

	// The favicon.ico
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadFile("favicon.ico")
		w.Header().Set("Content-Type", "image/x-icon")
		_, _ = fmt.Fprintln(w, string(bytes))
	})

	// Proxy main route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		target, err := url.Parse(query.Get("url"))
		if err != nil || target.IsAbs() == false {
			displayError(w, "The requested URL is invalid.")
			return
		}

		// Make reverse proxy
		proxy := &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				if referrer := query.Get("referrer"); referrer != "" {
					req.Header.Set("Referer", referrer)
				} else {
					req.Header.Set("Referer", target.Scheme+"://"+target.Host+"/")
				}

				req.Header.Set("Origin", target.Host)
				req.Host = target.Host
				req.URL = target

				log.Println(req)
			},

			ModifyResponse: func(response *http.Response) error {
				response.Header.Set("Access-Control-Allow-Origin", "*")

				// Handle redirection responses
				if response.Header.Get("Location") != "" {
					if query.Get("follow") != "false" {
						response.Header.Set("Location", proxyUrl+"/?url="+response.Header.Get("Location"))
					}
				}

				return nil
			},
		}

		proxy.ServeHTTP(w, r)
	})

	log.Println("Started: " + proxyUrl)
	log.Fatal(http.ListenAndServe(":"+proxyPort, nil))
}

// Create an error response
func displayError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_, err := w.Write([]byte("Error: " + message))
	if err != nil {
		panic("Cannot respond to the request.")
	}
}
