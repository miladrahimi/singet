# SingleFetch
A simple web proxy that is able to fetch (stream) a single URL.

It is written in Go language and ready to deploy on [Heroku](https://heroku.com).

## How to use

After running the application, you can fetch a single url this way:

```
[host]:[port]/?url=[url]
```

For example:

```
https://proxy.com/?url=https://url.com/video.mp4
```

The proxy pass the request body and headers (http method, referrer, auth, etc.) to the requested url.

### Redirection

The requested url might return a http redirection response.
In this case the proxy behave based on the `redirection` parameter.

```
[host]:[port]/?url=[url]&redirection=[value]
```

#### Default behaviour
Without the `redirection` parameter the default behaviour is returning the response with no manipulation.
For example, if the request was `http://proxy.com/?url=https://google.com`
the response would be redirection to `https://www.google.com`.
It might not be your expected result, because if you make this request on using your browser,
it redirects you to the real Google address (`https://www.google.com`) with no proxy.

#### follow
With the `follow` value for the `redirection` parameter the behaviour is redirecting through the proxy.
For example, if the request was `http://proxy.com/?url=https://google.com&redirection=follow`
the response would be redirection to `http://proxy.com/?url=https://www.google.com&redirection=follow`.
As you can see, if your client handle redirection it redirects you to the new url through the proxy.

#### stop
With the `stop` value for the `redirection` parameter the behaviour is returning a JSON response contains location.
For example, if the request was `http://proxy.com/?url=https://google.com&redirection=stop`
the response would be:

```
HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Content-Length: ...
Content-Type: application/json
Date: ...

{"location":"https://www.google.com/"}
```

## License
OpenURL is initially created by [Milad Rahimi](http://miladrahimi.com)
and released under the [MIT License](http://opensource.org/licenses/mit-license.php).
