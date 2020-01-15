package ezbunt

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

func setup(t *testing.T) string {
	path, err := ioutil.TempFile(os.TempDir(), "ezbunt_test.db")
	if err != nil {
		t.FailNow()
	}
	absPath, err := filepath.Abs(path.Name())

	return absPath
}

func teardown(dbFilePath string, t *testing.T) {
	err := os.Remove(dbFilePath)
	if err != nil {
		t.Fail()
	}
}

func TestNewDB(t *testing.T) {
	dbFilePath := setup(t)
	defer teardown(dbFilePath, t)
	ez := NewEzbunt(dbFilePath)

	is := is.New(t)

	err := ez.WriteKeyVal("this", "that")
	is.NoErr(err) // expect no error
}

func TestWriteKeyVal(t *testing.T) {
	dbFilePath := setup(t)
	defer teardown(dbFilePath, t)
	ez := NewEzbunt(dbFilePath)

	kStr := "Gummer"
	wantVal := "Gormley"
	err := ez.WriteKeyVal(kStr, wantVal)
	is := is.New(t)
	is.NoErr(err) // expect no error

	got, err := ez.GetVal(kStr)
	is.NoErr(err)
	is.Equal(got, wantVal) // expect to be equal

}

func TestWriteKeyValTTL(t *testing.T) {
	dbFilePath := setup(t)
	defer teardown(dbFilePath, t)
	ez := NewEzbunt(dbFilePath)

	kStr := "Gummer"
	vStr := "Gormley"
	err := ez.WriteKeyValTTL(kStr, vStr, 0)
	is := is.New(t)
	is.NoErr(err) // expect no error

	_, err = ez.GetVal(kStr)
	if err == nil {
		t.Fail()
	}
}
