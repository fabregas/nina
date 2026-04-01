package nn

import (
	"strings"
	"syscall/js"
)

var currentPath string

func init() {
	currentPath = global.Get("location").Get("pathname").String()

	// listen Back/Forward buttons in browser
	global.Call("addEventListener", "popstate", js.FuncOf(func(this js.Value, args []js.Value) any {
		currentPath = js.Global().Get("location").Get("pathname").String()

		nina.scheduleUpdate(nil)
		return nil
	}))
}

func CurrentPath() string {
	return currentPath
}

func Navigate(path string) {
	if currentPath == path {
		return
	}

	js.Global().Get("history").Call("pushState", nil, "", path)
	currentPath = path

	nina.scheduleUpdate(nil)
}

type RouteDef struct {
	path  string
	exact bool
	node  Node
}

func Route(path string, node Node) RouteDef {
	return RouteDef{path: path, exact: true, node: node}
}

func RoutePrefix(path string, node Node) RouteDef {
	return RouteDef{path: path, exact: false, node: node}
}

func Router(fallback Node, routes ...RouteDef) Node {
	current := CurrentPath()

	for _, r := range routes {
		if r.exact {
			if current == r.path {
				return r.node
			}
		} else {
			if current == r.path || strings.HasPrefix(current, r.path+"/") {
				return r.node
			}
		}
	}

	return fallback
}
