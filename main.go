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

	f, err := os.Create(*filename)
	if err != nil {
		fmt.Printf("failed to open %q: %s\n", err)
		os.Exit(1)
	}

	remaining := *maxSize * 1024 * 1024
	scn := bufio.NewScanner(os.Stdin)
	for scn.Scan() {
		if len(scn.Text()) > remaining {
			f, err = replaceF(*filename, f)
			if err != nil {
				fmt.Printf("failed to cycle to %q: %s\n", err)
				os.Exit(1)
			}
			remaining = *maxSize * 1024 * 1024
		}
		n, err := fmt.Fprintln(f, scn.Text())
		if err != nil {
			fmt.Printf("failed to write to %q: %s\n", err)
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
		fmt.Printf("failed to sync %q: %s\n", err)
	}
	err = handle.Close()
	if err != nil {
		fmt.Printf("error closing %q: %s\n", err)
	}
	stamp := time.Now().Format("2006-01-02_150405")
	ext := filepath.Ext(filename)
	name := fmt.Sprintf("%s_%s%s", filename[:len(filename)-len(ext)], stamp, ext)
	err = os.Rename(filename, name)
	if err != nil {
		fmt.Printf("error renaming %q: %s\n", err)
	}
	return os.Create(filename)
}
