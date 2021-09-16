package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

func isGZipAcceptable(r *http.Request) bool {
	// e.g. `Accept-Encoding: deflate, gzip;q=1.0, *;q=0.5`
	for _, v := range r.Header["Accept-Encoding"] {
		if strings.Contains(v, "gzip") {
			return true
		}
	}
	return false
}

func main() {
	l, err := net.Listen("tcp4", ":8080")
	if err != nil {
		panic(err)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go func() {
			defer c.Close()
			// ループさせることで HTTP/1.1 の Keep-Alive にする
			for {
				err := c.SetReadDeadline(time.Now().Add(5 * time.Second))
				if err != nil {
					panic(err)
				}
				req, err := http.ReadRequest(
					bufio.NewReader(c))
				if err != nil {
					if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
						break
					}
					// 相手がソケットをクローズした
					if err == io.EOF {
						break
					}
					panic(err)
				}

				// Response.Write() は HTTP/1.1 より古いバージョンが使われる場合
				// もしくは長さがわからない場合は Connection: close ヘッダを付与してしまう
				resp := http.Response{
					StatusCode: 200,
					ProtoMajor: 1,
					ProtoMinor: 1,
					Header:     make(http.Header),
				}
				if isGZipAcceptable(req) {
					buff := &bytes.Buffer{}
					writer := gzip.NewWriter(buff)

					content := "Hello World (gzipped)\n"
					if _, err := io.WriteString(writer, content); err != nil {
						panic(err)
					}
					// Close しないと書き込まれない
					if err = writer.Close(); err != nil {
						panic(err)
					}
					resp.Body = ioutil.NopCloser(buff)
					resp.ContentLength = int64(buff.Len())
					// ヘッダは圧縮されない
					// 少量のデータを通信するほど効率が悪くなる
					resp.Header.Set("Content-Encoding", "gzip")
				} else {
					content := "Hello World\n"
					resp.Body = ioutil.NopCloser(strings.NewReader(content))
					resp.ContentLength = int64(len(content))
				}
				err = resp.Write(c)
				if err != nil {
					panic(err)
				}
			}
		}()
	}
}
