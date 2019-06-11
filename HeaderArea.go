package main

import (
	"github.com/aarzilli/nucular"
)

var firstHeader = HttpHeaderDisplay{
	true,
	"Content-Type",
	"application/json",
}

var secondHeader = HttpHeaderDisplay{
	false,
	"Accept",
	"application/json",
}

var displayHeadersInputs []*nucular.TextEditor
var displayHeadersList = []*HttpHeaderDisplay{}

type HttpHeaderDisplay struct {
	enabled bool
	name    string
	value   string
}

// func renderRequestHeaders(w *nucular.Window,
// 	displayHeadersList *[]*HttpHeaderDisplay,
// 	displayHeadersInputList *[]*nucular.TextEditor) {
// 	w.Row(50).Dynamic(2)

// 	displayHeaderStructs := *displayHeadersList
// 	displayHeaderInputs := *displayHeadersInputList

// 	for index, headerDisplay := range displayHeaderStructs {
// 		w.CheckboxText(headerDisplay.name, &headerDisplay.enabled)
// 		displayHeaderInputs[index].Edit(w)
// 	}
// }

var addingHeader = false

func renderRequestHeaders(w *nucular.Window) {

	w.Row(50).Dynamic(2)

	for index, headerDisplay := range displayHeadersList {
		w.CheckboxText(headerDisplay.name, &headerDisplay.enabled)
		displayHeadersInputs[index].Edit(w)
	}

	if addingHeader {
		w.Row(50).Dynamic(2)
		var textEditor = nucular.TextEditor{}
		textEditor.Buffer = []rune("key")
		textEditor.Flags = nucular.EditSelectable | nucular.EditClipboard
		textEditor.Edit(w)

		var anotherTextEditor = nucular.TextEditor{}
		anotherTextEditor.Buffer = []rune("value")
		anotherTextEditor.Flags = nucular.EditSelectable | nucular.EditClipboard
		anotherTextEditor.Edit(w)

	}
}

func setupHeaderAreaFromHistory(headers *map[string]string) {

	asValue := *headers

	// Update header fields with values from history
	if len(asValue) > 0 {

		for key, value := range asValue {

			asHeaderDisplay := HttpHeaderDisplay{
				false,
				key,
				value,
			}

			displayHeadersList = append(displayHeadersList, &asHeaderDisplay)
		}
	} else {
		// Add some default headers
		displayHeadersList = append(displayHeadersList, &firstHeader)
		displayHeadersList = append(displayHeadersList, &secondHeader)
	}
}

func getEnabledHeaders() map[string]string {

	var resultingRequestHeaders = make(map[string]string)

	for index, el := range displayHeadersList {

		val := string(displayHeadersInputs[index].Buffer)

		// For some reason the active/enabled value of nucular.CheckboxText is inverted
		// Flip it
		if !el.enabled {
			resultingRequestHeaders[el.name] = val
		}
	}

	return resultingRequestHeaders
}

// Function to display input fields for adding a header in a flow
func startAddHeaderProcess() {

	// Add one to displayHeadersList
	// Add one to displayHeadersInputs

	vutffac := &HttpHeaderDisplay{
		enabled: false,
		name:    "myface",
		value:   "hiface",
	}

	displayHeadersList = append(displayHeadersList, vutffac)

	var textEditor = nucular.TextEditor{}
	textEditor.Buffer = []rune(vutffac.value)
	textEditor.Flags = nucular.EditSelectable | nucular.EditClipboard
	displayHeadersInputs = append(displayHeadersInputs, &textEditor)

	addingHeader = true

}
