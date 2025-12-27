package provider

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestClientLoadSave(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "store.json")

	client := &urlShortenerClient{
		Store:     make(map[string]string),
		mu:        &sync.Mutex{},
		storeFile: f,
	}
	client.Store["short.ly/1"] = "https://example.com/1"
	if err := client.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// create new client and load
	client2 := &urlShortenerClient{
		Store:     make(map[string]string),
		mu:        &sync.Mutex{},
		storeFile: f,
	}
	if err := client2.Load(); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if got := client2.Store["short.ly/1"]; got != "https://example.com/1" {
		t.Fatalf("unexpected value: %s", got)
	}

	// cleanup file
	os.Remove(f)
}
