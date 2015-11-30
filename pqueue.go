// Copyright 2015 By Jash. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// mail: shijian0912@163.com

package evloop

type Interface interface {
	Prior(Interface) bool
}

type Pqueue struct {
	data  []Interface
	count int
}

func NewPqueue() *Pqueue {
	var pq = new(Pqueue)
	pq.data = make([]Interface, 1, 100)
	return pq
}

func (pq *Pqueue) Empty() bool {
	return pq.data == nil || len(pq.data) <= 1
}

func (pq *Pqueue) Top() Interface {
	if pq.Empty() {
		return nil
	}
	return pq.data[1]
}

func (pq *Pqueue) Push(elem Interface) {
	pq.data = append(pq.data, elem)

	var ia = len(pq.data) - 1
	var ib = int(ia / 2)

	// root node is at index 1.
	for ia != 1 && pq.data[ia].Prior(pq.data[ib]) {
		pq.data[ia], pq.data[ib] = pq.data[ib], pq.data[ia]
		ia, ib = ib, ib/2
	}

	pq.count++
}

func (pq *Pqueue) Pop() {
	var n = len(pq.data) - 1
	pq.data[1] = pq.data[n]

	pq.data = pq.data[0:n]
	n--

	// p: parent, l: left, r: right.
	var p = 1
	for {
		var l = p * 2
		if l >= n {
			break
		}

		var i, r = l, l + 1
		if r < n && pq.data[r].Prior(pq.data[l]) {
			i = r
		}

		if pq.data[i].Prior(pq.data[p]) {
			pq.data[p], pq.data[l] = pq.data[l], pq.data[p]
			p = i
		} else {
			break
		}
	}
	pq.count--
}

func (pq *Pqueue) Count() int {
	return pq.count
}

func (pq *Pqueue) IsValid() bool {
	if len(pq.data) == pq.count {
		return true
	}
	return false
}
