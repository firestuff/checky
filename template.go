package main

type Template struct {
	Id   string
	Test string
}

func (t *Template) GetType() string {
	return "template"
}

func (t *Template) GetId() string {
	return t.Id
}
