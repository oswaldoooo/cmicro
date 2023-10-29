package aio

import (
	"bytes"
	"crypto/sha256"
	"io"
	"os"
	"runtime"
	"sync"
	"syscall"
)

var (
	//design by system
	DEFAULT_BUFFER_SIZE           = os.Getpagesize()
	OUTOFRANGE          str_error = "out of range"
)

// write content to buffer, and wirte to disk when it close
type Aio struct {
	buffer, sha_val []byte
	data_len        uint64
	fd              int
	linkname        string
	off             uint64
	off_lock        sync.Mutex
	lock            sync.RWMutex
	ReadEn          func([]byte) (int, error)
}
type str_error string

func (s *str_error) Error() string {
	return string(*s)
}
func (s *Aio) SetSeek(off int64) error {
	runtime.KeepAlive(s.fd)
	s.lock.Lock()
	defer s.lock.Unlock()
	if uint64(off) >= s.data_len {
		return &OUTOFRANGE
	}
	s.off = uint64(off)
	return nil
}
func (s *Aio) Read(p []byte) (n int, err error) {
	runtime.KeepAlive(s.fd)
	s.off_lock.Lock()
	defer s.off_lock.Unlock()
	s.lock.RLock()
	defer s.lock.RUnlock()
	var datalen uint64 = uint64(len(p))
	if s.off >= s.data_len {
		return 0, &OUTOFRANGE
	}
	if datalen > s.data_len-s.off { //data len is readable length
		datalen = s.data_len - s.off
	}
	copy(p[s.off:s.off+datalen], s.buffer[:datalen])
	s.off += datalen
	n = int(datalen)
	return
}

func (s *Aio) Write(p []byte) (n int, err error) {
	runtime.KeepAlive(s.fd)
	s.lock.Lock()
	defer s.lock.Unlock()
	var datalen uint64 = uint64(len(p))
	if datalen+s.data_len > uint64(len(s.buffer)) {
		err = &OUTOFRANGE
		return
	}
	copy(s.buffer[s.data_len:s.data_len+datalen], p)
	s.data_len += datalen
	n = int(datalen)
	return
}

func (s *Aio) Close() error {
	var (
		strerr  str_error
		content []byte
		err     error
	)
	content, err = os.ReadFile(s.linkname[1 : len(s.linkname)-5])
	if err == nil { //compare buffer whether change
		hsh := sha256.New()
		hsh.Reset()
		_, err = hsh.Write(s.buffer[:s.data_len])
		if err == nil {
			if bytes.Equal(s.sha_val, hsh.Sum(nil)) {
				goto closestart
			}
			hsh.Reset()
			_, err = hsh.Write(content)
			if err == nil {
				if !bytes.Equal(s.sha_val, hsh.Sum(nil)) {
					strerr += "content changed by other process;"
				}
			}
		}
	}
	if err != nil {
		strerr += str_error(err.Error())
	}
	if len(strerr) == 0 {
		//write buffer to disk
		_, err = syscall.Seek(s.fd, 0, io.SeekStart)
		if err == nil {
			_, err = syscall.Write(s.fd, s.buffer[:s.data_len])
		}
		if err != nil {
			strerr += str_error(err.Error())
		}
	}
closestart:
	runtime.KeepAlive(s.fd)
	err = syscall.Close(s.fd)
	if err != nil {
		strerr += str_error(err.Error())
	}
	err = syscall.Unlink(s.linkname)
	if err != nil {
		strerr += str_error(err.Error())
	}
	if len(strerr) > 0 {
		return &strerr
	}
	return nil
}

func OpenFile(pathname string, flag int, perm os.FileMode) (*Aio, error) {
	var (
		err      error
		linkname string = "." + pathname + ".swap"
	)

	err = syscall.Link(pathname, linkname)
	if err == nil {
		var fd int
		fd, err = syscall.Open(linkname, flag, uint32(perm))
		if err == nil {
			aioptr := &Aio{fd: fd, buffer: make([]byte, DEFAULT_BUFFER_SIZE), linkname: linkname}
			var lang int
			lang, err = syscall.Read(fd, aioptr.buffer)
			if err == nil {
				aioptr.data_len = uint64(lang)
				//make sha sign
				hsh := sha256.New()
				hsh.Reset()
				_, err = hsh.Write(aioptr.buffer[:lang])
				if err == nil {
					aioptr.sha_val = hsh.Sum(nil)
				}
			}
			if err != nil {
				syscall.Close(fd)
				aioptr = nil
			}
			return aioptr, err
		} else {
			syscall.Unlink(linkname)
		}
	}
	return nil, err
}
