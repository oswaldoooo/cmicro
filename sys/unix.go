package sys

/*
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
	"fmt"
	"os"
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
