# Getting Started

Welcome to **Nina UI**! This guide will help you set up your first WebAssembly project.

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.21+** (for optimized Wasm compilation)
- **A local web server** (e.g., `goexec`, `python -m http.server`, or other)

## 1. Initialize your project

Create a new directory for your project and initialize the Go module:

```bash
mkdir my-nina-app
cd my-nina-app
go mod init my-nina-app
```


## 2. Prepare the WebAssembly Environment

Go requires a small JavaScript glue file to run WebAssembly. Copy it from your Go installation to your project folder:

```bash
# for Go 1.24+
cp $(go env GOROOT)/lib/wasm/wasm_exec.js .

# for older versions
cp $(go env GOROOT)/misc/wasm/wasm_exec.js .

```


## 3. Create the HTML Entry Point

Create `dist/index.html`:

```html
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    
    <script src="https://unpkg.com/@tailwindcss/browser@4"></script>
    
    <style type="text/tailwindcss">
      @theme {
        --color-primary: #0f172a;
        --color-primary-foreground: #f8fafc;
      }
    </style>

    <script src="wasm_exec.js"></script>
    <script>
      const go = new Go();
      WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
        go.run(result.instance);
      });
    </script>
  </head>
  <body>
    <div id="app"></div>
  </body>
</html>
```

## 4. Write your first Nina application

Create `main.go`. This is where the magic happens using Nina's declarative API:

```go
package main

import (
    "https://github.com/fabregas/nina/ui"
    "https://github.com/fabregas/nina/nn"
)

func main() {
    // Create a beautiful button
    app := nn.Div().
        Class("flex flex-col items-center justify-center h-screen gap-4").
        Children(
            nn.H1().Class("text-4xl font-bold").Text("Welcome to Nina UI"),
            ui.Button().
                SizeLg().
                Text("Click me!")
        )

    // Mount to the #app div
    nn.Mount("#app", app)
    
    // Keep the Wasm app running
    select {}
}
```

## 5. Build and Run

Compile your Go code to WebAssembly:

```bash
GOOS=js GOARCH=wasm go build -o ./dist/main.wasm main.go
```

Now, start your local web server. For example, using Python:

```bash
python3 -m http.server 8080
```

Open `http://localhost:8080` in your browser. You should see a perfectly styled, interactive UI rendered entirely by Go!


## Next Steps

- Explore the [Architecture](/guide/architecture) to understand the 4-layer system.
- Check out the [Component Gallery](/components/button) for ready-to-use widgets.
- Learn about [State Management](/guide/state) in Nina UI.
