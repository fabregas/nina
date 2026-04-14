# Button

Displays a button or a component that looks like a button. It is a fundamental building block of **Nina UI**, supporting multiple variants, sizes, and the powerful Slot pattern for semantic HTML.

## Basic Usage

By default, the `Button` builder generates a standard HTML `<button>` element.

```go
import (
    "github.com/fabregas/nina/ui"
)

// Inside your component or main function:
ui.Button().Text("Save Changes")
```

### Generated HTML
```html
<button class="inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground shadow hover:bg-primary/90 h-9 px-4 py-2">
  Save Changes
</button>
```

---

## Variants

Nina UI buttons come with several pre-defined style variants inspired by shadcn/ui. Use the `Variant()` method to change the appearance.

  

<Preview height="100px">

<div class="flex gap-4">
<button class="group/button text-sm gap-1.5 transition-all items-center border-transparent disabled:opacity-50 [&amp;_svg]:pointer-events-none bg-primary px-3 outline-none focus-visible:border-ring whitespace-nowrap [&amp;_svg]:shrink-0 hover:bg-primary/80 font-medium text-primary-foreground has-data-[icon=inline-start]:pl-2.5 inline-flex rounded-4xl aria-invalid:border-destructive justify-center focus-visible:ring-3 select-none active:not-aria-[haspopup]:translate-y-px dark:aria-invalid:border-destructive/50 h-9 focus-visible:ring-ring/30 dark:aria-invalid:ring-destructive/40 border bg-clip-padding shrink-0 [&amp;_svg:not([class*='size-'])]:size-4 aria-invalid:ring-3 aria-invalid:ring-destructive/20 has-data-[icon=inline-end]:pr-2.5 disabled:pointer-events-none" data-slot="button">Default</button><button class="group/button bg-clip-padding outline-none text-secondary-foreground transition-all select-none bg-secondary px-3 [&amp;_svg]:shrink-0 aria-expanded:bg-secondary shrink-0 font-medium items-center disabled:pointer-events-none disabled:opacity-50 rounded-4xl hover:bg-secondary/80 gap-1.5 has-data-[icon=inline-end]:pr-2.5 focus-visible:border-ring active:not-aria-[haspopup]:translate-y-px dark:aria-invalid:ring-destructive/40 h-9 aria-invalid:ring-destructive/20 [&amp;_svg:not([class*='size-'])]:size-4 has-data-[icon=inline-start]:pl-2.5 whitespace-nowrap focus-visible:ring-3 justify-center border-transparent text-sm dark:aria-invalid:border-destructive/50 border aria-invalid:ring-3 aria-invalid:border-destructive [&amp;_svg]:pointer-events-none inline-flex focus-visible:ring-ring/30 aria-expanded:text-secondary-foreground" data-slot="button">Secondary</button><button class="group/button text-sm [&amp;_svg:not([class*='size-'])]:size-4 bg-destructive/10 h-9 font-medium outline-none whitespace-nowrap active:not-aria-[haspopup]:translate-y-px aria-invalid:border-destructive dark:bg-destructive/20 select-none dark:aria-invalid:border-destructive/50 [&amp;_svg]:shrink-0 dark:focus-visible:ring-destructive/40 rounded-4xl focus-visible:ring-destructive/20 gap-1.5 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 has-data-[icon=inline-end]:pr-2.5 items-center border-transparent disabled:pointer-events-none text-destructive dark:hover:bg-destructive/30 aria-invalid:ring-3 px-3 bg-clip-padding disabled:opacity-50 hover:bg-destructive/20 focus-visible:border-destructive/40 border has-data-[icon=inline-start]:pl-2.5 focus-visible:ring-3 inline-flex justify-center [&amp;_svg]:pointer-events-none shrink-0 transition-all" data-slot="button">Destructive</button><button class="group/button aria-invalid:ring-destructive/20 bg-background has-data-[icon=inline-end]:pr-2.5 shrink-0 [&amp;_svg]:shrink-0 [&amp;_svg:not([class*='size-'])]:size-4 font-medium disabled:opacity-50 items-center border-border text-sm aria-expanded:text-foreground focus-visible:ring-3 [&amp;_svg]:pointer-events-none rounded-4xl border select-none hover:bg-muted dark:hover:bg-input/30 transition-all aria-invalid:border-destructive hover:text-foreground dark:bg-transparent whitespace-nowrap bg-clip-padding outline-none dark:aria-invalid:border-destructive/50 h-9 inline-flex gap-1.5 px-3 disabled:pointer-events-none aria-expanded:bg-muted aria-invalid:ring-3 focus-visible:ring-ring/30 has-data-[icon=inline-start]:pl-2.5 justify-center focus-visible:border-ring active:not-aria-[haspopup]:translate-y-px dark:aria-invalid:ring-destructive/40" data-slot="button">Outline</button><button class="group/button border-transparent bg-clip-padding hover:text-foreground focus-visible:ring-3 aria-invalid:ring-destructive/20 aria-invalid:border-destructive dark:aria-invalid:ring-destructive/40 dark:hover:bg-muted/50 font-medium text-sm has-data-[icon=inline-end]:pr-2.5 select-none [&amp;_svg:not([class*='size-'])]:size-4 aria-expanded:bg-muted h-9 has-data-[icon=inline-start]:pl-2.5 justify-center outline-none focus-visible:border-ring dark:aria-invalid:border-destructive/50 hover:bg-muted rounded-4xl border active:not-aria-[haspopup]:translate-y-px [&amp;_svg]:shrink-0 items-center aria-expanded:text-foreground gap-1.5 shrink-0 aria-invalid:ring-3 focus-visible:ring-ring/30 px-3 inline-flex whitespace-nowrap disabled:opacity-50 disabled:pointer-events-none [&amp;_svg]:pointer-events-none transition-all" data-slot="button">Ghost</button><button class="group/button rounded-4xl focus-visible:ring-ring/30 dark:aria-invalid:border-destructive/50 has-data-[icon=inline-end]:pr-2.5 aria-invalid:ring-destructive/20 border-transparent text-sm hover:underline aria-invalid:border-destructive px-3 aria-invalid:ring-3 justify-center bg-clip-padding disabled:pointer-events-none items-center border select-none focus-visible:border-ring active:not-aria-[haspopup]:translate-y-px disabled:opacity-50 [&amp;_svg]:pointer-events-none text-primary gap-1.5 font-medium whitespace-nowrap outline-none underline-offset-4 h-9 has-data-[icon=inline-start]:pl-2.5 dark:aria-invalid:ring-destructive/40 shrink-0 focus-visible:ring-3 [&amp;_svg:not([class*='size-'])]:size-4 transition-all [&amp;_svg]:shrink-0 inline-flex" data-slot="button">Link</button>
</div>

