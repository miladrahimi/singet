# SingleFetch
A simple web proxy that is able to fetch (stream) a single URL.

It is written in Go language and ready to deploy on [Heroku](https://heroku.com).

## How to use

After running the application, you can fetch a single url this way:

```
[your-server]/?url=[url]
```

For example:

```
https://example.com/?https://target.com/video.mp4
```

You may set referrer this way:

```
[your-server]/?url=[url]&referrer=[another-url]
```

And if you don't want the proxy to follow redirection:

```
[your-server]/?url=[url]&follow=false
```

## License
OpenURL is initially created by [Milad Rahimi](http://miladrahimi.com)
and released under the [MIT License](http://opensource.org/licenses/mit-license.php).