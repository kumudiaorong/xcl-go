package xcl

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Xcl struct {
	Section
	file_abs_path   string
	last_write_time time.Time
}

func NewXcl(path string) (*Xcl, error) {
	abs_path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	var xcl = &Xcl{file_abs_path: abs_path}
	xcl.Section = *newSec("", "")
	return xcl,nil
}

var kvreg = regexp.MustCompile(`([^\s=]+)\s*=\s*(?:(s|b|i|u|f)')?([^']*)`)

func (xcl *Xcl) prase_file() error {
	ifs, err := os.Open(xcl.file_abs_path)
	if err != nil {
		return err
	}
	defer ifs.Close()
	buf := bufio.NewScanner(ifs)
	var sec *Section = &xcl.Section
	for buf.Scan() {
		line := strings.TrimSpace(buf.Text())
		if line == "" {
			continue
		}
		var l = len(line)
		if line[0] == '[' && line[l-1] == ']' {
			sec, _, err = sec.TryInsertSec(line[1 : l-1])
			if err != nil {
				continue
			}
		} else {
			sec.prase_kv(line)
		}
	}
	if buf.Err() != nil {
		return buf.Err()
	}
	return nil
}
func (xcl *Xcl) Load(path string) error {
	new_xcl, err := Load(path)
	if err != nil {
		return err
	}
	*xcl = *new_xcl
	return nil
}
func (xcl *Xcl) Reload(force bool) error {
	info, err := os.Stat(xcl.file_abs_path)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return errors.New("not a regular file")
	}
	if force || xcl.last_write_time != info.ModTime() {
		new_xcl, err := load(xcl.file_abs_path)
		if err != nil {
			return err
		}
		*xcl = *new_xcl
		return nil
	}
	return nil
}
func (xcl *Xcl) Save(force bool) error {
	if force || xcl.NeedUpdate() {
		par := filepath.Dir(xcl.file_abs_path)
		info, err := os.Stat(par)
		if os.IsNotExist(err) {
			os.MkdirAll(par, 0755)
		} else if err == nil && info.IsDir() {
			ofs, err := os.OpenFile(xcl.file_abs_path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				return err
			}
			defer ofs.Close()
			fmt.Fprint(ofs, xcl.String())
		}
	}
	return nil
}

func load(abs_path string) (*Xcl, error) {
	var xcl = &Xcl{file_abs_path: abs_path}
	xcl.kvs = make(map[string]interface{})
	xcl.secs = make(map[string]*Section)
	err := xcl.prase_file()
	if err != nil {
		return nil, err
	}
	return xcl, nil
}

func Load(path string) (*Xcl, error) {
	abs_path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(abs_path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, errors.New("not a regular file")
	}
	xcl, err := load(abs_path)
	if err != nil {
		return nil, err
	}
	return xcl, nil
}
