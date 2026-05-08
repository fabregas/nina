# State: Persistent Component Data

While [Signals](./signals.md) are reactive primitives perfect for granular updates and shared data, sometimes a component needs to manage a complex, purely local data structure. 

In a Virtual DOM framework, component structs are often recreated from scratch when their parent re-renders. If you use standard Go variables inside your component struct, their values will be lost during this recreation. 

To solve this, the framework provides the `nn.State[T]` embedding. It ensures that your local data **survives** component recreation and persists across the entire lifecycle of the component in the DOM.

## How It Works

1. **Define a Data Struct:** Create a standard Go struct to hold your state fields.
2. **Embed `nn.State[T]`:** Embed this type into your component, passing your data struct as the generic type parameter.
3. **Access via `c.Data`:** Read or mutate your fields using the `.Data` property provided by the embedding.

Under the hood, when the framework's reconciliation engine diffs the Virtual DOM, it recognizes the `nn.State[T]` embedding and automatically copies the state from the previous component instance into the newly created one.

---

## Usage Example

Here is an example of a form component that preserves user input, even if the parent component triggers a full re-render of the page.

```go
// 1. Define the shape of your local state
type ProfileFormState struct {
    Username string
    Age      int
    HasError bool
}

// 2. Embed nn.State with your custom type
type ProfileForm struct {
    *nn.BaseComponent
    nn.State[ProfileFormState]
}

// Constructor
func NewProfileForm() *ProfileForm {
    c := &ProfileForm{}
    // Initialize default state values
    c.Data = &ProfileFormState{
        Username: "Guest",
        Age:      18,
    }
    return c
}

func (c *ProfileForm) View() nn.Node {
    return nn.Div().
        Class("flex flex-col gap-4 p-6").
        Children(
            nn.Text(fmt.Sprintf("Editing profile: %s", c.Data.Username)),

            nn.Input().
                Type("text").
                Value(c.Data.Username). // Read from preserved state
                OnInput(func(e nn.Event) {
                    // Mutate the state
                    c.Data.Username = e.TargetValue()
                    
                    // component will be updated automatically
                }),
        )
}
```

## State vs. Signals Comparison

| Feature | `nn.State[T]` | `nn.Signal[T]` |
| :--- | :--- | :--- |
| **Primary Use Case** | Bundling multiple related fields of strictly local component data (e.g., form fields). | Shared state, deep reactivity, or individual reactive values. |
| **Reactivity** | Passive. You mutate `c.Data` and must call `c.Update()` to trigger a render. | Active. Calling `.Set()` automatically triggers renders for subscribers. |
| **Scope** | Tied strictly to the component's DOM lifecycle. | Can be passed around, shared via Context, or declared globally. |
| **Data Access** | Direct field access (`c.Data.Name = "John"`). | Method access (`sig.Get/Set`). |
