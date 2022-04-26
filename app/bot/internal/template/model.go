package template

type MStartWelcome struct {
	ID       int64
	Username string
	Notice   string
}

type MMyDesc struct {
	Current int
	BindMax int
}

type MMyView struct {
	Alias        string
	ClientID     string
	ClientSecret string
	UpdatedAt    string
}
