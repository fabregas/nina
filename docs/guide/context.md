# Context: Dependency Injection

In complex applications, passing data from a top-level component down to deeply nested children through standard properties is cumbersome and leads to "prop drilling". 

The **Context API** solves this by allowing parent components to seamlessly inject data, state, or signals into their entire sub-tree. Because our framework leverages Go Generics (`[T]`), Context is **100% type-safe** and does not rely on fragile string keys.

## How It Works

Context relies on two main concepts:
1. **Providing:** A parent component attaches a value of a specific type `T` to itself.
2. **Getting:** Any child, grandchild, or deeply nested component can ask the framework to traverse up the tree and find the nearest provided value of type `T`.

---

## Core API

The framework exposes three global functions for managing Context:

### `ProvideContext[T](c Component, value T)`
Publishes a static value of type `T` to the specified component `c`. All descendants of this component will be able to retrieve this value.

### `GetContext[T](c Component) T`
Retrieves the nearest Context of type `T`. It starts at component `c` and walks up the virtual component tree until it finds a component that provided type `T`. 

### `ProvideContextDefer[T](c Component, provider func() T)`
A deferred (lazy) provider. Instead of publishing a static value immediately, you provide a function that returns the value. This function is evaluated during the rendering phase. 
**Why is this useful?** It is extremely helpful when the Context is meant to share the component's internal `State` or a value that is not fully initialized during the component's constructor phase.

---

## Usage Examples

### 1. Basic Static Context (Theming)
A classic use case is providing a global configuration, like a UI theme, from the root layout down to individual buttons.

```go
// 1. Define a specific type for your context to ensure type safety
type ThemeContext string

// --- Parent Component ---
type RootLayout struct {
    *nn.BaseComponent
}

func NewRootLayout() *RootLayout {
    c := &RootLayout{}
    // Provide the context early in the lifecycle
    nn.ProvideContext[ThemeContext](c, "dark-mode")
    return c
}

func (c *RootLayout) View() nn.Node {
    return nn.Div().Children(
        nn.Comp(NewDeeplyNestedButton()),
    )
}

// --- Child Component ---
type NestedButton struct {
    *nn.BaseComponent
}

func (c *NestedButton) View() nn.Node {
    // Retrieve the context from anywhere in the tree!
    theme := nn.GetContext[ThemeContext](c)
    
    return nn.Button().
        Class(string(theme)). // applies "dark-mode"
        Text("Submit")
}
```

### 2. Deferred Context (Sharing State)
Sometimes, a parent component needs to share its internal state (via nn.State[T]) with its children. Because state might be mutated after the component object is created, ProvideContextDefer ensures the children always get the correct, up-to-date reference during the render cycle.

```go
// The shared state data
type TableState struct {
    SortColumn string
    Ascending  bool
}

// --- Parent Component ---
type DataTable struct {
    *nn.BaseComponent
    nn.State[TableState]
}

func NewDataTable() *DataTable {
    c := &DataTable{}
    c.InitState(func() *TableState {
        return &TableState{SortColumn: "id", Ascending: true}
    })

    // Use Defer! The func will be called when the framework actually
    // renders the children, ensuring it captures the latest c.Data.
    nn.ProvideContextDefer[TableState](c, func() TableState {
        return *c.Data
    })

    return c
}

func (c *DataTable) View() nn.Node {
    return ui.Div().Children(
        // Renders children that will consume the deferred TableConfig
        nn.Comp(NewTableHeader()), 
    )
}
```

::: tip Combining Signals and Context
While Context is great for static configurations or structured state, you can also Provide a Signal through Context (nn.ProvideContext[*nn.Signal[int]](...)). This gives you the ultimate architecture: deeply injected, globally available, granularly reactive state!
:::
