package internal

import (
	reflect "github.com/goccy/go-reflect"
	"unsafe"
)

func Btoa(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Atob(s string) []byte {
	sp := unsafe.Pointer(&s)
	b := *(*[]byte)(sp)
	(*reflect.SliceHeader)(unsafe.Pointer(&b)).Cap = (*reflect.StringHeader)(sp).Len
	return b
}
