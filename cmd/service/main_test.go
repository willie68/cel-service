package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	log "github.com/willie68/cel-service/internal/logging"
	"github.com/willie68/cel-service/pkg/model"
	"github.com/willie68/cel-service/pkg/protofiles"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	serverAddr = "localhost:50051"
	client     protofiles.EvalServiceClient
	conn       *grpc.ClientConn
)

func TestGRPCJson(t *testing.T) {
	go initServer()
	initClient(t)

	ast := assert.New(t)

	celModels := readJson("../../test/data/data1.json", t)

	for _, cm := range celModels {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		grpcContext, err := structpb.NewStruct(cm.Request.Context)
		ast.Nil(err)
		celRequest := protofiles.CelRequest{
			Context:    grpcContext,
			Expression: cm.Request.Expression,
		}

		result, err := client.Evaluate(ctx, &celRequest)
		ast.Nil(err)
		ast.NotNil(result)

		ast.Equal(cm.Result, result.Result)
	}
	closeClient()
}

func readJson(filename string, t *testing.T) []model.TestCelModel {
	ast := assert.New(t)
	ya, err := ioutil.ReadFile(filename)
	ast.Nil(err)
	var celModels []model.TestCelModel
	err = json.Unmarshal(ya, &celModels)
	ast.Nil(err)
	return celModels
}

func initServer() {
	lis, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Logger.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	grpcServer = grpc.NewServer(opts...)
	protofiles.RegisterEvalServiceServer(grpcServer, newCelServer())
	grpcServer.Serve(lis)
}

func initClient(t *testing.T) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	var err error
	conn, err = grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Logger.Fatalf("fail to dial: %v", err)
	}
	client = protofiles.NewEvalServiceClient(conn)
}

func closeClient() {
	conn.Close()
	grpcServer.Stop()
}
