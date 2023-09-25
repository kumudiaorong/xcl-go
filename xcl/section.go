package xcl

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Section struct {
	full_name   string
	name        string
	kvs         map[string]interface{}
	secs        map[string]*Section
	update_flag bool
}

func newSec(path string, name string) *Section {
	sec := &Section{}
	if path == "" {
		sec.full_name = name
	} else {
		sec.full_name = path + name
	}
	sec.name = name
	sec.kvs = make(map[string]interface{})
	sec.secs = make(map[string]*Section)
	sec.update_flag = false
	return sec
}
func (sec *Section) SetName(name string) {
	sec.name = name
	sec.update_flag = false
}

func (sec *Section) Name() string {
	return sec.name
}

func (sec *Section) String() string {
	var sb strings.Builder
	if sec.name != "" {
		sb.WriteString("[" + sec.full_name + "]\n")
	}
	for key, value := range sec.kvs {
		sb.WriteString(key + " = ")
		switch value := value.(type) {
		case string:
			sb.WriteString("s'" + value)
		case bool:
			sb.WriteString("b'" + strconv.FormatBool(value))
		case int:
			sb.WriteString("i'" + strconv.Itoa(value))
		case uint:
			sb.WriteString("u'" + strconv.Itoa(int(value)))
		case float64:
			sb.WriteString("f'" + strconv.FormatFloat(value, 'f', -1, 64))
		default:
			sb.WriteString("s'" + fmt.Sprintf("%v", value))
		}
		sb.WriteString("\n")
	}
	sb.WriteString("\n")
	for _, value := range sec.secs {
		sb.WriteString(value.String())
	}
	sec.update_flag = true
	return sb.String()
}

func (sec *Section) Clear() {
	if len(sec.kvs) == 0 && len(sec.secs) == 0 {
		return
	}
	sec.kvs = make(map[string]interface{})
	sec.secs = make(map[string]*Section)
	sec.update_flag = false
}

func (sec *Section) NeedUpdate() bool {
	if sec.update_flag {
		return true
	}
	for _, value := range sec.secs {
		if value.NeedUpdate() {
			return true
		}
	}
	return false
}

func (sec *Section) prase_kv(line string) bool {
	var matches = kvreg.FindStringSubmatch(line)
	if matches[0] != line {
		return false
	}
	var value interface{}
	if len(matches) == 4 {
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
	} else if len(matches) == 3 {
		value = matches[2]
	} else {
		value = ""
	}
	sec.kvs[matches[1]] = value
	return true
}

var secreg = regexp.MustCompile(`([^\[\]']+)(?:'([^\[\]']+))*`)

func (sec *Section) insertSec(name_sub []string) *Section {
	fmt.Printf("insertSec: %v\n", name_sub)
	new_sec := newSec(sec.full_name, name_sub[0])
	sec.secs[name_sub[0]] = new_sec
	if len(name_sub) == 1 {
		return new_sec
	}
	return new_sec.insertSec(strings.SplitN(name_sub[1], "'", 1))
}
func (sec *Section) tryInsertSec(full_name string) (*Section, bool) {
	fmt.Printf("tryInsertSec: %v\n", full_name)
	names := strings.SplitN(full_name, "'", 1)
	sub, ok := sec.secs[names[0]]
	if !ok {
		return sec.insertSec(names), true
	}
	return sub.tryInsertSec(names[1])
}

func (sec *Section) TryInsertSec(path string) (*Section, bool, error) {
	if len(path) == 0 {
		return sec, true, nil
	}
	if !secreg.MatchString(path) {
		return nil, false, errors.New("wrong format")
	}
	sec, suc := sec.tryInsertSec(path)
	return sec, suc, nil
}


func (sec *Section) TryInsert(path string, value interface{}) (bool, error) {
	fmt.Printf("TryInsert: %v | %v\n", path, value)
	switch value.(type) {
	case string:
	case bool:
	case int:
	case uint:
	case float64:
	default:
		return false, errors.New("wrong type")
	}
	if !secreg.MatchString(path) {
		return false, errors.New("wrong format")
	}
	idx := strings.LastIndex(path, "'")
	if idx != -1 {
		sec, _ = sec.tryInsertSec(path[:idx])
	}
	_, ok := sec.kvs[path[idx+1:]]
	if !ok {
		sec.kvs[path[idx+1:]] = value
	}
	return !ok, nil
}
