package http

import (
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/toledoom/gork/internal/app/command"
	"github.com/toledoom/gork/internal/ports/grpc/proto/battle"
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

// TODO: Implement the actual handlers
// Create the http requests and responses for each handler
// ////////////////////////////////////
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

	err = api.app.HandleCommand(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (api *Api) FinishBattleHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(data)
}

func (api *Api) GetRankHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(data)
}

func (api *Api) GetTopPlayersHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(data)
}

func (api *Api) CreatePlayerHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(data)
}

func (api *Api) GetPlayerByIDHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(data)
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
