package cobrautil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"syscall"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func panicf(format string, args ...any) {
	panic(fmt.Sprintf(format, args...))
}

func ensureDir(dir string) (err error) {
	err = os.MkdirAll(dir, os.ModePerm)
	if errors.Is(err, syscall.EEXIST) {
		err = nil
		goto end
	}
end:
	return err
}

func fileExists(file string) (exists bool, err error) {
	// TODO Differentiate that from dir
	_, err = os.Stat(file)
	if errors.Is(err, fs.ErrNotExist) {
		err = nil
		exists = false
		goto end
	}
	if err != nil {
		goto end
	}
	exists = true
end:
	return exists, err
}

func marshalJSONFile(_ Context, file string, object any) error {
	content, err := json.MarshalIndent(object, "", "\t")
	if err != nil {
		goto end
	}
	err = os.WriteFile(file, content, os.ModePerm)
end:
	return err
}

var DefaultLanguage = language.English

var defaultLanguage language.Tag
var titler cases.Caser

func title(s string) string {
	if defaultLanguage != DefaultLanguage {
		defaultLanguage = DefaultLanguage
		titler = cases.Title(defaultLanguage)
	}
	return titler.String(s)
}
