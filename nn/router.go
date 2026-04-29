package nn

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"syscall/js"
)

type RouteContext struct {
	Path   string
	Params map[string]string // for /users/:id
	Query  url.Values        // for ?page=2
}

type Handler func(ctx RouteContext) Component

type Middleware func(next Handler) Handler

type route struct {
	segments []string
	factory  Handler
}

func (r *route) match(path string) (map[string]string, bool) {
	path = strings.Trim(path, "/")

	if path == "" && len(r.segments) == 1 && r.segments[0] == "" {
		return make(map[string]string), true
	}

	pathSegments := strings.Split(path, "/")

	if len(pathSegments) != len(r.segments) {
		return nil, false
	}

	params := make(map[string]string)

	for i, seg := range r.segments {
		if strings.HasPrefix(seg, ":") {
			// this is param, for example ":id"
			paramName := seg[1:]
			params[paramName] = pathSegments[i]
		} else if seg != pathSegments[i] {
			// static segments are not the same (for example "users" != "posts")
			return nil, false
		}
	}

	return params, true
}

func (r *route) Use(middlewares ...Middleware) *route {
	for i := len(middlewares) - 1; i >= 0; i-- {
		r.factory = middlewares[i](r.factory)
	}
	return r
}

type Router struct {
	BaseComponent

	routes   []*route
	notFound func(RouteContext) Component

	// current 'page' that we render
	current Component

	initialized bool
}

// --------------------- global state -----
var (
	activeRouters []*Router
	historyMu     sync.Mutex
)

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) OnMount() {
	historyMu.Lock()
	activeRouters = append(activeRouters, r)
	historyMu.Unlock()
}

func (r *Router) OnDestroy() {
	historyMu.Lock()
	for i, act := range activeRouters {
		if act == r {
			activeRouters = append(activeRouters[:i], activeRouters[i+1:]...)
			break
		}
	}
	historyMu.Unlock()
}

func (r *Router) View() Node {
	if !r.initialized {
		loc := js.Global().Get("window").Get("location")
		fullPath := loc.Get("pathname").String() + loc.Get("search").String()

		r.resolve(fullPath)
		r.initialized = true
	}

	if r.current == nil {
		return Div()
	}

	return Div().Children(Comp(r.current))
}

func (r *Router) NotFound(f func(RouteContext) Component) {
	r.notFound = f
}

func (r *Router) Add(pattern string, factory func(RouteContext) Component) *route {
	segments := strings.Split(strings.Trim(pattern, "/"), "/")
	rt := &route{
		segments: segments,
		factory:  factory,
	}

	r.routes = append(r.routes, rt)
	return rt
}

func (r *Router) handlePath(rawURL string) {
	r.resolve(rawURL)
	// local re-render
	Update(r)
}

func (r *Router) resolve(rawURL string) {
	// 1. parse URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println("[Nina Router] Error parsing URL:", err)
		return
	}
	cleanPath := parsedURL.Path      // "/users/42"
	queryParams := parsedURL.Query() // map["sort"]:["asc"]

	var matchedComponent Component
	var ctx RouteContext

	// 2. find matched route
	for _, route := range r.routes {
		params, isMatch := route.match(cleanPath)
		if isMatch {
			ctx = RouteContext{
				Path:   cleanPath,
				Params: params,
				Query:  queryParams,
			}

			// call route's factory for receive component object
			matchedComponent = route.factory(ctx)
			break
		}
	}

	// 3. 404 handling
	if matchedComponent == nil {
		if r.notFound != nil {
			ctx = RouteContext{
				Path:  cleanPath,
				Query: queryParams,
			}
			matchedComponent = r.notFound(ctx)
		} else {
			matchedComponent = &defaultNotFound{path: cleanPath}
		}
	}

	// 4. update router state
	if matchedComponent != nil {
		r.current = matchedComponent
	}
}

func Navigate(path string) {
	global.Get("history").Call("pushState", nil, "", path)
	notifyRouters(path)
}

// should be called once (for example in nn.Mount)
func initHistory() {
	cb := js.FuncOf(func(this js.Value, args []js.Value) any {
		loc := js.Global().Get("window").Get("location")
		fullPath := loc.Get("pathname").String() + loc.Get("search").String()
		notifyRouters(fullPath)
		return nil
	})
	global.Get("window").Call("addEventListener", "popstate", cb)
}

func notifyRouters(path string) {
	historyMu.Lock()
	defer historyMu.Unlock()
	for _, r := range activeRouters {
		r.handlePath(path)
	}
}

type defaultNotFound struct {
	BaseComponent

	path string
}

func (d *defaultNotFound) View() Node {
	return Div().
		Style("font-family: sans-serif; text-align: center; padding: 50px; color: #333;").
		Children(
			H1().Text("404"),
			P().Text("Page not found"),
			Div().
				Style("margin-top: 20px; color: #666; font-style: italic;").
				Children(
					Span().Text("Path: "),
					Code().
						Style("background: #eee; padding: 2-4px; border-radius: 4px;").
						Text(d.path),
				),
			Div().
				Style("margin-top: 30px;").
				Children(
					A().
						Href("/").
						Style("color: #007bff; text-decoration: none; font-weight: bold;").
						Text("← Go to main page").
						OnClick(func(ev Event) {
							ev.PreventDefault()
							Navigate("/")
						}),
				),
		)
}

type redirectComp struct {
	BaseComponent

	Target string
}

func (c *redirectComp) View() Node {
	return Div()
}

func (c *redirectComp) OnMount() {
	go func() {
		Navigate(c.Target)
	}()
}

func Redirect(target string) *redirectComp {
	return &redirectComp{Target: target}
}
