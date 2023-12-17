package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/toledoom/gork/internal/app/command"
	"github.com/toledoom/gork/internal/app/query"
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
	c := &command.StartBattle{
		BattleID:  battleID,
		Player1ID: sbr.PlayerId1,
		Player2ID: sbr.PlayerId2,
	}

	err := gork.HandleCommand[*command.StartBattle](s.app, c)
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

	err := gork.HandleCommand[*command.FinishBattle](s.app, c)
	if err != nil {
		return nil, err
	}

	q := &query.GetBattleResult{
		BattleID: fbr.BattleId,
	}
	gbrr, err := gork.HandleQuery[*query.GetBattleResult, *query.GetBattleResultResponse](s.app, q)
	if err != nil {
		return nil, err
	}

	return &battle.FinishBattleResponse{
		Player1Score: gbrr.Player1Score,
		Player2Score: gbrr.Player2Score,
	}, nil
}

func (s *GameServer) GetRank(ctx context.Context, grr *leaderboard.GetRankRequest) (*leaderboard.GetRankResponse, error) {
	q := &query.GetRank{
		PlayerID: grr.PlayerId,
	}

	getRankResponse, err := gork.HandleQuery[*query.GetRank, *query.GetRankResponse](s.app, q)
	if err != nil {
		return nil, err
	}

	return &leaderboard.GetRankResponse{
		Rank: getRankResponse.Rank,
	}, err
}

func (s *GameServer) GetTopPlayers(ctx context.Context, gtp *leaderboard.GetTopPlayersRequest) (*leaderboard.GetTopPlayersResponse, error) {
	q := &query.GetTopPlayers{
		NumPlayers: gtp.NumPlayers,
	}
	getTopPlayersResponse, err := gork.HandleQuery[*query.GetTopPlayers, *query.GetTopPlayersResponse](s.app, q)
	if err != nil {
		return nil, err
	}

	var protoMemberList []*leaderboard.Member
	for _, m := range getTopPlayersResponse.MemberList {
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
	c := &command.CreatePlayer{
		PlayerID: cpr.Id,
		Name:     cpr.Name,
	}
	err := gork.HandleCommand[*command.CreatePlayer](s.app, c)
	if err != nil {
		return nil, err
	}

	return &player.CreatePlayerResponse{}, err
}

func (s *GameServer) GetPlayerById(ctx context.Context, cpr *player.GetPlayerByIdRequest) (*player.GetPlayerByIdResponse, error) {
	q := &query.GetPlayerByID{
		PlayerID: cpr.Id,
	}
	getPlayerByIDResponse, err := gork.HandleQuery[*query.GetPlayerByID, *query.GetPlayerByIDResponse](s.app, q)
	if err != nil {
		return nil, err
	}

	return &player.GetPlayerByIdResponse{
		Name:  getPlayerByIDResponse.Player.Name,
		Score: getPlayerByIDResponse.Player.Score,
	}, err
}
