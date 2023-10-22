package sys

// import "C"
import (
	"fmt"
	"os"
	"sync"
	"syscall"
	"unsafe"
)

// common signal value
const (
	OPEN   = 0b00000001
	CLOSE  = 0b00000010
	MODIFY = 0b00000100
)

func init() {
	fmt.Println("system is", sysinfo)
}

type Posix interface {
	ShmOpen(pathname string, size int64) (unsafe.Pointer, error)
	ShmClose(ptr unsafe.Pointer)
	ShmUnlink(pathname string)
}
type SystemV interface {
	ShmAt(string, int64) (unsafe.Pointer, error)
	ShmDt(unsafe.Pointer)
	ShmDel(unsafe.Pointer)
}
type UnixEn struct {
	resource_map sync.Map
}
type Semaphore struct {
	count uint8 //just count
	mutex sync.Mutex
}

func (s *UnixEn) ShmOpen(pathname string, size int64) (unsafe.Pointer, error) {
	var (
		fd  int
		err error
	)
	var rwmod, mod int = syscall.O_CREAT | syscall.O_RDWR, 0740
	if sysinfo == "linux" {
		syscall.Syscall(syscall.SYS_SHM_OPEN, uintptr(unsafe.Pointer(&pathname)), uintptr(unsafe.Pointer(&rwmod)), uintptr(unsafe.Pointer(&mod)))
	} else if sysinfo != "windows" {
		fd, err = syscall.Open(pathname, syscall.O_CREAT|syscall.O_RDWR, 0740)
	}
	if err == nil {
		defer syscall.Close(fd)
		err = syscall.Ftruncate(fd, size)
		if err == nil {
			var ans []byte
			ans, err = syscall.Mmap(fd, 0, int(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
			if err == nil {
				ptr := unsafe.Pointer(&ans[0])
				s.resource_map.Store(ptr, ans)
				return ptr, nil
			}
		}
	}
	return nil, err
}

func (s *UnixEn) ShmClose(ptr unsafe.Pointer) {
	ans, ok := s.resource_map.LoadAndDelete(ptr)
	if ok {
		syscall.Munmap(ans.([]byte))
	} else {
		fmt.Fprintln(os.Stderr, "shm close failed")
	}
}
func (s *UnixEn) ShmUnlink(pathname string) {
	var err error
	if sysinfo == "linux" {
		_, _, err = syscall.Syscall(syscall.SYS_SHM_UNLINK, uintptr(unsafe.Pointer(&pathname)), 0, 0)
	} else if sysinfo != "windows" {
		err = os.Remove(pathname)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "shm unlink error", err.Error())
	}
}

// semaphore
func (s *Semaphore) Wait() uint8 {
	for s.wait() == -1 {
	}
	return s.count
}
func (s *Semaphore) wait() int8 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.count < 1 {
		return -1
	} else {
		s.count--
		return int8(s.count)
	}
}
func (s *Semaphore) Post() uint8 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.count++
	return s.count
}
