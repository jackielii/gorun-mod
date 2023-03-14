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

	td := try1(os.MkdirTemp("", "gorun-mod"))
	sh(td, "cp %s %s", fp, td)

	ext := filepath.Ext(fp)
	base := filepath.Base(fp)
	fn := base[:len(base)-len(ext)]

	sh(td, "go mod init %s", fn)
	sh(td, "go mod tidy")

	for _, m := range []string{"go.mod", "go.sum"} {
		fp := filepath.Join(td, m)
		if _, err := os.Stat(fp); err != nil {
			continue
		}
		content := try1(os.ReadFile(fp))
		fmt.Printf("// %s >>>\n", m)
		for _, line := range bytes.Split(content, []byte("\n")) {
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

func try1[T any](v T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return v
}
