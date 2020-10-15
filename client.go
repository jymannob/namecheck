package namecheck

import "net/http"

type Client interface {
	Get(url string) (*http.Response, error)
}
