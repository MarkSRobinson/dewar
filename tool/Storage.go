package tool

import (
	"github.com/op/go-logging"
	"github.com/boltdb/bolt"
	"os"
	"encoding/json"
)

var log = logging.MustGetLogger("storage")

type Institution struct {
	Name string
}

type Storage struct {
	db *bolt.DB
	homedirectory string
}

func (s *Storage) AddRecord(accessToken string, institution Institution) {
	s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		bolB, _ := json.Marshal(institution)
		b.Put([]byte(accessToken), []byte(bolB))
		return nil
	})
}

func (s *Storage) GetRecords() map[string]Institution {
	records := make(map[string]Institution)

	s.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("MyBucket"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			res := Institution{}
			json.Unmarshal(v, &res)
			records[string(k)] = res
		}

		return nil
	})


	return records
}

func (s *Storage) init() error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
		if err != nil {
			log.Errorf("create bucket failure: %s", err)
			return err
		} else {
			log.Debug("create bucket success")
		}
		return err
	})

	return err
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func OpenStorage(directory string) (*Storage, error) {
	homeDir := directory + "/.plaid/"

	err := os.Mkdir(homeDir, os.ModeDir)

	if(err != nil && !os.IsExist(err)) {
		return nil, err
	}

	db, err := bolt.Open(homeDir + "my.db", 0600, nil)
	if err != nil {
		log.Errorf("Could not open DB %#v",err)
		return nil, err
	}
	s := Storage{db, ""}

	err = s.init()


	if ( err != nil) {
		s.Close()
		return nil, err
	}

	return &s, nil
}
