package main

import (
	"embed"
	"encoding/json"
)

//go:embed translations
var translations embed.FS

type Translation struct {
	Lang         string
	Translations map[string]string
}

var translationsMap = make(map[string]*Translation)

func loadTranslation(lang string) (*Translation, error) {
	if translation, ok := translationsMap[lang]; ok {
		return translation, nil
	}
	data, err := translations.ReadFile("translations/" + lang + ".json")
	if err != nil {
		return nil, err
	}
	var translation Translation
	err = json.Unmarshal(data, &translation.Translations)
	if err != nil {
		return nil, err
	}
	translationsMap[lang] = &translation
	return &translation, nil
}

func GetTranslation(lang, key string) string {
	translation, err := loadTranslation(lang)
	if err != nil {
		translation, err = loadTranslation("en")
		if err != nil {
			return key
		}
	}
	value, ok := translation.Translations[key]
	if !ok {
		jsonData, _ := json.Marshal(translation.Translations)
		return string(jsonData)
	}
	return value
}
