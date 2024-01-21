package optimiser

import (
	"fmt"
	"sort"
	"strings"
)

type Router struct {
	Parent     *Router
	Subrouters []*Router

	PathSegment string
	Service     *string
}

func (r *Router) String() string {
	var lines []string

	r.Walk(func(path, service string) {
		lines = append(lines, path+" -> "+service)
	})

	return strings.Join(lines, "\n")
}

func (r *Router) Walk(f func(path, service string)) {
	r.sort()

	for _, subrouter := range r.Subrouters {
		subrouter.Walk(f)
	}

	if r.Service != nil {
		f(r.Path(), *r.Service)
	}
}

func (r *Router) sort() {
	for _, subrouter := range r.Subrouters {
		subrouter.sort()
	}

	sort.Slice(r.Subrouters, func(i, j int) bool {
		return r.Subrouters[i].PathSegment > r.Subrouters[j].PathSegment
	})
}

func (r *Router) Path() string {
	if r.Parent == nil {
		return r.PathSegment
	}

	return r.Parent.Path() + "/" + r.PathSegment
}

func (r *Router) Subrouter(path string) *Router {
	segments := strings.Split(strings.TrimPrefix(path, "/"), "/")

	for _, subrouter := range r.Subrouters {
		if subrouter.PathSegment == segments[0] {
			if len(segments) == 1 {
				return subrouter
			}

			return subrouter.Subrouter("/" + strings.Join(segments[1:], "/"))
		}
	}

	subrouter := &Router{
		Parent:      r,
		Subrouters:  []*Router{},
		PathSegment: segments[0],
		Service:     nil,
	}
	r.Subrouters = append(r.Subrouters, subrouter)

	if len(segments) > 1 {
		return subrouter.Subrouter("/" + strings.Join(segments[1:], "/"))
	}
	return subrouter
}

func (r *Router) Bind(service string) error {
	if r.Service != nil {
		return fmt.Errorf("path already bound to %q", *r.Service)
	}

	r.Service = &service
	return nil
}

func (r *Router) Prune() {
	for _, subrouter := range r.Subrouters {
		subrouter.Prune()
	}

	// Subroutes that only have one service bound to them can be pruned and replaced with a single "*" route
	subServices := r.subServices()
	if len(subServices) == 1 {
		r.Subrouters = []*Router{}
		r.Subrouter("*").Bind(subServices[0])
	}
}

func (r *Router) subServices() []string {
	subServices := map[string]bool{}
	for _, subrouter := range r.Subrouters {
		if subrouter.Service != nil {
			subServices[*subrouter.Service] = true
		}

		for _, subService := range subrouter.subServices() {
			subServices[subService] = true
		}
	}

	rtn := []string{}
	for subService := range subServices {
		rtn = append(rtn, subService)
	}
	return rtn
}

func NewRouter() *Router {
	return &Router{
		Parent:      nil,
		Subrouters:  []*Router{},
		PathSegment: "",
		Service:     nil,
	}
}
