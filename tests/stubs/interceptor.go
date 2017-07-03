package stubs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/pubnub/go/tests/helpers"
)

type Interceptor struct {
	Transport *interceptTransport
}

func NewInterceptor() *Interceptor {
	return &Interceptor{
		Transport: &interceptTransport{},
	}
}

func (i *Interceptor) AddStub(stub *Stub) {
	i.Transport.AddStub(stub)
}

func (i *Interceptor) GetClient() *http.Client {
	return &http.Client{
		Transport: i.Transport,
	}
}

type Stub struct {
	Method             string
	Path               string
	Query              string
	ResponseBody       string
	ResponseStatusCode int
	MixedPathPositions []int
	IgnoreQueryKeys    []string
	MixedQueryKeys     []string
}

func (s *Stub) Match(req *http.Request) bool {
	if s.Method != req.Method {
		log.Printf("Methods are not equal: %s != %s\n", s.Method, req.Method)
		return false
	}

	if !helpers.PathsEqual(s.Path, req.URL.EscapedPath(), s.MixedPathPositions) {
		return false
	}

	expectedQuery, _ := url.ParseQuery(s.Query)
	actualQuery := req.URL.Query()

	fmt.Println(expectedQuery, "\n", actualQuery)
	if !helpers.QueriesEqual(&expectedQuery,
		&actualQuery,
		s.IgnoreQueryKeys,
		s.MixedQueryKeys) {
		return false
	}

	return true
}

type interceptTransport struct {
	Stubs []*Stub
}

func (i *interceptTransport) RoundTrip(req *http.Request) (*http.Response,
	error) {

	for _, v := range i.Stubs {
		if v.Match(req) {
			var statusString string

			switch v.ResponseStatusCode {
			case 200:
				statusString = "200 OK"
			case 403:
				statusString = "403 Forbidden"
			default:
				statusString = ""
			}

			return &http.Response{
				Status:           statusString,
				StatusCode:       v.ResponseStatusCode,
				Proto:            "HTTP/1.0",
				ProtoMajor:       1,
				ProtoMinor:       0,
				Request:          req,
				Header:           http.Header{"Content-Length": {"256"}},
				TransferEncoding: nil,
				Close:            true,
				Body:             ioutil.NopCloser(bytes.NewBufferString(v.ResponseBody)),
				ContentLength:    256,
			}, nil
		}
	}

	// Nothing was found
	return &http.Response{
		Status:           "530 No stub matched",
		StatusCode:       530,
		Proto:            "HTTP/1.0",
		ProtoMajor:       1,
		ProtoMinor:       0,
		Request:          req,
		TransferEncoding: nil,
		Body:             ioutil.NopCloser(bytes.NewBufferString("No Stub Matched")),
		Close:            true,
		ContentLength:    256,
	}, nil
}

func (i *interceptTransport) AddStub(stub *Stub) {
	i.Stubs = append(i.Stubs, stub)
}