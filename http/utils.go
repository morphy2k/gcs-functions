package http

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func bodyToStderr(r io.ReadCloser) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintln(os.Stderr, string(data)); err != nil {
		return err
	}

	return nil
}
