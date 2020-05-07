package main

import "unsafe"

// []byte 转 string
func bytesToStr(value []byte) string {
	return *(*string)(unsafe.Pointer(&value)) // nolint
}

// string 转 []byte
func strToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s)) // nolint
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h)) // nolint
}
