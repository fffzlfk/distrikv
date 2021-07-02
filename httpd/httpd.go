package httpd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fffzlfk/distrikv/config"
	"github.com/fffzlfk/distrikv/db"
	"github.com/fffzlfk/distrikv/replication"
	"github.com/fffzlfk/distrikv/utils"
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
	// fmt.Fprintf(w, "redirecting from shard %d at shard %d\n (%q)\n", s.shards.Index, shard, url)

	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error redirecting the request: %v", err)
		return
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}

// GetHandler get the value of key
func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	shard := s.shards.GetIndex(key)

	if shard != s.shards.Index {
		s.redirect(w, r, shard)
		return
	}

	value, err := s.db.GetKey(key)
	resp := &utils.Resp{
		Shard:    shard,
		CurShard: s.shards.Index,
		Addr:     s.shards.Addrs[shard],
		Value:    string(value),
		Err:      err,
	}
	json.NewEncoder(w).Encode(resp)
	// fmt.Fprintf(w, "shard=%d current-shard=%d addr=%q value=%q error = %v\n", shard, s.shards.Index, s.shards.Addrs[shard], value, err)
}

// SetHandler puts key-values to db
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
	resp := &utils.Resp{
		Shard:    shard,
		CurShard: s.shards.Index,
		Addr:     s.shards.Addrs[shard],
		Err:      err,
	}
	json.NewEncoder(w).Encode(resp)
	// fmt.Fprintf(w, "shard=%d current-shard=%d addr=%q error = %v\n", shard, s.shards.Index, s.shards.Addrs[shard], err)
}

// DeleteHandler deletes key-values to db
func (s *Server) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	shard := s.shards.GetIndex(key)

	if shard != s.shards.Index {
		s.redirect(w, r, shard)
		return
	}

	err := s.db.DeleteKey(key)
	resp := &utils.Resp{
		Shard:    shard,
		CurShard: s.shards.Index,
		Addr:     s.shards.Addrs[shard],
		Err:      err,
	}
	json.NewEncoder(w).Encode(resp)
}

// DeleteExtraKeysHandler
func (s *Server) DeleteExtraKeysHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Fprintf(w, "Error = %v", s.db.DeleteExtraKeys(func(key string) bool {
		return s.shards.GetIndex(key) != s.shards.Index
	}))
}

func (s *Server) GetNextForReplicationHandler(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	k, v, err := s.db.GetNextForReplication()
	enc.Encode(replication.NextKeyValue{
		Key:   string(k),
		Value: string(v),
		Err:   err,
	})
}

func (s *Server) DeleteReplicationKeyHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")

	if err := s.db.DeleteReplicationKey([]byte(key), []byte(value)); err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		fmt.Fprint(w, "error:", err)
		return
	}
	fmt.Fprint(w, "ok")
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, nil)
}
