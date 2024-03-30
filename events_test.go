package annkesdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testGenerateMockConnector() Connector {
	return Connector{
		User:     "mock-user",
		Password: "mock-password",
		Secure:   false,
	}
}

func TestConnector_makeUpdateRequestBadBody(t *testing.T) {
	connector := testGenerateMockConnector()
	err := connector.makeUpdateRequest("/", make(chan int))
	assert.ErrorContains(t, err, "xml: unsupported type: chan int")
}

func TestConnector_makeUpdateRequestInvalidUrl(t *testing.T) {
	connector := testGenerateMockConnector()
	connector.Host = "mock host"
	err := connector.makeUpdateRequest("/", "body")
	assert.ErrorContains(t, err, "parse \"http://mock host/\": invalid character \" \" in host name")
}

func TestConnector_makeUpdateRequestNonExistingHost(t *testing.T) {
	connector := testGenerateMockConnector()
	connector.Host = "localhost:9999"
	err := connector.makeUpdateRequest("/", "body")
	assert.ErrorContains(t, err, "Put \"http://localhost:9999/\": dial tcp 127.0.0.1:9999: connect: connection refused")
}
