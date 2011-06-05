package gosigner 

import (
	"http"
	"crypto/hmac"
	"json"
	"encoding/base64"
	"template"
        "os"
        "strings"
        "log"
)

var (
        indexTemplate = template.MustParseFile("index.html", nil)
        indexJST   *template.Template
)
type Signature struct {
	Signature, Content, Key string
}

func signHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
        var data Signature
        dec.Decode(&data)
        keyBytes := []byte(data.Key)
	content := data.Content
	mac := hmac.NewSHA1(keyBytes)
	mac.Write([]byte(content))

	out := make([]byte, base64.StdEncoding.EncodedLen(len(mac.Sum())))
	base64.StdEncoding.Encode(out, mac.Sum())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&Signature{string(out), content, data.Key})
}
type JSTTemplate struct{
        Template string
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	sw := new(StringWritter)
        indexJST := template.New(nil)
        //indexJST.SetDelims("{","}")
        if err := indexJST.ParseFile("index.jst"); err != nil {
                panic("Unable to parse template")
        }
        indexJST.Execute(sw, nil)
        log.Print(sw.s)
        indexTemplate.Execute(w, &JSTTemplate{strings.Replace(sw.s, "\n", "" , -1)})
}

func init() {
	//Add singHandler as handler of /sign
	http.HandleFunc("/sign", signHandler)
        http.HandleFunc("/", indexHandler)
}


type StringWritter struct {
	s string
}

//Writes the template as string
func (self *StringWritter) Write(p []byte) (n int, err os.Error) {
	self.s += string(p)
	return len(self.s), nil
}

