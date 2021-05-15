// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fogo
// Date: 2021/5/14

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fsgo/iconvs/convert"
)

const version = "0.1 2021-05-15"

var fromCode = flag.String("f", "", "from encoding")
var toCode = flag.String("t", "", "to encoding")
var write = flag.Bool("w", true, "write file")
var list = flag.Bool("l", false, "list all available encodings")

func init() {
	flag.Usage = func() {
		cmd := os.Args[0]
		fmt.Fprintf(os.Stderr, "usage: %s [flags] [files ...]\n", cmd)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nsite :    https://github.com/fsgo/iconvs\n")
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

	if *fromCode == *toCode {
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
		if len(names[i]) == len(names[j]) {
			return names[i] < names[j]
		}
		return len(names[i]) < len(names[j])
	})
	fmt.Println(strings.Join(names, "\n"))
}

func convertGlob(fname string) {
	ms, err := filepath.Glob(fname)
	if err != nil {
		log.Fatalf("Glob(%q) failed: %v", fname, err)
	}
	for _, sf := range ms {
		doConvert(sf)
	}
}

var errIgnore = fmt.Errorf("ignore")
var convertFail int

func doConvert(fName string) (ret error) {
	defer func() {
		if ret == errIgnore {
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

	bf, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	out, err := convert.Convert(*fromCode, *toCode, bf)
	if err != nil {
		return err
	}
	if !*write {
		fmt.Print(string(out))
		return errIgnore
	}
	return os.WriteFile(fName, out, fs.Mode())
}
