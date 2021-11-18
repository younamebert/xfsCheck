package common

import "reflect"

func Equal(x, y interface{}) bool {
	return reflect.DeepEqual(x, y)
}
