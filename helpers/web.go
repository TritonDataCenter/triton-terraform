package helpers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	local "github.com/joyent/gosdc/localservices/cloudapi"
	"github.com/joyent/gosign/auth"
	"github.com/julienschmidt/httprouter"
)

// Server creates a local test double for testing API responses
type Server struct {
	Server     *httptest.Server
	oldHandler http.Handler
	Mux        *httprouter.Router
	API        *local.CloudAPI
	Creds      *auth.Credentials
}

// NewServer returns a Server
func NewServer() (*Server, error) {
	s := new(Server)

	s.Server = httptest.NewServer(nil)
	s.oldHandler = s.Server.Config.Handler
	s.Mux = httprouter.New()
	s.Server.Config.Handler = s.Mux

	key, err := ioutil.ReadFile(TestKeyFile)
	if err != nil {
		return nil, err
	}

	authentication, err := auth.NewAuth(TestAccount, string(key), "rsa-sha256")

	s.Creds = &auth.Credentials{
		UserAuthentication: authentication,
		SdcKeyId:           TestKeyID,
		SdcEndpoint:        auth.Endpoint{URL: s.Server.URL},
	}

	s.API = local.New(s.Server.URL, TestAccount)
	s.API.SetupHTTP(s.Mux)

	return s, nil
}

// URL returns the URL of the server
func (s *Server) URL() string {
	return s.Server.URL
}

// Stop stops the server for teardown
func (s *Server) Stop() {
	s.Mux = nil
	s.Server.Config.Handler = s.oldHandler
	s.Server.Close()
}
