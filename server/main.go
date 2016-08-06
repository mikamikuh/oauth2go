package main

import (
  "strconv"
  "log"
  "net/http"
  "html/template"
  "database/sql"
  _ "github.com/lib/pq"
  "github.com/RangelReale/osin"
  "github.com/ory-am/osin-storage/storage/postgres"
)

type AuthPage struct {
  Client string
  Message string
  URL string
}

func main() {

  db, err := sql.Open("postgres", "user=postgres dbname=mikamikuh sslmode=disable")

  if err != nil {
    log.Fatal(err)
    return
  }

  s := postgres.New(db)

  err = s.CreateSchemas()

  if err != nil {
    log.Fatal(err)
    return
  }

  // TestStorage implements the "osin.Storage" interface
  server := osin.NewServer(osin.NewServerConfig(), s)

  http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
    resp := server.NewResponse()
    defer resp.Close()

    id := r.URL.Query().Get("client_id")
    log.Print(id)
    c, err := resp.Storage.GetClient(id)
    if err != nil {
      log.Print(err)
      return
    }

    p := &AuthPage {Client: c.GetId(), Message: "View your email address", URL: "/approval"}
    t, err := template.ParseFiles("server/html/auth.html")

    if err != nil {
      log.Print(err)
      return
    }
    t.Execute(w, p)
  })

  http.HandleFunc("/approval", func(w http.ResponseWriter, r *http.Request) {
    resp := server.NewResponse()
    defer resp.Close()

    if ar := server.HandleAuthorizeRequest(resp, r); ar != nil {
      has := r.Form.Get("submit_access")

      ar.Authorized, err = strconv.ParseBool(has)

      if err != nil {
        log.Print(err)
        return
      }

      server.FinishAuthorizeRequest(resp, r, ar)
    }

    osin.OutputJSON(resp, w, r)
  })

  http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
    resp := server.NewResponse()
    defer resp.Close()

    if ar := server.HandleAccessRequest(resp, r); ar != nil {
      ar.Authorized = true
      server.FinishAccessRequest(resp, r, ar)
    }
    osin.OutputJSON(resp, w, r)
  })

  http.ListenAndServe(":14000", nil)
}
