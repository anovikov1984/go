package pubnub

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pubnub/go/pnerr"
)

const TIME_PATH = "/time/0"

var emptyTimeResp *TimeResponse

type timeBuilder struct {
	opts *timeOpts
}

func newTimeBuilder(pubnub *PubNub) *timeBuilder {
	builder := timeBuilder{
		opts: &timeOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newTimeBuilderWithContext(pubnub *PubNub, context Context) *timeBuilder {
	builder := timeBuilder{
		opts: &timeOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *timeBuilder) Transport(tr http.RoundTripper) *timeBuilder {
	b.opts.Transport = tr
	return b
}

func (b *timeBuilder) Execute() (*TimeResponse, StatusResponse, error) {
	rawJson, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyTimeResp, status, err
	}

	return newTimeResponse(rawJson, status)
}

type timeOpts struct {
	pubnub *PubNub

	Transport http.RoundTripper

	ctx Context
}

func (o *timeOpts) config() Config {
	return *o.pubnub.Config
}

func (o *timeOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *timeOpts) context() Context {
	return o.ctx
}

func (o *timeOpts) validate() error {
	return nil
}

func (o *timeOpts) buildPath() (string, error) {
	return TIME_PATH, nil
}

func (o *timeOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.Uuid)

	return q, nil
}

func (o *timeOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *timeOpts) httpMethod() string {
	return "GET"
}

func (o *timeOpts) isAuthRequired() bool {
	return false
}

func (o *timeOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *timeOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *timeOpts) operationType() OperationType {
	return PNTimeOperation
}

type TimeResponse struct {
	Timetoken int64
}

func newTimeResponse(jsonBytes []byte, status StatusResponse) (*TimeResponse, StatusResponse, error) {
	resp := &TimeResponse{}

	var value interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyTimeResp, status, e
	}

	if parsedValue, ok := value.([]interface{}); ok {
		if tt, ok := parsedValue[0].(float64); ok {
			resp.Timetoken = int64(tt)
		}
	}

	return resp, status, nil
}
