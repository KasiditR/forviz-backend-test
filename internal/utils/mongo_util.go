package utils

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"reflect"
	"regexp"
	"strings"
)

func toSnakeCase(str string) string {
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}

func GetName(a interface{}) string {
	t := reflect.TypeOf(a)
	name := getName(t)

	name = toSnakeCase(name)

	if !strings.HasSuffix(name, "s") {
		name += "s"
	}

	return name
}

func getName(t reflect.Type) string {
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Ptr || t.Kind() == reflect.Array || t.Kind() == reflect.Map || t.Kind() == reflect.Chan {
		return getName(t.Elem())
	}
	return t.Name()
}

func GetID(a interface{}) bson.ObjectID {
	tv := reflect.ValueOf(a)
	mVal := tv.FieldByName("ID")
	if reflect.Value.IsZero(mVal) {
		return bson.NilObjectID
	}
	aV := mVal.Interface().(bson.ObjectID)
	return aV
}
