package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	glog "github.com/golang/glog"
	commonproto "online_consultant/proto"
	proto "online_consultant/proto/template_server"
	nameresolver "online_consultant/server/template_client/service/name_resolver"
)

const (
	address     = "local:///OnlineConsultant.TemplateServer.TemplateService"
	resolverConfig = "./resolver.conf"
)

func main() {
	flag.Parse()
	// Init resolver
	if nameresolver.GetResolverConfig().Init(resolverConfig) == false {
        glog.Fatal("Load resolver config failed!")
    }
    glog.Info("Load resolver config success")
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithBalancerName("round_robin"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := proto.NewTemplateServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Test(ctx, &proto.TestRequest{
		Header: &commonproto.RequestHeader{
			RequestId: "R00001",
			ClientId: "C00001",
			ClientType: commonproto.ClientType_H5Client,
			Version: "V1",
			TestFlag: true,
		},
		Data: &proto.TestRequestData{
			Dummy: 1, 
		},
	})
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("rpc success: %d %s", r.Header.Retcode, r.Header.Retmsg)
}