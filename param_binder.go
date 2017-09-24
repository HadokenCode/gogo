package gogo

import (
	"reflect"
)

type Decoder struct {
	binder func(v string)
}

type Binder struct {
	types map[reflect.Type]*Decoder
}
