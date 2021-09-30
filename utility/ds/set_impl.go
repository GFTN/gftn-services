// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package ds

import (
	"errors"
)

var void Void

type Set struct {
	data map[interface{}]Void
}

type SetIterator struct {
	*Set
	iteratedSet Set
}

func NewSet() Set {
	s := Set{}
	s.data = make(map[interface{}]Void)
	return s
}

func (s *Set) Add(e interface{}) {
	s.data[e] = void
}

func (s *Set) Remove(e interface{}) {
	delete(s.data, e)
}

func (s *Set) Contains(e interface{}) bool {
	_, ok := s.data[e]
	return ok
}

func (s Set) Size() int {
	return len(s.data)
}

func (iter SetIterator) HasMore() bool {
	if len(iter.Set.data) > len(iter.iteratedSet.data) {
		return true
	}
	return false
}

func (iter *SetIterator) Next() (interface{}, error) {
	if len(iter.Set.data) == len(iter.iteratedSet.data) {
		return nil, errors.New("No more data")
	}
	/*
		Find the next key which currently does not exist in iteratedSet
		where the keys have been visited from Set;
		Then, mark the next key as visited by inserting it into iteratedSet and return it.
	*/
	for k := range iter.Set.data {
		if _, ok := iter.iteratedSet.data[k]; !ok {
			iter.iteratedSet.data[k] = void
			return k, nil
		}
	}

	return nil, errors.New("No more data")
}

func (s *Set) NewIterator() (*SetIterator, error) {
	if len(s.data) == 0 {
		return nil, errors.New("Set is empty")
	}
	return &SetIterator{s, NewSet()}, nil
}
