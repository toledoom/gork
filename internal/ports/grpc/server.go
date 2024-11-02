package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/toledoom/gork/internal/app/usecases"
	"github.com/toledoom/gork/internal/ports/grpc/proto/battle"
	"github.com/toledoom/gork/internal/ports/grpc/proto/leaderboard"
	"github.com/toledoom/gork/internal/ports/grpc/proto/player"
	"github.com/toledoom/gork/pkg/gork"
)

type GameServer struct {
	app *gork.App

	battle.UnimplementedBattleServer
	leaderboard.UnimplementedLeaderboardServer
	player.UnimplementedPlayerServer
}

func NewGameServer(app *gork.App) *GameServer {
	return &GameServer{
		app: app,
	}
}

func (s *GameServer) StartBattle(ctx context.Context, sbr *battle.StartBattleRequest) (*battle.StartBattleResponse, error) {
	battleID := uuid.New().String()
	input := usecases.StartBattleInput{
		BattleID:  battleID,
		Player1ID: sbr.PlayerId1,
		Player2ID: sbr.PlayerId2,
	}
	_, err := gork.ExecuteUseCase[usecases.StartBattleInput, usecases.StartBattleOutput](s.app, input)
	if err != nil {
		return nil, err
	}

	return &battle.StartBattleResponse{
		BattleId: battleID,
	}, nil
}

func (s *GameServer) FinishBattle(ctx context.Context, fbr *battle.FinishBattleRequest) (*battle.FinishBattleResponse, error) {
	input := usecases.FinishBattleInput{
		BattleID: fbr.BattleId,
		WinnerID: fbr.WinnerId,
	}
	ucOutput, err := gork.ExecuteUseCase[usecases.FinishBattleInput, usecases.FinishBattleOutput](s.app, input)
	if err != nil {
		return nil, err
	}

	return &battle.FinishBattleResponse{
		Player1Score: ucOutput.Player1Score,
		Player2Score: ucOutput.Player2Score,
	}, nil
}

func (s *GameServer) GetRank(ctx context.Context, grr *leaderboard.GetRankRequest) (*leaderboard.GetRankResponse, error) {
	input := usecases.GetRankInput{
		PlayerID: grr.PlayerId,
	}

	output, err := gork.ExecuteUseCase[usecases.GetRankInput, usecases.GetRankOutput](s.app, input)
	if err != nil {
		return nil, err
	}

	return &leaderboard.GetRankResponse{
		Rank: output.Rank,
	}, err
}

func (s *GameServer) GetTopPlayers(ctx context.Context, gtp *leaderboard.GetTopPlayersRequest) (*leaderboard.GetTopPlayersResponse, error) {
	input := usecases.GetTopPlayersInput{
		NumPlayers: gtp.NumPlayers,
	}

	output, err := gork.ExecuteUseCase[usecases.GetTopPlayersInput, usecases.GetTopPlayersOutput](s.app, input)
	if err != nil {
		return nil, err
	}

	var protoMemberList []*leaderboard.Member
	for _, m := range output.MemberList {
		protoMember := &leaderboard.Member{
			Id:    m.PlayerID,
			Score: m.Score,
		}
		protoMemberList = append(protoMemberList, protoMember)
	}

	return &leaderboard.GetTopPlayersResponse{
		MemberList: protoMemberList,
	}, err
}

func (s *GameServer) CreatePlayer(ctx context.Context, cpr *player.CreatePlayerRequest) (*player.CreatePlayerResponse, error) {
	input := usecases.CreatePlayerInput{
		PlayerID: cpr.Id,
		Name:     cpr.Name,
	}

	_, err := gork.ExecuteUseCase[usecases.CreatePlayerInput, usecases.CreatePlayerOutput](s.app, input)
	if err != nil {
		return nil, err
	}

	return &player.CreatePlayerResponse{}, err
}

func (s *GameServer) GetPlayerById(ctx context.Context, cpr *player.GetPlayerByIdRequest) (*player.GetPlayerByIdResponse, error) {
	input := usecases.GetPlayerByIDInput{
		PlayerID: cpr.Id,
	}

	output, err := gork.ExecuteUseCase[usecases.GetPlayerByIDInput, usecases.GetPlayerByIDOutput](s.app, input)
	if err != nil {
		return nil, err
	}

	return &player.GetPlayerByIdResponse{
		Name:  output.Player.Name,
		Score: output.Player.Score,
	}, err
}
