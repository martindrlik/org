// Package store represents simple in memory storage
// with rolling - old entries are overwritten by new
// ones.
//
package store

import (
	"sync"
	"time"
)

type Value struct {
	Name string
	Time time.Time
	Data []byte
}

type Store struct {
	idx int
	mux sync.Mutex
	val []Value

	byName map[string]map[int]struct{}
}

func NewStore(max int) *Store {
	return &Store{
		val:    make([]Value, 0, max),
		byName: make(map[string]map[int]struct{}),
	}
}

func (s *Store) Add(name string, data []byte) {
	v := Value{name, time.Now(), data}
	s.mux.Lock()
	defer s.mux.Unlock()
	if cap(s.val) > len(s.val) {
		m := make(map[int]struct{})
		m[len(s.val)] = struct{}{}
		s.byName[name] = m
		s.val = append(s.val, v)
		return
	}
	if s.idx == len(s.val) {
		s.idx = 0
	}
	if is, ok := s.byName[s.val[s.idx].Name]; ok {
		delete(is, s.idx)
	}
	s.val[s.idx] = v
	is, ok := s.byName[name]
	if !ok {
		is = make(map[int]struct{})
		s.byName[name] = is
	}
	s.byName[name][s.idx] = struct{}{}
	s.idx++
}

func (s *Store) All() []Value {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.val
}

func (s *Store) ByName(name string) ([]Value, bool) {
	s.mux.Lock()
	defer s.mux.Unlock()
	is, ok := s.byName[name]
	if !ok {
		return nil, false
	}
	v := make([]Value, len(is))
	i := -1
	for intoVal := range is {
		i++
		v[i] = s.val[intoVal]
	}
	return v, true
}
