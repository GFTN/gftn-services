// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package ds

type Void struct{}

type Iterator interface {
	HasMore() bool
	Next() (interface{}, error)
}

type Container interface {
	Add(interface{})
	Remove(interface{})
	Contains(interface{}) bool
	Size() int
}
