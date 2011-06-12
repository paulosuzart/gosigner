package gosigner

import (
	"appengine"
	"http"
	"log"
	"json"
)

type Resource struct {
	Accepts, Renders string
	GET, POST, PUT, DELETE   *Path
}

type Context struct {
	ctx      appengine.Context
	w        http.ResponseWriter
	r        *http.Request
	resource *Resource
}

type Path struct {
	pattern string
	handler SuGoHandler
}

type SuGoHandler func(c Context)

var _sugo *SuGo
        
func  Add(r *Resource) {
	if r.GET != nil {
		_sugo.routes[r.GET.pattern] = r
	}
	if r.POST != nil {
		_sugo.routes[r.POST.pattern] = r
	}
}

func (self *Context) ReadJSON(data interface{}) {
	dec := json.NewDecoder(self.r.Body)
	dec.Decode(data)
}

func (self *Context) RenderJSON(data interface{}) {
	self.w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(self.w).Encode(data)
}

func (self *Context) Render(data interface{}) {

	self.w.Header().Set("Content-Type", self.resource.Renders)
	if self.resource.Renders == "application/json" {
		json.NewEncoder(self.w).Encode(data)
	}
}

type SuGo struct {
        root string
        routes map[string]*Resource
}

func MakeSuGo(root string) {
        _sugo = &SuGo{root, map[string]*Resource{}}
}

func Start() {
        //TODO if started, not start again
        http.HandleFunc(_sugo.root, httpHandler)
}
func httpHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var ctx Context
	ctx.ctx = c
	ctx.w = w
	ctx.r = r
	resource := _sugo.findHandler(r.URL.Path)
	ctx.resource = resource
	switch r.Method {
	case "GET":
		resource.GET.handler(ctx)
	case "POST":
		resource.POST.handler(ctx)
	}
}
func (self *SuGo) findHandler(path string) *Resource {
	log.Print(path)
	return self.routes[path]
}
func ensurePOST(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		h(w, r)
	}
}

func ensureGET(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		h(w, r)
	}
}
