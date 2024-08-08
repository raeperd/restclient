package reqxml

import "github.com/raeperd/restclient"

// Error is a ValidatorHandler that applies DefaultValidator
// and decodes the response as an XML object
// if the DefaultValidator check fails.
func Error(v any) restclient.ResponseHandler {
	return restclient.ValidatorHandler(restclient.DefaultValidator, To(v))
}
