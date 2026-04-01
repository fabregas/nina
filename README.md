# Nina Framework 🚀

Nina is a lightweight, fast, and declarative UI framework for building Single Page Applications (SPAs) in Go using WebAssembly.

Instead of writing your frontend in JavaScript/TypeScript and struggling to keep types synchronized with your backend, Nina allows you to build robust, type-safe, and highly responsive interfaces entirely in pure Go.

> ⚠️ **Disclaimer: Work in Progress (Alpha)**
> Nina is currently in early active development. While the core concepts and engine are functional, the API is highly volatile and subject to breaking changes without notice. It is **not yet recommended for production use**.


## ✨ Key Features

* **100% Go:** No JavaScript, no Webpack, no Vite. Your entire UI, state management, and business logic are written in pure Go.
* **Custom Virtual DOM:** Features a smart `patch` algorithm that surgically updates only the DOM elements that have actually changed.
* **Keyed Diffing:** High-performance rendering for large lists and data tables. It preserves element focus and avoids expensive, cascading DOM manipulations.
* **Smart Engine (Local Rendering):** Update isolated components locally (`nn.Update(comp)`) without triggering a full page re-render. Perfect for high-frequency data streams like WebSockets or background timers.
* **Lifecycle Hooks:** Built-in `OnMount` and `OnDestroy` interfaces give you safe, predictable control over goroutines, HTTP requests, and memory cleanup.
* **Built-in SPA Router:** Declarative, nested client-side routing out-of-the-box (`nn.Router`, `nn.RoutePrefix`) without relying on third-party libraries.
* **Opt-in Memoization:** Granular control over the rendering of heavy components via the `nn.Pure` interface (custom hash-based `isDirty` checks).

## 🎯 Ideal Use Cases

* Internal Admin Panels and Backoffice tools.
* Infrastructure monitoring dashboards (e.g., tracking PostgreSQL replication, managing EKS clusters, or viewing real-time metrics).
* Any application that benefits from strict typing and sharing identical DTOs (Data Transfer Objects) between the frontend and backend.

## 📦 Quick Start

### 1. Creating a Component

This example demonstrates both the automatic global state updates (via UI events) and the highly optimized local rendering engine (via background goroutines).

```go
package main

import (
	"fmt"
	"time"

    // Import the exact package folder
	"[github.com/fabregas/nina/nn](https://github.com/fabregas/nina/nn)"
)

type Dashboard struct {
	clicks int
	ticks  int
}

// OnMount is called automatically when the component appears on the screen
func (d *Dashboard) OnMount() {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			d.ticks++
			
			// Triggers a highly optimized local re-render for THIS component only!
			// The rest of the application tree is completely ignored.
			nn.Update(d) 
		}
	}()
}

func (d *Dashboard) View() *nn.Element {
	return nn.Div().Class("p-8").Children(
		nn.H1().Text("Nina Framework Demo"),
		
		nn.Div().Text(fmt.Sprintf("Background ticks (Local Update): %d", d.ticks)),
		nn.Div().Text(fmt.Sprintf("Button clicks (Auto Update): %d", d.clicks)),
		
		nn.Button().
			Class("bg-blue-500 text-white px-4 py-2 mt-4 rounded").
			Text("Click Me").
			OnClick(func() {
				d.clicks++
				// No need to call nn.Update here! 
				// Nina automatically schedules a batched UI update after any event handler.
			}),
	)
}
```


### 2. Mounting the Application


```go
func main() {
	app := &Dashboard{}
	
	// Mount the app to the HTML container with id="app"
	nn.Mount("app", app)
	
	// Block the main goroutine (required for WebAssembly execution)
	select {}
}
```

### 3. Building the project

```console
# GOOS=js GOARCH=wasm go build -o main.wasm
```

To run the application in a browser, you will need the standard wasm_exec.js file (provided by the Go installation) and a basic index.html file to load your compiled binary.


## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
