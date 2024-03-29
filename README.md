# Singet

Singet means single get.
It's a simple web proxy to fetch (get, download, or stream) a single URL.
It's written in the Go programming language.

## How to use

After running, you can fetch a single URL this way:

```
[host]:[port]/?url=[url]
```

For example:

```
https://example.com/?url=https://target.com/video.mp4
```

The proxy passes the request body and headers (HTTP method, referrer, auth, etc.) to the target URL.

### Base64 Encoding

If you want to use this proxy on a censored network, you might need to encode URLs.
To do so, you can use the `base64` parameter instead of the `url` like in this example:

```
// Using "url" parameter:
https://example.com/?url=https://miladrahimi.com

// Using "base64" parameter:
https://example.com/?base64=aHR0cHM6Ly9taWxhZHJhaGltaS5jb20=
```

### Redirection

The target URL might return an HTTP redirection response.
In this case, the proxy behaves based on the `r` parameter or its default behavior when the parameter is not present.
You can set this parameter this way:

```
[host]:[port]/?url=[url]&r=[value]
```

#### Default

When `r` parameter is `default` or missing, it returns the response without manipulation.
For example, if the request was `http://example.com/?url=https://google.com` the response would be a standard redirection
to `https://www.google.com`.

#### Follow

When `r` parameter is `follow`, it redirects through the proxy.
For example, if the request was `http://example.com/?r=follow&url=https://google.com` the response would be a
redirection to `http://example.com/?r=follow&url=https://www.google.com`.

#### Stop

When `r` parameter is `stop`, it returns a JSON response contains the new location.
For example, if the request was `http://example.com/?r=stop&url=https://google.com` the response would be:

```
HTTP/1.1 200 OK
Content-Type: application/json
...

{"location":"https://www.google.com/"}
```

### HTTP Header Manipulation

By default, Singet passes the request headers without any change.
You may want to manipulate some HTTP headers.
In this case, you can pass the related header in the query string with the prefix `h__` like in the example below.

```
https://example.com/?url=https://www.google.com&h__referer=https://images.google.com
```

It will set the `REFERER` header to `https://images.google.com`.

## License
Singet is initially created by [Milad Rahimi](http://miladrahimi.com)
and released under the [MIT License](http://opensource.org/licenses/mit-license.php).
