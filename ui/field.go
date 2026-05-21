package ui

import "github.com/fabregas/nina/nn"

func FieldGroup() *simpleBuilder {
	return simple("div").
		Attr("data-slot", "field-group").
		Class("group/field-group @container/field-group flex w-full flex-col gap-7 data-[slot=checkbox-group]:gap-3 *:data-[slot=field-group]:gap-4")
}

// ---------------------------------------------

type fieldBuilder struct {
	baseBuilder[*fieldBuilder]

	orientationAttr string
}

func Field() *fieldBuilder {
	baseClass := "group/field flex w-full gap-3 data-[invalid=true]:text-destructive"

	b := &fieldBuilder{}
	b.baseBuilder = base(b, "div")
	b.Attr("role", "group").
		Attr("data-slot", "field").
		Class(baseClass)

	return b
}

func (f *fieldBuilder) OrientationVertical() *fieldBuilder {
	f.orientationAttr = "vertical"
	return f
}

func (f *fieldBuilder) OrientationHorizontal() *fieldBuilder {
	f.orientationAttr = "horizontal"
	return f
}

func (f *fieldBuilder) OrientationResponsive() *fieldBuilder {
	f.orientationAttr = "responsive"
	return f
}

func (f *fieldBuilder) build(ctx *buildContext) {
	var orientationClass string
	switch f.orientationAttr {
	case "horizontal":
		orientationClass = "flex-row items-center has-[>[data-slot=field-content]]:items-start *:data-[slot=field-label]:flex-auto has-[>[data-slot=field-content]]:[&>[role=checkbox],[role=radio]]:mt-px"
	case "responsive":
		orientationClass = "flex-col *:w-full @md/field-group:flex-row @md/field-group:items-center @md/field-group:*:w-auto @md/field-group:has-[>[data-slot=field-content]]:items-start @md/field-group:*:data-[slot=field-label]:flex-auto [&>.sr-only]:w-auto @md/field-group:has-[>[data-slot=field-content]]:[&>[role=checkbox],[role=radio]]:mt-px"
	default: // vertical
		orientationClass = "flex-col *:w-full [&>.sr-only]:w-auto"
	}

	ctx.Props.Class(orientationClass)
	ctx.Props.Attr("data-orientation", f.orientationAttr)
}

// ---------------------------------------------

func FieldLabel() *simpleBuilder {
	baseClasses := "group/field-label peer/field-label flex w-fit gap-2 leading-snug group-data-[disabled=true]/field:opacity-50 has-data-checked:bg-input/30 has-[>[data-slot=field]]:rounded-2xl has-[>[data-slot=field]]:border *:data-[slot=field]:p-4 has-[>[data-slot=field]]:w-full has-[>[data-slot=field]]:flex-col"

	return simple("label").
		Attr("data-slot", "field-label").
		Class(baseClasses)
}

// ---------------------------------------------

func FieldTitle() *simpleBuilder {
	baseClasses := "flex w-fit items-center gap-2 text-sm leading-snug font-medium group-data-[disabled=true]/field:opacity-50"

	return simple("div").
		Attr("data-slot", "field-label").
		Class(baseClasses)
}

// ---------------------------------------------

func FieldDescription() *simpleBuilder {
	baseClasses := "text-left text-sm leading-normal font-normal text-muted-foreground group-has-data-horizontal/field:text-balance [[data-variant=legend]+&]:-mt-1.5 last:mt-0 nth-last-2:-mt-1 [&>a]:underline [&>a]:underline-offset-4 [&>a:hover]:text-primary"

	return simple("p").
		Attr("data-slot", "field-description").
		Class(baseClasses)
}

// ---------------------------------------------

type fieldErrorBuilder struct {
	baseBuilder[*fieldErrorBuilder]

	errors []string
}

func FieldError() *fieldErrorBuilder {
	b := &fieldErrorBuilder{}
	b.baseBuilder = base(b, "div")
	b.Attr("role", "alert").
		Attr("data-slot", "field-error").
		Class("text-sm font-normal text-destructive")

	return b
}

func (e *fieldErrorBuilder) Errors(errs ...string) *fieldErrorBuilder {
	e.errors = errs
	return e
}

func (e *fieldErrorBuilder) build(ctx *buildContext) {
	if len(ctx.Children) > 0 {
		return
	}

	var uniqueErrors []string
	seen := make(map[string]bool)

	for _, err := range e.errors {
		if err != "" && !seen[err] {
			seen[err] = true
			uniqueErrors = append(uniqueErrors, err)
		}
	}

	if len(uniqueErrors) == 0 {
		return
	}

	if len(uniqueErrors) == 1 {
		ctx.Children = append(ctx.Children, nn.Text(uniqueErrors[0]))
		return
	}

	var listItems []nn.AsNode
	for _, err := range uniqueErrors {
		listItems = append(listItems, nn.Li().Text(err))
	}

	ul := nn.Ul().
		Class("ml-4 flex list-disc flex-col gap-1").
		Children(listItems...)

	ctx.Children = append(ctx.Children, ul)
}
