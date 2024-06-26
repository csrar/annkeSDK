package annkesdk

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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
		{
			name:          "non existing host url",
			inUser:        "mock-user",
			inHost:        "localhost:9999",
			expectedError: "Error executing login request Get",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := NewConnector(tc.inHost, tc.inUser, "mock-password", false)
			assert.Nil(t, cfg)
			assert.ErrorContains(t, err, tc.expectedError)
		})
	}
}

func testLoginHanlerHelper(loginResponse, sessionResponse string, loginStatus, sessionStatus int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case loginPath:
			{
				w.WriteHeader(loginStatus)
				fmt.Fprintln(w, loginResponse)
			}
		case sessionPath:
			{
				w.WriteHeader(sessionStatus)
				fmt.Fprintln(w, sessionResponse)
			}
		default:
			fmt.Fprintln(w, "mock-server")
		}
	})
}

func TestConnector_NewConnector(t *testing.T) {
	cases := []struct {
		name            string
		expectedError   string
		loginResponse   string
		sessionResponse string
		loginStatus     int
		sessionStatus   int
	}{
		{
			name:          "invalid login response",
			loginResponse: "mock-error",
			loginStatus:   http.StatusInternalServerError,
			expectedError: "received unexpected response from: /ISAPI/Security/sessionLogin/capabilities status: 500 payload: mock-error",
		},
		{
			name:          "invalid login response",
			loginResponse: "mock-response",
			loginStatus:   http.StatusOK,
			expectedError: "Error unmarshaling login response",
		},
		{
			name:            "invalid session status",
			loginResponse:   "<?xml version=\"1.0\" encoding=\"UTF-8\" ?><SessionLoginCap version=\"1.0\" xmlns=\"http://www.std-cgi.com/ver20/XMLSchema\"><sessionID>123</sessionID><challenge>123</challenge><iterations>100</iterations><isIrreversible>true</isIrreversible><salt>123</salt><isSessionIDValidLongTerm opt=\"true,false\">false</isSessionIDValidLongTerm><sessionIDVersion>2</sessionIDVersion></SessionLoginCap>",
			loginStatus:     http.StatusOK,
			sessionStatus:   http.StatusInternalServerError,
			sessionResponse: "mock-response",
			expectedError:   "received unexpected response from: /ISAPI/Security/sessionLogin status: 500 payload: mock-response",
		},
		{
			name:            "valid session response",
			loginResponse:   "<?xml version=\"1.0\" encoding=\"UTF-8\" ?><SessionLoginCap version=\"1.0\" xmlns=\"http://www.std-cgi.com/ver20/XMLSchema\"><sessionID>123</sessionID><challenge>123</challenge><iterations>100</iterations><isIrreversible>true</isIrreversible><salt>123</salt><isSessionIDValidLongTerm opt=\"true,false\">false</isSessionIDValidLongTerm><sessionIDVersion>2</sessionIDVersion></SessionLoginCap>",
			loginStatus:     http.StatusOK,
			sessionStatus:   http.StatusOK,
			sessionResponse: "mock-response",
			expectedError:   "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(testLoginHanlerHelper(tc.loginResponse, tc.sessionResponse, tc.loginStatus, tc.sessionStatus))
			defer ts.Close()
			cfg, err := NewConnector(ts.URL[7:], "mock-user", "mock-password", false)
			if tc.expectedError != "" {
				assert.ErrorContains(t, err, tc.expectedError)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, cfg)
			}
		})
	}
}

func TestConnector_newSesion(t *testing.T) {
	cases := []struct {
		name               string
		conn               Connector
		expectedError      string
		inputLoginResponse models.LoginResponse
	}{
		{
			name:               "invalid connector URL",
			conn:               Connector{User: "mock-User", Host: "mock host"},
			inputLoginResponse: mockLoginResponse(),
			expectedError:      "invalid character \" \" in host name",
		},
		{
			name:               "non existing host url",
			conn:               Connector{User: "mock-User", Host: "localhost:9999"},
			inputLoginResponse: mockLoginResponse(),
			expectedError:      "dial tcp 127.0.0.1:9999: connect: connection refused",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.conn.newSession(models.LoginResponse{})
			if tc.expectedError != "" {
				assert.ErrorContains(t, err, tc.expectedError)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
