package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"gopkg.in/yaml.v3"

	"github.com/ghkadim/highload_architect/internal/app/mysql"
	"github.com/ghkadim/highload_architect/internal/logger"
	"github.com/ghkadim/highload_architect/internal/models"
)

const (
	syncMode    = "sync"
	cleanupMode = "cleanup"
)

func main() {
	conf := flag.String("config", "", "path to config")
	debug := flag.Bool("debug", false, "enable debug log")
	dryRun := flag.Bool("dry-run", false, "dry run")
	mode := flag.String("mode", syncMode, fmt.Sprintf("available options '%s', '%s'", syncMode, cleanupMode))
	flag.Parse()

	l := logger.Init(*debug)
	defer func() { _ = l.Sync() }()

	if *mode != syncMode && *mode != cleanupMode {
		logger.Fatal("Invalid mode value %s", *mode)
	}

	data, err := os.ReadFile(*conf)
	if err != nil {
		logger.Fatal("Failed to read config from %s: %v", *conf, err)
	}

	var cnf config
	err = yaml.Unmarshal(data, &cnf)
	if err != nil {
		logger.Fatal("Failed to parse config from %s: %v", *conf, err)
	}

	// trap Ctrl+C and call cancel on the context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	r := newResharder(cnf, *dryRun)
	switch *mode {
	case syncMode:
		maxID, err := r.sync(ctx)
		if err != nil {
			logger.Fatal("Failed to sync shards: %v", err)
		}
		logger.Info("Syncing shards completed with max id=%d", maxID)
	case cleanupMode:
		deleted, err := r.cleanup(ctx)
		if err != nil {
			logger.Fatal("Failed to cleanup shards: %v", err)
		}
		logger.Info("Shards cleanup completed with %d rows deleted", deleted)
	default:
		logger.Fatal("Unexpected mode: %s", *mode)
	}
}

type resharder struct {
	clients         map[string]storage
	dedicatedShards map[models.UserID]string
	conf            config
	dryRun          bool
}

func newResharder(conf config, dryRun bool) *resharder {
	return &resharder{
		clients:         make(map[string]storage),
		dedicatedShards: make(map[models.UserID]string),
		conf:            conf,
		dryRun:          dryRun,
	}
}

func (r resharder) sync(ctx context.Context) (models.DialogMessageID, error) {
	maxID := models.DialogMessageID(r.conf.FromID)
	for _, beforeDb := range r.conf.Before {
		beforeCl, err := r.getClient(beforeDb)
		if err != nil {
			return 0, err
		}

		for _, afterDb := range r.conf.After {
			if beforeDb.Address == afterDb.Address {
				continue
			}

			logger.Debug("Syncing shards [%s]->[%s] started", beforeDb.Address, afterDb.Address)
			afterCl, err := r.getClient(afterDb)
			if err != nil {
				return 0, err
			}
			messageCount := 0
			id := models.DialogMessageID(r.conf.FromID)
			for {
				logger.Debug("Fetching [%s] messages from id=%d matching shard '%s'",
					beforeDb.Address, id, afterDb.ShardMatchRegexp)
				messages, err := beforeCl.DialogMatchingShard(ctx, afterDb.ShardMatchRegexp, id, r.conf.Limit)
				if err != nil {
					return 0, err
				}
				if len(messages) == 0 {
					if maxID < id {
						maxID = id
					}
					logger.Info("Syncing [%s]->[%s] completed with max id=%d, %d messages copied", beforeDb.Address, afterDb.Address, id, messageCount)
					break
				}
				messageCount += len(messages)

				logger.Debug("Copying [%s]->[%s] %d messages", beforeDb.Address, afterDb.Address, len(messages))
				if !r.dryRun {
					err = afterCl.DialogBulkInsert(ctx, messages)
					if err != nil {
						return 0, err
					}
				}
				id = messages[len(messages)-1].ID
			}
		}
	}
	return maxID, nil
}

func (r resharder) cleanup(ctx context.Context) (int64, error) {
	var deletedCount int64
	for _, db := range r.conf.After {
		cl, err := r.getClient(db)
		if err != nil {
			return 0, err
		}
		logger.Info("Deleting [%s] messages not matching %s", db.Address, db.ShardMatchRegexp)
		if !r.dryRun {
			count, err := cl.DialogsNotMatchingShardDelete(ctx, db.ShardMatchRegexp)
			if err != nil {
				return deletedCount, err
			}
			logger.Debug("Deleted [%s] %d messages", db.Address, count)
			deletedCount += count
		}
	}
	return deletedCount, nil
}

func (r resharder) getClient(dbCnf db) (storage, error) {
	if cl, ok := r.clients[dbCnf.Address]; ok {
		return cl, nil
	}
	cl, err := mysql.NewStorage(
		dbCnf.User, dbCnf.Password, dbCnf.Address, dbCnf.Database, r.dedicatedShards)
	if err != nil {
		return nil, err
	}
	r.clients[dbCnf.Address] = cl
	return cl, nil
}

type storage interface {
	DialogBulkInsert(ctx context.Context, messages []models.DialogMessage) error
	DialogsNotMatchingShardDelete(ctx context.Context, matchExpr string) (int64, error)
	DialogMatchingShard(ctx context.Context, matchExpr string, fromID models.DialogMessageID, limit int64) ([]models.DialogMessage, error)
}

type db struct {
	User             string `yaml:"user"`
	Password         string `yaml:"password"`
	Address          string `yaml:"address"`
	Database         string `yaml:"database"`
	ShardMatchRegexp string `yaml:"shardMatchRegexp"`
}

type config struct {
	Before []db  `yaml:"before"`
	After  []db  `yaml:"after"`
	FromID int64 `yaml:"fromId"`
	Limit  int64 `yaml:"limit"`
}
