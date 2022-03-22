package main

import "time"

type Template struct {
	Id    string  `json:"id"`
	Title string  `json:"title"`
	Items []*Item `json:"items"`
}

type Item struct {
	Id string `json:"id"`

	Check *Check `json:"check,omitempty"`
}

type Check struct {
	Text      string     `json:"text"`
	Owner     string     `json:"owner"`
	Completed *time.Time `json:"completed"`
}

func NewTemplate() *Template {
	return &Template{
		Items: []*Item{},
	}
}

func (t *Template) GetType() string {
	return "template"
}

func (t *Template) GetId() string {
	return t.Id
}

func (t *Template) SetId(id string) {
	t.Id = id
}

func (t *Template) IsValid() bool {
	return true
}
