package lru

import (
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestCache(t *testing.T) {
	data := map[string]String{
		"k1": "v1", "k2": "v2", "k3": "v3",
	}

	t.Run("null", func(t *testing.T) {
		cache := New(0)
		for k, v := range data {
			cache.Set(k, v)
			if val, ok := cache.Get(k); ok {
				t.Fatal("should not value, but get:", val)
			}
		}
	})

	t.Run("normal", func(t *testing.T) {
		cache := New(int64(len(data) * 4))
		for k, v := range data {
			cache.Set(k, v)
			if val, ok := cache.Get(k); !ok || val != data[k] {
				t.Fatalf("should be value:%s, but get: %s", data[k], val)
			}
		}
	})

	t.Run("full", func(t *testing.T) {
		cache := New(int64((len(data) - 1) * 4))
		cache.Set("k1", data["k1"])
		if val, ok := cache.Get("k1"); !ok || val != data["k1"] {
			t.Fatalf("should be value:%s, but get: %s", data["k1"], val)
		}

		cache.Set("k2", data["k2"])
		if val, ok := cache.Get("k2"); !ok || val != data["k2"] {
			t.Fatalf("should be value:%s, but get: %s", data["k2"], val)
		}

		cache.Set("k3", data["k3"])
		if val, ok := cache.Get("k3"); !ok || val != data["k3"] {
			t.Fatalf("should be value:%s, but get: %s", data["k3"], val)
		}

		if val, ok := cache.Get("k1"); ok {
			t.Fatal("should not value, but get:", val)
		}
	})
}
