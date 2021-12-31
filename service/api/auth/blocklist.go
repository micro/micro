package auth

import (
	"context"
	"fmt"
	"sync"
)

var DefaultBlockList BlockList = &InMemoryBlockList{
	blocked: map[string]bool{},
}

// BlockList manages a list of blocked IDs.
type BlockList interface {
	// IsBlocked returns true if the id is blocked
	IsBlocked(ctx context.Context, id, namespace string) (bool, error)
	// Add adds an ID to the blocklist
	Add(ctx context.Context, id, namespace string) error
	// Remove removes an ID from the blocklist
	Remove(ctx context.Context, id, namespace string) error
	// List returns the entire list of blocked IDs
	List(ctx context.Context) ([]string, error)
}

// InMemoryBlockList is an implementation of BlockList that maintains the list in memory only.
type InMemoryBlockList struct {
	sync.RWMutex
	blocked map[string]bool
}

func (i *InMemoryBlockList) IsBlocked(ctx context.Context, id, namespace string) (bool, error) {
	i.RLock()
	defer i.RUnlock()
	return i.blocked[fmt.Sprintf("%s:%s", namespace, id)], nil
}

func (i *InMemoryBlockList) Add(ctx context.Context, id, namespace string) error {
	i.Lock()
	defer i.Unlock()
	i.blocked[fmt.Sprintf("%s:%s", namespace, id)] = true
	return nil
}

func (i *InMemoryBlockList) Remove(ctx context.Context, id, namespace string) error {
	i.Lock()
	defer i.Unlock()
	delete(i.blocked, fmt.Sprintf("%s:%s", namespace, id))
	return nil
}

func (i *InMemoryBlockList) List(ctx context.Context) ([]string, error) {
	i.RLock()
	defer i.RUnlock()
	res := make([]string, len(i.blocked))
	idx := 0
	for k, _ := range i.blocked {
		res[idx] = k
		idx++
	}
	return res, nil
}
