package main

import (
	"fmt"
	"log"
	"net"

	"github.com/toledoom/gork/internal/app"
	grpcport "github.com/toledoom/gork/internal/ports/grpc"
	"github.com/toledoom/gork/internal/ports/grpc/proto/battle"
	"github.com/toledoom/gork/internal/ports/grpc/proto/leaderboard"
	"github.com/toledoom/gork/internal/ports/grpc/proto/player"
	"github.com/toledoom/gork/pkg/gork"
	"google.golang.org/grpc"
)

func main() {
	a := gork.NewApp(app.SetupUseCases, app.SetupCommandHandlers, app.SetupQueryHandlers)
	a.Start(app.SetupServices)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 50051))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gameServer := grpcport.NewGameServer(a)
	grpcServer := grpc.NewServer()

	battle.RegisterBattleServer(grpcServer, gameServer)
	leaderboard.RegisterLeaderboardServer(grpcServer, gameServer)
	player.RegisterPlayerServer(grpcServer, gameServer)
	log.Print("Server started")
	grpcServer.Serve(lis)
}
