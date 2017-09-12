package utils

import "github.com/spf13/viper"

//LangCheckFields проверяет поля на пустые значения, возввращает true или false
func LangCheckFields(dt map[string]interface{}) bool {
	if len(dt) == 0 {
		return false
	}
	for _, value := range dt {
		v := value.(string)
		if v == "" {
			return false
		}
	}
	return true
}

//LangGetField ищет в объекте требуемое текстовое поле, если его нет, выдаёт по умолчанию обязательное (ru).
func LangGetField(maps map[string]interface{}, lang string) string {
	l := viper.GetString("services.email.lang")
	if maps[lang] == nil {
		return maps[l].(string)
	}
	return maps[lang].(string)
}
