package core

import (
	"fmt"
	"strconv"
	"strings"
	"bufio"
	"io"
	"log"
)

func Str2Int(s string) (int, error) {
	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return -1, err
	}

	return int(n), nil
}

func Str2IntSlice(s string) ([]int, error) {
	strSlice := strings.Split(s, ",")
	intSlice := make([]int, len(strSlice))
	for _, v := range strSlice {
		n, err := Str2Int(v)
		if err != nil {
			return nil, err
		}

		intSlice = append(intSlice, int(n))
	}

	return intSlice, nil
}

func Str2Int32Slice(s string) ([]int32, error) {
	strSlice := strings.Split(s, ",")
	intSlice := make([]int32, len(strSlice))
	for _, v := range strSlice {
		n, err := Str2Int(v)
		if err != nil {
			return nil, err
		}

		intSlice = append(intSlice, int32(n))
	}

	return intSlice, nil
}

func ConcatStrings(ss []string, sep string) string {
	var b strings.Builder
	for _, s := range ss {
		b.WriteString(s)
		b.WriteString(sep)
	}

	return b.String()
}

func ConcatRunes(rss [][]rune) string {
	var b strings.Builder
	for _, rs := range rss {
		bs := []byte(string(rs))
		b.Write(bs)
	}

	fmt.Println(b.String())

	return b.String()
}

func ScanRunes(text string) {
	r := bufio.NewReader(strings.NewReader(text))
    for {
        if c, sz, err := r.ReadRune(); err != nil {
            if err == io.EOF {
                break
            } else {
                log.Fatal(err)
            }
        } else {
            fmt.Printf("%q [%d]\n", string(c), sz)
        }
    }
}
