package restclient_test

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/raeperd/restclient"
)

func ExampleValidatorHandler() {
	var (
		regularBody string
		errBody     string
	)

	// If we fail validation because the response is a 404,
	// we handle the body with errBody instead of regularBody
	// for separate processing.
	err := restclient.
		URL("http://example.com/404").
		ToString(&regularBody).
		AddValidator(
			restclient.ValidatorHandler(
				restclient.DefaultValidator,
				restclient.ToString(&errBody),
			)).
		Fetch(context.Background())
	switch {
	case errors.Is(err, restclient.ErrInvalidHandled):
		fmt.Println("got errBody:",
			strings.Contains(errBody, "Example Domain"))
	case err != nil:
		fmt.Println("unexpected error", err)
	case err == nil:
		fmt.Println("unexpected success")
	}

	fmt.Println("got regularBody:", strings.Contains(regularBody, "Example Domain"))
	// Output:
	// got errBody: true
	// got regularBody: false
}
