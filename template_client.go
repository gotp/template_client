package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	//glog "github.com/golang/glog"
	common "github.com/gotp/proto"
	svrproto "github.com/gotp/proto/template_server"
	nameresolver "github.com/gotp/template_client/service/name_resolver"
)

const (
	address     = "local:///gotp.TemplateServer.TemplateService"
	resolverConfig = "./resolver.conf"
)

func main() {
	flag.Parse()
	// Init resolver
	if nameresolver.GetResolverConfig().Init(resolverConfig) == false {
        log.Fatal("Load resolver config failed!")
    }
    log.Printf("Load resolver config success")
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithBalancerName("round_robin"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := svrproto.NewTemplateServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Test(ctx, &svrproto.TestRequest{
		Header: &common.RequestHeader{
			RequestId: "R00001",
			ClientId: "C00001",
			ClientType: common.ClientType_H5Client,
			Version: "V1",
			TestFlag: true,
		},
		Data: &svrproto.TestRequestData{
			Dummy: 1, 
		},
	})
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("rpc success: %d %s", r.Header.Retcode, r.Header.Retmsg)
}
