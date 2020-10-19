package mockclient

import (
	"io/ioutil"
	"net/http"

	"github.com/jub0bs/namecheck"
)

type clientFunc func(url string) (*http.Response, error)

func (f clientFunc) Get(url string) (*http.Response, error) {
	return f(url)
}

func WithError(err error) namecheck.Client {
	get := func(_ string) (*http.Response, error) {
		return nil, err
	}
	return clientFunc(get)
}

func WithStatusCode(sc int) namecheck.Client {
	get := func(_ string) (*http.Response, error) {
		res := http.Response{
			StatusCode: sc,
			Body:       ioutil.NopCloser(nil),
		}
		return &res, nil
	}
	return clientFunc(get)
}
