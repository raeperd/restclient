package restclient_test

import (
	"context"
	"fmt"

	"github.com/raeperd/restclient"
)

func ExampleHasStatusErr() {
	err := restclient.
		URL("http://example.com/404").
		CheckStatus(200).
		Fetch(context.Background())
	if restclient.HasStatusErr(err, 404) {
		fmt.Println("got a 404")
	}
	// Output:
	// got a 404
}
