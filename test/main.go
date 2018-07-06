package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if !strings.HasSuffix(path, ".txt") {
			return nil
		}
		fmt.Println(path)
		return nil
	})
}
