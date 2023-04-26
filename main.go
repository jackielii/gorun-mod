package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage: gorun-mod <file.go>")
		os.Exit(1)
	}
	fp := try1(filepath.Abs(args[0]))
	ext := filepath.Ext(fp)
	base := filepath.Base(fp)
	module := base[:len(base)-len(ext)]

	td := try1(os.MkdirTemp("", "gorun-mod"))
	source := try1(os.ReadFile(fp))
	source = bytes.Replace(source, []byte("//go:build ignore\n"), nil, 1)
	try0(os.WriteFile(filepath.Join(td, base), source, 0644))

	sh(td, "go mod init %s", module)
	sh(td, "go mod tidy")

	for _, m := range []string{"go.mod", "go.sum"} {
		fp := filepath.Join(td, m)
		if _, err := os.Stat(fp); err != nil {
			continue
		}
		mod := try1(os.ReadFile(fp))
		fmt.Printf("// %s >>>\n", m)
		for _, line := range bytes.Split(mod, []byte("\n")) {
			if len(line) != 0 {
				fmt.Println("// " + string(line))
			}
		}
		fmt.Printf("// <<< %s\n", m)
	}

	os.RemoveAll(td)
}

func sh(cwd string, format string, args ...any) string {
	cmd := exec.Command("sh", "-c", fmt.Sprintf(format, args...))
	cmd.Dir = cwd
	out, err := cmd.Output()
	if err != nil {
		log.Printf("sh:\n"+format+"\n", args...)
		log.Fatal(err)
	}
	return strings.TrimSpace(string(out))
}

func try0(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func try1[T any](v T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return v
}
