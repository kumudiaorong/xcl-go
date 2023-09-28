package xcl

import (
	"fmt"
	"reflect"
	"testing"
)

func assert(t *testing.T, stat any, estr string) {
	T := reflect.TypeOf(stat)
	if T.Kind() == reflect.Bool && !stat.(bool) {
		t.Fatalf("Assert : %s", estr)
	}
}

func TestNewSec(t *testing.T) {
	sec := NewSec("path", "name")
	assert(t, sec.full_name == "path'name", "New Sec Wrong Field full_name")
	assert(t, sec.kvs != nil, "New Sec Wrong Field kvs")
	assert(t, sec.sub_section != nil, "New Sec Wrong Field secs")
	assert(t, sec.name == "name", "New Sec Wrong Field name")
}
func TestSecName(t *testing.T) {
	sec := NewSec("path", "name")
	assert(t, sec.Name() == sec.name, fmt.Sprintf("Sec Name() Wrong, Want %s", sec.name))
}
func TestSecSetName(t *testing.T) {
	sec := NewSec("path", "name")
	sec.SetName("new")
	assert(t, sec.name == "new", "Sec SetName() Wrong, Want name")
}
func TestSecTryInsertValue(t *testing.T) {
	sec := NewSec("", "")
	f := func(k string, v any) {
		assert(t, sec.tryInsertValue(k, v), `Sec TryInsertValue("key","str") Fail`)
		assert(t, sec.kvs[k] == v, `Sec TryInsertValue("key","str") Insert Wrong Value`)
		assert(t, !sec.tryInsertValue(k, v), `Sec TryInsertValue("key","str") Should not Success`)
	}
	f("key", "str")
	f("bool", true)
	f("int", int(1))
	f("uint", uint(1))
	f("float64", float64(1))
}
func TestSecFind(t *testing.T) {
	sec := NewSec("", "")
	sec.tryInsertValue("key", "str")
	v, err := sec.Find("key")
	assert(t, err == nil, fmt.Sprintf(`Sec Find("key","str") Unexpect Error %v`, err))
	assert(t, v.(string) == "str", fmt.Sprintf(`Sec Find("key","str") Wrong Value %s, Want "str"`, v))
}
func TestSecTryInsert(t *testing.T) {
	sec := NewSec("", "")
	suc, err := sec.TryInsert("key", "str")
	assert(t, err == nil, fmt.Sprintf(`Sec TryInsert("key", "str") Unexpect Error %v`, err))
	assert(t, suc, `Sec TryInsert("key", "str") Fail`)
	v, ok := sec.kvs["key"]
	assert(t, ok, `Sec TryInsert("key", "str"), Fail`)
	assert(t, v.(string) == "str", `Sec TryInsert("key", "str"), Wrong Value`)
	suc, err = sec.TryInsert("key'sub", "str")
	assert(t, err == nil, fmt.Sprintf(`Sec TryInsert("key'sub", "str") Unexpect Error %v`, err))
	assert(t, suc, `Sec TryInsert("key'sub", "str") Fail`)
	sub, ok := sec.sub_section["key"]
	assert(t, ok, `Sec TryInsert("key'sub", "str") Fail`)
	v, ok = sub.kvs["sub"]
	assert(t, ok, `Sec TryInsert("key", "str"), Fail`)
	assert(t, v.(string) == "str", `Sec TryInsert("key", "str"), Wrong Value`)
}
func TestSecDecode(t *testing.T) {
	sec := NewSec("", "")
	type T struct {
		PKey *string
		Key  string
	}
	type TT struct {
		PSub *T
		Sub  T
		PKey *string
		Key  string
	}
	sec.tryInsertValue("Key", "value")
	sec.tryInsertValue("PKey", "pvalue")
	sec.TryInsert("Sub'Key", "svalue")
	sec.TryInsert("Sub'PKey", "spvalue")
	sec.TryInsert("PSub'Key", "psvalue")
	sec.TryInsert("PSub'PKey", "pspvalue")
	obj := &TT{Sub: T{Key: ""}, Key: ""}
	err := sec.Decode(obj)
	assert(t, err == nil, fmt.Sprintf(`Sec Decode(...) Unexpect Error %v`, err))
	assert(t, obj.Key == "value", fmt.Sprintf(`Sec Decode(...) Wrong Value %s, Want "value"`, obj.Key))
	assert(t, *obj.PKey == "pvalue", fmt.Sprintf(`Sec Decode(...) Wrong Pointer Value %s, Want "value"`, *obj.PKey))
	assert(t, obj.Sub.Key == "svalue", fmt.Sprintf(`Sec Decode(...) Wrong Sub Value %s, Want "svalue"`, obj.Sub.Key))
	assert(t, *obj.Sub.PKey == "spvalue", fmt.Sprintf(`Sec Decode(...) Wrong Sub Pointer Value %s, Want "svalue"`, *obj.Sub.PKey))
	assert(t, obj.PSub.Key == "psvalue", fmt.Sprintf(`Sec Decode(...) Wrong Pointer Sub Value %s, Want "svalue"`, obj.PSub.Key))
	assert(t, *obj.PSub.PKey == "pspvalue", fmt.Sprintf(`Sec Decode(...) Wrong Pointer Sub Pointer Value %s, Want "pspvalue"`, *obj.PSub.PKey))
}
