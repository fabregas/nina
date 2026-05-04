//go:build js && wasm

package nn

import (
	"encoding/json"
	"sync"
	"syscall/js"
)

// wrapper over local storage js api
type domStorage struct {
	// key: name in localStorage
	// value: map where key is pointer to component, value - it callback
	listeners map[string]map[Component]func(string)
	mu        sync.Mutex
}

func newDomStorage() *domStorage {
	s := &domStorage{}

	window := js.Global().Get("window")

	cb := js.FuncOf(func(this js.Value, args []js.Value) any {
		event := args[0]

		key := event.Get("key").String()
		newValue := event.Get("newValue")

		valStr := ""
		if !newValue.IsNull() {
			valStr = newValue.String()
		}

		s.notify(key, valStr)

		return nil
	})

	window.Call("addEventListener", "storage", cb)

	return s
}

func (s *domStorage) raw() js.Value {
	return js.Global().Get("localStorage")
}

func (s *domStorage) Set(key, value string) {
	s.raw().Call("setItem", key, value)

	s.notify(key, value)
}

func (s *domStorage) Get(key string) (string, bool) {
	val := s.raw().Call("getItem", key)
	if val.IsNull() {
		return "", false
	}
	return val.String(), true
}

func (s *domStorage) Remove(key string) {
	s.raw().Call("removeItem", key)
	s.notify(key, "")
}

func (s *domStorage) Clear() {
	s.raw().Call("clear")
}

func (s *domStorage) SetJSON(key string, value any) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	s.Set(key, string(bytes))
	return nil
}

func (s *domStorage) GetJSON(key string, value any) error {
	valStr, exists := s.Get(key)
	if !exists {
		return NotFound
	}

	return json.Unmarshal([]byte(valStr), value)
}

func (s *domStorage) Watch(comp Component, key string, callback func(newValue string)) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.listeners == nil {
		s.listeners = make(map[string]map[Component]func(string))
	}
	if s.listeners[key] == nil {
		s.listeners[key] = make(map[Component]func(string))
	}

	s.listeners[key][comp] = callback
}

func (s *domStorage) unwatchAll(comp Component) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for key := range s.listeners {
		delete(s.listeners[key], comp)
	}
}

func (s *domStorage) notify(key, value string) {
	s.mu.Lock()
	callbacks := s.listeners[key]
	s.mu.Unlock()

	for _, cb := range callbacks {
		cb(value)
	}
}
