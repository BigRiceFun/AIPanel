package terminal

import "io"

type session interface {
	io.Reader
	io.Writer
	Close() error
	Wait() error
}
