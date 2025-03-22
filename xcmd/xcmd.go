package xcmd

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

func Execute(dir, name string, args ...string) error {
	return ExecuteWithWriter(os.Stdout, dir, name, args...)
}

func ExecuteWithResult(dir, name string, args ...string) (string, error) {
	writer := bytes.Buffer{}
	ExecuteWithWriter(&writer, dir, name, args...)
	return writer.String(), nil
}

func ExecuteWithWriter(writer io.Writer, dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir // Set the working directory
	cmd.Stdout = writer
	cmd.Stderr = writer
	err := cmd.Run()
	return err
}
