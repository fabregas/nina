---
layout: home

hero:
  name: "Nina UI"
  text: "Build Modern Web Apps in Pure Go."
  tagline: "A declarative, WebAssembly-native frontend framework. Ditch JavaScript, forget HTML templates, and unleash the full power of Go in the browser."
  actions:
    - theme: brand
      text: Get Started
      link: /guide/getting-started
    - theme: alt
      text: View Components
      link: /ui/button
---

## Why Choose Nina UI?

### 🚫 Zero JavaScript. 100% Go.
Leave the JavaScript ecosystem behind. Nina UI compiles directly to WebAssembly, allowing you to write your entire frontend exclusively in Go. No more context switching between languages, no Webpack, and no `node_modules`.

### ⚡️ Unleash Native Go Features
Because your app runs in WebAssembly, you have the full power of the Go standard library at your fingertips. Fire up **Goroutines** for concurrent tasks, use channels for state management, and make network requests using the standard `net/http` client directly from the browser.

### 🎨 Declarative & Template-Free
Say goodbye to clunky HTML templates and string concatenation. Nina UI uses a fluent, deeply integrated Builder API. Construct your UI programmatically with type safety and IDE autocomplete:

```go
nn.Div().
    Class("flex items-center gap-4").
    Attr("data-role", "card").
    Children(
        nn.Text("Hello, WebAssembly!"),
    )
```

### 🧠 Smart, Granular Rendering
Performance is a feature. Nina UI features a highly optimized, intelligent rendering engine. When your application state changes, Nina doesn't re-render the entire page. It isolates the exact component that triggered the update.

### 🎯 Surgical DOM Patching
Actual DOM mutations are expensive. Nina UI performs an in-memory diff and touches the real browser DOM *only* when absolutely necessary. If only a single class name or attribute changes, only that specific node is updated.

---

## 🏗️ The Architecture
Nina UI is designed with strict boundaries and a clear separation of concerns, ensuring your codebase remains maintainable as it scales.

#### 1. The Core Layer (DOM & Lifecycle)
The foundation. This layer handles direct communication with the browser via `syscall/js`. It manages raw HTML elements, memory mapping, node lifecycle, and the core diffing/patching algorithm.

#### 2. The Unified Component Layer
The core of your design system. This layer seamlessly combines Tailwind-styled structural primitives with reactive logic. Whether you need a simple, stateless building block or a fully encapsulated, self-re-rendering widget (like a Dropdown or Dialog), everything shares the same cohesive API. Components dynamically scale from purely visual elements to state-aware, event-driven widgets powered by Signals.

#### 4. The Application Layer
Your domain. This is where you compose the primitives and smart components to build your actual application screens. Focus entirely on your business logic, knowing the lower layers are handling the heavy lifting.



