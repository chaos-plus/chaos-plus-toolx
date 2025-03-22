package xpkg

import (
	"bytes"
	"os/exec"
	"strings"
)

func GetPkgPath(module string) string {
	v, _ := GetPkgPathE(module)
	return v
}

func GetPkgPathE(module string) (string, error) {
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", module)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}
