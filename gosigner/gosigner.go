package gosigner

import (
	"crypto/hmac"
	"encoding/base64"
        "sugo"
)

const (
	version = "0.0.2"
)

type Signature struct {
	Signature, Content, Key string
}

type Version struct{
        Version string
}

func signHandler(c sugo.Context) {
	var data Signature
	c.ReadJSON(&data)
	keyBytes := []byte(data.Key)
	content := data.Content
	mac := hmac.NewSHA1(keyBytes)
	mac.Write([]byte(content))

	out := make([]byte, base64.StdEncoding.EncodedLen(len(mac.Sum())))
	base64.StdEncoding.Encode(out, mac.Sum())

	c.Render(&Signature{string(out), content, data.Key})
}

func init() {

	signResource := &sugo.Resource{
		Accepts: "application/json",
		Renders: "application/json",
		POST:     &sugo.Path{"/api/sign", signHandler},
	}
        
        versionResource := &sugo.Resource{
                Renders: "application/json",
                GET : &sugo.Path{"/api/ver", func(c sugo.Context) {
                                                c.Render(&Version{version})
                                            }},
        }
        sugo.Make("/") 
	sugo.Add(signResource)
        sugo.Add(versionResource)
        sugo.Start()
}
