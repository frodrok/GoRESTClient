package main

import (
	"GoRESTClient/configHandler"
	"GoRESTClient/httpClient"
	"bytes"
	"encoding/json"
	"fmt"
	"image"

	//"io"
	"os"
	"strings"
	"time"

	//	"strconv"

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

var status = "Status"
var responseContentType string = "Content Type"
var tabbableFields []*nucular.TextEditor

type Pair struct {
	a, b interface{}
}

// Default values unless history is found
var auth = Pair{"admin", "admin"}

var httpMethod = "GET"

var cHandler configHandler.ConfigHandlerActions

var splith = &nucular.ScalableSplit{}

// Construct the render function
func textEditorDemo() func(w *nucular.Window) {

	var urlEditorField nucular.TextEditor
	urlEditorField.Flags = nucular.EditSelectable | nucular.EditClipboard

	urlEditorField.Maxlen = 1000
	urlEditorField.Active = true

	// Instantiate configHandler
	cHandler, _ = configHandler.NewConfigHandler("/home/fredrik/.gostman")

	// Load history
	var history = cHandler.GetRequestHistory()

	// Load config
	config := cHandler.GetConfig()

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

		urlEditorField.Buffer = []rune(lastHistoryEntry.Url)

		usernameEditorField.Buffer = []rune(lastHistoryEntry.Username)

		passwordEditorField.Buffer = []rune(lastHistoryEntry.Password)

		setupHeaderAreaFromHistory(&lastHistoryEntry.Headers)

	}

	tabbableFields = append(tabbableFields, &passwordEditorField)

	var responseBodyEditor nucular.TextEditor
	responseBodyEditor.Flags = nucular.EditSelectable | nucular.EditMultiline | nucular.EditClipboard | nucular.EditReadOnly
	responseBodyEditor.Buffer = []rune("response body")

	tabbableFields = append(tabbableFields, &responseBodyEditor)

	splith.MinSize = 80
	splith.Size = 120
	splith.Spacing = 5

	fmt.Println("config", config)

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
			&requestBodyEditorField, &displayHeadersList, displayHeadersInputs,
			cHandler)

		w.Row(30).Static(50, 500, 125, 125)

		w.LabelColored(httpMethod, "LT", colornames.Aquamarine)

		urlEditorField.Edit(w)
		usernameEditorField.Edit(w)
		passwordEditorField.Edit(w)

		// Only display the HTTP request body editor if the selected
		// request is PUT or POST
		if httpMethod == "PUT" || httpMethod == "POST" {
			w.Row(100).Dynamic(1)

			requestBodyEditorField.Edit(w)
		}

		renderRequestHeaders(w)

		w.Row(30).Dynamic(4)

		w.LabelWrap("Wut")
		w.LabelWrap("Wut2")
		w.LabelColored(status, "RT", colornames.Aquamarine)
		w.LabelColored(responseContentType, "RT", colornames.Aquamarine)

		// Calculate 60% of the window height and set the height to that
		var fifty_percent float64 = float64(w.Bounds.H) * float64(0.6)
		w.Row(int(fifty_percent)).Dynamic(1)

		responseBodyEditor.Maxlen = 100000
		responseBodyEditor.Edit(w)

	}
}

func cycleSelectedInputFieldForward() {

	var foldOver = len(tabbableFields) - 1

	// Find the active and set the next element to active
	for e := range tabbableFields {
		element := tabbableFields[e]

		if element.Active {

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

		if element.Active {

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

		if element.Active {

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
	displayHeadersInputs []*nucular.TextEditor,
	cHandler configHandler.ConfigHandlerActions) {

	mw := w.Master()

	if in := w.Input(); in != nil {
		k := in.Keyboard
		for _, e := range k.Keys {
			scaling := mw.Style().Scaling

			fmt.Sprint(e.Code.String())

			LOG(e.Code.String())
			LOG(e.String())

			fmt.Println(string(e.Rune) == "+")

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

				var request httpClient.HttpRequest

				requestHeaders := getEnabledHeaders()

				// Don't send a body in GET or DELETE requests
				if httpMethod == "GET" || httpMethod == "DELETE" {

					request = httpClient.HttpRequest{
						httpMethod,
						string(urlField.Buffer),
						string(usernameField.Buffer),
						string(passwordField.Buffer),
						"",
						requestHeaders,
					}

				} else {

					request = httpClient.HttpRequest{
						httpMethod,
						string(urlField.Buffer),
						string(usernameField.Buffer),
						string(passwordField.Buffer),
						string(requestBodyField.Buffer),
						requestHeaders,
					}
				}

				responses := make(chan httpClient.HttpResponse)

				go func() {
					var response = httpClient.CallHttp(&request)
					responses <- response
				}()

				response := <-responses

				// Check if we find any rune's with code 13 and clean the body if so
				var clean = removeCRs(response.Body)

				formattedBody := formatBody(clean, response.ContentType)

				responseField.Buffer = []rune(formattedBody)

				status = response.StatusString
				responseContentType = response.ContentType

				// Store url, headers and credentials in history
				cHandler.AddRequestToHistory(request)

			case (e.Modifiers == key.ModShift && (e.Code == key.CodeTab)):
				cycleSelectedInputFieldBackward()

			case (e.Code == key.CodeTab):
				cycleSelectedInputFieldForward()

			case (e.Modifiers == key.ModControl && (e.Code == key.CodeW)):
				selectWordInCurrentEditField()

			// Special case here, golang's event has no code for
			// the symbol '+', so string compare it on the rune
			case (e.Modifiers == key.ModControl && (string(e.Rune) == "+")):
				startAddHeaderProcess()
			}

		}
	}
}

func formatBody(body string, contentType string) string {

	isHtml := strings.Contains(contentType, "html")
	isXml := strings.Contains(contentType, "xml")
	isJson := strings.Contains(contentType, "json")

	if isHtml {
		return gohtml.Format(body)
	}

	if isXml {
		return xmlfmt.FormatXML(body, "", "  ")
	}

	if isJson {

		var prettyJson bytes.Buffer
		err := json.Indent(&prettyJson, []byte(body), "", "  ")

		if err != nil {
			println("json parse error: ", err)
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
