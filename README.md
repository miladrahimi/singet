# SingleFetch

A simple web proxy that can fetch (or stream) a single URL. It is written in Go language and ready to deploy on [Heroku](https://heroku.com).

## How to use

After running, you can fetch a single URL this way:

```
[host]:[port]/?url=[url]
```

For example:

```
https://domain.com/?url=https://url.com/video.mp4
```

The proxy passes the request body and headers (HTTP method, referrer, auth, etc.) to the requested URL.

### Redirection

The requested URL might return an HTTP redirection response. In this case, the proxy behaves based on the `redirection` parameter or its default behavior when the parameter is not present. You can pass this parameter to the app this way:

```
[host]:[port]/?url=[url]&redirection=[value]
```

#### Default behaviour

The default behavior is returning the response with no manipulation. For example, if the request was `http://proxy.com/?url=https://google.com` the response would be a redirection to `https://www.google.com`.

#### follow

When `redirection` is `follow`, It will be redirected through the proxy. For example, if the request was `http://proxy.com/?url=https://google.com&redirection=follow` the response would be a redirection to `http://proxy.com/?url=https://www.google.com&redirection=follow`.

#### stop

When `redirection` is `stop`, It returns a JSON response contains the taregt location. For example, if the request was `http://proxy.com/?url=https://google.com&redirection=stop` the response would be:

```http
HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Content-Length: ...
Content-Type: application/json
Date: ...

{"location":"https://www.google.com/"}
```

### HTTP Header Manipulation

In default, SingleFetch passes the request headers without any change, but you may want to manipulate some headers like referer or any other header. In this case, you can pass the related header in the query string with the prefix `h_` like this example:

```
http://proxy.com/?url=https://google.com&h_referer=http://images.google.com
```

## License
SingleFetch is initially created by [Milad Rahimi](http://miladrahimi.com)
and released under the [MIT License](http://opensource.org/licenses/mit-license.php).
