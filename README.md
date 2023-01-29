# commonlib

# etcd grpc client
```go
package main

import (
	"context"
	"fmt"
	"github.com/newgate8820/commonlib/etcd_server"
	"github.com/newgate8820/commonlib/testlib/testrpc/client/protocol/helloworld"
	"github.com/newgate8820/commonlib/transport"
	"go.uber.org/zap"
	"google.golang.org/grpc/resolver"
	"time"
)

const App = "base-hello"

func main() {

	addrs := []string{"etcd-server:2379"}
	r := etcd_server.NewResolver(addrs, zap.NewNop())
	resolver.Register(r)

	conn := transport.NewGrpcConn("etcd:///"+App, 10)

	helloClient := helloworld.NewGreeterClient(conn)
	tk := time.NewTicker(time.Second * 3)
	for range tk.C {
		res, err := helloClient.SayHello(context.Background(), &helloworld.HelloRequest{
			Name: "xiaoming",
		})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	}

}
```
# server node1
```go
package main

import (
	"fmt"
	"github.com/newgate8820/commonlib/etcd_server"
	"github.com/newgate8820/commonlib/testlib/testrpc/node1/protocol/helloworld"
	"github.com/newgate8820/commonlib/testlib/testrpc/node1/service"
	"github.com/newgate8820/commonlib/transport"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"

	"time"
)

const (
	app         = "base-hello"
	grpcAddress = "node1"
	port        = ":8084"
	etcdAddress = "etcd-server:2379"
)

func main() {
	addrs := []string{etcdAddress}
	etcdRegister := etcd_server.NewRegister(addrs, zap.NewNop())
	node := etcd_server.Server{
		Name: app,
		Addr: grpcAddress + port,
	}

	server, err := Start()
	if err != nil {
		panic(fmt.Sprintf("start server failed : %v", err))
	}

	if _, err := etcdRegister.Register(node, 10); err != nil {
		panic(fmt.Sprintf("server register failed: %v", err))
	}

	fmt.Println("service started listen on", grpcAddress)

	select {}
	server.Stop()
	etcdRegister.Stop()
	time.Sleep(time.Second)
}

func Start() (*grpc.Server, error) {
	s := transport.NewGrpcServer()

	helloworld.RegisterGreeterServer(s, &service.Service{})
	reflection.Register(s)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()

	return s, nil
}


```

# server node2
```go
package main

import (
	"fmt"
	"github.com/newgate8820/commonlib/etcd_server"
	"github.com/newgate8820/commonlib/testlib/testrpc/node2/protocol/helloworld"
	"github.com/newgate8820/commonlib/testlib/testrpc/node2/service"
	"github.com/newgate8820/commonlib/transport"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"

	"time"
)

const (
	app         = "base-hello"
	grpcAddress = "node2"
	port        = ":8085"
	etcdAddress = "etcd-server:2379"
)

func main() {
	addrs := []string{etcdAddress}
	etcdRegister := etcd_server.NewRegister(addrs, zap.NewNop())
	node := etcd_server.Server{
		Name: app,
		Addr: grpcAddress + port,
	}

	server, err := Start()
	if err != nil {
		panic(fmt.Sprintf("start server failed : %v", err))
	}

	if _, err := etcdRegister.Register(node, 10); err != nil {
		panic(fmt.Sprintf("server register failed: %v", err))
	}

	fmt.Println("service started listen on", grpcAddress)
	select {}
	server.Stop()
	etcdRegister.Stop()
	time.Sleep(time.Second)
}

func Start() (*grpc.Server, error) {
	s := transport.NewGrpcServer()

	helloworld.RegisterGreeterServer(s, &service.Service{})
	reflection.Register(s)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()

	return s, nil
}

```