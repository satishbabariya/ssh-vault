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
	// Open a PostgreSQL database.
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// maintains a pool of idle connections. To maximize pool performance
	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	pgdb.SetMaxOpenConns(maxOpenConns)
	pgdb.SetMaxIdleConns(maxOpenConns)

	// create db
	db := bun.NewDB(pgdb, pgdialect.New(), bun.WithDiscardUnknownColumns())

	// Print all queries to stdout.
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

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

func (store *Store) IdentityExists(ctx context.Context, github_id string) (bool, error) {
	exists, err := store.db.NewSelect().Model((*model.Identity)(nil)).Where("github_id = ?", github_id).Exists(ctx)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (store *Store) CreateIdentity(ctx context.Context, github_id string) error {
	identity := &model.Identity{
		GithubID: github_id,
	}
	_, err := store.db.NewInsert().Model(identity).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
