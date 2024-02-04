package ps

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type OnBeforeExecute func(stdin io.WriteCloser, cmd string, flags *[]*string) error
type OnExecute func(stdin io.WriteCloser, cmd string, flags *[]*string) error
type OnBeforeStart func(cmd string, flags []string) []string
type OnScanner func(scanner *bufio.Scanner, cmd string)
type OnStop func(stdin io.WriteCloser, cmder *exec.Cmd, cmd string)

type EventHandler struct {
	OnBeforeExecute OnBeforeExecute
	OnExecute       OnExecute
	OnBeforeStart   OnBeforeStart
	OnStop          OnStop
	OnScanner       OnScanner
	StopTimeoutMs   int64
}

// Stayopen abstracts running cmd with `-stay_open` to greatly improve
// performance. Remember to call Stayopen.Stop() to signal cmd to shutdown
// to avoid zombie perl processes
type Stayopen struct {
	l       sync.Mutex
	cmd     *exec.Cmd
	command string
	flags   []string
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	scanner *bufio.Scanner
	event   EventHandler
}

func (s *Stayopen) execute(flags ...string) ([]byte, error) {
	s.l.Lock()
	defer s.l.Unlock()
	if s.cmd == nil {
		return nil, errors.New("Stopped")
	}
	args := []*string{}
	for _, v := range flags {
		args = append(args, &v)
	}
	err := s.event.OnBeforeExecute(s.stdin, s.command, &args)
	if err != nil {
		return nil, errors.Wrap(err, "Failed onBeforeExecute")
	}
	for _, f := range args {
		if f != nil {
			fmt.Fprintln(s.stdin, *f)
		}
	}
	err = s.event.OnExecute(s.stdin, s.command, &args)
	if err != nil {
		return nil, errors.Wrap(err, "Failed onExecute")
	}
	if !s.scanner.Scan() {
		return nil, errors.New("Failed to read output")
	} else {
		results := s.scanner.Bytes()
		sendResults := make([]byte, len(results))
		copy(sendResults, results)
		return sendResults, nil
	}
}

func (s *Stayopen) Stop() {
	s.l.Lock()
	defer s.l.Unlock()
	cmd := s.cmd
	ps := cmd.ProcessState
	s.cmd = nil
	s.event.OnStop(s.stdin, cmd, s.command)
	ms := s.event.StopTimeoutMs
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ms)*time.Millisecond) // 设置5秒超时
	defer cancel()
	<-ctx.Done()
	if ps != nil && !ps.Exited() {
		cmd.Cancel()
	} else if ps == nil && cmd.Process != nil {
		cmd.Process.Kill()
	}
}

func (s *Stayopen) start() error {
	s.l.Lock()
	defer s.l.Unlock()
	flags := s.event.OnBeforeStart(s.command, s.flags)
	s.cmd = exec.Command(s.command, flags...)
	stdin, err := s.cmd.StdinPipe()
	if err != nil {
		return errors.Wrap(err, "Failed getting stdin pipe")
	}
	stdout, err := s.cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "Failed getting stdout pipe")
	}
	s.stdin = stdin
	s.stdout = stdout
	s.scanner = bufio.NewScanner(stdout)
	s.event.OnScanner(s.scanner, s.command)
	if err := s.cmd.Start(); err != nil {
		return errors.Wrap(err, "Failed starting cmd in stay_open mode")
	}
	return nil
}

func NewStayOpen(h EventHandler, cmd string, flags ...string) (*Stayopen, error) {
	stayopen := &Stayopen{
		command: cmd,
		flags:   flags,
		event:   h,
	}
	return stayopen, stayopen.start()
}
