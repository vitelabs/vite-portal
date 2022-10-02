package service

// UpdateStatus tries to update the local status of a subset of all nodes specified by the limit parameter.
// Once all nodes have been updated, it starts from the beginning.
func (s *Service) UpdateStatus(chain string, limit int) {
	if chain == "" || limit <= 0 {
		return
	}
	e := s.store.GetEnumerator(chain)
	processed := *s.status.GetRequestedSet(chain)
	count := 0
	for e.MoveNext() {
		n, found := e.Current()
		if !found {
			continue
		}
		if processed.Contains(n.Id) {
			continue
		}
		processed.Add(n.Id)

		count++
		if count > limit {
			break
		}
	}
	if count == 0 {
		processed.Clear()
	}
}

// SendStatus sends the local status information about every node to Apache Kafka
func (s *Service) SendStatus(chain string) {
	// round := time.Now().UnixMilli() / 1000 / 60
}