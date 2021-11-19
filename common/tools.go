package common

import (
	"reflect"
	"sort"
)

func Equal(x, y interface{}) bool {
	return reflect.DeepEqual(x, y)
}

func IsHave(target string, str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, target)
	return index < len(str_array) && str_array[index] == target
}
