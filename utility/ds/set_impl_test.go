// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package ds

import (
	"testing"
)

func TestAddElement(t *testing.T) {
	s := NewSet()
	e := "dummy"
	s.Add(e)
}

func TestRemoveElement(t *testing.T) {
	s := NewSet()
	e := "dummy"

	s.Add(e)
	if s.Size() <= 0 {
		t.Fail()
	}

	s.Remove(e)
	if s.Size() > 0 {
		t.Fail()
	}
}

func TestContainsElement(t *testing.T) {
	s := NewSet()
	e := "dummy"
	f := "empty"

	s.Add(e)
	if !s.Contains(e) || s.Contains(f) {
		t.Fail()
	}

	s.Add(f)
	if !s.Contains(e) || !s.Contains(f) {
		t.Fail()
	}
}

func TestHasMoreElement(t *testing.T) {
	s := NewSet()
	iter, err := s.NewIterator()

	if iter != nil || err == nil {
		t.Fail()
	}

	e := "dummy"
	s.Add(e)

	iter, err = s.NewIterator()

	if iter == nil || err != nil {
		t.Fail()
	}

	if !iter.HasMore() {
		t.Fail()
	}

	key, err := iter.Next()
	if key == nil || err != nil {
		t.Fail()
	}

	s.Remove(e)

	if iter.HasMore() {
		t.Fail()
	}

	key, err = iter.Next()
	if key != nil || err == nil {
		t.Fail()
	}

}

func TestNextElement(t *testing.T) {
	s := NewSet()
	e := "dummy"
	s.Add(e)
	iter, err := s.NewIterator()

	if err != nil {
		t.Fail()
	}

	key, err := iter.Next()
	if key == nil || err != nil {
		t.Fail()
	}

	if iter.HasMore() {
		t.Fail()
	}

}
