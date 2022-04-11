package store

import (
	"context"
	"database/sql"
	"runtime"
	"ssh-vault/internal/model"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type Store struct {
	db *bun.DB
}

func NewStore(dsn string) (*Store, error) {

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// maintains a pool of idle connections. To maximize pool performance
	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	sqldb.SetMaxOpenConns(maxOpenConns)
	sqldb.SetMaxIdleConns(maxOpenConns)

	// create db
	db := bun.NewDB(sqldb, pgdialect.New(), bun.WithDiscardUnknownColumns())

	// log queries
	db.AddQueryHook(bundebug.NewQueryHook())
	bundebug.NewQueryHook(bundebug.WithVerbose(true))

	return &Store{db: db}, nil
}

func (store *Store) CreateTable(ctx context.Context, model interface{}) error {
	res, err := store.db.NewCreateTable().IfNotExists().Model(model).Exec(ctx)
	if err != nil {
		return err
	}

	if _, err := res.RowsAffected(); err != nil {
		return err
	}

	return nil
}

func (store *Store) Init(ctx context.Context) error {
	models := []interface{}{
		(*model.Identity)(nil),
		(*model.Permission)(nil),
		(*model.Remote)(nil),
		(*model.Credential)(nil),
	}

	// register models
	for _, model := range models {
		store.db.RegisterModel(model)
	}

	// create tables
	for _, model := range models {
		if err := store.CreateTable(ctx, model); err != nil {
			return err
		}
	}

	return nil
}

func (store *Store) Close() error {
	return store.db.Close()
}
