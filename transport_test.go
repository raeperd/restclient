package restclient_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/raeperd/restclient"
	"github.com/raeperd/restclient/internal/be"
)

func TestUserAgentTransport(t *testing.T) {
	// Wrap an existing transport or use nil for http.DefaultTransport
	baseTrans := http.DefaultClient.Transport
	trans := restclient.UserAgentTransport(baseTrans, "my-user/agent")

	var headers postman
	err := restclient.
		URL("https://postman-echo.com/get").
		Transport(trans).
		ToJSON(&headers).
		Fetch(context.Background())
	be.NilErr(t, err)
	be.Equal(t, "my-user/agent", headers.Headers["user-agent"])
}
