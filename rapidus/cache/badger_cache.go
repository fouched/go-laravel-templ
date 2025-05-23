package cache

import (
	"github.com/dgraph-io/badger/v4"
	"time"
)

type BadgerCache struct {
	Conn   *badger.DB
	Prefix string
}

func (b *BadgerCache) Has(str string) (bool, error) {
	_, err := b.Get(str)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (b *BadgerCache) Get(str string) (interface{}, error) {
	var fromCache []byte

	err := b.Conn.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(str))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			fromCache = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	decoded, err := decode(string(fromCache))
	if err != nil {
		return nil, err
	}
	item := decoded[str]

	return item, nil
}

// Set creates an entry in the cache with an optional expiry time in seconds
func (b *BadgerCache) Set(str string, value interface{}, expireSecs ...int) error {
	entry := Entry{}
	entry[str] = value
	encoded, err := encode(entry)
	if err != nil {
		return err
	}

	if len(expireSecs) > 0 {
		err = b.Conn.Update(func(txn *badger.Txn) error {
			e := badger.NewEntry([]byte(str), encoded).WithTTL(time.Second * time.Duration(expireSecs[0]))
			err = txn.SetEntry(e)
			if err != nil {
				return err
			}
			//err = txn.Commit()
			return err
		})
	} else {
		err = b.Conn.Update(func(txn *badger.Txn) error {
			e := badger.NewEntry([]byte(str), encoded)
			err = txn.SetEntry(e)
			if err != nil {
				return err
			}
			//err = txn.Commit()
			return err
		})
	}

	return nil
}

func (b *BadgerCache) Forget(str string) error {
	err := b.Conn.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(str))
		return err
	})
	return err
}

func (b *BadgerCache) EmptyByMatch(str string) error {
	return b.emptyByMatch(str)
}

func (b *BadgerCache) Empty() error {
	return b.emptyByMatch("")
}

func (b *BadgerCache) emptyByMatch(str string) error {
	deleteKeys := func(keysForDelete [][]byte) error {
		if err := b.Conn.Update(func(txn *badger.Txn) error {
			for _, key := range keysForDelete {
				if err := txn.Delete(key); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
		return nil
	}

	collectSize := 100000
	err := b.Conn.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.AllVersions = false
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		keysForDelete := make([][]byte, 0, collectSize)
		keysCollected := 0

		for it.Seek([]byte(str)); it.ValidForPrefix([]byte(str)); it.Next() {
			key := it.Item().KeyCopy(nil)
			keysForDelete = append(keysForDelete, key)
			keysCollected++
			if keysCollected == collectSize {
				if err := deleteKeys(keysForDelete); err != nil {
					return err
				}
			}
		}

		if keysCollected > 0 {
			if err := deleteKeys(keysForDelete); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
