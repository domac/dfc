package store

import (
	"testing"
)

func TestDBFunctions(t *testing.T) {
	resourceDB, err := OpenResourceDB("/tmp/dfc")

	if err != nil {
		t.Fatal(err)
	}
	defer resourceDB.Close()

	err = resourceDB.Set([]byte("key1"), []byte("value1"))
	if err != nil {
		t.Fail()
	}

	value1, _ := resourceDB.Get([]byte("key1"))
	if string(value1) != "value1" {
		t.Fail()
	}

	err = resourceDB.Set([]byte("key2"), []byte("value2"))
	if err != nil {
		t.Fail()
	}

	value2, _ := resourceDB.Get([]byte("key2"))
	if string(value2) != "value2" {
		t.Fail()
	}

	keys := resourceDB.Keys()
	if len(keys) != 2 {
		t.Fail()
	}

	resourceDB.Remove([]byte("key1"))
	keys = resourceDB.Keys()
	if len(keys) != 1 {
		t.Fail()
	}
}
