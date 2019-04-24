package main

import (
	"encoding/base64"
	"fmt"
	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"golang.org/x/image/colornames"
	"golang.org/x/mobile/event/key"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var Wnd nucular.MasterWindow

var Title = "Gostman"

func main() {
	fmt.Println("hello sir\n");
	Wnd = nucular.NewMasterWindow(0, Title, textEditorDemo())

	if Wnd == nil {
		_, _ = fmt.Fprintf(os.Stderr, "unknown demo %q\n", "WUT")
		_, _ = fmt.Fprintf(os.Stderr, "ggggg\n")
		os.Exit(1)
	}

	Wnd.Main()
}

var status = 0
var responseContentType string
var tabbableFields [3]*nucular.TextEditor

type Pair struct {
	a, b string
}

var auth = Pair{"admin", "admin"}

var client = &http.Client{}
var httpMethod = "GET"

func textEditorDemo() func(w *nucular.Window) {

	var urlEditorField nucular.TextEditor
	urlEditorField.Flags = nucular.EditSelectable
	urlEditorField.Buffer = []rune("https://google.se/API/version")
	urlEditorField.Maxlen = 1000
	urlEditorField.Active = true
	tabbableFields[0] = &urlEditorField

	var requestBodyEditorField nucular.TextEditor
	requestBodyEditorField.Flags = nucular.EditSelectable | nucular.EditMultiline | nucular.EditClipboard | nucular.EditIbeamCursor
	requestBodyEditorField.Buffer = []rune("request body")
	requestBodyEditorField.Maxlen = 150
	tabbableFields[1] = &requestBodyEditorField


	var responseBodyEditor nucular.TextEditor
	responseBodyEditor.Flags = nucular.EditSelectable | nucular.EditMultiline | nucular.EditClipboard | nucular.EditIbeamCursor

	responseBodyEditor.Buffer = []rune("response body")

	tabbableFields[2] = &responseBodyEditor

	return func(w *nucular.Window) {

		keybindings(w, &urlEditorField, &responseBodyEditor)

		// w.Row(30).Dynamic(1)
		w.Row(30).Dynamic(3)
		// w.Row(30).Static(180)
		w.LabelColored(httpMethod, "LT", colornames.Aquamarine)

		urlEditorField.Edit(w)

		var pressed = w.Button(label.T("send"), false)

		if (httpMethod == "PUT" || httpMethod == "POST") {
			w.Row(100).Dynamic(1)

			requestBodyEditorField.Edit(w)
		}

		w.Row(30).Dynamic(1)
		w.LabelColored(responseContentType, "LT", colornames.Aquamarine)
		w.LabelColored(strconv.Itoa(status), "RT", colornames.Aquamarine)

		w.Row(500).Dynamic(1)
		responseBodyEditor.Maxlen = 1000
		responseBodyEditor.Edit(w)

		if pressed {
			fmt.Printf("my press %s\n", string(urlEditorField.Buffer))

			var response = callHttp(string(urlEditorField.Buffer), httpMethod)
			status = response.status

			responseBodyEditor.Buffer = []rune(response.body)
			responseContentType = response.contentType
		}
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

func keybindings(w *nucular.Window, urlField *nucular.TextEditor, responseField *nucular.TextEditor) {
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

			case (e.Modifiers == key.ModControl && (e.Code == key.CodeQ)):
				httpMethod = "GET"
			case (e.Modifiers == key.ModControl && (e.Code == key.CodeW)):
				httpMethod = "PUT"
			case (e.Modifiers == key.ModControl && (e.Code == key.CodeE)):
				httpMethod = "POST"
			case (e.Modifiers == key.ModControl && (e.Code == key.CodeR)):
				httpMethod = "DELETE"

			case (e.Modifiers == key.ModControl && (e.Code == key.CodeReturnEnter)):

				var response = callHttp(string(urlField.Buffer), httpMethod)

				var clean = removeCRs(response.body)

				responseField.Buffer = []rune(clean)

				status = response.status
				responseContentType = response.contentType

			case (e.Modifiers == key.ModShift && (e.Code == key.CodeTab)):
				cycleSelectedInputFieldBackward()

			case (e.Code == key.CodeTab):
				cycleSelectedInputFieldForward()



			}

		}
	}
}

func removeCRs(s string) string {
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

func basicAuth(username string, password string) string {
	var auth = username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

type HttpResponse struct {
	status int
	body string
	contentType string
	contentLength string
}

func callHttp(s string, method string) HttpResponse {

	fmt.Println("callHttp called")
	req, err := http.NewRequest(method, s, nil)
	req.Header.Add("Authorization", "Basic "+basicAuth(auth.a, auth.b))

	resp, err := client.Do(req)

	if err != nil {
		return HttpResponse{
			500,
			"",
			"error",
			"0",
		}
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	for k, v := range resp.Header {
		fmt.Printf("key[%s] value[%s]\n", k, v)
	}

	return HttpResponse{
		status: resp.StatusCode,
		body: string(body),
		contentType: resp.Header.Get("Content-Type"),
		contentLength: resp.Header.Get("Content-Length"),
	}

}
