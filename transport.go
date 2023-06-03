package main

import (
	"context"
	"encoding/json"
	"net/http"
)

const (
	CONTENT_TYPE   = "Content-Type"
	TransPort_Type = "application/json"
)

// func MakeHttpHandler(ctx context.Context) http.Handler {
// 	router := mux.NewRouter()
// 	newroute := router.Methods("GET")
// 	srv := ServiceImp{}
// 	healend := MakeHealthCheckEndpoint(&srv)
// 	discend := MakeDiscoverServiceEndpoint(&srv)
// 	newroute.Path("/health").Handler(kithttp.NewServer(
// 		healend,
// 		decodeHealthCheck,
// 		encodeAnything,
// 	))
// 	newroute.Path("/discoverservice").Handler(kithttp.NewServer(discend, decodeDiscoverService, encodeAnything))

// 	return router
// }

//	func decodeGreet(ctx context.Context, req *http.Request) (any, error) {
//		return GreetRequest{}, nil
//	}
func decodeHealthCheck(ctx context.Context, req *http.Request) (any, error) {
	return HealthCheckRequest{}, nil
}
func decodeDiscoverService(ctx context.Context, req *http.Request) (any, error) {
	return DiscoverServiceRequest{}, nil
}
func encodeAnything(ctx context.Context, repw http.ResponseWriter, rep any) error {
	repw.Header().Set(CONTENT_TYPE, TransPort_Type)
	return json.NewEncoder(repw).Encode(rep)
}
