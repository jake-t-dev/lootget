package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type channelStore struct {
	mu       sync.RWMutex
	path     string
	channels map[string]string
}

func newChannelStore(path string) (*channelStore, error) {
	s := &channelStore{path: path, channels: make(map[string]string)}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &s.channels); err != nil {
		return nil, err
	}
	if s.channels == nil {
		return nil, fmt.Errorf("channel settings must be a JSON object")
	}
	return s, nil
}

func (s *channelStore) output(guild string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.channels[guild]
}

func (s *channelStore) set(guild, channel string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	next := make(map[string]string, len(s.channels)+1)
	for k, v := range s.channels {
		next[k] = v
	}
	next[guild] = channel
	data, err := json.Marshal(next)
	if err != nil {
		return err
	}
	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	file, err := os.CreateTemp(dir, ".channels-*")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	if _, err := file.Write(data); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := os.Rename(file.Name(), s.path); err != nil {
		return err
	}
	s.channels = next
	return nil
}
