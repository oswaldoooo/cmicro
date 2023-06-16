## **CMIRCO**
### **console_api**
```go
// example
cl, err := console_api.NewClient("origin", "localhost", 9000)
if err != nil {
    fmt.Println("connect to console failed")
    os.Exit(1)
}
err = client.RegisterService("serviceName", "serviceId", console_api.RELOAD, &console_api.Address{Add: "localhost", Port: 8000, TypeName: "tcp"})
if err != nil && err != io.EOF {
    fmt.Println(err.Error())
    os.Exit(1)
}


//deregister service from console
client.DeregisterService("serviceId")
```