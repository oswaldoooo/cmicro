package rdtree

import "bytes"

// this is red black tree map
type node struct {
	key []byte
	val []byte
}
type RdTree struct {
	d    []node
	size int
}

var (
	Default_RdTree_Capsize = 100
)

func NewRdTree() *RdTree {
	return &RdTree{d: make([]node, Default_RdTree_Capsize)}
}

type Bytor interface {
	Bytes() []byte
}

func find_suitable_position(ndqueue []node, size int, key Bytor) (int, bool) {
	if size == 0 {
		return 0, true
	}
	var (
		mid, ok     int
		left, right int = 0, size
	)
	for right-left > 1 {
		mid = (right + left) / 2
		ok = bytes.Compare(key.Bytes(), ndqueue[mid].key)
		if ok == 0 {
			return mid, false
		} else if ok == -1 {
			return mid, true
		} else {
			//move right
			left = mid
		}
	}
	return size, true
}
func find_position(ndqueue []node, size int, key Bytor) int {
	if size == 0 {
		return -1
	}
	var (
		mid, ok     int
		left, right int = 0, size
	)
	for right-left > 1 {
		mid = (right + left) / 2
		ok = bytes.Compare(key.Bytes(), ndqueue[mid].key)
		if ok == 0 {
			return mid
		} else if ok == -1 {
			return -1
		}
	}
	return -1
}

func (s *RdTree) Insert(key, val Bytor) {
	pos, isinert := find_suitable_position(s.d, s.size, key)
	if pos == -1 {
		panic("get position failed")
	}
	if isinert {
		if s.size >= Default_RdTree_Capsize {
			//extend
			newd := make([]node, 3*s.size/2)
			copy(newd, s.d)
			s.d = newd
		}
		copy(s.d[pos+1:s.size+1], s.d[pos:s.size])
		s.size++
	}
	s.d[pos].key = key.Bytes()
	s.d[pos].val = val.Bytes()
}
func RdTreeFind[T any](s *RdTree, key Bytor) *T {

	return nil
}

func RdTreeErase[T any](s *RdTree, key Bytor) {

}
