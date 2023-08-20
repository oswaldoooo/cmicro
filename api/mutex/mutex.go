package mutex

import (
	"context"
	"fmt"
	"strconv"

	"github.com/oswaldoooo/cmicro/pkg/mutex"
	"google.golang.org/grpc"
)

var client mutex.MutexClient
var rwclient mutex.RwmutexClient
var NULLCLIENT = fmt.Errorf("client is not create")

func DialMutexCloudBackGround(host string, port int) error {
	clients, err := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	if err == nil {
		clients.Connect()
		client = mutex.NewMutexClient(clients)
		rwclient = mutex.NewRwmutexClient(clients)
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
	info, err := client.Freelock(context.Background(), &mutex.MutexInfo{MutexId: mutexid})
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
	info, err := client.Islock(context.Background(), &mutex.MutexInfo{MutexId: mutexid})
	if err == nil {
		return info.Ok, nil
	}
	return false, err
}
func TryRLock(mutexid int32) (bool, error) {
	if rwclient == nil {
		return false, NULLCLIENT
	}
	info, err := rwclient.Rlock(context.Background(), &mutex.MutexInfo{MutexId: mutexid})
	if err == nil {
		return info.Ok, nil
	}
	return false, err
}
func RUnlock(mutexid int32) error {
	if rwclient == nil {
		return NULLCLIENT
	}
	_, err := rwclient.Runlock(context.Background(), &mutex.MutexInfo{MutexId: mutexid})
	return err
}
func TryWLock(mutexid int32) (bool, error) {
	if rwclient == nil {
		return false, NULLCLIENT
	}
	info, err := rwclient.Lock(context.Background(), &mutex.MutexInfo{MutexId: mutexid})
	if err == nil {
		return info.Ok, nil
	}
	return false, err
}
func UnWlock(mutexid int32) error {
	if rwclient == nil {
		return NULLCLIENT
	}
	_, err := rwclient.Unlock(context.Background(), &mutex.MutexInfo{MutexId: mutexid})
	return err
}

func GetRWMutex() (int32, error) {
	if rwclient == nil {
		return -1, NULLCLIENT
	}
	info, err := rwclient.Getlock(context.Background(), &mutex.MutexInfo{})
	if err == nil {
		return info.MutexId, nil
	}
	return -1, err
}

func Release_RWLock(mutexid int32) error {
	if rwclient == nil {
		return NULLCLIENT
	}
	info, err := rwclient.Freelock(context.Background(), &mutex.MutexInfo{MutexId: mutexid})
	if err == nil {
		if !info.Ok {
			err = fmt.Errorf("free mutex %d failed", mutexid)
		}
	}
	return err
}

// [rislock,wislock]
func RWIsLock(mutexid int32) ([2]bool, error) {
	ans := [2]bool{false, false}
	if rwclient == nil {
		return ans, NULLCLIENT
	}
	info, err := rwclient.Islock(context.Background(), &mutex.MutexInfo{MutexId: mutexid})
	if err == nil {
		ans[0] = info.Rok
		ans[1] = info.Ok
		return ans, nil
	}
	return ans, err
}
