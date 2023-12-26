package main

import (
	"fmt"
	"unsafe"
)

var usless [0]struct {
	i [0]struct {
		j [0]struct {
			k [0]struct {
				m [0]struct {
					n [0]struct{ o [0]struct{ p [0]struct{} } }
				}
			}
		}
	}
}

func main() {
	varInfo("usless", usless)
}
func varInfo[T any](name string, v T) {
	fmt.Printf("var %s: type %T, size %d\n", name, v, unsafe.Sizeof(v))
}
