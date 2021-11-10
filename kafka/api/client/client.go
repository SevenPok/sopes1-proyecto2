package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	pb "github.com/sevenpok/api-rabbit/gen/proto"
	"google.golang.org/grpc"
)

type Game struct {
	ID       int    `json:"id"`
	GameName string `json:"gameName"`
	Players  int    `json:"players"`
}

func squidGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	//Extraer valores de los parametros
	vars := mux.Vars(r)
	//Conversion de string a integer
	id, errId := strconv.Atoi(vars["id"])
	players, errPlayers := strconv.Atoi(vars["players"])

	if errId != nil || errPlayers != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error unmarshalling the request"}`))
		return
	}
	//crear struct con los valores
	game := Game{id, vars["gameName"], players}
	w.WriteHeader(http.StatusOK)
	result, _ := json.Marshal(game)

	CreateGame(&game)

	w.Write(result)
}

func server() {
	const port string = "4000"
	router := mux.NewRouter()

	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/game/{id}/gamename/{gameName}/players/{players}", squidGame).Methods("POST")

	log.Println("Server Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my api")
}

func main() {
	server()
}

func CreateGame(game *Game) {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Println(err)
	}
	client_ := pb.NewTestApiClient(conn)

	resp, err := client_.CreateGame(context.Background(), &pb.Game{Id: int64(game.ID), GameName: game.GameName, Players: int64(game.Players)})
	if err != nil {
		log.Println(err)
	}

	fmt.Println(resp)
	fmt.Println(resp.Msg)
}
