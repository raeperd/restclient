package restclient_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/raeperd/restclient"
)

func ExampleNewCookieJar() {
	// Create a client that preserve cookies between requests
	myClient := *http.DefaultClient
	myClient.Jar = restclient.NewCookieJar()
	// Use the client to make a request
	err := restclient.
		URL("http://httpbin.org/cookies/set/chocolate/chip").
		Client(&myClient).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("could not connect to httpbin.org:", err)
	}
	// Now check that cookies we got
	for _, cookie := range myClient.Jar.Cookies(&url.URL{
		Scheme: "http",
		Host:   "httpbin.org",
	}) {
		fmt.Println(cookie)
	}
	// And we'll see that they're reused on subsequent requests
	var cookies struct {
		Cookies map[string]string
	}
	err = restclient.
		URL("http://httpbin.org/cookies").
		Client(&myClient).
		ToJSON(&cookies).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("could not connect to httpbin.org:", err)
	}
	fmt.Println(cookies)

	// And we can manually add our own cookie values
	// without overriding existing ones
	err = restclient.
		URL("http://httpbin.org/cookies").
		Client(&myClient).
		Cookie("oatmeal", "raisin").
		ToJSON(&cookies).
		Fetch(context.Background())
	if err != nil {
		fmt.Println("could not connect to httpbin.org:", err)
	}
	fmt.Println(cookies)

	// Output:
	// chocolate=chip
	// {map[chocolate:chip]}
	// {map[chocolate:chip oatmeal:raisin]}
}
