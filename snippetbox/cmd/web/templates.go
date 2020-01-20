package main

import "ellipe.party/snippetbox/pkg/models"

// templateData : holds the information required by the template, since only one dynamic element can be sent
// to the HTML Template you can use an struct to hold data comming from different sources as one.
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
