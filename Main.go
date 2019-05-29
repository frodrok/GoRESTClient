package main

import (
	"GoRESTClient/configHandler"
	"bytes"
	"encoding/json"
	"fmt"
	"image"

	//"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aarzilli/nucular"
	"github.com/go-xmlfmt/xmlfmt"
	"github.com/yosssi/gohtml"
	"golang.org/x/image/colornames"
	"golang.org/x/mobile/event/key"
)

func LOG(s string) {
	fmt.Printf("%s - %s\n", time.Now().Format("2009-09-01T15:04"), s)
}

var Wnd nucular.MasterWindow

var Title = "Gostman"

func main() {

	Wnd = nucular.NewMasterWindowSize(0, Title, image.Point{1000, 1000}, textEditorDemo())

	if Wnd == nil {
		_, _ = fmt.Fprintf(os.Stderr, "unknown demo %q\n", "WUT")
		_, _ = fmt.Fprintf(os.Stderr, "ggggg\n")
		os.Exit(1)
	}

	Wnd.Main()
}

var status = "http status"
var responseContentType string = "content type"
var tabbableFields []*nucular.TextEditor

type Pair struct {
	a, b interface{}
}

// Load these default values from config file
var auth = Pair{"admin", "admin"}

var client = &http.Client{}
var httpMethod = "GET"
var HistoryStruct, _ = configHandler.GetHistoryOrCreateHistoryFile()

var history []configHandler.HistoryEntry = HistoryStruct.Content

var splith = &nucular.ScalableSplit{}

type HttpHeaderDisplay struct {
	enabled bool
	name    string
	value   string
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
	"application/json",
}

