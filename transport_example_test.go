package restclient_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/raeperd/restclient"
)

func ExampleReplayString() {
	const res = `HTTP/1.1 200 OK

An example response.`

	var s string
	const expected = `An example response.`
	if err := restclient.
		URL("http://response.example").
		Transport(restclient.ReplayString(res)).
		ToString(&s).
		Fetch(context.Background()); err != nil {
		panic(err)
	}
	fmt.Println(s == expected)
	// Output:
	// true
}

func ExamplePermitURLTransport() {
	// Wrap an existing transport or use nil for http.DefaultTransport
	baseTrans := http.DefaultClient.Transport
	trans := restclient.PermitURLTransport(baseTrans, `^http://example\.com/`)
	var s string
	if err := restclient.
		URL("http://example.com/").
		Transport(trans).
		ToString(&s).
		Fetch(context.Background()); err != nil {
		panic(err)
	}
	fmt.Println(strings.Contains(s, "Example Domain"))

	if err := restclient.
		URL("http://unauthorized.example.com/").
		Transport(trans).
		ToString(&s).
		Fetch(context.Background()); err != nil {
		fmt.Println(err) // unauthorized subdomain not allowed
	}
	// Output:
	// true
	// ErrTransport: Get "http://unauthorized.example.com/": requested URL not permitted by regexp: ^http://example\.com/
}

func ExampleRoundTripFunc() {
	// Wrap an underlying transport in order to add request middleware
	baseTrans := http.DefaultClient.Transport
	var checksumTransport restclient.RoundTripFunc = func(req *http.Request) (res *http.Response, err error) {
		// Read and checksum the body
		b, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		h := md5.New()
		h.Write(b)
		checksum := fmt.Sprintf("%X", h.Sum(nil))
		// Must clone requests before modifying them
		req2 := *req
		req2.Header = req.Header.Clone()
		// Add header and body to the clone
		req2.Header.Add("Checksum", checksum)
		req2.Body = io.NopCloser(bytes.NewBuffer(b))
		return baseTrans.RoundTrip(&req2)
	}
	var data postman
	err := restclient.
		URL("https://postman-echo.com/post").
		BodyBytes([]byte(`Hello, World!`)).
		ContentType("text/plain").
		Transport(checksumTransport).
		ToJSON(&data).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("Error!", err)
	}
	fmt.Println(data.Headers["checksum"])
	// Output:
	// 65A8E27D8879283831B664BD8B7F0AD4
}

func ExampleLogTransport() {
	logger := func(req *http.Request, res *http.Response, err error, d time.Duration) {
		fmt.Printf("method=%q url=%q err=%v status=%q duration=%v\n",
			req.Method, req.URL, err, res.Status, d.Round(1*time.Second))
	}
	// Wrap an existing transport or use nil for http.DefaultTransport
	baseTrans := http.DefaultClient.Transport
	trans := restclient.LogTransport(baseTrans, logger)
	var s string
	if err := restclient.
		URL("http://example.com/").
		Transport(trans).
		ToString(&s).
		Fetch(context.Background()); err != nil {
		fmt.Println("Error!", err)
	}
	// Works for bad responses too
	baseTrans = restclient.ErrorTransport(errors.New("can't connect"))

	trans = restclient.LogTransport(baseTrans, logger)

	if err := restclient.
		URL("http://example.com/").
		Transport(trans).
		ToString(&s).
		Fetch(context.Background()); err != nil {
		fmt.Println("Error!", err)
	}
	// Output:
	// method="GET" url="http://example.com/" err=<nil> status="200 OK" duration=0s
	// method="GET" url="http://example.com/" err=can't connect status="" duration=0s
	// Error! ErrTransport: Get "http://example.com/": can't connect
}
