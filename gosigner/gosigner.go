package gosigner

import (
	"http"
	"crypto/hmac"
	"json"
	"encoding/base64"
	"appengine"
	"appengine/user"
	"log"
)

const (
	version = "0.0.2"
)

type Signature struct {
	Signature, Content, Key string
}

type Version struct {
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
func versionHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(&Version{version})
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.String(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}
}

func auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)
		u := user.Current(c)
		if u == nil {
			url, err := user.LoginURL(c, r.URL.String())
			if err != nil {
				http.Error(w, err.String(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Location", url)
			w.WriteHeader(http.StatusFound)
			return
		}
		log.Printf("User %s logged.", u.Email)

		h(w, r)
	}
}

func keyHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var data Key
	json.NewDecoder(r.Body).Decode(&data)
	data.Save(c)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}


func keysHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	var owner string
	if u == nil {
		owner = "test@example.com"
	} else {
		owner = u.Email
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AllKeys(c, owner))
}

func POST(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		h(w, r)
	}
}

func GET(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		h(w, r)
	}
}

func init() {
	http.HandleFunc("/api/key", POST(keyHandler))
	http.HandleFunc("/api/keys", GET(keysHandler))
	http.HandleFunc("/api/enter", auth(indexHandler))
	http.HandleFunc("/api/ver", versionHandler)
	http.HandleFunc("/api/sign", signHandler)
}
