package httpd

import (
	"fmt"
	"hash/fnv"
	"io"
	"net/http"

	"github.com/fffzlfk/distrikv/db"
)

// Server contains HTTP method handlers to be used for the database
type Server struct {
	db         *db.Database
	shardIndex int
	shardCount int
	addrs      map[int]string
}

// NewServer creates a new Server instance with HTTP handlers
func NewServer(db *db.Database, shardIndex, shardCount int, addrs map[int]string) *Server {
	return &Server{
		db:         db,
		shardIndex: shardIndex,
		shardCount: shardCount,
		addrs:      addrs,
	}
}

func (s *Server) getShard(key string) int {
	h := fnv.New64()
	h.Write([]byte(key))
	return int(h.Sum64() % uint64(s.shardCount))
}

func (s *Server) redirect(w http.ResponseWriter, r *http.Request, shard int) {
	url := "http://" + s.addrs[shard] + r.RequestURI
	fmt.Fprintf(w, "redirecting from shard %d at shard %d\n (%q)\n", s.shardIndex, shard, url)

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
	shard := s.getShard(key)

	if shard != s.shardIndex {
		s.redirect(w, r, shard)
		return
	}

	value, err := s.db.GetKey(key)
	fmt.Fprintf(w, "shard=%d current-shard=%d addr=%q value=%q error = %v\n", shard, s.shardIndex, s.addrs[shard], value, err)
}

// SetHandler
func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")
	shard := s.getShard(key)

	if shard != s.shardIndex {
		s.redirect(w, r, shard)
		return
	}

	err := s.db.SetKey(key, []byte(value))
	fmt.Fprintf(w, "shard=%d current-shard=%d error = %v\n", shard, s.shardIndex, err)
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, nil)
}
