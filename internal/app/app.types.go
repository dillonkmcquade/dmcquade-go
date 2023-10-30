package app

type Project struct {
	Src         string `json:"src"`
	Url         string `json:"url"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Github      string `json:"github"`
	Youtube     string `json:"youtube"`
	Id          int    `json:"id"`
	NoInfo      bool   `json:"noInfo"`
}
type AppData struct {
	Projects       [4]Project
	ProjectDetails [3]ProjectDetails
	Skills         [10]string
}

type ProjectDetailSection struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type ProjectDetails struct {
	Id       int                    `json:"id"`
	Title    string                 `json:"title"`
	Image    string                 `json:"image"`
	Sections []ProjectDetailSection `json:"sections"`
}