</Preview>


```go
nn.Div().Class("flex gap-4").Children(
    ui.Button().Text("Default"),
    ui.Button().Secondary().Text("Secondary"),
    ui.Button().Destructive().Text("Destructive"),
    ui.Button().Outline().Text("Outline"),
    ui.Button().Ghost().Text("Ghost"),
    ui.Button().Link().Text("Link"),
)
```

---

## Sizes

You can easily adjust the size and padding of the button using the `Size*()` methods.

<Preview height="100px">
<div class="gap-4 flex items-center"><button class="group/button outline-none text-primary-foreground focus-visible:ring-ring/30 [&amp;_svg:not([class*='size-'])]:size-4 justify-center border font-medium bg-primary has-data-[icon=inline-end]:pr-2 disabled:pointer-events-none disabled:opacity-50 transition-all has-data-[icon=inline-start]:pl-2 aria-invalid:ring-3 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 active:not-aria-[haspopup]:translate-y-px border-transparent focus-visible:ring-3 gap-1 items-center select-none dark:aria-invalid:ring-destructive/40 rounded-4xl bg-clip-padding text-sm aria-invalid:border-destructive aria-invalid:ring-destructive/20 px-3 inline-flex whitespace-nowrap focus-visible:border-ring shrink-0 hover:bg-primary/80 h-8 dark:aria-invalid:border-destructive/50" data-slot="button">Small</button><button class="group/button aria-invalid:ring-destructive/20 has-data-[icon=inline-start]:pl-2.5 outline-none select-none aria-invalid:ring-3 focus-visible:border-ring aria-invalid:border-destructive border whitespace-nowrap items-center focus-visible:ring-ring/30 disabled:opacity-50 active:not-aria-[haspopup]:translate-y-px dark:aria-invalid:border-destructive/50 [&amp;_svg:not([class*='size-'])]:size-4 focus-visible:ring-3 has-data-[icon=inline-end]:pr-2.5 rounded-4xl [&amp;_svg]:pointer-events-none text-primary-foreground gap-1.5 transition-all h-9 inline-flex shrink-0 hover:bg-primary/80 px-3 [&amp;_svg]:shrink-0 bg-primary font-medium dark:aria-invalid:ring-destructive/40 bg-clip-padding text-sm justify-center border-transparent disabled:pointer-events-none" data-slot="button">Default</button><button class="group/button rounded-4xl select-none focus-visible:ring-3 h-10 has-data-[icon=inline-end]:pr-3 border-transparent text-sm aria-invalid:border-destructive dark:aria-invalid:ring-destructive/40 inline-flex font-medium outline-none [&amp;_svg:not([class*='size-'])]:size-4 aria-invalid:ring-destructive/20 aria-invalid:ring-3 bg-primary shrink-0 whitespace-nowrap gap-1.5 justify-center text-primary-foreground focus-visible:border-ring [&amp;_svg]:shrink-0 dark:aria-invalid:border-destructive/50 disabled:pointer-events-none disabled:opacity-50 border transition-all focus-visible:ring-ring/30 active:not-aria-[haspopup]:translate-y-px items-center bg-clip-padding [&amp;_svg]:pointer-events-none hover:bg-primary/80 px-4 has-data-[icon=inline-start]:pl-3" data-slot="button">Large</button><button class="group/button font-medium aria-invalid:ring-destructive/20 bg-background focus-visible:ring-ring/30 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive disabled:opacity-50 aria-expanded:bg-muted bg-clip-padding aria-expanded:text-foreground dark:hover:bg-input/30 justify-center focus-visible:border-ring active:not-aria-[haspopup]:translate-y-px aria-invalid:ring-3 [&amp;_svg]:pointer-events-none [&amp;_svg:not([class*='size-'])]:size-3 border-border select-none dark:aria-invalid:border-destructive/50 outline-none hover:text-foreground border hover:bg-muted inline-flex text-sm whitespace-nowrap disabled:pointer-events-none items-center size-6 shrink-0 dark:bg-transparent rounded-4xl focus-visible:ring-3 transition-all [&amp;_svg]:shrink-0" data-slot="button"><svg class="shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12h14"></path> <path d="m12 5 7 7-7 7"></path></svg></button><button class="group/button border-border hover:bg-muted bg-clip-padding focus-visible:border-ring transition-all select-none focus-visible:ring-ring/30 aria-expanded:bg-muted rounded-4xl dark:bg-transparent whitespace-nowrap bg-background hover:text-foreground aria-expanded:text-foreground dark:hover:bg-input/30 focus-visible:ring-3 dark:aria-invalid:ring-destructive/40 [&amp;_svg]:shrink-0 aria-invalid:ring-destructive/20 [&amp;_svg]:pointer-events-none [&amp;_svg:not([class*='size-'])]:size-4 border disabled:pointer-events-none items-center outline-none disabled:opacity-50 text-sm inline-flex justify-center aria-invalid:ring-3 dark:aria-invalid:border-destructive/50 size-8 active:not-aria-[haspopup]:translate-y-px shrink-0 font-medium aria-invalid:border-destructive" data-slot="button"><svg class="shrink-0" stroke-linejoin="round" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M5 12h14"></path> <path d="m12 5 7 7-7 7"></path></svg></button><button class="group/button dark:aria-invalid:border-destructive/50 select-none aria-invalid:ring-destructive/20 dark:hover:bg-input/30 shrink-0 disabled:opacity-50 bg-clip-padding font-medium [&amp;_svg]:pointer-events-none size-9 disabled:pointer-events-none hover:bg-muted inline-flex items-center focus-visible:ring-3 dark:aria-invalid:ring-destructive/40 transition-all outline-none focus-visible:border-ring active:not-aria-[haspopup]:translate-y-px bg-background [&amp;_svg]:shrink-0 aria-expanded:text-foreground text-sm border aria-expanded:bg-muted whitespace-nowrap aria-invalid:border-destructive border-border focus-visible:ring-ring/30 aria-invalid:ring-3 hover:text-foreground dark:bg-transparent justify-center rounded-4xl [&amp;_svg:not([class*='size-'])]:size-4" data-slot="button"><svg class="shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" xmlns="http://www.w3.org/2000/svg"><path d="M5 12h14"></path> <path d="m12 5 7 7-7 7"></path></svg></button><button class="group/button font-medium aria-expanded:text-foreground dark:bg-transparent focus-visible:ring-ring/30 bg-background justify-center border text-sm border-border active:not-aria-[haspopup]:translate-y-px [&amp;_svg]:shrink-0 rounded-4xl transition-all aria-invalid:ring-destructive/20 dark:hover:bg-input/30 disabled:opacity-50 inline-flex outline-none disabled:pointer-events-none dark:aria-invalid:ring-destructive/40 shrink-0 aria-invalid:ring-3 [&amp;_svg:not([class*='size-'])]:size-4 aria-invalid:border-destructive size-10 items-center bg-clip-padding focus-visible:border-ring focus-visible:ring-3 whitespace-nowrap [&amp;_svg]:pointer-events-none hover:bg-muted dark:aria-invalid:border-destructive/50 hover:text-foreground aria-expanded:bg-muted select-none" data-slot="button"><svg class="shrink-0" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M5 12h14"></path> <path d="m12 5 7 7-7 7"></path></svg></button></div>
</Preview>



