package main

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/go-kit/kit/sd/consul"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

type discoverClient interface {
	Register(serviceName, instanceid, registeradd string, registerport int, helathcheckurl string) bool
	Deregister(instanceId string) bool
	DiscoverServices(serviceName string) ([]any, error)
}
type Client struct {
	ConsulAdd      string
	ConsulPort     int
	HealthCheckUrl string
	consul.Client
	config      consulapi.Config
	instanceMap sync.Map
	mux         sync.Mutex
}

func (s *Client) Register(serviceName string, instanceid string, registeradd string, registerport int, helathcheckurl string) bool {
	info := consulapi.AgentServiceRegistration{
		Name:    serviceName,
		ID:      instanceid,
		Address: registeradd,
		Port:    registerport,
		Meta:    nil,
		Check:   &consulapi.AgentServiceCheck{HTTP: "http://" + registeradd + ":" + strconv.Itoa(registerport) + helathcheckurl, DeregisterCriticalServiceAfter: "30s", Interval: "15s"},
	}
	err := s.Client.Register(&info)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err == nil
}

func (s *Client) Deregister(instanceId string) bool {
	err := s.Client.Deregister(&consulapi.AgentServiceRegistration{ID: instanceId})
	return err == nil
}

func (s *Client) DiscoverServices(serviceName string) ([]any, error) {
	instances, ok := s.instanceMap.Load(serviceName)
	if ok {
		return instances.([]any), nil

	}
	s.mux.Lock()
	instances, ok = s.instanceMap.Load(serviceName)
	if ok {
		return instances.([]any), nil

	}
	go func() {
		//使用consul服务实例监控某个服务名服务实例列表变化
		params := make(map[string]any)
		params["type"] = "service"
		params["service"] = serviceName
		plan, _ := watch.Parse(params)
		plan.Handler = func(u uint64, i interface{}) {
			if i == nil {
				return
			}
			v, ok := i.([]*consulapi.ServiceEntry)
			if !ok {
				return
			}
			if len(v) == 0 { //没有服务实例在线
				s.instanceMap.Store(serviceName, []any{})
			}
			var healthService []any
			for _, service := range v {
				if service.Checks.AggregatedStatus() == consulapi.HealthPassing {
					healthService = append(healthService, service.Service)
				}
			}
			s.instanceMap.Store(serviceName, healthService)
		}
		defer plan.Stop()
		plan.Run(s.config.Address)
	}()
	defer s.mux.Unlock()
	entries, _, err := s.Client.Service(serviceName, "", false, nil)
	if err == nil && len(entries) > 0 {
		instances := make([]any, len(entries))
		if len(entries) == 0 {
			return nil, errors.New("entries is null")
		}
		for i := 0; i < len(entries); i++ {
			instances[i] = entries[i].Service
		}
		s.instanceMap.Store(serviceName, instances)
		return instances, nil
	} else {
		s.instanceMap.Store(serviceName, []any{})
		return nil, err
	}
}

func NewClient(ConsulAdd string, ConsulPort int) *Client {
	config := consulapi.Config{Address: ConsulAdd + ":" + strconv.Itoa(ConsulPort)}
	cl, err := consulapi.NewClient(&config)
	if err == nil {
		return &Client{Client: consul.NewClient(cl), config: config}
	}
	return nil
}
