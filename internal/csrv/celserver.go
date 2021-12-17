package csrv

import (
	"context"

	"github.com/willie68/cel-service/internal/celproc"
	log "github.com/willie68/cel-service/internal/logging"
	"github.com/willie68/cel-service/pkg/protofiles"
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

func NewCelServer() *celServer {
	s := &celServer{}
	return s
}
