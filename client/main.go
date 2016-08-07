package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type LoginPage struct {
	URL string
}

type TokenResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Expires      int    `json:"expires_in"`
}

func main() {
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		p := &LoginPage{URL: "http://localhost:14000/auth?client_id=1111&response_type=code"}
		t, err := template.ParseFiles("client/html/login.html")

		if err != nil {
			log.Print(err)
			return
		}

		t.Execute(w, p)
	})

	http.HandleFunc("/appauth", func(w http.ResponseWriter, r *http.Request) {

		c := r.URL.Query().Get("code")

		if c != "" {
			data := url.Values{}
			data.Add("client_id", "1111")
			data.Add("client_secret", "aabbccdd")
			data.Add("grant_type", "authorization_code")
			data.Add("code", c)
			urlStr := "http://localhost:14000/token"

			client := &http.Client{}
			r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			resp, err := client.Do(r)
			defer resp.Body.Close()

			if err != nil {
				log.Print(err)
				return
			}

			tr := &TokenResponse{}
			err = json.NewDecoder(resp.Body).Decode(tr)

			if err != nil {
				log.Print(err)
				return
			}

			t, _ := template.ParseFiles("client/html/approved.html")
			t.Execute(w, tr)
		} else {
			t, _ := template.ParseFiles("client/html/denied.html")
			t.Execute(w, nil)
		}
	})

	log.Print("Client started...")
	http.ListenAndServe(":14001", nil)
}
