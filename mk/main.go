package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"sync"
)

func makePath(path string, logErrors bool) {
	_, err := os.Stat(path)
	switch {
	case errors.Is(err, os.ErrNotExist):
		index := strings.LastIndex(path, "/")
		if index != -1 {
			err := os.MkdirAll(path[:index], fs.ModeDir)
			if logErrors && err != nil {
				fmt.Println(err)
			}
		}
		if index != len(path)-1 {
			_, err := os.Create(path)
			if logErrors && err != nil {
				fmt.Println(err)
			}
		}
	case errors.Is(err, os.ErrExist), err == nil:
	default:
		if logErrors {
			message, _ := strings.CutPrefix(err.Error(), "CreateFile ")
			fmt.Println(message)
		}
	}
}

func makePaths(paths []string, logErrors bool) {
	home := os.Getenv("HOME")
	if len(home) == 0 {
		user := os.Getenv("USERPROFILE")
		if len(user) != 0 {
			home = user
		}
	}
	home = strings.ReplaceAll(home, "\\", "/")
	os.Setenv("HOME", home)
	home += "/"

	var wg sync.WaitGroup
	for _, path := range paths {
		wg.Add(1)
		defer wg.Wait()
		go func(path string) {
			defer wg.Done()
			path = strings.ReplaceAll(path, "\\", "/")
			path = strings.TrimLeft(path, "/")
			path, isHomePath := strings.CutPrefix(path, "~/")
			if isHomePath {
				path = home + path
			}
			makePath(path, logErrors)
		}(path)
	}
}

func help() string {
	return `
  NAME
      mk - create a new file or directory

  SYNOPSIS
      mk file ...
      mk path/to/file ...
      mk directory/ ...
      mk path/to/directory/ ...

  DESCRIPTION
      mk creates a new file or directory if the file or directory does not
      exist. mk will create a new directory if the path provided ends with
      a forward slash (/); otherwise, mk will create a new file. mk will
      also create any intermediate directories needed.

  OPTIONS
      -e   Log errors to console.`
}

func main() {
	showErrors := flag.Bool("e", false, "show errors")
	flag.Parse()
	args := flag.Args()

	if len(args) > 0 {
		makePaths(args, *showErrors)
	} else {
		fmt.Println(help())
	}
}
