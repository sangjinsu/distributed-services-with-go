package server

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ProductRequest struct {
	Record Record `json:"record"`
}

type ProductResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
	Record Record `json:"record"`
}

type httpServer struct {
	Log *Log
}

func NewHttpServer(addr string) *http.Server {
	httpSvr := newHttpServer()
	r := http.NewServeMux()
	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/" {
			http.NotFound(writer, request)
			return
		}

		if request.Method == "GET" {
			httpSvr.handleConsume(writer, request)
			return
		}

		if request.Method == "POST" {
			httpSvr.handleProduce(writer, request)
			return
		}

		if request.Method != "POST" && request.Method != "GET" {
			http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
	})
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

func newHttpServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

func (s *httpServer) handleProduce(w http.ResponseWriter, r *http.Request) {
	var req ProductRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	offset, err := s.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := ProductResponse{Offset: offset}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	var req ConsumeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	record, err := s.Log.Read(req.Offset)
	if errors.Is(err, ErrOffsetNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := ConsumeResponse{Record: record}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
