package must

import (
	"bytes"
	"io"
	"log"
)

// Read is a helper function that asserts that an error is nil and returns the reader.
func Read(reader io.ReadCloser, err error) io.Reader {
	if err != nil {
		reader.Close()
		log.Fatalf("%v: unexpected error", err)
	}

	buffer := &bytes.Buffer{}
	_, err = io.Copy(buffer, reader)
	reader.Close()
	if err != nil {
		log.Fatalf("%v: failed to read from reader", err)
	}

	return buffer
}
