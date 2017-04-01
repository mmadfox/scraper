package scraper

import "github.com/boltdb/bolt"

var bucketName []byte = []byte("urls")

type boltdbVisits struct {
	db *bolt.DB
}

func (v *boltdbVisits) fix(u string) string {
	if len(u) == 0 {
		return "__#__"
	}
	return u
}

func (v *boltdbVisits) Visit(u string) bool {
	u = v.fix(u)
	var ok bool = false
	err := v.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		v := b.Get([]byte(u))
		if len(v) > 0 {
			ok = true
			return nil
		}
		return b.Put([]byte(u), []byte("1"))
	})
	if err != nil {
		panic(err)
	}
	return ok
}

func (v *boltdbVisits) ResetVisit(u string) error {
	u = v.fix(u)
	err := v.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		return b.Delete([]byte(u))
	})
	return err
}

func (v *boltdbVisits) Drop() error {
	err := v.db.Update(func(tx *bolt.Tx) error {
		if err := tx.DeleteBucket(bucketName); err != nil {
			return err
		}
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
	return err
}

func (v *boltdbVisits) Close() error {
	return v.db.Close()
}

func NewBoltDbVisits(dbpath string) (Visiter, error) {
	db, err := bolt.Open(dbpath, 0600, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
	return &boltdbVisits{
		db: db,
	}, err
}
