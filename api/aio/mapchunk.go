package aio
import (
        "errors"
        "fmt"
        "os"
        "syscall"
        "unsafe"
        "io"
)

type Chunk struct{
        fd int
        fdcloser io.Closer
        b []byte
        realsize int64
}
func mapchunk(fd int,chunksize int64)(b []byte,err error){
        err=syscall.Ftruncate(fd,chunksize)
        if err==nil{
                b,err=syscall.Mmap(fd,0,int(chunksize),syscall.PROT_READ|syscall.PROT_WRITE,syscall.MAP_SHARED)
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
        if s.fdcloser==nil{
                err=syscall.Close(s.fd)
                if err==nil{

                }
        }else{
                err=s.fdcloser.Close()
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
        //var err error
        //ans.fd,ans.b,err=mapchunk(filename,size)
        f,err:=Open(filename,syscall.O_CREAT|syscall.O_RDWR)
        if err==nil{
                ans.b,err=mapchunk(f.Fd(),size)
                if err==nil{
                        ans.fd=f.Fd()
                        ans.fdcloser=f
                        return ans
                }
                f.Close()
        }
        if err!=nil{
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
