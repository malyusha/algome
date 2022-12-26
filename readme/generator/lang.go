package generator

import (
	"path/filepath"
	"strings"
)

const (
	LangGo         = "Go"
	LangPython     = "Python"
	LangJavascript = "JS"
	LangCpp        = "C++"
	LangPhp        = "PHP"

	ExtGo     = "go"
	ExtPython = "py"
	ExtJS     = "js"
	ExtCpp    = "cpp"
	ExtPhp    = "php"
)

var extToLang = map[string]string{
	ExtGo:     LangGo,
	ExtPython: LangPython,
	ExtJS:     LangJavascript,
	ExtCpp:    LangCpp,
	ExtPhp:    LangPhp,
}

func langFromPath(filename string) (bool, string) {
	ext := strings.TrimPrefix(filepath.Ext(filename), ".")

	if lang, ok := extToLang[ext]; ok {
		return true, lang
	}

	return false, ext
}