```go
nn.Div().Class("flex items-center gap-4").Children(
    ui.Button().SizeSm().Text("Small"),
    ui.Button().Text("Default"),
    ui.Button().SizeLg().Text("Large"),
    
    // The "icon" size creates a perfect square, ideal for Lucide icons
    ui.Button().Outline().SizeIconXs().Children(
        icons.ArrowRight(),
    ),
    ui.Button().Outline().SizeIconSm().Children(
        icons.ArrowRight(),
    ),
    ui.Button().Outline().SizeIcon().Children(
        icons.ArrowRight(),
    ),
    ui.Button().Outline().SizeIconLg().Children(
        icons.ArrowRight(),
    ),
)
```

---

## The Slot Pattern (AsChild)

One of the most powerful features of Nina UI is the **Slot Pattern**. Sometimes you need a semantic link `<a>` for SEO or routing, but you want it to look exactly like a button. 

Instead of wrapping a link inside a button (which is invalid HTML), use the `AsChild()` method. This will merge all the button's Tailwind classes onto your custom element.

```go
ui.Button().
    Outline()
    SizeLg().
    AsChild(
        nn.A().Href("https://github.com").Attr("target", "_blank"),
    ).
    Children(
        icons.Github(),
        nn.Text("View on GitHub"),
    )
```

### Generated HTML (Slot)
Notice how the output is a clean `<a>` tag, completely preserving the button's design system:
```html
<a href="[https://github.com](https://github.com)" target="_blank" class="inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-colors border border-input bg-transparent shadow-sm hover:bg-accent hover:text-accent-foreground h-10 rounded-md px-8">
  <svg class="mr-2 size-4" ...></svg>
  View on GitHub
</a>
```

---

## API Reference

TBD
