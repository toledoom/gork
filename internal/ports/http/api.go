package http

import (
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/toledoom/gork/internal/app/command"
	"github.com/toledoom/gork/internal/app/query"
	"github.com/toledoom/gork/internal/ports/grpc/proto/battle"
	"github.com/toledoom/gork/internal/ports/grpc/proto/leaderboard"
	"github.com/toledoom/gork/internal/ports/grpc/proto/player"
	"github.com/toledoom/gork/pkg/application"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Api struct {
	app *application.App
}

func NewApi(app *application.App) *Api {
	return &Api{
		app: app,
	}
}

func (api *Api) StartBattleHandler(w http.ResponseWriter, r *http.Request) {
	battleID := uuid.New().String()

	httpReq := &battle.StartBattleRequest{}
	startBattleReq, err := decodeHttpRequest[*battle.StartBattleRequest](r, w, httpReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c := &command.StartBattle{
		BattleID:  battleID,
		Player1ID: startBattleReq.PlayerId1,
		Player2ID: startBattleReq.PlayerId2,
	}

	err = application.HandleCommand[*command.StartBattle](api.app, c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (api *Api) FinishBattleHandler(w http.ResponseWriter, r *http.Request) {
	httpReq := &battle.FinishBattleRequest{}
	finishBattleReq, err := decodeHttpRequest[*battle.FinishBattleRequest](r, w, httpReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c := &command.FinishBattle{
		BattleID: finishBattleReq.BattleId,
		WinnerID: finishBattleReq.WinnerId,
	}

	err = application.HandleCommand[*command.FinishBattle](api.app, c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	q := &query.GetBattleResult{
		BattleID: finishBattleReq.BattleId,
	}
	gbrr, err := application.HandleQuery[*query.GetBattleResult, *query.GetBattleResultResponse](api.app, q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := &battle.FinishBattleResponse{
		Player1Score: gbrr.Player1Score,
		Player2Score: gbrr.Player2Score,
	}
	marshalledResp, _ := protojson.Marshal(resp)
	w.Write(marshalledResp)
}

func (api *Api) GetRankHandler(w http.ResponseWriter, r *http.Request) {
	httpReq := &leaderboard.GetRankRequest{}
	getRankReq, err := decodeHttpRequest[*leaderboard.GetRankRequest](r, w, httpReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	q := &query.GetRank{
		PlayerID: getRankReq.PlayerId,
	}

	getRankResponse, err := application.HandleQuery[*query.GetRank, *query.GetRankResponse](api.app, q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := &leaderboard.GetRankResponse{
		Rank: getRankResponse.Rank,
	}
	marshalledResp, _ := protojson.Marshal(resp)
	w.Write(marshalledResp)
}

func (api *Api) GetTopPlayersHandler(w http.ResponseWriter, r *http.Request) {
	httpReq := &leaderboard.GetTopPlayersRequest{}
	getTopPlayersReq, err := decodeHttpRequest[*leaderboard.GetTopPlayersRequest](r, w, httpReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	q := &query.GetTopPlayers{
		NumPlayers: getTopPlayersReq.NumPlayers,
	}

	getTopPlayersResponse, err := application.HandleQuery[*query.GetTopPlayers, *query.GetTopPlayersResponse](api.app, q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var memberList []*leaderboard.Member
	for _, m := range getTopPlayersResponse.MemberList {
		member := &leaderboard.Member{
			Id:    m.PlayerID,
			Score: m.Score,
		}
		memberList = append(memberList, member)
	}
	resp := &leaderboard.GetTopPlayersResponse{
		MemberList: memberList,
	}
	marshalledResp, _ := protojson.Marshal(resp)
	w.Write(marshalledResp)
}

func (api *Api) CreatePlayerHandler(w http.ResponseWriter, r *http.Request) {
	httpReq := &player.CreatePlayerRequest{}
	createPlayerReq, err := decodeHttpRequest[*player.CreatePlayerRequest](r, w, httpReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c := &command.CreatePlayer{
		PlayerID: createPlayerReq.Id,
		Name:     createPlayerReq.Name,
	}

	err = application.HandleCommand[*command.CreatePlayer](api.app, c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	marshalledResp, _ := protojson.Marshal(&player.CreatePlayerResponse{})
	w.Write(marshalledResp)
}

func (api *Api) GetPlayerByIDHandler(w http.ResponseWriter, r *http.Request) {
	httpReq := &player.GetPlayerByIdRequest{}
	getPlayerByIDReq, err := decodeHttpRequest[*player.GetPlayerByIdRequest](r, w, httpReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	q := &query.GetPlayerByID{
		PlayerID: getPlayerByIDReq.Id,
	}

	getPlayerByIDResponse, err := application.HandleQuery[*query.GetPlayerByID, *query.GetPlayerByIDResponse](api.app, q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := &player.GetPlayerByIdResponse{
		Name:  getPlayerByIDResponse.Player.Name,
		Score: getPlayerByIDResponse.Player.Score,
	}
	marshalledResp, _ := protojson.Marshal(resp)
	w.Write(marshalledResp)
}

func decodeHttpRequest[T protoreflect.ProtoMessage](r *http.Request, w http.ResponseWriter, httpReq T) (T, error) {
	defer r.Body.Close()
	body, err := io.ReadAll(io.Reader(r.Body))
	if err != nil {
		return httpReq, err
	}
	err = protojson.Unmarshal(body, httpReq)
	if err != nil {
		return httpReq, err
	}
	return httpReq, err
}
