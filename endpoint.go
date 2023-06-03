package main

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type HealthCheckRequest struct {
}
type HealthCheckResponse struct {
	Status bool `json:"status"`
}
type DiscoverServiceRequest struct {
	ServiceName string `json:"servicename"`
}
type DiscoverServiceResponse struct {
	Instances any    `json:"instances"`
	Error     string `json:"error"`
}

// type GreetRequest struct {
// 	First_name string `json:"first_name"`
// 	Last_name  string `json:"last_name"`
// }
// type GreetResponse struct {
// 	Result any `json:"result"`
// }

//	func MakeGreetEndpoint(srv service) endpoint.Endpoint {
//		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
//			gr := request.(GreetRequest)
//			return GreetResponse{Result: srv.Greet(gr.First_name, gr.Last_name)}, nil
//		}
//	}
func MakeHealthCheckEndpoint(srv normal_service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return HealthCheckResponse{Status: srv.HealthCheck()}, nil
	}
}
func MakeDiscoverServiceEndpoint(srv normal_service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DiscoverServiceRequest)
		ins, err := srv.DiscoverService(req.ServiceName)
		if err != nil {
			return DiscoverServiceResponse{}, err
		}
		return DiscoverServiceResponse{Instances: ins}, nil
	}
}
