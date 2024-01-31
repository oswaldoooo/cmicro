package sys

import "unsafe"

type Pointer[T any] struct {
	T *T
	b []byte
}

func NewPointer[T any]() *Pointer[T] {
	var (
		ans = &Pointer[T]{}
		t   T
	)
	ans.b = make([]byte, unsafe.Sizeof(t))
	ans.T = (*T)(unsafe.Pointer(&ans.b[0]))
	return ans
}
func NewPointer2[T any](b []byte) *Pointer[T] {
	var (
		t T
	)
	if len(b) >= int(unsafe.Sizeof(t)) {
		return &Pointer[T]{b: b, T: (*T)(unsafe.Pointer(&b[0]))}
	}
	return nil
}
func (s *Pointer[T]) Bytes() []byte {
	return s.b
}
