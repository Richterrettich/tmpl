package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	rawArgs := os.Args[1:]

	if len(rawArgs) == 0 {
		fatal("tmpl needs at least one argument")
	}

	flags, args := parseFlags(rawArgs)

	if len(args) == 0 {
		fatal("tmpl needs at least one argument")
	}

	mainFile, args := args[0], args[1:]

	templateFiles := []string{mainFile}
	for _, arg := range args {
		if isDirectory(arg) {
			filepath.Walk(arg, func(path string, info os.FileInfo, err error) error {
				handleError(err)
				if !isDirectory(path) && strings.HasSuffix(path, ".tmpl") {
					templateFiles = append(templateFiles, path)
				}
				return nil
			})
		} else if strings.HasSuffix(arg, ".tmpl") {
			templateFiles = append(templateFiles, arg)
		}
	}

	mainTemplate := template.New("main")
	for _, templateFile := range templateFiles {
		content, err := ioutil.ReadFile(templateFile)
		handleError(err)
		mainTemplate, err = mainTemplate.Parse(fmt.Sprintf("{{ define \"%s\" }}%s\n{{ end }}", templateFile, content))
		handleError(err)
	}
	handleError(mainTemplate.ExecuteTemplate(os.Stdout, mainFile, flags))
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	handleError(err)
	return fileInfo.IsDir()
}

func handleError(err error) {
	if err != nil {
		fatal(err.Error())
	}
}

func fatal(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func parseFlags(rawArgs []string) (map[string]string, []string) {
	cleanArgs := make([]string, 0)
	flags := make(map[string]string)
	i := 0
	for {
		if i >= len(rawArgs) {
			return flags, cleanArgs
		}
		arg := rawArgs[i]
		if strings.HasPrefix(arg, "-") {
			if len(rawArgs)-1 < i+1 {
				fatal("no value provided for flag %s", arg)
			}
			val := rawArgs[i+1]
			flags[arg[1:len(arg)]] = val
			i = i + 2
		} else {
			cleanArgs = append(cleanArgs, arg)
			i = i + 1
		}
	}
}
