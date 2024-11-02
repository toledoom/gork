package http

import (
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/toledoom/gork/internal/app/usecases"
	"github.com/toledoom/gork/internal/ports/grpc/proto/battle"
	"github.com/toledoom/gork/internal/ports/grpc/proto/leaderboard"
	"github.com/toledoom/gork/internal/ports/grpc/proto/player"
	"github.com/toledoom/gork/pkg/gork"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Api struct {
	app *gork.App
}

func NewApi(app *gork.App) *Api {
	return &Api{
		app: app,
	}
}

func (api *Api) StartBattleHandler(w http.ResponseWriter, r *http.Request) {
	battleID := uuid.New().String()

	httpReq := &battle.StartBattleRequest{}
	startBattleReq, err := decodeHttpRequest[*battle.StartBattleRequest](r, httpReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := usecases.StartBattleInput{
		BattleID:  battleID,
		Player1ID: startBattleReq.PlayerId1,
		Player2ID: startBattleReq.PlayerId2,
	}
	_, err = gork.ExecuteUseCase[usecases.StartBattleInput, usecases.StartBattleOutput](api.app, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp := &battle.StartBattleResponse{
		BattleId: battleID,
	}
	marshalledResp, _ := protojson.Marshal(resp)
	w.WriteHeader(http.StatusCreated)
	w.Write(marshalledResp)
}

func (api *Api) FinishBattleHandler(w http.ResponseWriter, r *http.Request) {
	httpReq := &battle.FinishBattleRequest{}
	finishBattleReq, err := decodeHttpRequest[*battle.FinishBattleRequest](r, httpReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := usecases.FinishBattleInput{
		BattleID: finishBattleReq.BattleId,
		WinnerID: finishBattleReq.WinnerId,
	}
	ucOutput, err := gork.ExecuteUseCase[usecases.FinishBattleInput, usecases.FinishBattleOutput](api.app, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := &battle.FinishBattleResponse{
		Player1Score: ucOutput.Player1Score,
		Player2Score: ucOutput.Player2Score,
	}
	marshalledResp, _ := protojson.Marshal(resp)
	w.Write(marshalledResp)
}

func (api *Api) GetRankHandler(w http.ResponseWriter, r *http.Request) {
	httpReq := &leaderboard.GetRankRequest{}
	getRankReq, err := decodeHttpRequest[*leaderboard.GetRankRequest](r, httpReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := usecases.GetRankInput{
		PlayerID: getRankReq.PlayerId,
	}

	output, err := gork.ExecuteUseCase[usecases.GetRankInput, usecases.GetRankOutput](api.app, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := &leaderboard.GetRankResponse{
		Rank: output.Rank,
	}
	marshalledResp, _ := protojson.Marshal(resp)
	w.Write(marshalledResp)
}

func (api *Api) GetTopPlayersHandler(w http.ResponseWriter, r *http.Request) {
	httpReq := &leaderboard.GetTopPlayersRequest{}
	getTopPlayersReq, err := decodeHttpRequest[*leaderboard.GetTopPlayersRequest](r, httpReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := usecases.GetTopPlayersInput{
		NumPlayers: getTopPlayersReq.NumPlayers,
	}

	output, err := gork.ExecuteUseCase[usecases.GetTopPlayersInput, usecases.GetTopPlayersOutput](api.app, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var memberList []*leaderboard.Member
	for _, m := range output.MemberList {
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
	createPlayerReq, err := decodeHttpRequest[*player.CreatePlayerRequest](r, httpReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := usecases.CreatePlayerInput{
		PlayerID: createPlayerReq.Id,
		Name:     createPlayerReq.Name,
	}

	_, err = gork.ExecuteUseCase[usecases.CreatePlayerInput, usecases.CreatePlayerOutput](api.app, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	marshalledResp, _ := protojson.Marshal(&player.CreatePlayerResponse{})
	w.WriteHeader(http.StatusCreated)
	w.Write(marshalledResp)
}

func (api *Api) GetPlayerByIDHandler(w http.ResponseWriter, r *http.Request) {
	httpReq := &player.GetPlayerByIdRequest{}
	getPlayerByIDReq, err := decodeHttpRequest[*player.GetPlayerByIdRequest](r, httpReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := usecases.GetPlayerByIDInput{
		PlayerID: getPlayerByIDReq.Id,
	}

	output, err := gork.ExecuteUseCase[usecases.GetPlayerByIDInput, usecases.GetPlayerByIDOutput](api.app, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := &player.GetPlayerByIdResponse{
		Name:  output.Player.Name,
		Score: output.Player.Score,
	}
	marshalledResp, _ := protojson.Marshal(resp)
	w.Write(marshalledResp)
}

func decodeHttpRequest[T protoreflect.ProtoMessage](r *http.Request, httpReq T) (T, error) {
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
