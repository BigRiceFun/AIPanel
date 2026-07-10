//go:build !windows

package terminal

import (
	"os"
	"os/exec"

	"github.com/creack/pty"
)

type ptySession struct {
	file *os.File
	cmd  *exec.Cmd
}

func newSession(shell string) (session, error) {
	cmd := exec.Command(shell)
	file, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}
	return &ptySession{file: file, cmd: cmd}, nil
}

func (s *ptySession) Read(p []byte) (int, error) {
	return s.file.Read(p)
}

func (s *ptySession) Write(p []byte) (int, error) {
	return s.file.Write(p)
}

func (s *ptySession) Close() error {
	return s.file.Close()
}

func (s *ptySession) Wait() error {
	return s.cmd.Wait()
}
