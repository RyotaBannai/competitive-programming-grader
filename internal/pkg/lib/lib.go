package lib

import (
	"fmt"
	"reflect"
	"sort"

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

// e.g. fs.FileInfo
type HasName interface {
	Name() string
}

func SortFilebyName[T HasName](s []T) {
	sort.Slice(s, func(i, j int) bool {
		return s[i].Name() < s[j].Name()
	})
}

func Contains[T comparable](slice []T, item T) bool {
	set := make(map[T]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
