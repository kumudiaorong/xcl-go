package xcl

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var kvreg = regexp.MustCompile(`([^\s=]+)\s*=\s*(s|b|i|u|f)('|\[)([^']*)`)

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

var secLineFmt = regexp.MustCompile(`\[([^\[\]']+(?:'[^\[\]']+)*)\](?:\[([^\[\]']+)\])?`)

func Decode(i io.Reader) (*Section, error) {
	buf := bufio.NewScanner(i)
	var sec *Section = NewSec("", "")
	var sub = sec
	for buf.Scan() {
		line := strings.TrimSpace(buf.Text())
		if line == "" {
			continue
		}
		// var l = len(line)
		if line[0] == '[' {
			matchs := secLineFmt.FindStringSubmatch(line)
			if matchs[0] != line {

			}
			// if matchs
			// path := line[1 : l-1]
			// if !secreg.MatchString(path) {
			// 	// return nil, false, errors.New("wrong format")
			// }
			sub = sec.insertSec(strings.SplitN(matchs[1], "'", 1))
			if len(matchs) > 2 {
				subset, ok := sub.section_array[matchs[2]]
				if !ok {
					subset = &[]*Section{}
					sub.section_array[matchs[2]] = subset
				}
				sub = NewSec(sub.full_name, matchs[2])
				*subset = append(*subset, sub)
			}
		} else {
			var matches = kvreg.FindStringSubmatch(line)
			if matches[0] != line {
				// return "", nil, errors.New("Wrong Format Key-Value Pair")
			}
			var value any = nil
			if matches[3] == "'" {
				switch matches[2] {
				case "s":
					value = matches[3]
				case "b":
					value = matches[3] == "true"
				case "i":
					if i, err := strconv.Atoi(matches[3]); err == nil {
						value = i
					}
				case "u":
					if i, err := strconv.Atoi(matches[3]); err == nil {
						value = uint(i)
					}
				case "f":
					if f, err := strconv.ParseFloat(matches[3], 32); err == nil {
						value = float64(f)
					}
				}
				if value == nil {
					// return "", nil, err
				}
			} else if matches[3] == "[" && matches[4] == "" {
				for buf.Scan() {
					line = strings.TrimSpace(buf.Text())
					l := len(line)
					if l == 0 || line[l-1] == ',' {
						
					}
				}
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
