package utils

import (
	"reflect"
	"runtime"
)

func GetFunctionName(i interface{}, nameOnly bool) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
