# State Management Guide: State, Signals & Context

When building complex applications with the Nina framework, you have three distinct tools for managing and distributing data: **State**, **Signals**, and **Context**. 

While they might seem to overlap, each serves a very specific architectural purpose. Understanding how they differ and how they synergize is the key to writing clean, performant, and maintainable WebAssembly applications.

---

## The Cheat Sheet

Here is a quick comparison of the three tools:

| Feature | `nn.State[T]` (State) | `nn.Signal[T]` (Signals) | `nn.Get/ProvideContext` (Context) |
| :--- | :--- | :--- | :--- |
| **What is it?** | A persistent data backpack for a single component. | A reactive, thread-safe data wire. | A dependency injection portal. |
| **Scope** | Strictly Local (Component level). | Anywhere (Local, Global, or Shared). | Hierarchical (Parent down to descendants). |
| **Reactivity** | **Passive:** Requires manual `c.Update()`. | **Active:** Automatically triggers `View()` on subscribers. | **None:** Just a delivery mechanism (unless delivering a Signal). |
| **Best For** | Form fields, local temporary UI states. | Data fetched from APIs, global states, counters. | Theming, routing, avoiding "Prop Drilling". |

---

## When to Use What? (Decision Guide)

### 1. Use `nn.State[T]` when:
- The data strictly belongs to **one specific component** and nowhere else.
- You have multiple related fields (e.g., a user filling out a `ProfileForm` with 10 different text inputs). 
- **Rule of thumb:** If the component is destroyed, should the data disappear forever? If yes, use State.

### 2. Use `nn.Signal[T]` when:
- The data dictates **what is visible on the screen** across multiple components.
- The data is updated asynchronously (e.g., a background goroutine fetching an API). Signals are protected by Mutexes, making them safe for concurrent updates.
- **Rule of thumb:** If changing this variable should instantly automatically update the UI without you having to write `c.Update()`, use a Signal.

### 3. Use Context when:
- You need to pass configuration or data to a deeply nested component, but you don't want to pass it through every intermediate component as a property.
- **Rule of thumb:** If you find yourself writing `NewComponent(data)` just to immediately pass that data into `NewChildComponent(data)`, you should use Context instead.

---

## The Ultimate Combo: Context + Signals

The true power of the Nina framework unlocks when you combine these tools. The most common and powerful pattern is **injecting a Signal through Context**.

This gives you the ultimate architecture: deeply injected, globally accessible, and granularly reactive state.

### Example: A Global Notification System
Imagine a root layout that holds the state of a "Toast" notification, but any button deep inside the application can trigger it.

**1. The Parent (Provides the Signal):**
```go
type NotificationContext *nn.Signal[string]

type RootLayout struct {
    *nn.BaseComponent
    notificationSig *nn.Signal[string]
}

func NewRootLayout() *RootLayout {
    c := &RootLayout{
        // Initialize the reactive signal
        notificationSig: nn.NewSignal[string](""),
    }
    
    // Inject the signal into the context tree
    nn.ProvideContext[NotificationContext](c, c.notificationSig)
    return c
}

func (c *RootLayout) View() nn.Node {
    // Read the signal for the UI. (Subscribes RootLayout to changes)
    msg := c.notificationSig.Get(c)

    return nn.Div().Children(
        // Render the notification if the signal isn't empty
        nn.If(msg != "", nn.Div().Class("toast").Text(msg)),
        
        // Render the rest of the app...
        nn.Comp(NewDeeplyNestedPage()),,
    )
}
```

**2. The Deeply Nested Child (Consumes and Updates):**

Notice how this child component doesn't need to know anything about RootLayout. It just grabs the signal from the ether and updates it.

```go
type SaveButton struct {
    *nn.BaseComponent
}

func (c *SaveButton) View() nn.Node {
    return nn.Button().
        OnClick(func(e nn.Event) {
            // 1. Retrieve the Signal from Context
            sig := nn.GetContext[NotificationContext](c)
            
            // 2. Set a new value.
            // This will AUTOMATICALLY trigger RootLayout.View() to show the toast!
            if sig != nil {
                sig.Set("Data saved successfully!")
            }
        }).
        Children(nn.Text("Save Data"))
}
```

## Summary

* Use State to remember things locally.

* Use Context to teleport data through the tree.

* Use Signals to make that data alive and reactive.
