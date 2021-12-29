package misc

import (
	"fmt"
	"reflect"

	"github.com/k0kubun/pp"
)

func Debug(c interface{}) {
	pp.Println(c)
}

type Zipped struct {
	Fst, Snd interface{}
}

func Zip(a, b interface{}) ([]Zipped, error) {
	v1 := reflect.ValueOf(a)
	v2 := reflect.ValueOf(b)

	if v1.Kind() == reflect.Ptr {
		v1 = v1.Elem()
	}
	if v2.Kind() == reflect.Ptr {
		v2 = v2.Elem()
	}

	if v1.Kind() != reflect.Slice || v2.Kind() != reflect.Slice {
		return nil, fmt.Errorf("expected slice type, found for\nfirst param [%v], and for\nsecond param [%v]", v1.Kind().String(), v2.Kind().String())
	}

	if v1.Len() != v2.Len() {
		return nil, fmt.Errorf("arguments must be the same length")
	}

	r := make([]Zipped, v1.Len())
	for i := 0; i < v1.Len(); i++ {
		r[i] = Zipped{v1.Index(i).Interface(), v2.Index(i).Interface()}
	}

	return r, nil
}
