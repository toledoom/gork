package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/toledoom/gork/internal/app/command"
	"github.com/toledoom/gork/internal/app/query"
	"github.com/toledoom/gork/internal/ports/grpc/proto/battle"
	"github.com/toledoom/gork/internal/ports/grpc/proto/leaderboard"
	"github.com/toledoom/gork/internal/ports/grpc/proto/player"
	"github.com/toledoom/gork/pkg/application"
)

type GameServer struct {
	app *application.App

	battle.UnimplementedBattleServer
	leaderboard.UnimplementedLeaderboardServer
	player.UnimplementedPlayerServer
}

func NewGameServer(app *application.App) *GameServer {
	return &GameServer{
		app: app,
	}
}

func (s *GameServer) StartBattle(ctx context.Context, sbr *battle.StartBattleRequest) (*battle.StartBattleResponse, error) {
	battleID := uuid.New().String()
	c := &command.StartBattle{
		BattleID:  battleID,
		Player1ID: sbr.PlayerId1,
		Player2ID: sbr.PlayerId2,
	}

	err := s.app.HandleCommand(c)
	if err != nil {
		return nil, err
	}

	return &battle.StartBattleResponse{
		BattleId: battleID,
	}, nil
}

func (s *GameServer) FinishBattle(ctx context.Context, fbr *battle.FinishBattleRequest) (*battle.FinishBattleResponse, error) {
	c := &command.FinishBattle{
		BattleID: fbr.BattleId,
		WinnerID: fbr.WinnerId,
	}

	err := s.app.HandleCommand(c)
	if err != nil {
		return nil, err
	}

	return &battle.FinishBattleResponse{
		Player1Score: 0, // TODO: Get the score using a query
		Player2Score: 0, // TODO: Get the score using a query
	}, nil
}

func (s *GameServer) GetRank(ctx context.Context, grr *leaderboard.GetRankRequest) (*leaderboard.GetRankResponse, error) {
	q := &query.GetRank{
		PlayerID: grr.PlayerId,
	}

	response, err := s.app.HandleQuery(q)
	if err != nil {
		return nil, err
	}
	getRankResponse := response.(*query.GetRankResponse)

	return &leaderboard.GetRankResponse{
		Rank: getRankResponse.Rank,
	}, err
}

func (s *GameServer) GetTopPlayers(ctx context.Context, gtp *leaderboard.GetTopPlayersRequest) (*leaderboard.GetTopPlayersResponse, error) {
	q := &query.GetTopPlayers{
		NumPlayers: gtp.NumPlayers,
	}
	response, err := s.app.HandleQuery(q)
	if err != nil {
		return nil, err
	}
	getTopPlayersResponse := response.(*leaderboard.GetTopPlayersResponse)
	return &leaderboard.GetTopPlayersResponse{
		MemberList: getTopPlayersResponse.MemberList,
	}, err
}

func (s *GameServer) CreatePlayer(ctx context.Context, cpr *player.CreatePlayerRequest) (*player.CreatePlayerResponse, error) {
	c := &command.CreatePlayer{
		PlayerID: cpr.Id,
		Name:     cpr.Name,
	}
	err := s.app.HandleCommand(c)
	if err != nil {
		return nil, err
	}

	return &player.CreatePlayerResponse{}, err
}

func (s *GameServer) GetPlayerById(ctx context.Context, cpr *player.GetPlayerByIdRequest) (*player.GetPlayerByIdResponse, error) {
	q := &query.GetPlayerByID{
		PlayerID: cpr.Id,
	}
	queryResponse, err := s.app.HandleQuery(q)
	if err != nil {
		return nil, err
	}
	getPlayerByIDResponse := queryResponse.(*player.GetPlayerByIdResponse)

	return &player.GetPlayerByIdResponse{
		Name:  getPlayerByIDResponse.Name,
		Score: getPlayerByIDResponse.Score,
	}, err
}
