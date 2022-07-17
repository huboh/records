package structType

import (
	"reflect"
)

type ForEachFunc func(f reflect.StructField, i int)
type FilterFunc func(structField reflect.StructField, i int) bool
type MapperFunc[T any] func(structField reflect.StructField, i int) T
type ReducerFunc[T any] func(initializer T, structField reflect.StructField, i int) T

func ForEach(s reflect.Type, f ForEachFunc) {
	for i := 0; i < s.NumField(); i++ {
		f(s.Field(i), i)
	}
}

func Map[T any](s reflect.Type, f MapperFunc[T]) []T {
	o := make([]T, s.NumField())

	for i := 0; i < s.NumField(); i++ {
		o[i] = f(s.Field(i), i)
	}

	return o
}

func Reduce[T any](s reflect.Type, initializer T, f ReducerFunc[T]) T {
	o := initializer

	for i := 0; i < s.NumField(); i++ {
		o = f(o, s.Field(i), i)
	}

	return o
}

func Filter(s reflect.Type, f FilterFunc) []reflect.StructField {
	o := make([]reflect.StructField, 0)

	for i := 0; i < s.NumField(); i++ {
		if f(s.Field(i), i) {
			o = append(o, s.Field(i))
		}
	}

	return o
}
