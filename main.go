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
			displayError(rw, "The requested URL is invalid.")
			return
		}

		// Make reverse proxy
		proxy := &httputil.ReverseProxy{
			Director: func(r *http.Request) {
				if referrer := query.Get("referrer"); referrer != "" {
					r.Header.Set("Referer", referrer)
				} else {
					r.Header.Set("Referer", target.Scheme+"://"+target.Host+"/")
				}

				r.Header.Set("Origin", target.Host)
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
						modifyResponseToRedirection(r, r.Header.Get("Location"))
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
		panic("Cannot respond to the request.")
	}
}

// modifyResponseToRedirection will modify the response to a location response
func modifyResponseToRedirection(r *http.Response, location string) {
	body, err := json.Marshal(map[string]string{"location": location})
	if err != nil {
		panic("Cannot respond to the request.")
	}

	r.Body = ioutil.NopCloser(bytes.NewReader(body))
	r.ContentLength = int64(len(body))
	r.StatusCode = http.StatusOK
	r.Header = http.Header{}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Content-Length", strconv.Itoa(len(body)))
}
