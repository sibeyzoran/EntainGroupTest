package main

import (
	"database/sql"
	"flag"
	"log"
	"net"

	"github.com/sibeyzoran/entaingrouptest/sports/db"
	"github.com/sibeyzoran/entaingrouptest/sports/service"
	"google.golang.org/grpc"
)

var (
	grpcEndpoint = flag.String("grpc-endpoint", "localhost:7000", "gRPC server endpoint")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("failed running grpc server: %s\n", err)
	}
}

func run() error {
	conn, err := net.Listen("tcp", ":7000")
	if err != nil {
		return err
	}

	sportsDB, err := sql.Open("sqlite3", "./db/sports.db")
	if err != nil {
		return err
	}

	sportsRepo := db.NewRacesRepo(sportsDB)
	if err := sportsRepo.Init(); err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	racing.RegisterRacingServer(
		grpcServer,
		service.NewRacingService(
			sportsRepo,
		),
	)

	log.Printf("gRPC server listening on: %s\n", *grpcEndpoint)

	if err := grpcServer.Serve(conn); err != nil {
		return err
	}

	return nil
}
