package sys

/*
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <pwd.h>
*/
import "C"

type pid_t = int32
type uid_t = int32
type gid_t = int32

func GetPid() pid_t {
	return int32(C.getpid())
}
func GetPPid() pid_t {
	return int32(C.getppid())
}
func GetUid() uid_t {
	return int32(C.getuid())
}
func GetEUid() uid_t {
	return int32(C.geteuid())
}
func GetGid() gid_t {
	return int32(C.getgid())
}
func GetEGid() gid_t {
	return int32(C.getegid())
}
