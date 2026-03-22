package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	pb "github.com/ricocynthia/botanica/proto"
	"github.com/ricocynthia/botanica/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Botanica API — gRPC backend + HTTP/REST wrapper (BFF pattern)
// Same architecture used at Alaska Airlines: HTTP requests from the
// frontend are translated into gRPC calls to backend services.

var grpcClient pb.BotanicaServiceClient

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// GET /
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		respondError(w, http.StatusNotFound, "route not found")
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"name":         "Botanica API",
		"description":  "A herbal remedies and foraging API by Cynthia Rico Cook — Earthy Mujer",
		"version":      "1.0.0",
		"architecture": "gRPC backend + HTTP/REST wrapper (BFF pattern)",
		"endpoints": map[string]string{
			"GET /remedies":               "List all remedies. Filter by ?type=tea or ?property=sleep",
			"GET /remedies/{id}":          "Get a remedy by ID",
			"GET /ingredients":            "List all unique remedy ingredients",
			"GET /forageables":            "List all plants and mushrooms. Filter by ?category=Plant or ?property=immune",
			"GET /forageables/{id}":       "Get a forageable by ID",
		},
	})
}

// GET /remedies and GET /remedies/{id}
func remediesHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) == 2 && parts[1] != "" {
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid remedy ID")
			return
		}
		remedy, err := grpcClient.GetRemedy(context.Background(), &pb.GetRemedyRequest{Id: int32(id)})
		if err != nil {
			respondError(w, http.StatusNotFound, "remedy not found")
			return
		}
		respondJSON(w, http.StatusOK, remedy)
		return
	}
	resp, err := grpcClient.GetRemedies(context.Background(), &pb.GetRemediesRequest{
		Type:     r.URL.Query().Get("type"),
		Property: r.URL.Query().Get("property"),
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch remedies")
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// GET /ingredients
func ingredientsHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := grpcClient.GetIngredients(context.Background(), &pb.GetIngredientsRequest{})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch ingredients")
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// GET /forageables and GET /forageables/{id}
func forageablesHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) == 2 && parts[1] != "" {
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid forageable ID")
			return
		}
		forageable, err := grpcClient.GetForageable(context.Background(), &pb.GetForageableRequest{Id: int32(id)})
		if err != nil {
			respondError(w, http.StatusNotFound, "forageable not found")
			return
		}
		respondJSON(w, http.StatusOK, forageable)
		return
	}
	resp, err := grpcClient.GetForageables(context.Background(), &pb.GetForageablesRequest{
		Category: r.URL.Query().Get("category"),
		Property: r.URL.Query().Get("property"),
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch forageables")
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

func main() {
	// Start gRPC server
	grpcPort := ":50051"
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterBotanicaServiceServer(grpcServer, &server.BotanicaServer{})
	log.Printf("🌿 Botanica gRPC server starting on %s", grpcPort)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Connect HTTP (BFF) layer to gRPC server
	conn, err := grpc.NewClient("localhost"+grpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()
	grpcClient = pb.NewBotanicaServiceClient(conn)

	// Start HTTP server
	httpPort := ":8080"
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/remedies", remediesHandler)
	mux.HandleFunc("/remedies/", remediesHandler)
	mux.HandleFunc("/ingredients", ingredientsHandler)
	mux.HandleFunc("/forageables", forageablesHandler)
	mux.HandleFunc("/forageables/", forageablesHandler)

	log.Printf("🍵 Botanica HTTP server starting on http://localhost%s", httpPort)
	log.Fatal(http.ListenAndServe(httpPort, mux))
}