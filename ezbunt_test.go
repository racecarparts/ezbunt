package ezbunt

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

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
	ez := New(dbFilePath)

	is := is.New(t)

	err := ez.WriteKeyVal("this", "that")
	is.NoErr(err) // expect no error
}

func TestWriteKeyVal(t *testing.T) {
	dbFilePath := setup(t)
	defer teardown(dbFilePath, t)
	ez := New(dbFilePath)

	kStr := "Gummer"
	wantVal := "Gormley"
	err := ez.WriteKeyVal(kStr, wantVal)
	is := is.New(t)
	is.NoErr(err) // expect no error

	got, err := ez.GetVal(kStr)
	is.NoErr(err)
	is.Equal(got, wantVal) // expect to be equal
}

func TestWriteKeyValAsInt(t *testing.T) {
	dbFilePath := setup(t)
	defer teardown(dbFilePath, t)
	ez := New(dbFilePath)

	kStr := "purpose"
	wantVal := 42
	err := ez.WriteKeyValAsInt(kStr, wantVal)
	is := is.New(t)
	is.NoErr(err) // expect no error

	got, err := ez.GetValAsInt(kStr)
	is.NoErr(err)
	is.Equal(got, wantVal) // expect to be equal

	got = ez.GetValAsIntDefault("horse", 42)
	is.Equal(got, 42) // expect to get default
}

func TestWriteKeyValBool(t *testing.T) {
	dbFilePath := setup(t)
	defer teardown(dbFilePath, t)
	ez := New(dbFilePath)

	kStr := "all"
	wantVal := true
	err := ez.WriteKeyValAsBool(kStr, wantVal)
	is := is.New(t)
	is.NoErr(err) // expect no error

	got, err := ez.GetValAsBool(kStr)
	is.NoErr(err)
	is.Equal(got, wantVal) // expect to be equal

	got = ez.GetValAsBoolDefault("horse", true)
	is.Equal(got, true) // expect to get default
}

func TestWriteKeyValTime(t *testing.T) {
	dbFilePath := setup(t)
	defer teardown(dbFilePath, t)
	ez := New(dbFilePath)

	now := time.Now()

	kStr := "Cindy"
	wantVal := now.UTC()
	err := ez.WriteKeyValAsTime(kStr, wantVal)

	is := is.New(t)

	is.NoErr(err) // expect no error on write

	got, err := ez.GetValAsTime(kStr)
	is.NoErr(err)          // expect no error on GetValAsTime
	is.Equal(got, wantVal) // expect to be equal

	got = ez.GetValAsTimeDefault("Lauper", now.UTC())
	is.Equal(got, now.UTC()) // expect to get default
}

func TestWriteKeyValTTL(t *testing.T) {
	dbFilePath := setup(t)
	defer teardown(dbFilePath, t)
	ez := New(dbFilePath)

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

func TestGetPairs(t *testing.T) {
	dbFilePath := setup(t)
	defer teardown(dbFilePath, t)
	ez := New(dbFilePath)

	pairs := make(map[string]string)
	pairs["dance:harlem"] = "shake"
	pairs["food:milk"] = "shake"
	pairs["dance:mashed"] = "potato"
	pairs["food:mashed"] = "potato"
	pairs["dance:cabbage"] = "patch"
	pairs["sensation:cabbage"] = "patch"
	pairs["sensation:rick"] = "astley"

	for k, v := range pairs {
		ez.WriteKeyVal(k, v)
	}

	is := is.New(t)
	dances, err := ez.GetPairs("dance")
	is.NoErr(err) // expect no error

	foods, err := ez.GetPairs("food")
	is.NoErr(err) // expect no error

	sensations, err := ez.GetPairs("sensation")
	is.NoErr(err) // expect no error

	is.Equal(len(dances), 3)     // expect 3 dances
	is.Equal(len(foods), 2)      // expect 2 foods
	is.Equal(len(sensations), 2) // expect 2 sensations
}

func TestDeleteKey(t *testing.T) {
	dbFilePath := setup(t)
	defer teardown(dbFilePath, t)
	ez := New(dbFilePath)

	ez.WriteKeyVal("oompa", "loompa")

	is := is.New(t)
	val, err := ez.GetVal("oompa")

	is.NoErr(err)           // expect no error
	is.Equal(val, "loompa") // expect equal

	val, err = ez.DeleteKey("oompa")
	is.NoErr(err)           // expect no error
	is.Equal(val, "loompa") // expect equal

	val, err = ez.GetVal("oompa")
	is.True(err != nil) // expect error
	is.True(val == "")  // expect zero val string

}
