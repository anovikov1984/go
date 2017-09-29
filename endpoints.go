package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/pubnub/go/utils"
)

type endpointOpts interface {
	config() Config
	client() *http.Client
	context() Context
	validate() error

	buildPath() (string, error)
	buildQuery() (*url.Values, error)
	// or bytes[]?
	buildBody() ([]byte, error)

	httpMethod() string
	operationType() OperationType
}

func defaultQuery(uuid string) *url.Values {
	v := &url.Values{}

	v.Set("pnsdk", "PubNub-Go/"+Version)
	v.Set("uuid", uuid)

	return v
}

func buildUrl(o endpointOpts) (*url.URL, error) {
	var stringifiedQuery string
	var signature string

	path, err := o.buildPath()
	if err != nil {
		return &url.URL{}, err
	}

	query, err := o.buildQuery()
	if err != nil {
		return &url.URL{}, err
	}

	if o.config().FilterExpression != "" {
		query.Set("filter-expr", o.config().FilterExpression)
	}

	if o.config().SecretKey != "" {
		timestamp := time.Now().Unix()
		query.Set("timestamp", strconv.Itoa(int(timestamp)))

		signedInput := o.config().SubscribeKey + "\n" + o.config().PublishKey + "\n"

		if o.operationType() == PNAccessManagerGrant ||
			o.operationType() == PNAccessManagerRevoke {
			signedInput += "grant\n"
		} else {
			signedInput += fmt.Sprintf("%s\n", path)
		}

		signedInput += utils.PreparePamParams(query)
		fmt.Println("signed input: ", signedInput)

		signature = utils.GetHmacSha256(o.config().SecretKey, signedInput)
	}

	if o.operationType() == PNPublishOperation {
		q, _ := o.buildQuery()
		v := q.Get("meta")
		if v != "" {
			query.Set("meta", v)
		}
	}

	if o.operationType() == PNSetStateOperation {
		q, _ := o.buildQuery()
		v := q.Get("state")
		query.Set("state", v)
	}

	if v := query.Get("uuid"); v != "" {
		query.Set("uuid", v)
	}

	if v := query.Get("auth"); v != "" {
		query.Set("auth", v)
	}

	stringifiedQuery = utils.PreparePamParams(query)

	if signature != "" {
		stringifiedQuery += fmt.Sprintf("&signature=%s", signature)
	}

	path = fmt.Sprintf("//%s%s", o.config().Origin, path)

	retUrl := &url.URL{
		Opaque:   path,
		Scheme:   "https",
		Host:     o.config().Origin,
		RawQuery: stringifiedQuery,
	}

	return retUrl, nil
}
