package main

import (
	"bytes"
	"flag"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

var header string = `
/**
 * This file generated automatically.  Do not modify.
 * Generated from workspace: {{.workspace}}
 */

package {{.package}}

import "github.com/BellerophonMobile/logberry"

var buildmetadata = &logberry.BuildMetadata{
  Host:     "{{.host}}",
  User:     "{{.user}}",
  Date:     "{{.date}}",

  Repositories: []logberry.RepositoryMetadata {
`

var repo string = `
    logberry.RepositoryMetadata{
      Repository: "{{.root}}",
      Branch:     "{{.branch}}",
      Commit:     "{{.commit}}",
      Dirty:      {{.dirty}},
      Path:       "{{.path}}",
    },
`

var footer string = `
  },
}
`

func main() {

	var workspace string
	flag.StringVar(&workspace, "workspace", ".", "Directory to scan.")

	var pkg string
	flag.StringVar(&pkg, "pkg", "main", "Package in which to identify code.")

	var target string
	flag.StringVar(&target, "target", "/dev/stdout", "File in which to generate signature.")

	var gofile string
	flag.StringVar(&gofile, "out", "", "If target is a Go file, root name using this parameter.  Suffix '.go' will be appended and it set as target.  This is a workaround for the go run command interpreting arguments ending in .go as files to run.")

	flag.Parse()

	if gofile != "" {
		target = gofile + ".go"
	}

	build := map[string]string{
		"workspace": workspace,
		"package":   pkg,
		"host":      exe("hostname", ""),
		"user":      exe("whoami", ""),
		"date":      time.Now().Format(time.RFC3339),
	}

	h := template.Must(template.New("header").Parse(header))
	r := template.Must(template.New("repo").Parse(repo))
	f := template.Must(template.New("footer").Parse(footer))

	var outb bytes.Buffer

	h.Execute(&outb, build)

	err := filepath.Walk(workspace,
		func(path string, info os.FileInfo, err error) error {

			if filepath.Base(path) == ".git" {

				dir := filepath.Dir(path)
				repo := map[string]interface{}{
					"root":   filepath.Base(exe("git rev-parse --show-toplevel", dir)),
					"branch": exe("git rev-parse --abbrev-ref HEAD", dir),
					"commit": exe("git rev-parse HEAD", dir),
					"dirty":  false,
					"path":   dir,
				}

				if strings.Contains(exe("git status -uno", dir), "modified") {
					repo["dirty"] = true
				}

				r.Execute(&outb, repo)

				return filepath.SkipDir
			}

			return nil
		})
	if err != nil {
		log.Panic(err)
	}

	f.Execute(&outb, build)

	gofmt, err := format.Source(outb.Bytes())
	if err != nil {
		log.Panic(err)
	}

	if err := ioutil.WriteFile(target, gofmt, 0644); err != nil {
		log.Panic(err)
	}
}

func exe(cmd string, wd string) string {

	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:]

	command := exec.Command(head, parts...)
	command.Dir = wd

	out, err := command.Output()
	if err != nil {
		log.Panic("Failed command", cmd, err)
	}

	return strings.TrimSpace(string(out))

}
