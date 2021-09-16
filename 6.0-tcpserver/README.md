# ソケット勉強用の HTTP サーバ
* HTTP/1.1 に対応
* Keep-Alive に対応
* レスポンスの gzip 圧縮に対応

## 実行例: gzip 圧縮なし
```
$ curl -v localhost:8080
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Length: 12
<
Hello World
* Connection #0 to host localhost left intact
```

## 実行例: gzip 圧縮あり
```
$ curl --compressed -v localhost:8080
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*
> Accept-Encoding: deflate, gzip, br
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Length: 46
< Content-Encoding: gzip
<
Hello World (gzipped)
* Connection #0 to host localhost left intact
```
