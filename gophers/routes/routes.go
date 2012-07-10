package routes

import (
	"net/http"
	"html/template"
)

func start() *RouteTable {
	rt := &RouteTable{}

	//      var rt = &Route {
	//          Name: "default",
	//          Pattern: "/{controller}/{action}/{id}",
	//			Method: "GET",
	//          Default: map[string]string { "controller": "home", "action": "index", "id": "0", },
	//          Constraint: map[string]string { "id": "\\d+" }
	//      }

	rt.AddRoute(&Route{
		Name:       "default",
		Pattern:    "/{controller}/{action}/{id}",
		Method:     "GET",
		Default:    map[string]string{"controller": "home", "action": "index", "id": "0"},
		Constraint: map[string]string{"id": "\\d+"},
	})

	return rt
}

type WebContext struct {
	w        http.ResponseWriter
	r        *http.Request
	rd       *RouteData
	vb       map[string]interface{}
	layout   string
	template string
}

func HandleRoute(w http.ResponseWriter, r *http.Request) {
	rt := start()
	rd, _ := rt.Match(w, r)
	ctx := &WebContext{w: w, r: r, rd: rd, template: "templates/" + rd.Controller + "/" + rd.Action + ".html"}
	DisplayTemplate(ctx)
}

func DisplayTemplate(ctx *WebContext) {
	if ctx.layout == "" {
		ctx.layout = "layout.html"
	}
	if ctx.vb == nil {
		ctx.vb = make(map[string]interface{})
	}
	t := template.Must(template.ParseFiles(ctx.layout, ctx.template))

	if err := t.Execute(ctx.w, ctx.vb); err != nil {
		http.Error(ctx.w, err.Error(), http.StatusInternalServerError)
	}
}

/*func Home(w http.ResponseWriter, r *http.Request) {
	helpers.DisplayTemplate("index", "templates/index.html", w, make(map[string]interface{}))
}*/
