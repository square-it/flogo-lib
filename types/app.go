package types

// App is the type for the App
type App struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	//Triggers    []TriggerConfig `json:"triggers"`
}
