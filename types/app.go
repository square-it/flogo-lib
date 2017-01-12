package types

// App is the type for the App
type App struct {
	Name        string     `json:"name"`
	Version     string     `json:"version"`
	Description string     `json:"description"`
	Triggers    []*Trigger `json:"triggers"`
	Actions     []*Action  `json:"actions"`
}

// Trigger is the type for the Trigger
type Trigger struct {
	Id   string      `json:"if"`
	Ref  string      `json:"ref"`
	Data interface{} `json:"data"`
}

// Action is the type for the Action
type Action struct {
	Id   string      `json:"if"`
	Ref  string      `json:"ref"`
	Data interface{} `json:"data"`
}
