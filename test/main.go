package main

import (
	"github.com/kumudiaorong/xcl-go/xcl"
)

func main() {
	sec := xcl.NewSec("", "")
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
	// sec.tryInsertValue("Key", "value")
	// sec.tryInsertValue("PKey", "pvalue")
	// sec.TryInsert("Sub'Key", "svalue")
	// sec.TryInsert("Sub'PKey", "spvalue")
	sec.TryInsert("PSub'Key", "psvalue")
	sec.TryInsert("PSub'PKey", "pspvalue")

	// assert(t, sec.secs["PSub"].kvs["PKey"].(string) == "pspvalue", "WWW")
	obj := &TT{Sub: T{Key: ""}, Key: ""}
	_ = sec.Decode(&obj)
	// assert(t, err == nil, fmt.Sprintf(`Sec Decode(...) Unexpect Error %v`, err))
	// assert(t, obj.Key == "value", fmt.Sprintf(`Sec Decode(...) Wrong Value %s, Want "value"`, obj.Key))
	// assert(t, *obj.PKey == "pvalue", fmt.Sprintf(`Sec Decode(...) Wrong Pointer Value %s, Want "value"`, *obj.PKey))
	// assert(t, obj.Sub.Key == "svalue", fmt.Sprintf(`Sec Decode(...) Wrong Sub Value %s, Want "svalue"`, obj.Sub.Key))
	// assert(t, *obj.Sub.PKey == "spvalue", fmt.Sprintf(`Sec Decode(...) Wrong Sub Pointer Value %s, Want "svalue"`, *obj.Sub.PKey))
	// assert(t, obj.PSub.Key == "psvalue", fmt.Sprintf(`Sec Decode(...) Wrong Pointer Sub Value %s, Want "svalue"`, obj.PSub.Key))
	// assert(t, *obj.PSub.PKey == "pspvalue", fmt.Sprintf(`Sec Decode(...) Wrong Pointer Sub Pointer Value %s, Want "pspvalue"`, *obj.PSub.PKey))
}
