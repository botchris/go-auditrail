package auditrail

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type fileDescriptor struct {
	fd *os.File
}

// NewFileLogger builds a new logger that writes log entries to the given file
// descriptor.
//
// The file descriptor must be writable. If it is not, an error will be
// returned. You can use os.Stdout or os.Stderr as file descriptors.
//
// The file descriptor must be closed by the caller.
func NewFileLogger(fd *os.File) (Logger, error) {
	if fd == nil {
		return nil, fmt.Errorf("file descriptor was nil")
	}

	st, err := fd.Stat()
	if err != nil {
		return nil, fmt.Errorf("%w: could not stat file descriptor", err)
	}

	mode := st.Mode()
	if mode.IsDir() {
		return nil, fmt.Errorf("file descriptor is a directory")
	}

	// check fi file has permissions for os.O_WRONLY or os.O_RDWR or os.O_APPEND
	perm := int(mode.Perm())
	isWritable := perm&os.O_WRONLY == 0 || perm&os.O_RDWR == 0 || perm&os.O_APPEND == 0

	if !isWritable {
		return nil, fmt.Errorf("file descriptor is not writable")
	}

	return &fileDescriptor{fd: fd}, nil
}

func (dsc *fileDescriptor) Log(_ context.Context, entry *Entry) error {
	log, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	_, wErr := dsc.fd.WriteString(string(log) + "\n")

	return wErr
}

// NewFilePathLogger builds a new logger that writes log entries to a file at
// the given path.
//
// If the file does not exist, it will be created.
func NewFilePathLogger(path string) (Logger, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, cErr := os.Create(path)
		if cErr != nil {
			return nil, cErr
		}

		if sErr := file.Close(); sErr != nil {
			return nil, sErr
		}
	}

	fd, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return NewFileLogger(fd)
}
