// Package str 字符串辅助方法
package str

import (
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

// Plural转为复数user->users
func Plural(word string) string {
	return pluralize.NewClient().Plural(word)
}

// Singular 转为单数 users ->user
func Singular(word string) string {
	return pluralize.NewClient().Singular(word)
}

// Snake转为snake_case,如TopicComment -> topic_comment
func Snake(s string) string {
	return strcase.ToSnake(s)
}

// Camel 转为CamelCase,如topic_comment -> TopicComment
func Camel(s string) string {
	return strcase.ToCamel(s)
}

// LowerCamel 转为lowerCamelCase,如TopicComment -> topicComment
func LowerCamel(s string) string {
	return strcase.ToLowerCamel(s)
}
