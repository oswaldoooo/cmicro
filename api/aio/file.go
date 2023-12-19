package aio
import (
        "errors"
        "os"
        "sync"
        "syscall"
)
var(
        ErrFileLock=errors.New("file locked")
        ErrEmptyName=errors.New("emtpy name")
)
//open file with lock
type File struct{
        mux sync.Mutex
        filename string
        fd int
}
func Open(filename string,mode int,perms ...uint32)(*File,error){
        if len(filename)==0{
                return nil,ErrEmptyName
        }
        var perm uint32=0644
        if len(perms)>0{
                perm=perms[0]
        }
        var fd int
        fd2,err:=syscall.Open(filename+".mutex",syscall.O_CREAT|syscall.O_EXCL,0640)
        if err!=nil{
                return nil,ErrFileLock
        }
        syscall.Close(fd2)
        //fmt.Println("create file lock success")
        fd,err=syscall.Open(filename,mode,perm)
        if err==nil{
                return &File{filename:filename,fd:fd},nil
        }
        return nil,err
}
func (s *File)Fd()int{
        return s.fd
}
func (s *File)Close()error{
        var errlist []error=make([]error,0,2)
        err:=syscall.Close(s.fd)
        if err!=nil{
                errlist=append(errlist,err)
        }
        err=os.Remove(s.filename+".mutex")
        if len(errlist)>0{
                return errors.Join(errlist...)
        }
        return nil
}
