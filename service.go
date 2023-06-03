package main

type normal_service interface {
	HealthCheck() bool
	DiscoverService(serviceName string) ([]any, error)
}

//	type service interface {
//		Greet(fisrt_name, last_name string) string
//		normal_service
//	}
type ServiceImp struct {
	discoverClient
}

// func (s *ServiceImp) Greet(fisrt_name string, last_name string) string {
// 	return fmt.Sprintf("Hello %s.%s\n", fisrt_name, last_name)
// }

func (s *ServiceImp) HealthCheck() bool {
	return true
}

func (s *ServiceImp) DiscoverService(serviceName string) ([]any, error) {
	return s.DiscoverServices(serviceName)
}
