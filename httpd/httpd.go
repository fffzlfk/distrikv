package httpd

import (
	"fmt"
	"io"
	"net/http"

	"github.com/fffzlfk/distrikv/config"
	"github.com/fffzlfk/distrikv/db"
)

// Server contains HTTP method handlers to be used for the database
type Server struct {
	db     *db.Database
	shards *config.Shards
}

// NewServer creates a new Server instance with HTTP handlers
func NewServer(db *db.Database, shards *config.Shards) *Server {
	return &Server{
		db:     db,
		shards: shards,
	}
}

func (s *Server) redirect(w http.ResponseWriter, r *http.Request, shard int) {
	url := "http://" + s.shards.Addrs[shard] + r.RequestURI
	fmt.Fprintf(w, "redirecting from shard %d at shard %d\n (%q)\n", s.shards.Index, shard, url)

	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error redirecting the request: %v", err)
		return
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}

// GetHandler
func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	shard := s.shards.GetIndex(key)

	if shard != s.shards.Index {
		s.redirect(w, r, shard)
		return
	}

	value, err := s.db.GetKey(key)
	fmt.Fprintf(w, "shard=%d current-shard=%d addr=%q value=%q error = %v\n", shard, s.shards.Index, s.shards.Addrs[shard], value, err)
}

// SetHandler
func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")
	shard := s.shards.GetIndex(key)

	if shard != s.shards.Index {
		s.redirect(w, r, shard)
		return
	}

	err := s.db.SetKey(key, []byte(value))
	fmt.Fprintf(w, "shard=%d current-shard=%d addr=%q error = %v\n", shard, s.shards.Index, s.shards.Addrs[shard], err)
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, nil)
}
