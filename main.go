package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	filename := flag.String("f", "", "log file name")
	maxSize := flag.Int("s", 10, "max log file size in MB")
	flag.Parse()
	if *filename == "" {
		fmt.Printf("filename is required\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// if file already exists, rename before overwiting
	_ = renameFile(*filename)

	f, err := os.Create(*filename)
	if err != nil {
		fmt.Printf("failed to open %q: %s\n", *filename, err)
		os.Exit(1)
	}

	remaining := *maxSize * 1024 * 1024
	scn := bufio.NewScanner(os.Stdin)
	for scn.Scan() {
		if len(scn.Text()) > remaining {
			f, err = replaceF(*filename, f)
			if err != nil {
				fmt.Printf("failed to cycle: %s\n", err)
				os.Exit(1)
			}
			remaining = *maxSize * 1024 * 1024
		}
		n, err := fmt.Fprintln(f, scn.Text())
		if err != nil {
			fmt.Printf("failed to write: %s\n", err)
			os.Exit(1)
		}
		remaining -= n
	}
	if scn.Err() != nil {
		fmt.Printf("scanner ended with %q\n", scn.Err())
	}
}

func replaceF(filename string, handle *os.File) (*os.File, error) {
	err := handle.Sync()
	if err != nil {
		fmt.Printf("failed to sync %q: %s\n", filename, err)
	}
	err = handle.Close()
	if err != nil {
		fmt.Printf("error closing %q: %s\n", filename, err)
	}
	renameFile(filename)
	return os.Create(filename)
}

func renameFile(filename string) error {
	stamp := time.Now().Format("2006-01-02_150405.00000")
	ext := filepath.Ext(filename)
	name := fmt.Sprintf("%s_%s%s", filename[:len(filename)-len(ext)], stamp, ext)
	return os.Rename(filename, name)
}
