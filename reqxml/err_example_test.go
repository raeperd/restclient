package reqxml_test

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/raeperd/restclient"
	"github.com/raeperd/restclient/reqxml"
)

func ExampleError() {
	type ErrorXML struct {
		XMLName   xml.Name `xml:"Error"`
		Code      string   `xml:"Code"`
		Message   string   `xml:"Message"`
		RequestID string   `xml:"RequestId"`
		HostID    string   `xml:"HostId"`
	}

	var errObj ErrorXML
	err := restclient.
		URL("http://example.s3.us-east-1.amazonaws.com").
		AddValidator(reqxml.Error(&errObj)).
		Fetch(context.Background())
	switch {
	case errors.Is(err, restclient.ErrInvalidHandled):
		fmt.Println(errObj.Message)
	case err != nil:
		fmt.Println("Error!", err)
	case err == nil:
		fmt.Println("unexpected success")
	}

	// Output:
	// Access Denied
}
