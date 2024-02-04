package ps

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

func NewExiftoolHandler() EventHandler {
	executeStr := "-execute"
	return EventHandler{
		OnBeforeExecute: func(stdin io.WriteCloser, cmd string, flags *[]*string) error {
			*flags = append(*flags, &executeStr)
			return nil
		},
		OnExecute: func(stdin io.WriteCloser, cmd string, flags *[]*string) error {
			return nil
		},
		OnBeforeStart: func(cmd string, flags []string) []string {
			flags = append([]string{"-stay_open", "True", "-@", "-", "-common_args"}, flags...)
			return flags
		},
		OnScanner: func(scanner *bufio.Scanner, cmd string) {
			scanner.Split(splitReadyToken)
		},
		OnStop: func(stdin io.WriteCloser, cmder *exec.Cmd, cmd string) {
			// write message telling it to close
			// but don't actually wait for the command to stop
			fmt.Fprintln(stdin, "-stay_open")
			fmt.Fprintln(stdin, "False")
			fmt.Fprintln(stdin, "-execute")
		},
	}
}

func splitReadyToken(data []byte, atEOF bool) (int, []byte, error) {
	delimPos := bytes.Index(data, []byte("{ready}\n"))
	delimSize := 8

	// maybe we are on Windows?
	if delimPos == -1 {
		delimPos = bytes.Index(data, []byte("{ready}\r\n"))
		delimSize = 9
	}

	if delimPos == -1 { // still no token found
		if atEOF {
			return 0, data, io.EOF
		} else {
			return 0, nil, nil
		}
	} else {
		if atEOF && len(data) == (delimPos+delimSize) { // nothing left to scan
			return delimPos + delimSize, data[:delimPos], bufio.ErrFinalToken
		} else {
			return delimPos + delimSize, data[:delimPos], nil
		}
	}
}
