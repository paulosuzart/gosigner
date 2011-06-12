package sugo 

import (
	"appengine"
	"http"
	"log"
	"json"
)

//A REST Resource.
type Resource struct {
	//Used by sugo.Context.Read and sugo.Context.Render
        //to use the correct parse for http body. Ex.:
        //"application/json"
        Accepts, Renders string
        //The Paths supported
	GET, POST, PUT, DELETE   *Path
}

//The SuGo Context. Single point of interaction with appengine api and
//Requst/ResponseWriter.
type Context struct {
	ctx      appengine.Context
	w        http.ResponseWriter
	r        *http.Request
	resource *Resource
}

//Represents the Path for a given resource
//handler by a given HTTP Verb.
type Path struct {
        //Regext to match request
	Pattern string
        //The actual function that handles request.
        Handler SuGoHandler
}

//Any function that receives sugo.Context as arguments.
type SuGoHandler func(c Context)

//The single SuGo instance
var _sugo *SuGo
        
//Adds a resource to SuGo
func  Add(r *Resource) {
	if r.GET != nil {
		_sugo.routes[r.GET.Pattern] = r
	}
	if r.POST != nil {
		_sugo.routes[r.POST.Pattern] = r
	}
}

//Forces SuGo to read a jSON body and put the resoult into
//data.
func (self *Context) ReadJSON(data interface{}) {
	dec := json.NewDecoder(self.r.Body)
	dec.Decode(data)
}

//Forces SuGo to write into w the data as JSON.
func (self *Context) RenderJSON(data interface{}) {
	self.w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(self.w).Encode(data)
}

//Renders the result based on Resouce.Renders attribute
func (self *Context) Render(data interface{}) {
	self.w.Header().Set("Content-Type", self.resource.Renders)
	if self.resource.Renders == "application/json" {
		json.NewEncoder(self.w).Encode(data)
	}
}

//The SuGo struct that holds the root context (usually '/')
//and the routes.
type SuGo struct {
        root string
        routes map[string]*Resource
}

func Make(root string) {
        _sugo = &SuGo{root, map[string]*Resource{}}
}

//Register the main httpHandler provided by SuGo.
func Start() {
        //TODO if started, not start again
        http.HandleFunc(_sugo.root, httpHandler)
}

//Receives the Http request and decides where to send it.        
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
		resource.GET.Handler(ctx)
	case "POST":
		resource.POST.Handler(ctx)
	}
}

//Search in the SuGo.routes the correct resource.
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
