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
	// "strings"
	"time"
)

var Wnd nucular.MasterWindow
var Wut = "hllo"

var Title = "Gostman"

func main() {
	fmt.Println("hello sir\n");
	Wnd = nucular.NewMasterWindow(0, Title, textEditorDemo())

	if Wnd == nil {
		_, _ = fmt.Fprintf(os.Stderr, "unknown demo %q\n", Wut)
		_, _ = fmt.Fprintf(os.Stderr, "ggggg\n")
		os.Exit(1)
	}

	Wnd.Main()
}

var compression int

type difficulty int

const (
	easy = difficulty(iota)
	hard
)

var op difficulty = easy
var status = 0

type Pair struct {
	a, b string
}

var auth = Pair{"admin", "admin"}

var client = &http.Client{}
var httpMethod = "GET"

func textEditorDemo() func(w *nucular.Window) {

	var textEditorEditor nucular.TextEditor
	textEditorEditor.Flags = nucular.EditSelectable
	textEditorEditor.Buffer = []rune("https://google.se/API/version")
	textEditorEditor.Maxlen = 150

	var responseBodyEditor nucular.TextEditor
	responseBodyEditor.Flags = nucular.EditSelectable | nucular.EditMultiline | nucular.EditClipboard | nucular.EditIbeamCursor
	responseBodyEditor.Buffer = []rune("responsebody")

	return func(w *nucular.Window) {

		keybindings(w, &textEditorEditor, &responseBodyEditor)

		// w.Row(30).Dynamic(1)
		w.Row(30).Dynamic(3)
		// w.Row(30).Static(180)
		w.LabelColored(httpMethod, "LT", colornames.Aquamarine)
		textEditorEditor.Edit(w)

		var pressed = w.Button(label.T("send"), false)

		w.Row(30).Dynamic(1)
		w.Label(strconv.Itoa(status), "RT")

		w.Row(500).Dynamic(1)
		responseBodyEditor.Maxlen = 1000
		responseBodyEditor.Edit(w)

		if pressed {
			fmt.Printf("my press %s\n", string(textEditorEditor.Buffer))

			var body, httpStatus = callHttp(string(textEditorEditor.Buffer), httpMethod)
			status = httpStatus

			responseBodyEditor.Buffer = []rune(body)
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

				var body, httpStatus = callHttp(string(urlField.Buffer), httpMethod)

				var clean = removeCRs(body)

				responseField.Buffer = []rune(clean)

				status = httpStatus

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

func callHttp(s string, method string) (string, int) {

	fmt.Println("called")
	req, err := http.NewRequest(method, s, nil)
	req.Header.Add("Authorization", "Basic "+basicAuth(auth.a, auth.b))

	resp, err := client.Do(req)

	if err != nil {
		return "error", 418
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return string(body), resp.StatusCode

}

func basicDemo(w *nucular.Window) {

	w.Row(30).Dynamic(1)
	w.Label(time.Now().Format("15:04:05"), "RT")

	w.Label("My gawd", "LT")

	var textEditorEditor nucular.TextEditor

	textEditorEditor.Flags = nucular.EditSelectable

	textEditorEditor.Buffer = []rune("prova")

	w.Row(30).Dynamic(1)
	textEditorEditor.Maxlen = 30
	textEditorEditor.Edit(w)

	w.Row(30).Static(80)
	if w.Button(label.T("button"), false) {
		fmt.Printf("button pressed! difficulty: %v compression: %d\n", op, compression)
	}
	w.Row(30).Dynamic(2)
	if w.OptionText("easy", op == easy) {
		op = easy
	}
	if w.OptionText("hard", op == hard) {
		op = hard
	}
	w.Row(25).Dynamic(1)
	w.PropertyInt("Compression:", 0, &compression, 100, 10, 1)
}
