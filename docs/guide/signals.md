# Signals: Reactive State Management

In traditional Virtual DOM frameworks, state is tightly coupled to the component's lifecycle. When a variable changes, the framework must re-render the entire component (and often its children) to figure out what changed in the DOM. 

Nina framework solves this performance bottleneck using **Signals**. 

A Signal is a reactive primitive that holds a value. It acts as a single source of truth that the UI can "subscribe" to. When the value inside the Signal changes, the framework surgically updates *only* the specific components that depend on it, without re-rendering the surrounding components.

## Why Use Signals?

1. **Granular Updates:** If a Signal controls the text on your cool button, only that text node is updated in the DOM when the state changes. The rest of the component tree is untouched.
2. **Decoupled from Components:** Signals can live outside of components. You can define a Signal globally, pass it through Context, or export it from a separate package.
3. **Concurrency Safe:** Built with Go's `sync.RWMutex` under the hood, Signals are thread-safe. You can safely mutate UI state from background goroutines, HTTP callbacks, or WebSockets without worrying about race conditions.

---

## Core API

Working with Signals involves three simple steps:

1. **`nn.NewSignal[T](initialValue)`**: Creates a new signal of any type `T`.
2. **`signal.Get(component)`**: Reads the current value **AND** registers the passed component as a subscriber. Every time this signal changes, this component's `View()` method will be re-evaluated. *(Note: Passing `nil` allows you to read the value without subscribing).*
3. **`signal.Set(newValue)`**: Updates the value and instantly triggers a re-render for all active subscribers.

---

## Usage Examples

### 1. Basic Local State (Counter)
The most common use case is managing simple local state within a component. We use `ui.Dynamic` to create a reactive boundary. Whenever `counterSig` changes, only the contents of `ui.Dynamic` will be re-evaluated.

```go

type counter struct {
    nn.BaseComponent

    counterSig *nn.Signal[int]
}

func Counter() *counter {
    return &counter{
        counterSig: nn.NewSignal[int](0), // initialize signal
    }
}

func (c *counter) View() nn.Node {
    return nn.Div().
        Class("p-4 border rounded shadow").
        Children(
            // Wrap the reactive part in ui.Dynamic
            nn.Text(fmt.Sprintf("Current count: %d", c.counterSig.Get(c))),
            
            nn.Button().
                Class("mt-2 bg-blue-500 text-white px-4 py-2 rounded").
                OnClick(func(e nn.Event) {
                    // Update the signal on user interaction
                    // We use Get(nil) here because we just want to read the current 
                    // value for the calculation, not subscribe the event handler.
                    c.counterSig.Set(c.counterSig.Get(nil) + 1)
                }).
                Children(ui.Text("Increment")),
        )
}
```


### 2. Conditional Rendering
Signals are perfect for toggling UI elements on and off, such as dropdowns, modals, or loading spinners.

```go
func (w *loadingWidget) View() nn.Node {
    return nn.Div().Children(
        nn.Button().
            OnClick(func(e nn.Event) {
                w.isLoading.Set(true)
                // Simulate a network request
                go func() {
                    time.Sleep(2 * time.Second)
                    w.isLoading.Set(false)
                }()
            }).
            Children(ui.Text("Fetch Data")),

        // Reactively show/hide content
        nn.IfElse(
            w.isLoading.Get(w),
            nn.Div().Class("spinner").Text("Loading..."),
            nn.Duv().Text("Data loaded successfully!"),
        ),
    )
}
```

### 3. Global / Shared State
Because Signals are just Go structs, they can be declared globally or shared across entirely different parts of your application without needing complex state managers (like Redux).

```go
// state.go
// Declare a global signal
var ActiveUser = nn.NewSignal[string]("Guest")

// header.go
func (h *Header) View() nn.Node {
    return nn.Header().Children(
            return nn.Text("Welcome, " + ActiveUser.Get(h))
    )
}

// login_form.go
func (f *LoginForm) View() nn.Node {
    return nn.Button().
        OnClick(func(e nn.Event) {
            // Updating this signal will instantly update the Header,
            // no matter where it is in the component tree.
            ActiveUser.Set("AdminUser")
        }).
        Children(nn.Text("Log In"))
}
```

::: tip Advanced Usage: Structs inside Signals
You can store complex structs or arrays inside a Signal. Just remember that updating a struct requires replacing it with a new instance (since Signals track the value assignment, not internal struct mutations).

```go
type UserData struct {
    Name  string
    Roles []string
}

userSig := nn.NewSignal[UserData](UserData{Name: "Alex"})

// To update, set a whole new struct:
userSig.Set(UserData{Name: "Alex", Roles: []string{"admin"}})
```
