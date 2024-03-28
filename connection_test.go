package annkesdk

import (
	"errors"
	"io"
	"testing"

	"github.com/csrar/annkeSDK/models"
	"github.com/stretchr/testify/assert"
)

func TestConnector_generateRandom(t *testing.T) {
	random := generateRandom()
	assert.GreaterOrEqual(t, random, 0)
	assert.Less(t, random, randomLenght)
}

func TestConnector_getProtocol(t *testing.T) {
	tests := []struct {
		name     string
		secure   bool
		expected string
	}{
		{"Secure", true, "https"},
		{"Non-Secure", false, "http"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Connector{Secure: tt.secure}
			protocol := c.getProtocol()
			assert.Equal(t, tt.expected, protocol)
		})
	}
}

func TestConnector_prepareLoginRequest(t *testing.T) {
	c := Connector{
		User:     "user",
		Password: "pass",
		Host:     "example.com",
		Secure:   true,
	}

	req, err := c.prepareLoginRequest()
	assert.NoError(t, err, "Unexpected error from prepareLoginRequest")

	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "https", req.URL.Scheme)
	assert.Equal(t, "https://user:pass@example.com"+loginPath+"?username=user&random=", req.URL.String()[:92])
}

func mockLoginResponse() models.LoginResponse {
	response := models.LoginResponse{
		Version:          "1.0",
		Xmlns:            "http://example.com",
		SessionID:        "abc123",
		Challenge:        "xyz456",
		Iterations:       1000,
		IsIrreversible:   true,
		Salt:             "salty",
		SessionIDVersion: 2,
	}
	response.IsSessionIDValidLongTerm.Text = true
	response.IsSessionIDValidLongTerm.Opt = "optional"
	return response
}

func TestConnector_prepareSessionRequest(t *testing.T) {
	c := Connector{
		User:     "user",
		Host:     "example.com",
		Secure:   true,
		Password: "pass",
	}
	expectedBody := "<SessionLogin><userName>user</userName><password>mock-pass</password><sessionID>abc123</sessionID><isSessionIDValidLongTerm>true</isSessionIDValidLongTerm><sessionIDVersion>2</sessionIDVersion></SessionLogin>"
	expectedUrl := "https://example.com/ISAPI/Security/sessionLogin?timeStamp="
	request, _ := c.prepareSessionRequest(mockLoginResponse(), "mock-pass")
	data, _ := io.ReadAll(request.Body)
	assert.Equal(t, expectedBody, string(data))
	assert.Equal(t, expectedUrl, request.URL.String()[:58])
	assert.Equal(t, "GET", request.Method)
}

func TestConnector_prepareSessionRequestError(t *testing.T) {
	c := Connector{
		User:     "user",
		Host:     "\\",
		Secure:   true,
		Password: "pass",
	}

	request, err := c.prepareSessionRequest(mockLoginResponse(), "mock-pass")
	assert.Nil(t, request)
	assert.Error(t, err, errors.New(`"parse "https://\\example.com/ISAPI/Security/sessionLogin?timeStamp=1711596731": invalid character "\\" in host name"`))
}

func TestHashPassword(t *testing.T) {
	connector := Connector{
		User:     "testuser",
		Password: "testpassword",
	}

	login := models.LoginResponse{
		Version:          "1.0",
		Xmlns:            "http://example.com",
		SessionID:        "abc123",
		Challenge:        "xyz456",
		Iterations:       10,
		IsIrreversible:   true,
		Salt:             "salty",
		SessionIDVersion: 2,
	}
	login.IsSessionIDValidLongTerm.Text = true
	login.IsSessionIDValidLongTerm.Opt = "optional"

	cases := []struct {
		name           string
		inIrreversible bool
		expectedHash   string
	}{
		{
			name:           "Is irreversible",
			inIrreversible: true,
			expectedHash:   "5bbb0d3b38d162f1fdcaa8a5e4bc53cf02914a113f1fe0b4590551c55002410f",
		},
		{
			name:           "Is not irreversible",
			inIrreversible: false,
			expectedHash:   "c489bf27321c295f8f34273a130bffac7021a6c1af85fdf409db62bd0cfee48d",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			login.IsIrreversible = tc.inIrreversible
			result := connector.hashPassword(login)
			assert.Equal(t, tc.expectedHash, result)
		})
	}

}

func TestConnector_NewConnectorErrors(t *testing.T) {
	cases := []struct {
		name          string
		inUser        string
		inHost        string
		expectedError string
	}{
		{
			name:          "missing user",
			expectedError: "error initializing Annke connection, missing parameter: User",
			inHost:        "mock",
		},
		{
			name:          "missing host",
			expectedError: "error initializing Annke connection, missing parameter: Host",
			inUser:        "mock",
		},
		{
			name:          "invalid URL",
			inUser:        "mock-user",
			inHost:        "\\",
			expectedError: "Error preparing login request parse \"http://mock-user:mock-password@\\\\/ISAPI/Security/sessionLogin/capabilities?username=mock-user&random=",
		},
	}

	for _, tc := range cases {
		t.Run(t.Name(), func(t *testing.T) {
			cfg, err := NewConnector(tc.inHost, tc.inUser, "mock-password", false)
			assert.Nil(t, cfg)
			assert.ErrorContains(t, err, tc.expectedError)
		})

	}
}
