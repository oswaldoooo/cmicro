package sys

/*
#include <sys/sem.h>
#include <sys/shm.h>
#include <stdlib.h>
#include <fcntl.h>
#include <sys/mman.h>
#include <stdio.h>
#include <unistd.h>
static int shmflag=IPC_CREAT|0640;
void* open_shm(char* pathname,int sizes,short sysinfo){
	int shmid=(sysinfo==1)?shm_open(pathname,O_RDWR|O_CREAT,0644):open(pathname,O_RDWR|O_CREAT,0644);
	if (shmid>0){
		if(ftruncate(shmid,sizes)!=-1){
			void* ans=mmap(NULL,sizes,PROT_READ|PROT_WRITE,MAP_SHARED,shmid,0);
			if(ans!=MAP_FAILED){
				return ans;
			}else{
				fprintf(stderr,"mmap failed\n");
			}
		}else{
			fprintf(stderr,"ftruncate failed\n");
		}
	}else{
		fprintf(stderr,"open shmid failed\n");
	}
	return 0;
}
void unmap(void* src,int sizes){
	msync(src,sizes,MS_SYNC);
	munmap(src,sizes);
}
*/
import "C"
import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

var sysinfo = strings.ToLower(runtime.GOOS)

// system v share memory open
func GetShare_Mem[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64](shmid int, dst_ptr **T) uintptr {
	shm, _, err := syscall.Syscall(syscall.SYS_SHMAT, uintptr(shmid), 0, 0)
	if len(err.Error()) < 1 {
		*dst_ptr = (*T)(unsafe.Pointer(shm))
		return 0
	}
	return shm
}
func CreateShare_Mem[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64](pathname string, gendid int, dst_ptr **T) uintptr {
	pn := C.CString(pathname)
	key := C.ftok(pn, (C.int)(gendid))
	C.free(unsafe.Pointer(pn))
	var te T
	shmid, _, err := syscall.Syscall(syscall.SYS_SHMGET, uintptr(key), unsafe.Sizeof(te), 01000|0640)
	if err != 0 {
		fmt.Println("[error]", err.Error())
		return 0
	}
	var sharemem uintptr
	sharemem, _, err = syscall.Syscall(syscall.SYS_SHMAT, shmid, 0, 0)
	if err != 0 {
		fmt.Println("[error]", err.Error())
		return 0
	}
	*dst_ptr = (*T)(unsafe.Pointer(sharemem))
	return shmid

}

// system v share memory close
func Close_Share_Mem(shm uintptr) error {
	_, _, err := syscall.Syscall(syscall.SYS_SHMDT, shm, 0, 0)
	if len(err.Error()) > 0 {
		return fmt.Errorf(err.Error())
	}
	return nil
}

// posix share memory interface
func Shm_Open[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64](pathname string, dst **T) unsafe.Pointer {
	pathinfo := C.CString(pathname)
	var te T
	var code int8
	if sysinfo == "linux" {
		code = 1
	} else {
		code = 0
	}
	ansptr := C.open_shm(pathinfo, (C.int)(unsafe.Sizeof(te)), (C.short)(code))
	if ansptr != nil {
		*dst = (*T)(unsafe.Pointer(ansptr))
	}
	C.free(unsafe.Pointer(pathinfo))
	return ansptr
}

// posix share memory close interface
func Shm_Close(ptr unsafe.Pointer, sizes int) {
	C.unmap(ptr, (C.int)(sizes))
}

// posix share memory delete interface
func Shm_Del(shm_name string) {
	shm_cname := C.CString(shm_name)
	if sysinfo == "linux" {
		C.shm_unlink(shm_cname)
	} else {
		os.Remove(shm_name)
	}
	C.free(unsafe.Pointer(shm_cname))
}

type ShmPosix[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64] interface {
	Open(pathname string, dst **T)
	Close(ptr unsafe.Pointer)
	Unlink(shm_name string)
}

