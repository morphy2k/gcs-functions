package http

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var client *http.Client

var (
	postURL  *url.URL
	trimPath bool
)

func init() {
	var err error
	if v := os.Getenv("POST_URL"); v != "" {
		postURL, err = url.Parse(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while parsing env \"URL\": %v", err)
		}
	}
	if v := os.Getenv("STRIP_PATH"); v != "" {
		trimPath, err = strconv.ParseBool(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while parsing env \"STRIP_PATH\": %v", err)
		}
	}

	client = http.DefaultClient
}
