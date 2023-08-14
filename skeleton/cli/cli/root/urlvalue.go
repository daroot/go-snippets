package root

import "net/url"

type urlValue struct {
	URL *url.URL
}

func (v urlValue) String() string {
	if v.URL != nil {
		return v.URL.String()
	}
	return ""
}

func (v urlValue) Set(s string) error {
	if u, err := url.Parse(s); err != nil {
		return err
	} else {
		*v.URL = *u
	}
	return nil
}
