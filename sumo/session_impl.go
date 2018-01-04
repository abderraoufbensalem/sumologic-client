package sumo

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type SessionImpl struct {
	address   string
	accessID  string
	accessKey string
}

func DefaultSession() Session {
	s := &SessionImpl{}
	s.SetAddress(ConfigKeys.SumoAddress)
	s.SetCredentials(ConfigKeys.SumoAccessId, ConfigKeys.SumoAccessKey)
	return s
}

func (s *SessionImpl) Discover() {
	client := &http.Client{
		Transport: s.CreateTransport(),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	executor := NewClientExecutor(s, client)

	req, _ := executor.NewRequest()
	req.SetEndpoint("/")
	res, _ := req.Get()

	if res != nil {
		location := res.Header("Location")
		if location != "" {
			s.SetAddress(strings.TrimRight(location, "/"))
		}
	}
}

func (s *SessionImpl) SetAddress(address string) {
	s.address = address
}

func (s *SessionImpl) SetCredentials(accessID, accessKey string) {
	s.accessID = accessID
	s.accessKey = accessKey
}

func (s *SessionImpl) Address() string {
	return s.address
}

func (s *SessionImpl) EndpointURL(endpoint string) *url.URL {
	uri, _ := url.Parse(fmt.Sprintf("%s/%s", strings.TrimRight(s.address, "/"), strings.TrimLeft(endpoint, "/")))
	return uri
}

func (s *SessionImpl) CreateTransport() http.RoundTripper {
	return NewAnonymousTransport(func(req *http.Request) (*http.Response, error) {
		req.SetBasicAuth(s.accessID, s.accessKey)
		return http.DefaultTransport.RoundTrip(req)
	})
}
