package restclient_test

import (
	"context"
	"fmt"
	"testing/fstest"

	"github.com/raeperd/restclient"
)

func ExampleReplayFS() {
	fsys := fstest.MapFS{
		"fsys.example - MKIYDwjs.res.txt": &fstest.MapFile{
			Data: []byte(`HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Date: Mon, 24 May 2021 18:48:50 GMT

An example response.`),
		},
	}
	var s string
	const expected = `An example response.`
	if err := restclient.
		URL("http://fsys.example").
		Transport(restclient.ReplayFS(fsys)).
		ToString(&s).
		Fetch(context.Background()); err != nil {
		panic(err)
	}
	fmt.Println(s == expected)
	// Output:
	// true
}
