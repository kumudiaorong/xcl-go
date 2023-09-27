package xcl

import (
	"errors"
	"strconv"
)

func prase_kv(line string) (string, interface{}, error) {
	var matches = kvreg.FindStringSubmatch(line)
	if matches[0] != line {
		return "", nil, errors.New("Wrong Format Key-Value Pair")
	}
	var value interface{}
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