// Construct the render function
func textEditorDemo() func(w *nucular.Window) {

	var urlEditorField nucular.TextEditor
	urlEditorField.Flags = nucular.EditSelectable | nucular.EditClipboard

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
	usernameEditorField.Buffer = []rune(auth.a.(string))
	usernameEditorField.Maxlen = 150

	tabbableFields = append(tabbableFields, &usernameEditorField)

	var passwordEditorField nucular.TextEditor
	passwordEditorField.Flags = nucular.EditSelectable | nucular.EditClipboard
	passwordEditorField.Buffer = []rune(auth.b.(string))
	passwordEditorField.Maxlen = 150

	// Set the buffer to the last value in the history or else empty
	if len(history) > 0 {

		var lastHistoryEntry = history[len(history)-1]

		urlFieldValue := lastHistoryEntry.Url
		urlEditorField.Buffer = []rune(urlFieldValue)

		usernameEditorField.Buffer = []rune(lastHistoryEntry.Username)

		passwordEditorField.Buffer = []rune(lastHistoryEntry.Password)

	}

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

	var displayHeadersInputs []*nucular.TextEditor

	// For each field in displayHeadersList create a nucular.TextEditor
	// and add it to displayHeadersInputsList
	for _, el := range displayHeadersList {
		var textEditor = nucular.TextEditor{}
		textEditor.Buffer = []rune(el.value)
		textEditor.Flags = nucular.EditSelectable | nucular.EditClipboard
		displayHeadersInputs = append(displayHeadersInputs, &textEditor)
		tabbableFields = append(tabbableFields, &textEditor)

	}

	// This is the render function, everything above it is called once upon initialization
	return func(w *nucular.Window) {

		handleKeybindings(w, &urlEditorField, &responseBodyEditor,
			&usernameEditorField, &passwordEditorField,
			&requestBodyEditorField, &displayHeadersList, displayHeadersInputs)

		w.Row(30).Static(50, 500, 125, 125)

		w.LabelColored(httpMethod, "LT", colornames.Aquamarine)

		urlEditorField.Edit(w)
		usernameEditorField.Edit(w)
		passwordEditorField.Edit(w)

		if httpMethod == "PUT" || httpMethod == "POST" {
			w.Row(100).Dynamic(1)

			requestBodyEditorField.Edit(w)
		}

		renderRequestHeaders(w, &displayHeadersList, &displayHeadersInputs)

		w.Row(30).Dynamic(4)

		w.LabelWrap("Wut")
		w.LabelWrap("Wut2")
		w.LabelColored(status, "RT", colornames.Aquamarine)
		w.LabelColored(responseContentType, "RT", colornames.Aquamarine)

		//w.Row(500).Dynamic(1)
		// Calculate 60% of the window height and set to that
		var fifty_percent float64 = float64(w.Bounds.H) * float64(0.6)
		w.Row(int(fifty_percent)).Dynamic(1)

		responseBodyEditor.Maxlen = 100000
		responseBodyEditor.Edit(w)

		//area := w.Row(0).SpaceBegin(0)

		//viewbounds, commitbounds := splith.Horizontal(w, rect.Rect{0, 0, 150, 150})
		//viewbounds, commitbounds := splith.Horizontal(w, area)

		//w.LayoutSpacePushRatio(0, 500, 1, 0.6)

		//responseBodyEditor.Maxlen = 100000
		//responseBodyEditor.Edit(w)

		//w.LayoutSpacePushScaled(viewbounds)
		//w.LabelColored("my face is yours", "LT", colornames.Blueviolet)
		//w.LayoutSpacePushScaled(commitbounds)
		//w.LabelColored("my face is yours", "LT", colornames.Blueviolet)

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

		if element.Active == true {

			if e+1 > foldOver {
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

		if element.Active == true {

			if e-1 < 0 {
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

func getSelectPoints(buffer []rune, startIndex int) *Pair {

	var start int
	var end int
	asRunes := buffer
	var MAX_ASCII_VAL rune = 123
	var MIN_ASCII_VAL rune = 64
	var NULL_ASCII_VAL rune = 0

	// Go forwards to get end
	if startIndex < len(asRunes) {
		for i := startIndex; i <= len(asRunes)-1; i++ {

			var currentRune = asRunes[i]
			if currentRune > MIN_ASCII_VAL && currentRune < MAX_ASCII_VAL && currentRune != NULL_ASCII_VAL {
				continue
			} else {
				end = i
				break
			}
		}
	} else {
		end = startIndex
	}

	// Go backwards to get start
	for i := startIndex - 1; i >= 0; i-- {
		var currentRune = asRunes[i]
		if currentRune > MIN_ASCII_VAL && currentRune < MAX_ASCII_VAL && currentRune != NULL_ASCII_VAL {
			continue
		} else {
			start = i + 1
			break
		}
	}

	return &Pair{start, end}
}

func selectWordInCurrentEditField() {
	for e := range tabbableFields {
		element := tabbableFields[e]

		if element.Active == true {

			// Find the previous non-alphabet character and the next one and
			// select everything in between that.
			// element.Cursor will give the current index of the cursor
			cursorIndex := element.Cursor
			currentBuffer := element.Buffer

			var points *Pair = getSelectPoints(currentBuffer, cursorIndex)

			element.SelectAll()
			element.SelectStart = points.a.(int)
			element.SelectEnd = points.b.(int)
		}
	}
}

func handleKeybindings(w *nucular.Window, urlField *nucular.TextEditor,
	responseField *nucular.TextEditor,
	usernameField *nucular.TextEditor,
	passwordField *nucular.TextEditor,
	requestBodyField *nucular.TextEditor,
	displayHeadersList *[]*HttpHeaderDisplay,
	displayHeadersInputs []*nucular.TextEditor) {

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
			case (e.Modifiers == key.ModControl && (e.Code == key.Code1)):
				httpMethod = "GET"
			case (e.Modifiers == key.ModControl && (e.Code == key.Code2)):
				httpMethod = "PUT"
			case (e.Modifiers == key.ModControl && (e.Code == key.Code3)):
				httpMethod = "POST"
			case (e.Modifiers == key.ModControl && (e.Code == key.Code4)):
				httpMethod = "DELETE"

				// Send HTTP request binds
			case (e.Modifiers == key.ModControl && (e.Code == key.CodeReturnEnter)):

				var request HttpRequest

				requestHeaders := createHeadersFromSelected(displayHeadersList,
					displayHeadersInputs)

				// Don't send a body in GET or DELETE requests
				if httpMethod == "GET" || httpMethod == "DELETE" {

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

				formattedBody := formatBody(clean, response.contentType)
				// Format some html
				//formatted := gohtml.Format(clean)

				responseField.Buffer = []rune(formattedBody)

				status = response.statusString
				responseContentType = response.contentType

				// Store url and credentials in history
				configHandler.SaveToHistory(string(urlField.Buffer),
					string(usernameField.Buffer),
					string(passwordField.Buffer))

			case (e.Modifiers == key.ModShift && (e.Code == key.CodeTab)):
				cycleSelectedInputFieldBackward()

			case (e.Code == key.CodeTab):
				cycleSelectedInputFieldForward()

			case (e.Modifiers == key.ModControl && (e.Code == key.CodeW)):
				selectWordInCurrentEditField()

			}

		}
	}
}

func createHeadersFromSelected(display *[]*HttpHeaderDisplay,
	inputFields []*nucular.TextEditor) map[string]string {

	hest := *display

	var resultingRequestHeaders = make(map[string]string)

	for index, el := range hest {

		val := string(inputFields[index].Buffer)

		// For some reason the active/enabled value of nucular.CheckboxText is inverted
		// Flip it
		if !el.enabled {
			resultingRequestHeaders[el.name] = val
		}
	}

	return resultingRequestHeaders
}

func formatBody(body string, contentType string) string {

	isHtml := strings.Index(contentType, "html") > -1
	isXml := strings.Index(contentType, "xml") > -1
	isJson := strings.Index(contentType, "json") > -1

	if isHtml {
		return gohtml.Format(body)
	}

	if isXml {
		return xmlfmt.FormatXML(body, "", "  ")
	}

	if isJson {

		var prettyJson bytes.Buffer
		error := json.Indent(&prettyJson, []byte(body), "", "  ")

		if error != nil {
			println("json parse error: ", error)
			return ""
		}

		return prettyJson.String()
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
