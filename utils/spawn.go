package utils

import (
	"bytes"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func ResolveCmd(cmd string) string {
	paths := make([]string, 0)
	if pwd, err := os.Getwd(); err == nil {
		paths = append(paths, pwd)
	}
	if exePath, err := os.Executable(); err == nil {
		paths = append(paths, filepath.Dir(exePath))
	}
	paths = append(paths, strings.Split(os.Getenv("PATH"), ":")...)
	cmds := make([]string, 0)
	cmds = append(cmds, cmd)
	if runtime.GOOS == "windows" {
		cmds = append(cmds, cmd+".cmd") // npm
	}
	for _, cmd := range cmds {
		for _, v := range paths {
			curr := path.Join(v, cmd)
			if info, err := os.Stat(curr); err == nil {
				if info.IsDir() {
					continue
				}
				return curr
			}
		}
	}
	return cmd
}

type SpawnOptions struct {
	// 使用程序std输入输入
	SystemStdout bool
}

type spawnOutput struct {
	p      *exec.Cmd
	writer *bytes.Buffer
}

func (s spawnOutput) String() string {
	if s.writer != nil {
		return s.writer.String()
	}
	return ""
}

func (s spawnOutput) Bytes() []byte {
	if s.writer != nil {
		return s.writer.Bytes()
	}
	return []byte{}
}

func SpawnSimple(cmd string, args []string) (spawnOutput, error) {
	return Spawn(cmd, args, nil)
}

func Spawn(cmd string, args []string, opt *SpawnOptions) (spawnOutput, error) {
	cmd = ResolveCmd(cmd)
	p := exec.Command(cmd, args...)
	writer := bytes.NewBufferString("")
	r := spawnOutput{
		p:      p,
		writer: writer,
	}
	if opt != nil && opt.SystemStdout {
		p.Stderr = os.Stderr
		p.Stdout = os.Stdout
	} else {
		p.Stdout = writer
		p.Stderr = writer
	}
	err := p.Start()
	if err != nil {
		return r, err
	}
	err = p.Wait()
	if err != nil {
		return r, err
	}
	return r, nil
}
