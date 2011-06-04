package gosigner 

import (
	"http"
	"crypto/hmac"
	"json"
	"encoding/base64"
	"template"
)

var (
        indexTemplate = template.MustParseFile("index.html", nil)
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	indexTemplate.Execute(w, 1)

}

func init() {
        
	//Add singHandler as handler of /sign
	http.HandleFunc("/sign", signHandler)
        http.HandleFunc("/", indexHandler)
}
