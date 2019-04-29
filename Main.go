package main

import (
	"GoRESTClient/configHandler"
	"encoding/json"
	"fmt"
	"github.com/aarzilli/nucular"
	"github.com/go-xmlfmt/xmlfmt"
	"github.com/yosssi/gohtml"
	"golang.org/x/image/colornames"
	"golang.org/x/mobile/event/key"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func LOG(s string) {
	fmt.Printf("%s - %s\n", time.Now().Format("2009-09-01T15:04"), s)
}

var Wnd nucular.MasterWindow

var Title = "Gostman"

func main() {

	Wnd = nucular.NewMasterWindow(0, Title, textEditorDemo())

	if Wnd == nil {
		_, _ = fmt.Fprintf(os.Stderr, "unknown demo %q\n", "WUT")
		_, _ = fmt.Fprintf(os.Stderr, "ggggg\n")
		os.Exit(1)
	}

	Wnd.Main()
}

var status = "django bango"
var responseContentType string = "ttt"
var tabbableFields [5]*nucular.TextEditor

type Pair struct {
	a, b string
}

// Load these default values from config file
var auth = Pair{"admin", "admin"}

var client = &http.Client{}
var httpMethod = "GET"
var HistoryStruct, _ = configHandler.GetHistoryOrCreateHistoryFile()

var history = HistoryStruct.Content

// Construct the render function
func textEditorDemo() func(w *nucular.Window) {

	var urlEditorField nucular.TextEditor
	urlEditorField.Flags = nucular.EditSelectable | nucular.EditClipboard

	// Set the buffer to the last value in the history or else empty
	urlEditorField.Buffer = []rune(history[len(history) -1])
	urlEditorField.Maxlen = 1000
	urlEditorField.Active = true

	tabbableFields[0] = &urlEditorField

	var requestBodyEditorField nucular.TextEditor
	// What determines if the field gets auto tabbed? Is it EditMultiline?
	requestBodyEditorField.Flags = nucular.EditSelectable | nucular.EditMultiline | nucular.EditClipboard
	requestBodyEditorField.Buffer = []rune("request body")
	requestBodyEditorField.Maxlen = 150

	tabbableFields[3] = &requestBodyEditorField

	var usernameEditorField nucular.TextEditor
	usernameEditorField.Flags = nucular.EditSelectable | nucular.EditClipboard
	usernameEditorField.Buffer = []rune(auth.a)
	usernameEditorField.Maxlen = 150

	tabbableFields[1] = &usernameEditorField

	var passwordEditorField nucular.TextEditor
	passwordEditorField.Flags = nucular.EditSelectable | nucular.EditClipboard
	passwordEditorField.Buffer = []rune(auth.b)
	passwordEditorField.Maxlen = 150

	tabbableFields[2] = &passwordEditorField

	var responseBodyEditor nucular.TextEditor
	responseBodyEditor.Flags = nucular.EditSelectable | nucular.EditMultiline | nucular.EditClipboard | nucular.EditReadOnly
	responseBodyEditor.Buffer = []rune("response body")

	tabbableFields[4] = &responseBodyEditor

	var history, _ = configHandler.GetHistoryOrCreateHistoryFile()

	fmt.Println("main history", history)

	return func(w *nucular.Window) {

		handleKeybindings(w, &urlEditorField, &responseBodyEditor,
			&usernameEditorField, &passwordEditorField, &requestBodyEditorField)


		w.Row(30).Static(50, 500, 125, 125)

		w.LabelColored(httpMethod, "LT", colornames.Aquamarine)

		urlEditorField.Edit(w)

		usernameEditorField.Edit(w)
		passwordEditorField.Edit(w)

		if (httpMethod == "PUT" || httpMethod == "POST") {
			w.Row(100).Dynamic(1)

			requestBodyEditorField.Edit(w)
		}

		w.Row(30).Dynamic(4)

		w.LabelWrap("Wut")
		w.LabelWrap("Wut2")
		w.LabelColored(status, "RT", colornames.Aquamarine)
		w.LabelColored(responseContentType, "RT", colornames.Aquamarine)

		//w.Row(500).Dynamic(1)
		w.RowScaled(500).Dynamic(1)

		responseBodyEditor.Maxlen = 100000
		responseBodyEditor.Edit(w)

	}
}

func cycleSelectedInputFieldForward() {

	var foldOver = len(tabbableFields) - 1

	// Find the active and set the next element to active
	for e := range tabbableFields {
		element := tabbableFields[e]

		if (element.Active == true) {

			if (e+1 > foldOver) {
				tabbableFields[0].Active = true
				element.Active = false
			} else {
				tabbableFields[e+1].Active = true
				element.Active = false
			}

			break
		}
	}

}

func cycleSelectedInputFieldBackward() {

	// Find the active and set the previous element to active
	var foldOver = len(tabbableFields) - 1

	// Find the active and set the next element to active
	for e := range tabbableFields {
		element := tabbableFields[e]

		if (element.Active == true) {

			if (e-1 < 0) {
				tabbableFields[foldOver].Active = true
				element.Active = false
			} else {
				tabbableFields[e-1].Active = true
				element.Active = false
			}

			break
		}
	}

}

func handleKeybindings(w *nucular.Window, urlField *nucular.TextEditor,
	responseField *nucular.TextEditor,
	usernameField *nucular.TextEditor,
	passwordField *nucular.TextEditor,
	requestBodyField *nucular.TextEditor) {

	mw := w.Master()

	if in := w.Input(); in != nil {
		k := in.Keyboard
		for _, e := range k.Keys {
			scaling := mw.Style().Scaling

			switch {
			case (e.Modifiers == key.ModControl || e.Modifiers == key.ModControl|key.ModShift) && (e.Code == key.CodeEqualSign):
				mw.Style().Scale(scaling + 0.1)
			case (e.Modifiers == key.ModControl || e.Modifiers == key.ModControl|key.ModShift) && (e.Code == key.CodeHyphenMinus):
				mw.Style().Scale(scaling - 0.1)
			case (e.Modifiers == key.ModControl) && (e.Code == key.CodeF):
				mw.SetPerf(!mw.GetPerf())

				// Change request method binds
			case (e.Modifiers == key.ModControl && (e.Code == key.CodeQ)):
				httpMethod = "GET"
			case (e.Modifiers == key.ModControl && (e.Code == key.CodeW)):
				httpMethod = "PUT"
			case (e.Modifiers == key.ModControl && (e.Code == key.CodeE)):
				httpMethod = "POST"
			case (e.Modifiers == key.ModControl && (e.Code == key.CodeR)):
				httpMethod = "DELETE"

				// Send HTTP request binds
			case (e.Modifiers == key.ModControl && (e.Code == key.CodeReturnEnter)):

				var request HttpRequest

				// Don't send a body in GET or DELETE requests
				if (httpMethod == "GET" || httpMethod == "DELETE") {

					request = HttpRequest{
						httpMethod,
						string(urlField.Buffer),
						string(usernameField.Buffer),
						string(passwordField.Buffer),
						"",
					}

				} else {
					request = HttpRequest{
						httpMethod,
						string(urlField.Buffer),
						string(usernameField.Buffer),
						string(passwordField.Buffer),
						string(requestBodyField.Buffer),
					}
				}

				var response = CallHttp(&request)

				// Check if we find any rune's with code 13 and clean the body if so :lenn:
				var clean = removeCRs(response.body)

				println(response.contentLength)
				println([]rune(clean))

				formattedBody := formatBody(clean, response.contentType)
				// Format some html
				//formatted := gohtml.Format(clean)

				responseField.Buffer = []rune(formattedBody)

				status = response.statusString
				responseContentType = response.contentType

			case (e.Modifiers == key.ModShift && (e.Code == key.CodeTab)):
				cycleSelectedInputFieldBackward()

			case (e.Code == key.CodeTab):
				cycleSelectedInputFieldForward()


			}

		}
	}
}

func formatBody(body string, contentType string) string {
	isHtml := strings.Index(contentType, "html") > -1
	isXml := strings.Index(contentType, "xml") > -1
	isJson := strings.Index(contentType, "json") > -1

	if (isHtml) {
		return gohtml.Format(body)
	}

	if (isXml) {
		return xmlfmt.FormatXML(body, "", "  ")
	}

	if (isJson) {
		var out io.Writer
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		if err := enc.Encode(body); err != nil {
			return body
		}
		return body
	}

	return body

}

func removeCRs(s string) string {
	/* Removes carriage return char's from a string */
	var asRune []rune = nil

	for _, char := range s {
		if char == 13 {
			continue
		} else {
			asRune = append(asRune, rune(char))
		}
	}

	return string(asRune)

}
