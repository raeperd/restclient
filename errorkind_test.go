package restclient_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/raeperd/restclient"
	"github.com/raeperd/restclient/internal/be"
)

func TestErrorKind(t *testing.T) {
	t.Parallel()
	var none restclient.ErrorKind = -1
	kinds := []restclient.ErrorKind{
		restclient.ErrURL,
		restclient.ErrRequest,
		restclient.ErrTransport,
		restclient.ErrValidator,
		restclient.ErrHandler,
	}
	ctx := context.Background()
	res200 := restclient.ReplayString("HTTP/1.1 200 OK\n\n")
	for _, tc := range []struct {
		ctx  context.Context
		want restclient.ErrorKind
		b    *restclient.Builder
	}{
		{ctx, none, restclient.
			URL("").
			Transport(res200),
		},
		{ctx, restclient.ErrURL, restclient.
			URL("http://%2020").
			Transport(res200),
		},
		{ctx, restclient.ErrURL, restclient.
			URL("hello world").
			Transport(res200),
		},
		{ctx, none, restclient.
			URL("http://world/#hello").
			Transport(res200),
		},
		{ctx, restclient.ErrRequest, restclient.
			URL("").
			Body(func() (io.ReadCloser, error) {
				return nil, errors.New("x")
			}).
			Transport(res200),
		},
		{ctx, restclient.ErrRequest, restclient.
			URL("").
			Method(" ").
			Transport(res200),
		},
		{nil, restclient.ErrRequest, restclient.
			URL("").
			Transport(res200),
		},
		{ctx, restclient.ErrTransport, restclient.
			URL("").
			Transport(restclient.ReplayString("")),
		},
		{ctx, restclient.ErrValidator, restclient.
			URL("").
			Transport(restclient.ReplayString("HTTP/1.1 404 Nope\n\n")),
		},
		{ctx, restclient.ErrHandler, restclient.
			URL("").
			Transport(res200).
			ToJSON(nil),
		},
	} {
		err := tc.b.Fetch(tc.ctx)
		for _, kind := range kinds {
			match := errors.Is(err, kind)
			be.Equal(t, kind == tc.want, match)
		}
		var askind = none
		be.Equal(t, tc.want != none, errors.As(err, &askind))
		be.Equal(t, tc.want, askind)
	}
}
