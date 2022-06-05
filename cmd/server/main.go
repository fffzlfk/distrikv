package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/fffzlfk/distrikv/config"
	"github.com/fffzlfk/distrikv/db"

	"github.com/fffzlfk/distrikv/httpd"
	"github.com/fffzlfk/distrikv/replica"
)

var (
	dbLocation     = flag.String("db-location", "", "the path to the bolt db database")
	httpAddr       = flag.String("http-addr", "", "set-addr")
	configFileName = flag.String("config-file", "sharding.toml", "set-config-file")
	shard          = flag.String("shard", "", "select the shard")
	isReplica      = flag.Bool("replica", false, "whether or not run as a replica")
)

func init() {
	flag.Parse()
	if *httpAddr == "" {
		log.Fatal("Must provide http-addr")
	}

	if *dbLocation == "" {
		log.Fatal("Must provide db-location")
	}

	if *shard == "" {
		log.Fatal("Must provide shard")
	}
}

func main() {
	cfg, err := config.ParseFile(*configFileName)
	if err != nil {
		log.Fatal(err)
	}

	shards, err := config.ParseShards(cfg.Shards, *shard)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Shard count = %d, current shard: %d\n", shards.Count, shards.Index)

	db, close, err := db.NewDatabase(*dbLocation, *isReplica)
	if err != nil {
		log.Fatalf("NewDataBase(%q): %v", *dbLocation, err)
	}
  defer func() {
    err := close()
    log.Fatal(err)
  }()

	// replication
	if *isReplica {
		masterAddrs, has := shards.Addrs[shards.Index]
		if !has {
			log.Fatal("master dose not exist:", err)
		}
		go replica.ClientLoop(db, masterAddrs, replica.Replication)
		go replica.ClientLoop(db, masterAddrs, replica.Deleted)
	}

	server := httpd.NewServer(db, shards)

	http.HandleFunc("/ping", server.PingHandler)

	http.HandleFunc("/get", server.GetHandler)

	http.HandleFunc("/set", server.SetHandler)

	http.HandleFunc("/delete", server.DeleteHandler)

	http.HandleFunc("/purge", server.DeleteExtraKeysHandler)

	http.HandleFunc("/next-replication-key", server.GetNextForReplicationHandler)

	http.HandleFunc("/delete-replication-key", server.DeleteReplicationKeyHandler)

	http.HandleFunc("/next-deleted-key", server.GetNextForDeletedHandler)

	http.HandleFunc("/delete-deleted-key", server.DeleteDeletedKeyHandler)

	// hash(key) % count = <current index>

	log.Fatal(server.ListenAndServe(*httpAddr))
}
