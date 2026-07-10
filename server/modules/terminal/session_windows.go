//go:build windows

package terminal

import (
	"io"
	"os/exec"
	"sync"
)

type pipeSession struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	reader *io.PipeReader
	writer *io.PipeWriter
	once   sync.Once
}

func newSession(shell string) (session, error) {
	cmd := exec.Command(shell)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	reader, writer := io.Pipe()
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go func() { _, _ = io.Copy(writer, stdout) }()
	go func() { _, _ = io.Copy(writer, stderr) }()

	return &pipeSession{cmd: cmd, stdin: stdin, reader: reader, writer: writer}, nil
}

func (s *pipeSession) Read(p []byte) (int, error) {
	return s.reader.Read(p)
}

func (s *pipeSession) Write(p []byte) (int, error) {
	return s.stdin.Write(p)
}

func (s *pipeSession) Close() error {
	s.once.Do(func() {
		_ = s.stdin.Close()
		_ = s.reader.Close()
		_ = s.writer.Close()
		if s.cmd.Process != nil {
			_ = s.cmd.Process.Kill()
		}
	})
	return nil
}

func (s *pipeSession) Wait() error {
	return s.cmd.Wait()
}
