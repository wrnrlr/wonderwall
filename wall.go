package main

import (
	"encoding/json"
	"io"
)

type Node struct {
	ID        string
	ClassName string                 `json:"className"`
	Attrs     map[string]interface{} `json:"attrs"`
}

type Action struct {
	Action   string `json:"action"`
	Node     Node   `json:"node"`
	Original *Node  `json:"original"`
}

func (a Action) Apply(s State) {
	switch a.Action {
	case "add":
		a.Add(s)
	case "remove":
		a.Remove(s)
	case "update":
		a.Update(s)
	}
}
func (a Action) Add(s State) {
	layer := s.Layer(a.Node.ClassName)
	*layer = append(*layer, a.Node)
}
func (a Action) Remove(s State) {
	layer := s.Layer(a.Node.ClassName)
	for i, n := range *layer {
		if n.ID == a.Node.ID {
			(*layer)[i].Attrs = a.Node.Attrs

			(*layer)[i] = (*layer)[len(*layer)-1] // Copy last element to index i.
			(*layer)[len(*layer)-1] = Node{}      // Erase last element (write zero value).
			*layer = (*layer)[:len(*layer)-1]

			return
		}
	}
}
func (a Action) Update(s State) {
	layer := s.Layer(a.Node.ClassName)
	for i, n := range *layer {
		if n.ID == a.Node.ID {
			(*layer)[i].Attrs = a.Node.Attrs
			return
		}
	}
}

func ReadAction(r io.Reader) (*Action, error) {
	var a Action
	err := json.NewDecoder(r).Decode(&a)
	return &a, err
}

type State struct {
	Image *[]Node `json:"image"`
	Text  *[]Node `json:"text"`
	Pen   *[]Node `json:"pen"`
}

func (s State) Layer(n string) *[]Node {
	if n == "Image" {
		return s.Image
	} else if n == "Text" {
		return s.Text
	} else if n == "Line" {
		return s.Pen
	} else {
		return nil
	}
}
