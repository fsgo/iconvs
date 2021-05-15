// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/5/14

package convert

import (
	"fmt"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
)

type Charset string

// Encodings 支持的所有的编码
var Encodings = map[string]encoding.Encoding{
	"UTF-8": unicode.UTF8,

	"UTF-16":   unicode.UTF16(unicode.BigEndian, unicode.UseBOM),
	"UTF-16BE": unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM),
	"UTF-16LE": unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM),

	"UTF-32":   utf32.UTF32(utf32.BigEndian, utf32.UseBOM),
	"UTF-32BE": utf32.UTF32(utf32.BigEndian, utf32.IgnoreBOM),
	"UTF-32LE": utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM),
}

func loadAll(list []encoding.Encoding) {
	for _, cm := range list {
		if ss, ok := cm.(fmt.Stringer); ok {
			name := strings.ReplaceAll(ss.String(), " ", "")
			Encodings[name] = cm
		}
	}
}

func init() {
	loadAll(charmap.All)
	loadAll(simplifiedchinese.All)
	loadAll(traditionalchinese.All)
	loadAll(korean.All)
	loadAll(japanese.All)

	// build alias
	for name := range Encodings {
		nameNew := strings.ReplaceAll(name, "-", "")
		nameNew = strings.ToUpper(nameNew)
		if name != nameNew {
			alias[nameNew] = name
		}
	}
}

var alias = map[string]string{
	"HZ-GB2312": "GB2312",
}

func charsetName(name string) string {
	name = strings.ToUpper(name)
	a, has := alias[name]
	if !has {
		return name
	}
	return a
}

func Convert(from string, to string, input []byte) ([]byte, error) {
	from = charsetName(from)
	to = charsetName(to)
	if from == to {
		return input, nil
	}
	ic := Encodings[from]
	if ic == nil {
		return nil, fmt.Errorf("from encoding not support %q", from)
	}
	oc := Encodings[to]
	if oc == nil {
		return nil, fmt.Errorf("to encoding not support %q", to)
	}
	bf, err := ic.NewDecoder().Bytes(input)
	if err != nil {
		return nil, fmt.Errorf("decode as %s failed, %w", from, err)
	}

	out, err := oc.NewEncoder().Bytes(bf)
	if err != nil {
		return nil, fmt.Errorf("encode to %s failed, %w", to, err)
	}
	return out, nil
}
