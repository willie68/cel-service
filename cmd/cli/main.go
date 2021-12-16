package main

import (
	"context"
	"flag"
	"time"

	"github.com/willie68/cel-service/internal/logging"
	log "github.com/willie68/cel-service/internal/logging"
	"github.com/willie68/cel-service/pkg/protofiles"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.example.com", "The server name used to verify the hostname returned by the TLS handshake")
)

func runEvaluate(client protofiles.EvalServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	jsonContext := map[string]interface{}{
		"data": map[string]interface{}{
			"value": 1,
		},
	}
	structValue, err := structpb.NewStruct(jsonContext)
	if err != nil {
		log.Logger.Fatalf("structpb new: %v", client, err)
	}

	celRequest := protofiles.CelRequest{
		Context:    structValue,
		Expression: "int(data.value) == 1",
	}

	celResponse, err := client.Evaluate(ctx, &celRequest)
	if err != nil {
		log.Logger.Fatalf("%v.Evaluate(_) = _, %v: ", client, err)
	}
	log.Logger.Infof("res: %v", celResponse)
}

func init() {
	initLogging()
}

func main() {
	log.Logger.Info("starting")

	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Logger.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := protofiles.NewEvalServiceClient(conn)

	// evaluate
	runEvaluate(client)

	log.Logger.Info("finished")
}

func initLogging() {
	log.Logger.SetLevel(logging.Debug)
	log.Logger.InitGelf()
}
