package generator

import (
	"path/filepath"
	"strings"
)

var langExtDefinitions = []string{
	"Go|go",
	"Python|py",
	"JS|js",
	"TS|ts",
	"C++|cpp",
	"PHP|php",
	"C|c",
	"C#|cs",
	"Swift|swift",
	"Java|java",
	"Kotlin|kt",
	"Rust|rs",
	"Ruby|rb",
	"Racket|rkt",
	"Dart|dart",
	"Scala|scala|sc",
}

var extToLang map[string]string

func init() {
	extToLang = buildLangsMap(langExtDefinitions)
}

func buildLangsMap(defs []string) map[string]string {
	out := make(map[string]string, len(defs))
	// load all languages and their possible extensions into map[extension]language_title.
	for _, def := range defs {
		parts := strings.Split(def, "|")
		for i := 1; i < len(parts); i++ {
			out[parts[i]] = parts[0]
		}
	}

	return out
}

func langFromPath(filename string) (bool, string) {
	ext := strings.TrimPrefix(filepath.Ext(filename), ".")

	if lang, ok := extToLang[ext]; ok {
		return true, lang
	}

	return false, ext
}
