package main

// Language is an enum with a string underlying type for the supported languages
type Language string

const (
	lua  Language = "lua"
	wren Language = "wren"
	moon Language = "moon"
	auto Language = "auto"
)

func isSupportedLanguage(lang Language) bool {
	if lang == lua {
		return true
	}
	if lang == wren {
		return true
	}
	if lang == moon {
		return true
	}
	return false
}
