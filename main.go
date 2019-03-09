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
		displayFavicon(w)
	})

	// Proxy main route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Parse requested url
		target, err := url.Parse(r.URL.RawQuery)
		if err != nil || target.IsAbs() == false {
			displayError(w, "Bad Request.")
			return
		}

		// Make reverse proxy
		proxy := &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.Header.Set("Referer", target.Scheme+"://"+target.Host+"/")
				req.Header.Set("Origin", target.Host)
				req.Host = target.Host
				req.URL = target

				log.Println(req)
			},

			ModifyResponse: func(response *http.Response) error {
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

	_, err := w.Write([]byte(message))
	if err != nil {
		panic("Cannot respond to the request.")
	}
}

// Create the favicon.ico
func displayFavicon(w http.ResponseWriter) {
	bytes, err := ioutil.ReadFile("resources/logo.png")
	if err != nil {
		_, _ = fmt.Fprintln(w, "")
	}

	w.Header().Set("Content-Type", "image/png")
	_, _ = fmt.Fprintln(w, string(bytes))
}
