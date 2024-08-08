package reqxml

import (
	"encoding/xml"

	"github.com/raeperd/restclient"
)

// Body is a BodyGetter that marshals a XML object.
func Body(v any) restclient.BodyGetter {
	return restclient.BodySerializer(xml.Marshal, v)
}

// BodyConfig sets the Builder's request body to the marshaled XML.
// It also sets ContentType to "application/xml"
// if it is not otherwise set.
func BodyConfig(v any) restclient.Config {
	return func(rb *restclient.Builder) {
		rb.
			Body(Body(v)).
			HeaderOptional("Content-Type", "application/xml")
	}
}
