/*
 * Copyright 2014 Google Inc. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package inmemory implements a simple in-memory GraphStore.
package inmemory

import (
	"sort"
	"sync"

	"kythe/go/storage"

	spb "kythe/proto/storage_proto"

	"code.google.com/p/goprotobuf/proto"
)

type store struct {
	entries []*spb.Entry
	mu      sync.RWMutex
}

// Create returns a new in-memory GraphStore
func Create() storage.GraphStore {
	return &store{}
}

// Close implements part of the GraphStore interface.
func (store) Close() error { return nil }

// Write implements part of the GraphStore interface.
func (s *store) Write(req *spb.WriteRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, u := range req.Update {
		s.insert(proto.Clone(&spb.Entry{
			Source:    req.Source,
			EdgeKind:  u.EdgeKind,
			Target:    u.Target,
			FactName:  u.FactName,
			FactValue: u.FactValue,
		}).(*spb.Entry))
	}
	return nil
}

func (s *store) insert(e *spb.Entry) {
	i := sort.Search(len(s.entries), func(i int) bool {
		return storage.EntryLess(e, s.entries[i])
	})
	if i == len(s.entries) {
		s.entries = append(s.entries, e)
	} else if i < len(s.entries) && storage.EntryCompare(e, s.entries[i]) == storage.EQ {
		s.entries[i] = e
	} else if i == 0 {
		s.entries = append([]*spb.Entry{e}, s.entries...)
	} else {
		s.entries = append(s.entries[:i], append([]*spb.Entry{e}, s.entries[i:]...)...)
	}
}

// Read implements part of the GraphStore interface.
func (s *store) Read(req *spb.ReadRequest, stream chan<- *spb.Entry) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	start := sort.Search(len(s.entries), func(i int) bool {
		comp := storage.VNameCompare(s.entries[i].Source, req.Source)
		return comp != storage.LT && (comp == storage.GT || req.GetEdgeKind() == "*" || s.entries[i].GetEdgeKind() >= req.GetEdgeKind())
	})
	end := sort.Search(len(s.entries), func(i int) bool {
		comp := storage.VNameCompare(s.entries[i].Source, req.Source)
		return comp == storage.GT || (req.GetEdgeKind() != "*" && s.entries[i].GetEdgeKind() > req.GetEdgeKind())
	})
	for i := start; i < end; i++ {
		stream <- s.entries[i]
	}
	return nil
}

// Scan implements part of the GraphStore interface.
func (s *store) Scan(req *spb.ScanRequest, stream chan<- *spb.Entry) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, e := range s.entries {
		if storage.EntryMatchesScan(req, e) {
			stream <- e
		}
	}
	return nil
}
