package ezbunt

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/buntdb"
)

// Ezbunt encapsulates and extends a buntdb instance
type Ezbunt struct {
	db *buntdb.DB
}

// New creates and returns an Ezbunt with a new buntdb (with a backing file at the path provided)
func New(dbFilePath string) *Ezbunt {
	newDB, err := buntdb.Open(dbFilePath)
	if err != nil {
		log.Fatal("Could not open data file path", err)
		return nil
	}
	return &Ezbunt{db: newDB}
}

// WriteKeyVal persists a key-value string pair indefinitley
func (ez *Ezbunt) WriteKeyVal(key string, val string) error {
	db := ez.db
	err := db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, val, nil)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

// WriteKeyValTTL persists a key-value string pair with a ttl in seconds
func (ez *Ezbunt) WriteKeyValTTL(key string, val string, ttlSeconds int) error {
	db := ez.db
	err := db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, val, &buntdb.SetOptions{Expires: true, TTL: time.Second * time.Duration(ttlSeconds)})
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

// WriteKeyValAsInt persists a key-value string/int pair indefinitley
func (ez *Ezbunt) WriteKeyValAsInt(key string, val int) error {
	valAsStr := strconv.Itoa(val)
	return ez.WriteKeyVal(key, valAsStr)
}

// WriteKeyValAsIntTTL persists a key-value string/int pair wit a ttl in seconds
func (ez *Ezbunt) WriteKeyValAsIntTTL(key string, val int, ttlSeconds int) error {
	valAsStr := strconv.Itoa(val)
	return ez.WriteKeyValTTL(key, valAsStr, ttlSeconds)
}

// WriteKeyValAsBool persists a key-value string/bool pair indefinitley
func (ez *Ezbunt) WriteKeyValAsBool(key string, val bool) error {
	valAsStr := strconv.FormatBool(val)
	return ez.WriteKeyVal(key, valAsStr)
}

// WriteKeyValAsBoolTTL persists a key-value string/bool pair indefinitley
func (ez *Ezbunt) WriteKeyValAsBoolTTL(key string, val bool, ttlSeconds int) error {
	valAsStr := strconv.FormatBool(val)
	return ez.WriteKeyValTTL(key, valAsStr, ttlSeconds)
}

// WriteKeyValAsJSON persists a key-value pair with a string key, and converts the val
// to a json []byte, and stores it as the value.  Useful for complex types.
func (ez *Ezbunt) WriteKeyValAsJSON(key string, val interface{}) error {
	valAsJSON, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return ez.WriteKeyVal(key, string(valAsJSON))
}

// WriteKeyValAsJSONTTL persists a key-value pair with a string key, and converts the val
// to a json []byte, and stores it as the value, with a TTL in seconds.  Useful for complex types.
func (ez *Ezbunt) WriteKeyValAsJSONTTL(key string, val interface{}, ttlSeconds int) error {
	valAsJSON, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return ez.WriteKeyValTTL(key, string(valAsJSON), ttlSeconds)
}

// GetVal retrieves the value, and possible error, for the corresponding key.
func (ez *Ezbunt) GetVal(key string) (string, error) {
	db := ez.db

	var theVal string
	err := db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(key)
		if err != nil {
			return err
		}
		theVal = val
		return nil
	})

	return theVal, err
}

// GetValAsInt retrieves the value as an `int` type, and possible error, for the corresponding key.
func (ez *Ezbunt) GetValAsInt(key string) (int, error) {
	val, err := ez.GetVal(key)
	if err != nil {
		return 0, err
	}
	valAsInt, convErr := strconv.Atoi(val)
	if convErr != nil {
		return 0, err
	}
	return valAsInt, nil
}

// GetValAsBool retrieves the value as a `bool` type, and possible error, for the corresponding key.
func (ez *Ezbunt) GetValAsBool(key string) (bool, error) {
	val, err := ez.GetVal(key)
	if err != nil {
		return false, err
	}
	valAsBool, convErr := strconv.ParseBool(val)
	if convErr != nil {
		return valAsBool, err
	}
	return valAsBool, nil
}

// GetValAsBytes retrieves the value as []byte, and possible error, for the corresponding key.  Useful for
// retrieving JSON objects.
func (ez *Ezbunt) GetValAsBytes(key string) ([]byte, error) {
	val, err := ez.GetVal(key)
	valBytes := []byte(val)
	if err != nil {
		return valBytes, err
	}
	return valBytes, nil
}

// GetValDefault retrieves the value for the corresponding key.
// If the key is not found, the value is not found, or an error is returned
// from the db, the provided default is returned.
func (ez *Ezbunt) GetValDefault(key string, defaultVal string) string {
	val, err := ez.GetVal(key)
	if err != nil {
		return defaultVal
	}
	return val
}

// GetValAsIntDefault retrieves the value for the corresponding key.
// If the key is not found, the value is not found, or an error is returned
// from the db, the provided default is returned.
func (ez *Ezbunt) GetValAsIntDefault(key string, defaultVal int) int {
	val, err := ez.GetValAsInt(key)
	if err != nil {
		return defaultVal
	}
	return val
}

// GetValAsBoolDefault retrieves the value for the corresponding key.
// If the key is not found, the value is not found, or an error is returned
// from the db, the provided default is returned.
func (ez *Ezbunt) GetValAsBoolDefault(key string, defaultVal bool) bool {
	val, err := ez.GetValAsBool(key)
	if err != nil {
		return defaultVal
	}
	return val
}

// GetPairs retrieves a map of key-value pairs based on a key prefix
func (ez *Ezbunt) GetPairs(keyPrefix string) (map[string]string, error) {
	db := ez.db
	pairs := make(map[string]string, 0)

	err := db.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("", func(key, value string) bool {
			if strings.HasPrefix(key, keyPrefix) {
				pairs[key] = value
			}
			return true
		})
		return err
	})

	return pairs, err
}

// DeleteKey removes a key-value pair and returns the value of the removed key
func (ez *Ezbunt) DeleteKey(key string) (string, error) {
	db := ez.db

	var theVal string
	err := db.Update(func(tx *buntdb.Tx) error {
		val, err := tx.Delete(key)
		if err != nil {
			return err
		}
		theVal = val
		return nil
	})

	return theVal, err
}
