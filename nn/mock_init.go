//go:build !js

package nn

import (
	"encoding/json"
	"sync"
	"time"
)

func init() {
	nina = &engine{
		registry:        make(map[Component]*componentNode),
		dirtyComponents: make(map[Component]bool),
	}

	nina.renderer = newMockRenderer()
	nina.storage = newMockStorage()
	nina.reqDomUpdate, _ = nina.renderer.initRequestAnimationFrame(nina.performUpdates)
}

///

type mockNode struct {
	Value any
}

func (b mockNode) isNative() {}

func (d mockNode) Equal(node NativeNode) bool {
	if node == nil {
		return false
	}
	n := node.(mockNode).Value

	return d.Value == n
}

func (b mockNode) Raw() any {
	return b.Value
}

type mockRenderer struct{}

func newMockRenderer() *mockRenderer {
	return &mockRenderer{}
}
func (d *mockRenderer) RootNode() NativeNode {
	return mockNode{"root"}
}

func (d *mockRenderer) Window() NativeNode {
	return mockNode{"window"}
}

func (d *mockRenderer) CreateElement(tag string) NativeNode {
	return mockNode{tag}
}
func (d *mockRenderer) CreateElementNS(ns, tag string) NativeNode {
	return mockNode{tag}
}
func (d *mockRenderer) CreateTextNode(text string) NativeNode {
	return mockNode{text}
}
func (d *mockRenderer) CreateComment(comment string) NativeNode {
	return mockNode{comment}
}
func (d *mockRenderer) CreateDocumentFragment() NativeNode {
	return mockNode{"fragment"}
}

func (d *mockRenderer) SetAttribute(node NativeNode, key, val string) {}

func (d *mockRenderer) RemoveAttribute(node NativeNode, key string) {
}

func (d *mockRenderer) HasAttribute(node NativeNode, key string) bool {
	return false
}

func (d *mockRenderer) GetAttribute(node NativeNode, key string) string {
	return ""
}

func (d *mockRenderer) AppendChild(parent, child NativeNode) {}

func (d *mockRenderer) InsertBefore(parent, child, anchor NativeNode) {}

func (d *mockRenderer) Remove(node NativeNode) {}

func (d *mockRenderer) AddEventListener(node NativeNode, event string, handler func(Event)) func() {
	return func() {}
}

func (d *mockRenderer) AddEventListenerWithCapture(node NativeNode, event string, handler func(Event)) func() {
	return func() {}
}
func (d *mockRenderer) AddResizeObserver(node NativeNode, handler func(Event)) func() {
	return func() {}
}

func (d *mockRenderer) NextSibling(node NativeNode) NativeNode {
	return nil
}

func (d *mockRenderer) FirstChild(node NativeNode) NativeNode {
	return nil
}

func (d *mockRenderer) SetInnerHTML(node NativeNode, html string) {}

func (d *mockRenderer) SetNodeValue(node NativeNode, val string) {}

func (d *mockRenderer) GetElementById(id string) NativeNode {
	return mockNode{id}
}

func (d *mockRenderer) Contains(node1, node2 NativeNode) bool {
	return false
}

func (d *mockRenderer) Closest(node NativeNode, selector string) NativeNode {
	return nil
}

func (d *mockRenderer) QuerySelector(node NativeNode, selector string) NativeNode {
	return nil
}

func (d *mockRenderer) QuerySelectorAll(node NativeNode, selector string) []NativeNode {
	return nil
}

func (d *mockRenderer) ScrollIntoView(node NativeNode, options map[string]any) {}

func (d *mockRenderer) Focus(node NativeNode) {}

func (d *mockRenderer) GetBoundingClientRect(node NativeNode) NativeNodeRect {
	return NativeNodeRect{}
}

func (d *mockRenderer) GetViewport() Viewport {
	return Viewport{}
}

func (d *mockRenderer) initRequestAnimationFrame(cb func()) (reqNext func(), cleaner func()) {
	reqNext = func() {
		go cb()
	}

	cleaner = func() {}

	return
}

func (d *mockRenderer) waitNextFrame() <-chan struct{} {
	ch := make(chan struct{})

	go func() {
		time.Sleep(time.Millisecond * 100)
		close(ch)
	}()

	return ch
}

func (d *mockRenderer) ToggleHTMLClass(class string) {}

func (d *mockRenderer) PushState(path string) {}

func (d *mockRenderer) OnPopState(handler func(path string)) func() {
	return func() {}
}

func (d *mockRenderer) GetCurrentPath() string {
	return "/"
}

/////////////////
// mock storage
////////////////

type mockStorage struct {
	listeners map[string]map[Component]func(string)
	mu        sync.Mutex
	data      map[string]string
}

func newMockStorage() *mockStorage {
	s := &mockStorage{data: make(map[string]string)}

	return s
}

func (s *mockStorage) Set(key, value string) {
	s.mu.Lock()
	s.data[key] = value
	s.mu.Unlock()

	s.notify(key, value)
}

func (s *mockStorage) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, ok := s.data[key]
	if !ok {
		return "", false
	}
	return val, true
}

func (s *mockStorage) Remove(key string) {
	s.mu.Lock()
	delete(s.data, key)
	s.mu.Unlock()

	s.notify(key, "")
}

func (s *mockStorage) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string]string)
}

func (s *mockStorage) SetJSON(key string, value any) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	s.Set(key, string(bytes))
	return nil
}

func (s *mockStorage) GetJSON(key string, value any) error {
	valStr, exists := s.Get(key)
	if !exists {
		return NotFound
	}

	return json.Unmarshal([]byte(valStr), value)
}

func (s *mockStorage) Watch(comp Component, key string, callback func(newValue string)) {
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

func (s *mockStorage) unwatchAll(comp Component) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for key := range s.listeners {
		delete(s.listeners[key], comp)
	}
}

func (s *mockStorage) notify(key, value string) {
	s.mu.Lock()
	callbacks := s.listeners[key]
	s.mu.Unlock()

	for _, cb := range callbacks {
		cb(value)
	}
}
