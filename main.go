// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fogo
// Date: 2021/5/14

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	// "sort"
	"sort"
	"strings"

	"github.com/fsgo/iconvs/convert"
)

const version = "0.3 2024-07-12"

var fromCode = flag.String("f", "", "from encoding")
var toCode = flag.String("t", "", "to encoding")
var write = flag.Bool("w", true, "write file when true, or print to stdout")
var list = flag.Bool("l", false, "list all available encodings")

func init() {
	flag.Usage = func() {
		cmd := os.Args[0]
		fmt.Fprintf(os.Stderr, "usage: %s [flags] [files ...]\n", cmd)
		flag.PrintDefaults()
		fmt.Fprint(os.Stderr, "\nsite :    https://github.com/fsgo/iconvs\n")
		fmt.Fprintf(os.Stderr, "version:  %s\n", version)
		os.Exit(2)
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	if *list {
		doList()
		return
	}
	if *fromCode == "" {
		log.Fatal("-f flag is required")
	}

	if *toCode == "" {
		log.Fatal("-t flag is required")
	}

	if *fromCode == *toCode {
		log.Fatal("-f and -t flags cannot be combined")
		return
	}
	if len(os.Args) == 0 {
		log.Fatal("no files")
	}
	for _, fname := range flag.Args() {
		convertGlob(fname)
	}

	if convertFail > 0 {
		log.Fatalf("%d files convert failed", convertFail)
	}
}

func doList() {
	names := make([]string, 0, len(convert.Encodings))
	for name := range convert.Encodings {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		return names[i] < names[j]
	})
	fmt.Println(strings.Join(names, "\n"))
}

func convertGlob(fname string) {
	ms, err := filepath.Glob(fname)
	if err != nil {
		log.Fatalf("Glob(%q) failed: %v", fname, err)
	}
	for _, sf := range ms {
		_ = doConvert(sf)
	}
}

var errIgnore = errors.New("ignore")
var convertFail int

func doConvert(fName string) (ret error) {
	defer func() {
		if errors.Is(ret, errIgnore) {
			return
		}

		if ret == nil {
			log.Printf("%s success\n", fName)
		} else {
			convertFail++
			log.Printf("%s failed, %s\n", fName, ret.Error())
		}
	}()

	f, err := os.Open(fName)
	if err != nil {
		return fmt.Errorf(" Open failed, %w", err)
	}
	defer f.Close()
	fs, err := f.Stat()
	if err != nil {
		return err
	}

	if fs.IsDir() {
		convertGlob(filepath.Join(fName, "*"))
		return errIgnore
	}

	bf, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	out, err := convert.Convert(bf, *fromCode, *toCode)
	if err != nil {
		return err
	}
	if bytes.Equal(bf, out) {
		return nil
	}
	if !*write {
		fmt.Print(string(out))
		return errIgnore
	}
	return os.WriteFile(fName, out, fs.Mode())
}
