package service

// UpdateStatus tries to update the local status of a subset of all nodes specified by the limit parameter.
// Once all nodes have been updated, it starts from the beginning.
func (s *Service) UpdateStatus(chain string, limit int) {
	
}

// SendStatus sends the local status information about every node to Apache Kafka
func (s *Service) SendStatus(chain string) {
	// round := time.Now().UnixMilli() / 1000 / 60
}