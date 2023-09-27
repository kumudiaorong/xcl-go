package xcl

import (
	// "errors"
	// "reflect"
	"bufio"
	"io"
	"strings"
)

type Decoder struct {
	sec *Section
}

func NewDecoder(i io.Reader) (*Decoder, error) {
	buf := bufio.NewScanner(i)
	var sec *Section = NewSec("", "")
	var sub = sec
	var err error
	for buf.Scan() {
		line := strings.TrimSpace(buf.Text())
		if line == "" {
			continue
		}
		var l = len(line)
		if line[0] == '[' && line[l-1] == ']' {
			sub, _, err = sec.TryInsertSec(line[1 : l-1])
			if err != nil {
				continue
			}
		} else {
			k, v, err := prase_kv(line)
			if err != nil {

			} else {
				sub.kvs[k] = v
			}
		}
	}
	if buf.Err() != nil {
		return nil, buf.Err()
	}
	return &Decoder{sec: sec}, nil
}

func (decoder *Decoder) Decode(s interface{}) error {
	return decoder.sec.Decode(s)
	// for i := 0; i < v.NumField(); i++ {
	// f := e.Field(i)
	// tf := t.Field(i)
	// sec.f tf.Name
	// if tf.Name == {

	// }
	// }
}
