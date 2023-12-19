package aio
import (
    "encoding/json"
    "errors"
    "fmt"
    "os"
    "syscall"
    "unsafe"
)

type Chunk struct{
    fd int
    b []byte
    realsize int64
}
func mapchunk(path string,chunksize int64)(fd int,b []byte,err error){
    //var fd int
    fd,err=syscall.Open(path,syscall.O_CREAT|syscall.O_RDWR,0640);
    if err==nil&&fd>0{
    //  defer syscall.Close(fd)
        err=syscall.Ftruncate(fd,chunksize)
        if err==nil{
            b,err=syscall.Mmap(fd,0,int(chunksize),syscall.PROT_READ|syscall.PROT_WRITE,syscall.MAP_SHARED)
        }
        if err!=nil{
            syscall.Close(fd)
            fd=0
        }
    }
    if fd<=0{
        err=errors.New("open "+path+" failed")
    }
    return
}
func(s *Chunk)unmap()error{
    syscall.Syscall(syscall.SYS_MSYNC,uintptr(unsafe.Pointer(&s.b[0])),syscall.MS_SYNC,0)
    err:=syscall.Munmap(s.b)
    if err==nil{
        s.realsize=0
        s.b=nil
    }
    return err
}
func (s *Chunk)Close()error{
    var err error
    err=syscall.Close(s.fd)
    if err==nil{

    }
    err=s.unmap()
    return err
}
//write to chunk;off_t is start
func (s *Chunk)Marshal(fun func(any)([]byte,error),v any,off_t uint)(length int,err error){
    var b []byte
    b,err=fun(v)
    if err==nil{
        if int(off_t)+len(b)>len(s.b){
            err=errors.New("out of range")
            return
        }
        copy(s.b[off_t:int(off_t)+len(b)],b)
        length=len(b)
    }
    return
}
func OpenChunk(filename string,size int64)*Chunk{
    var ans =new(Chunk)
    var err error
    ans.fd,ans.b,err=mapchunk(filename,size)
    if err==nil{
        return ans
    }else{
        fmt.Fprintln(os.Stderr,err.Error())
    }
    return nil
}
func CStrLen(p []byte)int{
    var ans int=0
    for _,e:=range p{
        if e!=0{
            ans++
        }else{
            break
        }
    }
  return ans
  }
func (s *Chunk)ReMap(newsize int64)error{
    err:=s.unmap()
    if err==nil{
        err=syscall.Ftruncate(s.fd,newsize)
        if err==nil{
            s.b,err=syscall.Mmap(s.fd,0,int(newsize),syscall.PROT_READ|syscall.PROT_WRITE,syscall.MAP_SHARED)
        }
    }
    return err
}
