package controllers

var (
	controllerlist = make(map[string]*Controller)
)

type Action struct {
	Name     string
	Run      map[string]interface{}
	Template string
	Layout   string
}

type Controller struct {
	Name    string
	Actions map[string]*Action
}

func GenerateControllers() map[string]*Controller {

	controllerlist["home"] = &Controller{
		Name:    "home",
		Actions: homeActions(),
	}

	return controllerlist
}
