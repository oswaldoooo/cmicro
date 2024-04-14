package delaymap

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/oswaldoooo/cmicro/kits"
	"github.com/oswaldoooo/cmicro/pkg/algorithm"
)

type _empty struct{}
type pair[Key, Val any] struct {
	Key  Key
	Val  Val
	call func() //callback function
	//internal field
	*time.Timer
	chanval reflect.Value
}

type RDelayMap[Key, Val any] struct {
	lock  sync.Mutex
	core  []pair[Key, Val]
	cmp   func(*Key, *Key) int
	top   int
	flush chan _empty
	//sets that element need ttl
	sets      *kits.Set[Key]
	sets_lock sync.Mutex
}

func NewRDelayMap[Key, Val any](cmp func(*Key, *Key) int) *RDelayMap[Key, Val] {
	if cmp == nil {
		panic("compare function can't be set to nil")
	}
	var rmap = RDelayMap[Key, Val]{
		sets:  kits.NewSet(cmp),
		flush: make(chan _empty, 1),
		cmp:   cmp,
		core:  make([]pair[Key, Val], 20),
	}
	return &rmap
}
func (r *RDelayMap[Key, Val]) Set(key Key, val Val, ttl time.Duration) {
	r.lock.Lock()
	defer r.lock.Unlock()
	var p = pair[Key, Val]{Key: key, Val: val}
	if ttl > 0 {
		p.Timer = time.NewTimer(ttl)
		p.chanval = reflect.ValueOf(p.Timer.C)
	}
	pos, need_insert := algorithm.Binary_Search(r.core, r.top, p, func(i int, p pair[Key, Val]) int {
		return r.cmp(&r.core[i].Key, &p.Key)
	})
	// fmt.Println("action ", pos, need_insert)
	if need_insert {
		corelen := len(r.core)
		if r.top+1 >= corelen {
			r.core = algorithm.Append(r.core)
		}
		copy(r.core[pos:r.top], r.core[pos+1:r.top+1])
		r.top++
	}
	r.core[pos] = p
	if ttl > 0 {
		r.sets_lock.Lock()
		r.sets.Set(key)
		r.sets_lock.Unlock()
		r.flush <- _empty{}
		// fmt.Printf("put key %v to sets. sets size %v\n", key, r.sets.Size())
		// r.sets.Range(func(k Key) error {
		// 	fmt.Println("key is", k)
		// 	return nil
		// })
	}
}
func (r *RDelayMap[Key, Val]) Delete(key Key) {
	r.lock.Lock()
	defer r.lock.Unlock()
	var p = pair[Key, Val]{Key: key}
	pos, need_insert := algorithm.Binary_Search(r.core, r.top, p, func(i int, p pair[Key, Val]) int {
		return r.cmp(&r.core[i].Key, &p.Key)
	})
	if !need_insert {
		r.sets.Delete(key)
		r.flush <- _empty{}
		copy(r.core[pos:r.top-1], r.core[pos+1:r.top])
		r.top--

	}
}

// SetCallBackWhenExpire implements DelayMap.
func (r *RDelayMap[Key, Val]) SetCallBackWhenExpire(k Key, callback func()) {
	pos, need_insert := algorithm.Binary_Search(r.core, r.top, pair[Key, Val]{Key: k}, func(i int, p pair[Key, Val]) int {
		return r.cmp(&r.core[i].Key, &p.Key)
	})
	if !need_insert {
		r.core[pos].call = callback
		r.sets_lock.Lock()
		r.sets.Set(k)
		r.sets_lock.Unlock()
		r.flush <- _empty{}
	}
}

// return nil k is not exist
func (r *RDelayMap[Key, Val]) Get(k Key) *Val {
	pos, need_insert := algorithm.Binary_Search(r.core, r.top, pair[Key, Val]{Key: k}, func(i int, p pair[Key, Val]) int {
		return r.cmp(&r.core[i].Key, &p.Key)
	})
	if !need_insert {
		return &r.core[pos].Val
	}
	return nil
}

// return nil k is not exist
func (r *RDelayMap[Key, Val]) getpair(k Key) (pair[Key, Val], bool) {
	pos, need_insert := algorithm.Binary_Search(r.core, r.top, pair[Key, Val]{Key: k}, func(i int, p pair[Key, Val]) int {
		return r.cmp(&r.core[i].Key, &p.Key)
	})
	if !need_insert {
		return r.core[pos], true
	}
	return pair[Key, Val]{}, false
}

// range function
func (r *RDelayMap[Key, Val]) Range(call func(Key, Val) error) (err error) {
	for i := 0; i < r.top; i++ {
		err = call(r.core[i].Key, r.core[i].Val)
		if err != nil {
			return
		}
	}
	return
}

// run with block
func (r *RDelayMap[Key, Val]) Run() {
	var (
		base_flush                           = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(r.flush)}
		select_list     []reflect.SelectCase = make([]reflect.SelectCase, 20)
		select_size     int                  = 1
		select_cap_size int                  = len(select_list)

		call_list []*pair[Key, Val] = make([]*pair[Key, Val], 20)
	)
	select_list[0] = base_flush
	for {
		//build select list
		select_size = 1
		r.sets_lock.Lock()
		r.sets.Range(func(k Key) error {
			pair, ok := r.getpair(k)
			if ok {
				if select_size >= select_cap_size {
					select_list = algorithm.Append(select_list)
					call_list = algorithm.Append(call_list)
					select_cap_size = len(select_list)
				}
				select_list[select_size] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: pair.chanval}
				call_list[select_size] = &pair
				// fmt.Printf("select %d is %+v\n", select_size, *call_list[select_size])
				select_size++
			} else {
				fmt.Printf("get key %v failed\n", k)
			}
			return nil
		})
		r.sets_lock.Unlock()
		//start
		pos, _, ok := reflect.Select(select_list[:select_size])
		// fmt.Printf("notify list actived %d %v sets size %d\n", pos, ok, r.sets.Size())
		if pos == 0 {
			if !ok {
				panic("logic error: can't close flush channel")
			}
			continue
		}
		if ok {
			if call_list[pos].call != nil {
				call_list[pos].call()
			}
		}
		r.Delete(call_list[pos].Key)
		r.sets_lock.Lock()
		r.sets.Delete(call_list[pos].Key)
		r.sets_lock.Unlock()
	}
}
