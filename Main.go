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
	"image"
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

	Wnd = nucular.NewMasterWindowSize(0, Title, image.Point{1000,1000}, textEditorDemo())

	if Wnd == nil {
		_, _ = fmt.Fprintf(os.Stderr, "unknown demo %q\n", "WUT")
		_, _ = fmt.Fprintf(os.Stderr, "ggggg\n")
		os.Exit(1)
	}

	Wnd.Main()
}

var status = "django bango"
var responseContentType string = "ttt"
var tabbableFields []*nucular.TextEditor

type Pair struct {
	a, b string
}

// Load these default values from config file
var auth = Pair{"admin", "admin"}

var client = &http.Client{}
var httpMethod = "GET"
var HistoryStruct, _ = configHandler.GetHistoryOrCreateHistoryFile()

var history = HistoryStruct.Content

var splith = &nucular.ScalableSplit{}

type HttpHeaderDisplay struct {
	enabled bool
	name string
	value string
}

var displayHeadersList = []*HttpHeaderDisplay{}

var wutface = HttpHeaderDisplay{
	true,
	"Content-Type",
	"application/json",
}

var wutface2 = HttpHeaderDisplay{
	false,
	"Accept",
	"application/smason",
}


// Construct the render function
func textEditorDemo() func(w *nucular.Window) {

	var urlEditorField nucular.TextEditor
	urlEditorField.Flags = nucular.EditSelectable | nucular.EditClipboard

	// Set the buffer to the last value in the history or else empty
	urlEditorField.Buffer = []rune(history[len(history) -1])
	urlEditorField.Maxlen = 1000
	urlEditorField.Active = true

	tabbableFields = append(tabbableFields, &urlEditorField)

	var requestBodyEditorField nucular.TextEditor
	// What determines if the field gets auto tabbed? Is it EditMultiline?
	requestBodyEditorField.Flags = nucular.EditSelectable | nucular.EditMultiline | nucular.EditClipboard
	requestBodyEditorField.Buffer = []rune("request body")
	requestBodyEditorField.Maxlen = 150

	tabbableFields = append(tabbableFields, &requestBodyEditorField)

	var usernameEditorField nucular.TextEditor
	usernameEditorField.Flags = nucular.EditSelectable | nucular.EditClipboard
	usernameEditorField.Buffer = []rune(auth.a)
	usernameEditorField.Maxlen = 150

	tabbableFields = append(tabbableFields, &usernameEditorField)

	var passwordEditorField nucular.TextEditor
	passwordEditorField.Flags = nucular.EditSelectable | nucular.EditClipboard
	passwordEditorField.Buffer = []rune(auth.b)
	passwordEditorField.Maxlen = 150

	tabbableFields = append(tabbableFields, &passwordEditorField)

	var responseBodyEditor nucular.TextEditor
	responseBodyEditor.Flags = nucular.EditSelectable | nucular.EditMultiline | nucular.EditClipboard | nucular.EditReadOnly
	responseBodyEditor.Buffer = []rune("response body")

	tabbableFields = append(tabbableFields, &responseBodyEditor)

	splith.MinSize = 80
	splith.Size = 120
	splith.Spacing = 5

	history, err := configHandler.GetHistoryOrCreateHistoryFile()

	if err != nil {
		println("could not load history")
	}

	fmt.Println("main history", history.Content)

	config, err := configHandler.LoadConfig()

	if err != nil {
		println("could not load config")
	}

	fmt.Println("config", config)

	displayHeadersList = append(displayHeadersList, &wutface)
	displayHeadersList = append(displayHeadersList, &wutface2)

	var displayHeadersInputsList []*nucular.TextEditor

	for _, el := range displayHeadersList {
		var textEditor = nucular.TextEditor{}
		textEditor.Buffer = []rune(el.value)
		textEditor.Flags = nucular.EditSelectable | nucular.EditClipboard
		displayHeadersInputsList = append(displayHeadersInputsList, &textEditor)
		tabbableFields = append(tabbableFields, &textEditor)

	}

	// This is the render function, everything above it is called once upon initialization
	return func(w *nucular.Window) {

		handleKeybindings(w, &urlEditorField, &responseBodyEditor,
			&usernameEditorField, &passwordEditorField, &requestBodyEditorField, &displayHeadersList)



		w.Row(30).Static(50, 500, 125, 125)

		w.LabelColored(httpMethod, "LT", colornames.Aquamarine)

		urlEditorField.Edit(w)
		usernameEditorField.Edit(w)
		passwordEditorField.Edit(w)

		if (httpMethod == "PUT" || httpMethod == "POST") {
			w.Row(100).Dynamic(1)

			requestBodyEditorField.Edit(w)
		}

		renderRequestHeaders(w, &displayHeadersList, &displayHeadersInputsList)

		w.Row(30).Dynamic(4)

		w.LabelWrap("Wut")
		w.LabelWrap("Wut2")
		w.LabelColored(status, "RT", colornames.Aquamarine)
		w.LabelColored(responseContentType, "RT", colornames.Aquamarine)

		//w.Row(500).Dynamic(1)
		w.RowScaled(500).Dynamic(1)

		responseBodyEditor.Maxlen = 100000
		responseBodyEditor.Edit(w)

		area := w.Row(0).SpaceBegin(0)

		//viewbounds, commitbounds := splith.Horizontal(w, rect.Rect{0, 0, 150, 150})
		viewbounds, commitbounds := splith.Horizontal(w, area)

		w.LayoutSpacePushScaled(viewbounds)
		w.LabelColored("my face is yours", "LT", colornames.Blueviolet)
		w.LayoutSpacePushScaled(commitbounds)
		w.LabelColored("my face is yours", "LT", colornames.Blueviolet)


	}
}

func renderRequestHeaders(w *nucular.Window,
	displayHeadersList *[]*HttpHeaderDisplay,
	displayHeadersInputList *[]*nucular.TextEditor) {
	w.Row(50).Dynamic(2)

	displayHeaderStructs := *displayHeadersList
	displayHeaderInputs := *displayHeadersInputList

	for index, headerDisplay := range displayHeaderStructs {
		w.CheckboxText(headerDisplay.name, &headerDisplay.enabled)
		displayHeaderInputs[index].Edit(w)
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
	requestBodyField *nucular.TextEditor,
	displayHeadersList *[]*HttpHeaderDisplay) {

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

				requestHeaders := createHeadersFromSelected(displayHeadersList)

				// Don't send a body in GET or DELETE requests
				if (httpMethod == "GET" || httpMethod == "DELETE") {

					request = HttpRequest{
						httpMethod,
						string(urlField.Buffer),
						string(usernameField.Buffer),
						string(passwordField.Buffer),
						"",
						requestHeaders,
					}

				} else {
					request = HttpRequest{
						httpMethod,
						string(urlField.Buffer),
						string(usernameField.Buffer),
						string(passwordField.Buffer),
						string(requestBodyField.Buffer),
						requestHeaders,
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

func createHeadersFromSelected(display *[]*HttpHeaderDisplay) map[string]string {
	hest := *display

	var resultingRequestHeaders = make(map[string]string)

	for _, el := range hest {

		// For some reason the active/enabled value of nucular.CheckboxText is inverted
		// Flip it
		if !el.enabled {
			resultingRequestHeaders[el.name] = el.value
		}
	}

	return resultingRequestHeaders
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
