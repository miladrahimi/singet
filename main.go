package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
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
		content, _ := ioutil.ReadFile("favicon.ico")
		w.Header().Set("Content-Type", "image/x-icon")
		_, _ = fmt.Fprintln(w, string(content))
	})

	// Proxy main route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		query := r.URL.Query()

		if query.Get("url") == "" {
			displayMessage(w, "Welcome!")
			return
		}

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

			ModifyResponse: func(r *http.Response) error {
				// Handle redirection responses
				if location := r.Header.Get("Location"); location != "" {
					if query.Get("follow") != "false" {
						r.Header.Set("Location", proxyUrl+"/?url="+r.Header.Get("Location"))
					} else {
						modifyResponseToRedirection(r, location)
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
	w.Header().Set("Content-type", "text/plain")

	_, err := w.Write([]byte("Error: " + message))
	if err != nil {
		panic("Cannot respond to the request.")
	}
}

// Create an error response
func displayMessage(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "text/plain")

	_, err := w.Write([]byte(message))
	if err != nil {
		panic("Cannot respond to the request.")
	}
}

// Modify response to location
func modifyResponseToRedirection(response *http.Response, location string) {
	body, err := json.Marshal(map[string]string{"location": location})
	if err != nil {
		panic("Cannot respond to the request.")
	}

	response.Status = "200 OK"
	response.StatusCode = http.StatusOK
	response.Header.Del("Location")
	response.Header.Set("Content-Type", "application/json")
	response.Body = ioutil.NopCloser(bytes.NewReader(body))
	response.ContentLength = int64(len(body))
	response.Header.Set("Content-Length", strconv.Itoa(len(body)))
}
