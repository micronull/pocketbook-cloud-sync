package donwload

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func Download(url, destination string) error {
	rsp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http GET %s: %w", url, err)
	}

	defer func() { _ = rsp.Body.Close() }()

	if rsp.StatusCode != http.StatusOK {
		return httpStatusError{rsp.StatusCode}
	}

	file, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("create file %s: %w", destination, err)
	}

	defer func() { _ = file.Close() }()

	_, err = io.Copy(file, rsp.Body)
	if err != nil {
		return fmt.Errorf("copy downloaded data to file %s: %w", destination, err)
	}

	return nil
}

type httpStatusError struct {
	code int
}

func (e httpStatusError) Error() string {
	return fmt.Sprintf("http status code: %d %s", e.code, http.StatusText(e.code))
}

func (e httpStatusError) Code() int {
	return e.code
}
