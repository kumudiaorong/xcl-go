package xcl

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	// "unsafe"
)

type Section struct {
	full_name     string
	name          string
	kvs           map[string]any
	sub_section   map[string]*Section
	section_array map[string]*[]*Section
}

func NewSec(path string, name string) *Section {
	sec := &Section{}
	if path == "" {
		sec.full_name = name
	} else {
		sec.full_name = path + "'" + name
	}
	sec.name = name
	sec.kvs = make(map[string]any)
	sec.sub_section = make(map[string]*Section)
	sec.section_array = make(map[string]*[]*Section)
	return sec
}

func (sec *Section) Name() string {
	return sec.name
}

func (sec *Section) SetName(name string) {
	sec.name = name
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
	for _, value := range sec.sub_section {
		sb.WriteString(value.String())
	}
	return sb.String()
}

func (sec *Section) Clear() {
	if len(sec.kvs) == 0 && len(sec.sub_section) == 0 {
		return
	}
	sec.kvs = make(map[string]any)
	sec.sub_section = make(map[string]*Section)
}

var secreg = regexp.MustCompile(`([^\[\]']+)(?:'([^\[\]']+))*`)

func (sec *Section) insertSec(name_sub []string) *Section {
	new_sec := NewSec(sec.full_name, name_sub[0])
	sec.sub_section[name_sub[0]] = new_sec
	if len(name_sub) > 1 {
		return new_sec.insertSec(strings.SplitN(name_sub[1], "'", 1))
	}
	return new_sec
}
func (sec *Section) tryInsertSec(full_name string) (*Section, bool) {
	names := strings.SplitN(full_name, "'", 1)
	sub, ok := sec.sub_section[names[0]]
	if !ok {
		return sec.insertSec(names), true
	}
	if len(names) > 1 {
		return sub.tryInsertSec(names[1])
	}
	return sub, false
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
func (sec *Section) tryInsertValue(key string, value any) bool {
	_, ok := sec.kvs[key]
	if !ok {
		sec.kvs[key] = value
	}
	return !ok
}

func (sec *Section) TryInsert(path string, value any) (bool, error) {
	if !secreg.MatchString(path) {
		return false, errors.New("wrong format")
	}
	idx := strings.LastIndex(path, "'")
	if idx != -1 {
		sec, _ = sec.tryInsertSec(path[:idx])
	}
	k := reflect.TypeOf(value).Kind()
	if k == reflect.String || k == reflect.Bool || k == reflect.Int || k == reflect.Uint || k == reflect.Float64 {
		return sec.tryInsertValue(path[idx+1:], value), nil
	}
	return false, nil
}
func (sec *Section) Find(path string) (any, error) {
	if len(path) == 0 {
		return nil, errors.New("Empty Path")
	}
	if !secreg.MatchString(path) {
		return nil, errors.New("wrong format")
	}
	names := strings.SplitN(path, "'", 1)
	if len(names) > 1 {
		sub, ok := sec.sub_section[names[0]]
		if !ok {
			return nil, nil
		} else {
			return sub.Find(names[1])
		}
	} else {
		v, ok := sec.kvs[names[0]]
		if !ok {
			return nil, nil
		} else {
			return v, nil
		}
	}
}
func decodePtr(ptr *reflect.Value) {
	k := ptr.Kind()
	for k == reflect.Ptr {
		if ptr.IsNil() && ptr.CanAddr() {
			ptr.Set(reflect.New(ptr.Type().Elem()))
		}
		*ptr = ptr.Elem()
		k = ptr.Kind()
	}
}
func (sec *Section) decode(e reflect.Value) error {
	t := e.Type()
	for i := 0; i < e.NumField(); i++ {
		ef := e.Field(i)
		efk := ef.Kind()
		// ef.IsValid()
		tf := t.Field(i)
		decodePtr(&ef)
		if !ef.IsValid() {
			return errors.New("Wrong Value")
		}
		if efk == reflect.Struct {
			sub, ok := sec.sub_section[tf.Name]
			if !ok {

			} else {
				sub.decode(ef)
			}
		} else if efk == reflect.String || efk == reflect.Bool || efk == reflect.Int || efk == reflect.Uint || efk == reflect.Float64 {
			v, ok := sec.kvs[tf.Name]
			if !ok {
			} else {
				ef.Set(reflect.ValueOf(v))
			}
		}
	}
	return nil
}
func (sec *Section) Decode(s any) error {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr {
		return errors.New("Not A Pointer")
	}
	if v.IsNil() {
		return errors.New("Nil Value")
	}
	v = v.Elem()
	for v.Kind() == reflect.Ptr {
		if v.IsNil() && v.CanAddr() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return errors.New(fmt.Sprintf("Not A Struct Pointer, Is A %+v", v.Type().Kind()))
	}
	return sec.decode(v)
}
func (sec *Section) encode(e reflect.Value) error {
	t := e.Type()
	for i := 0; i < e.NumField(); i++ {
		ef := e.Field(i)
		efk := ef.Kind()
		// ef.IsValid()
		tf := t.Field(i)
		decodePtr(&ef)
		if !ef.IsValid() {
			return errors.New("Wrong Value")
		}
		if efk == reflect.Struct {
			sub, _ := sec.tryInsertSec(tf.Name)
			sub.encode(ef)
		} else if efk == reflect.String || efk == reflect.Bool || efk == reflect.Int || efk == reflect.Uint || efk == reflect.Float64 {
			// sec.kvs[tf.Name]
			if ef.CanInterface() {
				sec.kvs[tf.Name] = ef.Interface()
			}
		} else if efk == reflect.Slice {
			efi := ef.Index(i)
			decodePtr(&efi)
			efk := efi.Kind()
			var f func(reflect.Value)
			if efk == reflect.Struct {
				sub, _ := sec.tryInsertSec(tf.Name)
				sub.encode(efi)
			}
			// if efk
			// for i := 0; i < ef.Len(); i++ {
			// 	ef.Index(i)
			// }
		}
	}
	return nil
}
func (sec *Section) Encode(s any) error {
	v := reflect.ValueOf(s)
	for v.Kind() == reflect.Ptr {
		if v.IsNil() && v.CanAddr() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return errors.New(fmt.Sprintf("Not A Struct Pointer, Is A %+v", v.Type().Kind()))
	}
	return sec.encode(v)
}
