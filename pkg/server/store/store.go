package store

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"

	_ "github.com/go-sql-driver/mysql"
	"github.com/satishbabariya/vault/pkg/server/model"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Store struct {
	db *bun.DB
}

func NewStore(database string, dsn string) (*Store, error) {
	if database == "postgres" {

		// check for valid dsn
		// TODO: fix first path segment in URL cannot contain colon goroutine 1 [running]:
		// Create PR to fix this
		// /go/pkg/mod/github.com/uptrace/bun/driver/pgdriver@v1.1.3/config.go line:195
		// export function ParseDSN(dsn string)
		// if _, err := pgdriver.WithDSN(dsn); err != nil {

		// }

		// Open a PostgreSQL database.
		pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

		// maintains a pool of idle connections. To maximize pool performance
		maxOpenConns := 4 * runtime.GOMAXPROCS(0)
		pgdb.SetMaxOpenConns(maxOpenConns)
		pgdb.SetMaxIdleConns(maxOpenConns)

		// create db
		db := bun.NewDB(pgdb, pgdialect.New(), bun.WithDiscardUnknownColumns())

		// Print all queries to stdout.
		// db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

		return &Store{db: db}, nil
	} else if database == "mysql" {
		// Open a MySQL database.
		sqldb, err := sql.Open("mysql", dsn)
		if err != nil {
			return nil, err
		}

		// maintains a pool of idle connections. To maximize pool performance
		maxOpenConns := 4 * runtime.GOMAXPROCS(0)
		sqldb.SetMaxOpenConns(maxOpenConns)
		sqldb.SetMaxIdleConns(maxOpenConns)

		// create db
		db := bun.NewDB(sqldb, mysqldialect.New())

		// Print all queries to stdout.
		// db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

		return &Store{db: db}, nil
	} else {
		return nil, fmt.Errorf("database %s not supported", database)
	}
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

func (store *Store) GetIdentity(ctx context.Context, github_id string) (*model.Identity, error) {
	identity := &model.Identity{}
	err := store.db.NewSelect().Model(identity).Where("github_id = ?", github_id).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return identity, nil
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

func (store *Store) ListRemotes(ctx context.Context, ids []int64) ([]*model.Remote, error) {
	remotes := []*model.Remote{}
	err := store.db.NewSelect().Model(&remotes).Where("id IN (?)", ids).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return remotes, nil
}

func (store *Store) GetPermissions(ctx context.Context, identity *model.Identity) ([]*model.Permission, error) {
	permissions := []*model.Permission{}
	err := store.db.NewSelect().Model(&permissions).Where("identity_id = ?", identity.ID).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}
