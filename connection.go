package annkesdk

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/csrar/annkeSDK/models"
)

type Connector struct {
	Host     string
	User     string
	Password string
	Secure   bool
	client   http.Client
}

func NewConnector(host, user, password string, secure bool) (*Connector, error) {
	cfg := Connector{
		Host:     host,
		User:     user,
		Password: password,
		Secure:   secure,
	}
	if err := cfg.validateParameters(); err != nil {
		return nil, err
	}

	jar, _ := cookiejar.New(nil)
	cfg.client = http.Client{Timeout: time.Duration(timeout) * time.Second, Jar: jar}
	loginResponse, err := cfg.login()

	if err != nil {
		return nil, err
	}
	return &cfg, cfg.newSession(loginResponse)
}
func (cfg Connector) validateParameters() error {
	if cfg.Host == "" {
		return NewAnnkeInitError("Host")
	}
	if cfg.User == "" {
		return NewAnnkeInitError("User")
	}
	return nil
}

func (c Connector) login() (models.LoginResponse, error) {
	loginResponse := models.LoginResponse{}
	req, err := c.prepareLoginRequest()
	if err != nil {
		return loginResponse, fmt.Errorf("Error preparing login request %w", err)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return loginResponse, fmt.Errorf("Error executing login request %w", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return loginResponse, fmt.Errorf("Error decoding login response %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return loginResponse, NewAnnkeError(resp.StatusCode, string(body), loginPath)
	}
	err = xml.Unmarshal(body, &loginResponse)
	if err != nil {
		return loginResponse, fmt.Errorf("Error unmarshaling login response %w", err)
	}
	return loginResponse, nil
}

func (c Connector) newSession(login models.LoginResponse) error {
	hashedPassword := c.hashPassword(login)
	req, err := c.prepareSessionRequest(login, hashedPassword)
	if err != nil {
		return fmt.Errorf("Error preparing session request %w", err)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("Error executing session request %w", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return NewAnnkeError(resp.StatusCode, string(body), sessionPath)
	}
	return nil
}

func (c Connector) hashPassword(login models.LoginResponse) string {
	var saltByte [32]byte
	if login.IsIrreversible {
		saltByte = sha256.Sum256([]byte(c.User + login.Salt + c.Password))
	} else {
		saltByte = sha256.Sum256([]byte(c.Password))
	}
	challenge := hex.EncodeToString(saltByte[:]) + login.Challenge
	saltByte = sha256.Sum256([]byte(challenge))

	for i := 2; login.Iterations > i; i++ {
		saltByte = sha256.Sum256([]byte(hex.EncodeToString(saltByte[:])))
	}
	return hex.EncodeToString(saltByte[:])
}

func (c Connector) prepareSessionRequest(login models.LoginResponse, hashedPassword string) (*http.Request, error) {
	session := models.Session{
		UserName:                 c.User,
		Password:                 hashedPassword,
		SessionID:                login.SessionID,
		IsSessionIDValidLongTerm: login.IsSessionIDValidLongTerm.Text,
		SessionIDVersion:         login.SessionIDVersion,
	}

	out, err := xml.Marshal(&session)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s://%s%s?timeStamp=%d", c.getProtocol(), c.Host, sessionPath, time.Now().Unix())
	req, err := http.NewRequest("GET", url, strings.NewReader(string(out)))
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (c Connector) prepareLoginRequest() (*http.Request, error) {
	url := fmt.Sprintf("%s://%s:%s@%s%s?username=%s&random=%d", c.getProtocol(), c.User, c.Password, c.Host, loginPath, c.User, generateRandom())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func generateRandom() int {
	return rand.Intn(randomLenght)
}

func (c Connector) getProtocol() string {
	protocol := "http"
	if c.Secure {
		protocol = "https"
	}
	return protocol
}
