package store

import (
	"encoding/json"

	"github.com/satishbabariya/vault/pkg/types"
	bolt "go.etcd.io/bbolt"
)

type Store struct {
	db *bolt.DB
}

func Open(path string) (*Store, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{
		ReadOnly: false,
	})
	if err != nil {
		return nil, err
	}

	return &Store{
		db: db,
	}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Stats() bolt.Stats {
	return s.db.Stats()
}

func (s *Store) Add(credential *types.Credential) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("credentials"))
		if err != nil {
			return err
		}

		if credential.Port == 0 {
			credential.Port = 22
		}

		data, err := json.Marshal(credential)
		if err != nil {
			return err
		}

		return b.Put([]byte(credential.Host), data)
	})
}

func (s *Store) Get(host string) (*types.Credential, error) {
	var credential types.Credential

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("credentials"))
		if b == nil {
			return nil
		}

		data := b.Get([]byte(host))
		if data == nil {
			return nil
		}

		return json.Unmarshal(data, &credential)
	})

	if err != nil {
		return nil, err
	}

	return &credential, nil
}

func (s *Store) Credentials() ([]types.Credential, error) {
	var credentials []types.Credential

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("credentials"))
		if b == nil {
			return nil
		}

		return b.ForEach(func(k, v []byte) error {
			var credential types.Credential
			if err := json.Unmarshal(v, &credential); err != nil {
				return err
			}

			credentials = append(credentials, credential)
			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	return credentials, nil
}
