package nn

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"syscall/js"
)

// wrapper over local storage js api
type storageWrapper struct {
	// key: name in localStorage
	// value: map where key is pointer to component, value - it callback
	listeners map[string]map[Component]func(string)
	mu        sync.Mutex
}

// Storage — gloval instrument for work with localStorage
var Storage = &storageWrapper{}

var NotFound = errors.New("not found")

func (s *storageWrapper) raw() js.Value {
	return js.Global().Get("localStorage")
}

func (s *storageWrapper) Set(key, value string) {
	s.raw().Call("setItem", key, value)

	s.notify(key, value)
}

func (s *storageWrapper) Get(key string) (string, bool) {
	val := s.raw().Call("getItem", key)
	if val.IsNull() {
		return "", false
	}
	return val.String(), true
}

func (s *storageWrapper) Remove(key string) {
	s.raw().Call("removeItem", key)
	s.notify(key, "")
}

func (s *storageWrapper) Clear() {
	s.raw().Call("clear")
}

func (s *storageWrapper) SetJSON(key string, value any) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	s.Set(key, string(bytes))
	return nil
}

func (s *storageWrapper) GetJSON(key string, value any) error {
	valStr, exists := Storage.Get(key)
	if !exists {
		return NotFound
	}

	return json.Unmarshal([]byte(valStr), value)
}

func (s *storageWrapper) Watch(comp Component, key string, callback func(newValue string)) {
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

func (s *storageWrapper) unwatchAll(comp Component) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for key := range s.listeners {
		delete(s.listeners[key], comp)
	}
}

func (s *storageWrapper) notify(key, value string) {
	fmt.Println("[notify]", key, value)
	s.mu.Lock()
	callbacks := s.listeners[key]
	s.mu.Unlock()

	for _, cb := range callbacks {
		cb(value)
	}
}

func initStorageListener() {
	window := global.Get("window")

	cb := js.FuncOf(func(this js.Value, args []js.Value) any {
		event := args[0]

		key := event.Get("key").String()
		newValue := event.Get("newValue")

		valStr := ""
		if !newValue.IsNull() {
			valStr = newValue.String()
		}

		Storage.notify(key, valStr)

		return nil
	})

	window.Call("addEventListener", "storage", cb)
}
