package server

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
)

type Server struct{ 
	sockPath string 
	listener net.Listener
}

func NewServer(sockPath string) *Server { 
	return &Server{ 
		sockPath: sockPath,
	}
}

func (s *Server) Start() error { 
	os.Remove(s.sockPath)

	l, err := net.Listen("unix", s.sockPath)
	if err != nil { 
		return err
	}
	s.listener = l 

	os.Chmod(s.sockPath, 0666)

	mux := http.NewServeMux()
	mux.HandleFunc("/status", s.handleStatus)

	go http.Serve(s.listener, mux)
	return nil
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) { 
	stats := map[string]any{
		"status": "running",
		"utun":   "utun0",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (s *Server) Stop() { 
	if s.listener != nil { 
		s.listener.Close()
	}
	os.Remove(s.sockPath)
}