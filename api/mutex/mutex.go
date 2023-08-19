package mutex

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"

	mutex "github.com/oswaldoooo/cmicro/pkg/mutex"
	"google.golang.org/grpc"
)

var client mutex.MutexClient
var NULLCLIENT = fmt.Errorf("client is not create")

func DialMutexCloudBackGround(host string, port int) error {
	clients, err := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	if err == nil {
		clients.Connect()
		client = mutex.NewMutexClient(clients)
	}
	return err
}
func GetMutex() (int32, error) {
	if client == nil {
		return -1, NULLCLIENT
	}
	info, err := client.Getlock(context.Background(), &mutex.MutexInfo{})
	if err == nil {
		return info.MutexId, nil
	}
	return -1, err
}
func TryLock(mutexid int32) (bool, error) {
	if client == nil {
		return false, NULLCLIENT
	}
	info, err := client.Lock(context.Background(), &mutex.MutexInfo{MutexId: mutexid})
	if err == nil {
		return info.Ok, nil
	}
	return false, err
}
func Unlock(mutexid int32) error {
	if client == nil {
		return NULLCLIENT
	}
	_, err := client.Unlock(context.Background(), &mutex.MutexInfo{MutexId: mutexid})
	return err
}
func Release_Lock(mutexid int32) error {
	if client == nil {
		return NULLCLIENT
	}
	info, err := client.Freelock(context.Background(), &mutex.MutexInfo{})
	if err == nil {
		if !info.Ok {
			err = fmt.Errorf("free mutex %d failed", mutexid)
		}
	}
	return err
}
func IsLock(mutexid int32) (bool, error) {
	if client == nil {
		return false, NULLCLIENT
	}
	info, err := client.Freelock(context.Background(), &mutex.MutexInfo{})
	if err == nil {
		return info.Ok, nil
	}
	return false, err
}

/*
锁id 列表 .....
初始化10把锁
锁住

*/

var mutex_list []mutex_mini = make([]mutex_mini, 10)
var available_mutex []int16 = []int16{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

type MutexCloud struct {
	mtx sync.Mutex
}

func (s *MutexCloud) Islock(ctx context.Context, muxinfo *mutex.MutexInfo) (*mutex.Response, error) {
	defer func() {
		if ans := recover(); ans != nil {
			log.Println("[panic error]", ans)
		}
	}()
	var resp mutex.Response
	if muxinfo.GetMutexId() < int32(len(mutex_list)) {
		resp.Ok = mutex_list[muxinfo.MutexId].islock
	}
	return &resp, nil
}

func (s *MutexCloud) Lock(ctx context.Context, muxinfo *mutex.MutexInfo) (*mutex.Response, error) {
	defer func() {
		if ans := recover(); ans != nil {
			log.Println("[panic error]", ans)
		}
	}()
	var resp mutex.Response
	log.Printf("accept lock request %d,%v", muxinfo.MutexId, mutex_list[muxinfo.MutexId].islock)
	if muxinfo.GetMutexId() < int32(len(mutex_list)) && !mutex_list[muxinfo.GetMutexId()].islock {
		mutex_list[muxinfo.GetMutexId()].islock = true
		resp.Ok = true
	} else {
		resp.Ok = false
	}
	return &resp, nil
}
func (s *MutexCloud) Unlock(ctx context.Context, muxinfo *mutex.MutexInfo) (*mutex.Response, error) {
	var resp mutex.Response
	log.Printf("accept unlock request %d,%v", muxinfo.MutexId, mutex_list[muxinfo.MutexId].islock)
	if muxinfo.GetMutexId() < int32(len(mutex_list)) && mutex_list[muxinfo.GetMutexId()].islock {
		mutex_list[muxinfo.GetMutexId()].islock = false
		resp.Ok = true
	} else {
		resp.Ok = false
	}
	return &resp, nil
}

func (s *MutexCloud) Getlock(ctx context.Context, muxinfo *mutex.MutexInfo) (*mutex.MutexInfo, error) {
	var info mutex.MutexInfo
	if len(available_mutex) > 0 {
		info.MutexId = int32(available_mutex[0])
		available_mutex = available_mutex[1:]
	} else {
		info.MutexId = -1
	}
	return &info, nil
}

func (s *MutexCloud) Freelock(ctx context.Context, muxinfo *mutex.MutexInfo) (*mutex.Response, error) {
	var resp mutex.Response
	resp.Ok = false
	if muxinfo.MutexId < int32(len(mutex_list)) {
		exist_mutex := false
		for _, mid := range available_mutex {
			if mid == int16(muxinfo.MutexId) {
				exist_mutex = true
				break
			}
		}
		if !exist_mutex {
			available_mutex = append(available_mutex, int16(muxinfo.MutexId))
			resp.Ok = true
		}
	}
	return &resp, nil
}

type mutex_mini struct {
	mtx    sync.Mutex //修改锁的状态，全靠这把锁
	islock bool
}
