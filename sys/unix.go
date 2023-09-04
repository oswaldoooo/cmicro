package sys

/*
#include <sys/shm.h>
#include <stdlib.h>
static int shmflag=IPC_CREAT|0640;
*/
import "C"
import (
	"fmt"
	"syscall"
	"unsafe"
)

func GetShare_Mem[T int | int32 | int64](shmid int, dst_ptr **T) uintptr {
	shm, _, err := syscall.Syscall(syscall.SYS_SHMAT, uintptr(shmid), 0, 0)
	if len(err.Error()) < 1 {
		*dst_ptr = (*T)(unsafe.Pointer(shm))
		return 0
	}
	return shm
}
func CreateShare_Mem[T int | int32 | int64](pathname string, gendid int, dst_ptr **T) uintptr {
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
func Close_Share_Mem(shm uintptr) error {
	_, _, err := syscall.Syscall(syscall.SYS_SHMDT, shm, 0, 0)
	if len(err.Error()) > 0 {
		return fmt.Errorf(err.Error())
	}
	return nil
}
