package xcmd

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

type Runner struct {
	WorkDir string
	Cmd     string
	Args    []string
	Env     []string
	Stdout  io.Writer
	Stderr  io.Writer
}

func New() *Runner {
	return &Runner{
		WorkDir: "",
		Cmd:     "",
		Args:    []string{},
		Env:     os.Environ(),
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}
}

func (r *Runner) SetWorkDir(workDir string) *Runner {
	r.WorkDir = workDir
	return r
}

func (r *Runner) SetCmd(cmd string) *Runner {
	r.Cmd = cmd
	return r
}

func (r *Runner) SetArgs(args []string) *Runner {
	r.Args = args
	return r
}

func (r *Runner) SetEnv(env []string) *Runner {
	r.Env = append(os.Environ(), env...)
	return r
}

func (r *Runner) AddEnv(env ...string) *Runner {
	r.Env = append(r.Env, env...)
	return r
}

func (r *Runner) SetStdout(stdout io.Writer) *Runner {
	r.Stdout = stdout
	return r
}

func (r *Runner) SetStderr(stderr io.Writer) *Runner {
	r.Stderr = stderr
	return r
}

func (r *Runner) Run() error {
	cmd := exec.Command(r.Cmd, r.Args...)
	cmd.Dir = r.WorkDir
	cmd.Env = r.Env
	cmd.Stdout = r.Stdout
	cmd.Stderr = r.Stderr
	return cmd.Run()
}

func Run(cmd string, args ...string) error {
	return New().SetCmd(cmd).SetArgs(args).Run()
}

func RunWithResult(cmd string, args ...string) (string, error) {
	buffer := bytes.Buffer{}
	err := New().
		SetCmd(cmd).
		SetArgs(args).
		SetStdout(&buffer).
		SetStderr(&buffer).
		Run()
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}
