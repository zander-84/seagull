package endpoint

import (
	"net/url"
)

// NewEndpoint new an Endpoint URL.
func NewEndpoint(scheme, host string) *url.URL {
	return &url.URL{Scheme: scheme, Host: host}
}

// ParseEndpoint parses an Endpoint URL.
func ParseEndpoint(endpoints []string, scheme string) (string, error) {
	for _, e := range endpoints {
		u, err := url.Parse(e)
		if err != nil {
			return "", err
		}

		// TODO: Compatibility processing
		// Function is to convert grpc:/127.0.0.1/?isSecure=true into grpcs:/127.0.0.1/
		// It will be deleted in about a month
		u = legacyURLToNew(u)
		if u.Scheme == scheme {
			return u.Host, nil
		}
	}
	return "", nil
}

func legacyURLToNew(u *url.URL) *url.URL {
	if u.Scheme == "https" || u.Scheme == "grpcs" {
		return u
	}
	return u
}

// Scheme is the scheme of endpoint url.
// examples: scheme="http",isSecure=true get "https"
func Scheme(scheme string, isSecure bool) string {
	if isSecure {
		return scheme + "s"
	}
	return scheme
}
