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
	DefaultHost = "http://0.0.0.0"
	DefaultPort = "8080"
)

func main() {
	proxyUrl := host() + ":" + port()

	http.HandleFunc("/favicon.ico", func(rw http.ResponseWriter, r *http.Request) {
		favicon(rw)
	})

	http.HandleFunc("/", func(rw http.ResponseWriter, request *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")

		query := request.URL.Query()

		if query.Get("url") == "" {
			displayError(rw, "Nothing requested.")
			return
		}

		target, err := url.Parse(query.Get("url"))
		if err != nil || target.IsAbs() == false {
			displayError(rw, "URL is invalid.")
			return
		}

		// Make a reverse proxy
		proxy := &httputil.ReverseProxy{
			Director: func(r *http.Request) {
				for name, value := range request.URL.Query() {
					if name[0:3] == "h__" {
						r.Header.Del(name[3:])

						for _, v := range value {
							r.Header.Set(name[3:], v)
						}
					}
				}

				r.Host = target.Host
				r.URL = target

				log.Println(r)
			},

			ModifyResponse: func(r *http.Response) error {
				r.Header.Del("Access-Control-Allow-Origin")

				// Handle redirection responses
				if r.StatusCode >= 300 && r.StatusCode < 400 && r.Header.Get("Location") != "" {
					if query.Get("redirection") == "follow" {
						r.Header.Set("Location", proxyUrl+"/?redirection=follow&url="+r.Header.Get("Location"))
					} else if query.Get("redirection") == "stop" {
						displayLocation(r, r.Header.Get("Location"))
					}
				}

				return nil
			},
		}

		proxy.ServeHTTP(rw, request)
	})

	log.Println("Serving on " + proxyUrl)
	log.Fatal(http.ListenAndServe(":"+port(), nil))
}

// host will return the proxy host
func host() string {
	if host := os.Getenv("HOST"); host == "" {
		return DefaultHost
	} else {
		return host
	}
}

// port will return the proxy port
func port() string {
	if port := os.Getenv("PORT"); port == "" {
		return DefaultPort
	} else {
		return port
	}
}

// favicon will return the favicon icon
func favicon(rw http.ResponseWriter) {
	content, _ := ioutil.ReadFile("favicon.ico")
	rw.Header().Set("Content-Type", "image/x-icon")
	_, _ = fmt.Fprintln(rw, string(content))
}

// displayError will create an error response
func displayError(rw http.ResponseWriter, error string) {
	rw.WriteHeader(http.StatusBadRequest)
	rw.Header().Set("Content-type", "application/json")

	body, err := json.Marshal(map[string]string{"error": error})
	if err != nil {
		panic("Cannot generate the JSON.")
	}

	_, err = rw.Write(body)
	if err != nil {
		panic("Cannot display the error.")
	}
}

// displayLocation will modify the response to a display location instead of redirecting
func displayLocation(r *http.Response, location string) {
	body, err := json.Marshal(map[string]string{"location": location})
	if err != nil {
		panic("Cannot display the location.")
	}

	r.Body = ioutil.NopCloser(bytes.NewReader(body))
	r.ContentLength = int64(len(body))
	r.StatusCode = http.StatusOK
	r.Header = http.Header{}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Content-Length", strconv.Itoa(len(body)))
}