func (s *UnixEn) ShmAt(pathname string, size int64) (unsafe.Pointer, error) {
	shmid, _, err := syscall.Syscall(syscall.SYS_SHMGET, uintptr(C.ftok(C.CString(pathname), C.int(rand.Intn(10)))), uintptr(size), 0100|0640)
	if len(err.Error()) == 0 {
		var mem uintptr
		mem, _, err = syscall.Syscall(syscall.SYS_SHMAT, shmid, 0, 0)
		if len(err.Error()) == 0 {
			ptr := unsafe.Pointer(mem)
			s.resource_map.Store(ptr, mem)
			return ptr, nil
		}
	}
	return nil, err
}
func (s *UnixEn) ShmDt(ptr unsafe.Pointer) {
	mem, ok := s.resource_map.Load(ptr)
	if ok {
		_, _, err := syscall.Syscall(syscall.SYS_SHMDT, mem.(uintptr), 0, 0)
		if len(err.Error()) != 0 {
			fmt.Fprintln(os.Stderr, "close systemv share memory failed", err.Error())
		}
	} else {
		fmt.Fprintln(os.Stderr, "invaild share memory address")
	}
}

func (s *UnixEn) ShmDel(_ unsafe.Pointer) {
	panic("not implemented") // TODO: Implement
}

type str_error string

func (s str_error) Error() string {
	return string(s)
}

//semaphore

const (
	SYSTEMV   = 0x01
	SETVAL    = 0x10
	IPC_RMID  = 0
	IPC_EXCL  = 02000
	IPC_CREAT = 01000
	SEM_UNDO  = 0x1000
)
const (
	SETVAL_DARWIN   = 0x08
	SEM_UNDO_DARWIN = 010000
)

type Semaphore interface {
	Post()
	Wait()
	Close() error
}

type systemv_sem int

func (s systemv_sem) Post() {
	var sembuf C.struct_sembuf
	sembuf.sem_num = 0
	sembuf.sem_op = 1
	sembuf.sem_flg = SEM_UNDO
	syscall.Syscall(syscall.SYS_SEMOP, uintptr(s), uintptr(unsafe.Pointer(&sembuf)), 1)
}

func (s systemv_sem) Wait() {
	var sembuf C.struct_sembuf
	sembuf.sem_num = 0
	sembuf.sem_op = -1
	sembuf.sem_flg = SEM_UNDO
	syscall.Syscall(syscall.SYS_SEMOP, uintptr(s), uintptr(unsafe.Pointer(&sembuf)), 1)
}

func (s systemv_sem) Close() error {
	ok, _, err := syscall.Syscall(syscall.SYS_SEMCTL, uintptr(s), 0, IPC_RMID)
	if ok == 0 {
		return nil
	}
	return err
}
func CreateSem(mod uint8, opts ...any) (sem Semaphore, err error) {
	var (
		mod_ uintptr = 0644
	)
	if len(opts) > 0 {
		for _, ele := range opts {
			switch reflect.TypeOf(ele).Kind() {
			case reflect.Int:
				mod_ = uintptr(ele.(int))
			case reflect.String:

			default:
				err = errors.New("unknow options")
				return
			}
		}
	}
	if mod == SYSTEMV {
		var semid, ok uintptr
		semid, _, err = syscall.Syscall(syscall.SYS_SEMGET, uintptr(C.ftok(C.CString("./"), C.int(8))), 1, mod_|01000)
		if semid > 0 {
			if sysinfo == "darwin" {
				ok, _, err = syscall.Syscall6(syscall.SYS_SEMCTL, semid, 0, SETVAL_DARWIN, 0, 0, 0)
			} else if sysinfo == "linux" {
				ok, _, err = syscall.Syscall6(syscall.SYS_SEMCTL, semid, 0, SETVAL, 0, 0, 0)
			} else {
				err = str_error("don't support system " + sysinfo)
				return
			}
			if ok == 0 {
				err = nil
				sem = systemv_sem(semid)
				return
			}
		}
	} else {
		err = errors.New("unkonw modtype")
	}
	return
}

// ternary expressions
func TernaryExpression[T any](ok bool, left, right T) T {
	if ok {
		return left
	}
	return right
}

// ternary expression for function
func TernaryExpressionFunc[T any](ok bool, left, right T, args ...any) {
	if reflect.TypeOf(left).Kind() != reflect.Func {
		panic("ternary expression func must be function template")
	}
	f := TernaryExpression(ok, left, right)
	fval := reflect.ValueOf(f)
	ftp := reflect.TypeOf(f)
	if ftp.NumIn() != len(args) {
		panic("function args required")
	}
	var (
		arg []reflect.Value = make([]reflect.Value, 0, len(args))
	)
	for i := 0; i < ftp.NumIn(); i++ {
		arg = append(arg, reflect.ValueOf(args[i]).Convert(ftp.In(i)))
	}
	fval.Call(arg)
}
