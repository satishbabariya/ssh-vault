package store

import (
	"context"
	"io/ioutil"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	connectionURL string
	connection    *pgxpool.Pool
}

func NewStore(connectionURL string) (*Store, error) {
	return &Store{
		connectionURL: connectionURL,
	}, nil
}

func (store *Store) Connect(ctx context.Context) error {
	connection, err := pgxpool.Connect(ctx, store.connectionURL)
	if err != nil {
		return err
	}
	store.connection = connection
	return nil
}

func (store *Store) Close() {
	store.connection.Close()
}

func (store *Store) Ping(ctx context.Context) error {
	return store.connection.Ping(ctx)
}

func (store *Store) AutoMigrate(ctx context.Context) error {

	file, err := ioutil.ReadFile("sql/postgresql.sql")
	if err != nil {
		return err
	}
	sql := string(file)

	_, err = store.connection.Exec(ctx, sql)
	if err != nil {
		return err
	}

	return nil
}

// type Store struct {
// 	db *bolt.DB
// }

// func Open(path string) (*Store, error) {
// 	db, err := bolt.Open(path, 0600, &bolt.Options{
// 		ReadOnly: false,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &Store{
// 		db: db,
// 	}, nil
// }

// func (s *Store) Close() error {
// 	return s.db.Close()
// }

// func (s *Store) Stats() bolt.Stats {
// 	return s.db.Stats()
// }

// func (s *Store) Add(credential model.Credential) error {
// 	return s.db.Update(func(tx *bolt.Tx) error {
// 		b, err := tx.CreateBucketIfNotExists([]byte("credentials"))
// 		if err != nil {
// 			return err
// 		}

// 		if credential.Port == 0 {
// 			credential.Port = 22
// 		}

// 		data, err := json.Marshal(credential)
// 		if err != nil {
// 			return err
// 		}

// 		return b.Put([]byte(credential.Host), data)
// 	})
// }

// func (s *Store) Get(host string) (*model.Credential, error) {
// 	var credential model.Credential

// 	err := s.db.View(func(tx *bolt.Tx) error {
// 		b := tx.Bucket([]byte("credentials"))
// 		if b == nil {
// 			return nil
// 		}

// 		data := b.Get([]byte(host))
// 		if data == nil {
// 			return nil
// 		}

// 		return json.Unmarshal(data, &credential)
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &credential, nil
// }

// func (s *Store) Remotes() ([]model.Remote, error) {
// 	var credentials []model.Credential

// 	err := s.db.View(func(tx *bolt.Tx) error {
// 		b := tx.Bucket([]byte("credentials"))
// 		if b == nil {
// 			return nil
// 		}

// 		return b.ForEach(func(k, v []byte) error {
// 			var credential model.Credential
// 			if err := json.Unmarshal(v, &credential); err != nil {
// 				return err
// 			}

// 			credentials = append(credentials, credential)
// 			return nil
// 		})
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	var remotes []model.Remote
// 	for _, credential := range credentials {
// 		remotes = append(remotes, model.Remote{
// 			Host: credential.Host,
// 			Port: credential.Port,
// 		})
// 	}

// 	return remotes, nil
// }
