package gapi

import (
	"fmt"

	db "github.com/santinofajardo/simpleBank/db/sqlc"
	"github.com/santinofajardo/simpleBank/pb"
	"github.com/santinofajardo/simpleBank/token"
	"github.com/santinofajardo/simpleBank/util"
	"github.com/santinofajardo/simpleBank/workers"
)

// Server servers gRPC request for our bancking service
type Server struct {
	pb.UnimplementedSimpleBankServer
	config      util.Config
	store       db.Store
	tokenMaker  token.Maker
	distributor workers.TaskDistributor
}

func NewServer(config util.Config, store db.Store, distributor workers.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSecret)
	if err != nil {
		return nil, fmt.Errorf("error creating the token maker: %w", err)
	}
	server := &Server{
		store:       store,
		tokenMaker:  tokenMaker,
		config:      config,
		distributor: distributor,
	}

	return server, nil
}
