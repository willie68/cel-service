package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	"github.com/willie68/cel-service/internal/celproc"
	"github.com/willie68/cel-service/internal/logging"
	log "github.com/willie68/cel-service/internal/logging"
	"github.com/willie68/cel-service/pkg/protofiles"
	"google.golang.org/grpc"
)

var (
	//tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	//certFile   = flag.String("cert_file", "", "The TLS cert file")
	//keyFile    = flag.String("key_file", "", "The TLS key file")
	//jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	grpcport = flag.Int("port", 50051, "The server port")
)

type celServer struct {
	protofiles.UnimplementedEvalServiceServer
}

func (c *celServer) Evaluate(ctx context.Context, req *protofiles.CelRequest) (*protofiles.CelResponse, error) {

	res, err := celproc.GRPCProcCel(req)
	log.Logger.Infof("req: %v, res: %v", req, res)

	if err != nil {
		log.Logger.Errorf("failed to listen: %v", err)
		return nil, err
	}
	return res, nil
}

func init() {
	initLogging()
}

func newServer() *celServer {
	s := &celServer{}
	return s
}

func main() {
	log.Logger.Info("starting")

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *grpcport))
	if err != nil {
		log.Logger.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	/*
		if *tls {
			if *certFile == "" {
				*certFile = data.Path("x509/server_cert.pem")
			}
			if *keyFile == "" {
				*keyFile = data.Path("x509/server_key.pem")
			}
			creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
			if err != nil {
				log.Logger.Fatalf("Failed to generate credentials %v", err)
			}
			opts = []grpc.ServerOption{grpc.Creds(creds)}
		}
	*/
	grpcServer := grpc.NewServer(opts...)
	protofiles.RegisterEvalServiceServer(grpcServer, newServer())
	grpcServer.Serve(lis)

	log.Logger.Info("finished")
}

func initLogging() {
	log.Logger.SetLevel(logging.Debug)
	log.Logger.InitGelf()
}
