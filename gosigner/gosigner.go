package gosigner 

import (
	"http"
	"crypto/hmac"
	"json"
	"encoding/base64"
)

const (
        version = "0.0.1"
)

type Signature struct {
	Signature, Content, Key string
}

type Version struct{
        Version string
}
// Signs the Content
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


//Used to get the current version of the app
func versionHandler(w http.ResponseWriter, r *http.Request){
        json.NewEncoder(w).Encode(&Version{version})
}        
func init() {
	http.HandleFunc("/ver", versionHandler)
        http.HandleFunc("/sign", signHandler)
}