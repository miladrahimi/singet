package main

import (
    "bytes"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "os"
    "strconv"
)

// main is the application entry point
func main() {
    http.HandleFunc("/favicon.ico", handleFavicon)
    http.HandleFunc("/redirect", handleRedirect)
    http.HandleFunc("/debug", handleDebug)
    http.HandleFunc("/", handleProxy)

    log.Println("Serving on " + port())
    log.Fatal(http.ListenAndServe(port(), nil))
}

// port returns the proxy port
func port() string {
    if port := os.Getenv("PORT"); port == "" {
        return ":8080"
    } else {
        return ":" + port
    }
}

// address returns the proxy address (URL)
func address() string {
    if address := os.Getenv("URL"); address == "" {
        return "http://127.0.0.1" + port()
    } else {
        return address
    }
}

// displayError will create an error response
func displayError(rw http.ResponseWriter, error string) {
    rw.WriteHeader(http.StatusBadRequest)
    rw.Header().Set("Content-type", "application/json")

    body, _ := json.Marshal(map[string]string{"error": error})

    if _, err := rw.Write(body); err != nil {
        panic(err)
    }
}

// displayLocation will modify the response to display a location instead of redirecting
func displayLocation(r *http.Response, location string) {
    body, _ := json.Marshal(map[string]string{"location": location})

    r.Body = io.NopCloser(bytes.NewReader(body))
    r.ContentLength = int64(len(body))
    r.StatusCode = http.StatusOK
    r.Header = http.Header{}
    r.Header.Set("Content-Type", "application/json")
    r.Header.Set("Content-Length", strconv.Itoa(len(body)))
}

// handleFavicon returns the favicon.ico file
func handleFavicon(rw http.ResponseWriter, r *http.Request) {
    content, _ := os.ReadFile("favicon.ico")
    rw.Header().Set("Content-Type", "image/x-icon")
    rw.Header().Set("X-Forwarded-Host", r.URL.Host)
    _, _ = fmt.Fprintln(rw, string(content))
}

func handleRedirect(rw http.ResponseWriter, r *http.Request) {
    http.Redirect(rw, r, "https://github.com/miladrahimi/singet", http.StatusSeeOther)
}

func handleDebug(rw http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        panic(err)
    }

    response, err := json.MarshalIndent(map[string]interface{}{
        "Host":    r.Host,
        "URI":     r.RequestURI,
        "headers": r.Header,
        "body":    body,
    }, "", " ")
    if err != nil {
        panic(err)
    }

    rw.WriteHeader(http.StatusOK)
    rw.Header().Set("Content-type", "application/json")
    rw.Header().Set("Access-Control-Allow-Origin", "www.google.com")
    _, _ = rw.Write(response)
}

func handleProxy(rw http.ResponseWriter, request *http.Request) {
    // Handle OPTIONS requests
    if request.Method == "OPTIONS" {
        rw.WriteHeader(204)
        return
    }

    var requestedUrl string

    // Receive requested URL
    query := request.URL.Query()
    if query.Get("url") != "" {
        requestedUrl = query.Get("url")
    } else if query.Get("base64") != "" {
        if u, err := base64.StdEncoding.DecodeString(query.Get("base64")); err != nil {
            displayError(rw, "Invalid Base64-encoded URL.")
            return
        } else {
            requestedUrl = string(u)
        }
    } else {
        displayError(rw, "Nothing requested.")
        return
    }

    // Validate the requested URL
    target, err := url.Parse(requestedUrl)
    if err != nil || target.IsAbs() == false {
        displayError(rw, "The URL is invalid.")
        return
    }

    // Setup proxy
    proxy := &httputil.ReverseProxy{
        Director: func(r *http.Request) {
            // Modify headers (inject h__ headers from query string)
            for parameterName := range request.URL.Query() {
                if parameterName[0:3] == "h__" {
                    r.Header.Del(parameterName[0:3])
                    r.Header.Set(parameterName[0:3], request.URL.Query().Get(parameterName))
                }
            }

            // Set target
            r.Host = target.Host
            r.URL = target
        },

        ModifyResponse: func(r *http.Response) error {
            r.Header.Set("Access-Control-Allow-Origin", "*")

            if r.Header.Get("Location") != "" {
                if query.Get("redirection") == "stop" {
                    displayLocation(r, r.Header.Get("Location"))
                } else {
                    r.Header.Set("Location", address()+"/?url="+r.Header.Get("Location"))
                }
            }

            return nil
        },
    }

    // Serve proxy
    proxy.ServeHTTP(rw, request)
}
