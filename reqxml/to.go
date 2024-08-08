package reqxml

import (
	"encoding/xml"

	"github.com/raeperd/restclient"
)

// To decodes a response as an XML object.
func To(v any) restclient.ResponseHandler {
	return restclient.ToDeserializer(xml.Unmarshal, v)
}
