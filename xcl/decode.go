package xcl

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var kvreg = regexp.MustCompile(`([^\s=]+)\s*=\s*(?:(s|b|i|u|f)')?([^']*)`)

func praseKV(line string) (string, any, error) {
	var matches = kvreg.FindStringSubmatch(line)
	if matches[0] != line {
		return "", nil, errors.New("Wrong Format Key-Value Pair")
	}
	var value any
	if len(matches) == 4 {
		var err error = nil
		switch matches[2] {
		case "s":
			value = matches[3]
		case "b":
			value = matches[3] == "true"
		case "i":
			var i int
			if i, err = strconv.Atoi(matches[3]); err == nil {
				value = i
			}
		case "u":
			var i int
			if i, err = strconv.Atoi(matches[3]); err == nil {
				value = uint(i)
			}
		case "f":
			var f float64
			if f, err = strconv.ParseFloat(matches[3], 32); err == nil {
				value = float64(f)
			}
		}
		if err != nil {
			return "", nil, err
		}
	} else if len(matches) == 3 {
		value = matches[2]
	} else {
		value = ""
	}
	return matches[1], value, nil
}

func Decode(i io.Reader) (*Section, error) {
	buf := bufio.NewScanner(i)
	var sec *Section = NewSec("", "")
	var sub = sec
	for buf.Scan() {
		line := strings.TrimSpace(buf.Text())
		if line == "" {
			continue
		}
		var l = len(line)
		if line[0] == '[' && line[l-1] == ']' {
			path := line[1 : l-1]
			if !secreg.MatchString(path) {
				// return nil, false, errors.New("wrong format")
			}
			sub = sec.insertSec(strings.SplitN(path, "'", 1))
		} else {
			k, v, err := praseKV(line)
			if err != nil {

			} else {
				sub.kvs[k] = v
			}
		}
	}
	if buf.Err() != nil {
		return nil, buf.Err()
	}
	return sec, nil
}
func Encode(sec *Section) string {
	return sec.String()
}
