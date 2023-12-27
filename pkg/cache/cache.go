package cache

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"
)

type TTLVal struct {
	t *time.Ticker
	v []string //key list
}
type KeyMap struct {
	m           map[string]any
	ttlmap      map[string]*TTLVal
	mux, ttlmux sync.Mutex
	chanlist    []chan bool
	refchan     chan bool
}

func backend(km *KeyMap) {
	var (
		chanlist []reflect.SelectCase
		chanmap  map[int]string = make(map[int]string)
		// refchan                 = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(km.refchan)}
		v     reflect.Value
		ve    any
		index int
		// ticker   *time.Ticker
		ok      bool
		chaname string
	)
	for {
		// fmt.Println("into cycle", len(km.ttlmap))
		// km.ttlmux.Lock()
		// for !km.ttlmux.TryLock() {
		// 	fmt.Println("backend try lock")
		// 	time.Sleep(20 * time.Millisecond)
		// }
		// fmt.Println("ttlmux locked")
		chanlist = make([]reflect.SelectCase, 1+len(km.ttlmap))
		chanlist[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(km.refchan)}
		i := 1
		for ttlname, ee := range km.ttlmap {
			chanlist[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ee.t.C)}
			chanmap[i] = ttlname
			i++
		}
		// km.ttlmux.Unlock()
		// fmt.Println("ttlmux unlocked")
		fmt.Println(len(chanlist))
		index, v, ok = reflect.Select(chanlist)
		if ok {
			ve = v.Interface()
			if _, ok = ve.(bool); ok {
				//refresh
				continue
			}
			chaname = chanmap[index]
			km.mux.Lock()
			//delete bind keys
			for _, ename := range km.ttlmap[chaname].v {
				delete(km.m, ename)
			}
			km.mux.Unlock()
			km.ttlmux.Lock()
			delete(km.ttlmap, chaname)
			km.ttlmux.Unlock()
			delete(chanmap, index)

		}
	}
}
func NewKeyMap() *KeyMap {
	ans := &KeyMap{m: make(map[string]any), refchan: make(chan bool, 1), ttlmap: make(map[string]*TTLVal)}
	go backend(ans)
	// fmt.Println("backend start")
	return ans
}
func (s *KeyMap) SetLease(leasename string, ttl int, bindkeys ...string) {
	for !s.ttlmux.TryLock() {
		// fmt.Println("set lease try lock")
		time.Sleep(10 * time.Millisecond)
	}
	// fmt.Println("set lease start")
	defer s.ttlmux.Unlock()
	// defer fmt.Println("set lease finised")
	if _, ok := s.ttlmap[leasename]; !ok {
		s.ttlmap[leasename] = &TTLVal{t: time.NewTicker(time.Duration(ttl) * time.Second)}
	}
	s.ttlmap[leasename].t.Reset(time.Duration(ttl) * time.Second)
	if len(bindkeys) > 0 {
		s.ttlmap[leasename].v = append(s.ttlmap[leasename].v, bindkeys...)
	}
	s.refchan <- true
}
func (s *KeyMap) SetKey(key string, val any, ttl int) {
	if ttl > 0 {
		s.SetLease(key, ttl, key)
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	s.m[key] = val
}
func (s *KeyMap) Range(fun func(string, any) error) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	var err error
	// fmt.Println("range total number", len(s.m))
	for k, v := range s.m {
		err = fun(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *KeyMap) BindKey2TTL(ttlname string, keys ...string) error {
	for !s.ttlmux.TryLock() {
		time.Sleep(10 * time.Millisecond)
	}
	defer s.ttlmux.Unlock()
	if info, ok := s.ttlmap[ttlname]; ok {
		for !s.mux.TryLock() {
			time.Sleep(10 * time.Millisecond)
		}
		defer s.mux.Unlock()
		for _, key := range keys {
			if _, ok := s.m[key]; ok {
				info.v = append(info.v, key)
			} else {
				return errors.New("key " + key + " is not existed")
			}
		}
	} else {
		return errors.New("ttl " + ttlname + " is not existed")
	}
	return nil
}
