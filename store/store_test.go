package store_test

import (
	"strconv"
	"testing"

	"github.com/martindrlik/org/store"
)

func TestAddAll(t *testing.T) {
	s := store.NewStore(2)
	s.Add("u", []byte("v"))
	s.Add("w", []byte("x"))
	v := s.All()
	if len(v) != 2 {
		t.Fatalf("All should return two entries; got %v", len(v))
	}
	if a := string(v[0].Data); a != "v" {
		t.Errorf("first entry should have v data; got %v", a)
	}
	if a := string(v[1].Data); a != "x" {
		t.Errorf("second entry should have x data; got %v", a)
	}
}

func TestByName(t *testing.T) {
	s := store.NewStore(2)
	s.Add("u", []byte("v"))
	s.Add("w", []byte("x"))
	t.Run("Found w", func(t *testing.T) {
		v, ok := s.ByName("w")
		if !ok {
			t.Errorf("ok should be true, w should be found")
		}
		if len(v) != 1 {
			t.Fatalf("number of entries should be 1; got %v", len(v))
		}
		if a := string(v[0].Data); a != "x" {
			t.Errorf("value for w should be x; got %v", a)
		}
	})
	t.Run("Not Found z", func(t *testing.T) {
		v, ok := s.ByName("z")
		if ok {
			t.Errorf("ok should be false, z should not be found")
		}
		if v != nil {
			t.Errorf("v should be nil")
		}
	})
}

func TestAddRolling(t *testing.T) {
	s := store.NewStore(1)
	s.Add("s", []byte("t")) // add 3 times so idx == len(val) { idx = 0 } is tested too
	s.Add("u", []byte("v"))
	s.Add("w", []byte("x"))
	v := s.All()
	if len(v) != 1 {
		t.Errorf("there should be only one entry; got %v", len(v))
	}
	if a := string(v[0].Data); a != "x" {
		t.Errorf("the value should be x; max capacity of store is one so first entry is overwritten by second and so on")
	}
}

func BenchmarkByName(b *testing.B) {
	const K = 20 * 1000
	s := store.NewStore(K)
	names := make([]string, K)
	for i := 0; i < K; i++ {
		name := strconv.FormatInt(int64(i), 36)
		s.Add(name, []byte("data"))
		names[i] = name
	}
	b.ResetTimer()
	for _, v := range names {
		s.ByName(v)
	}
}
