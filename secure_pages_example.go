package main


import ("net/http"
        "io"
        "io/ioutil"
        "log"
        "os"

        "encoding/json"

        "github.com/gorilla/sessions"
        )

//Config
var store *sessions.CookieStore
type Config struct {
    cookiePassword  string `json:"cookiePassword"`
    RootURL         string `json:"rootURL"`
    Port            string `json:"port"`
}
var config Config

func init() {
    file, err := ioutil.ReadFile("./config.json")
    if err != nil {
        log.Printf("File error: %v\n", err)
        os.Exit(1)
    }
    json.Unmarshal(file, &config)
    store = sessions.NewCookieStore([]byte(config.cookiePassword))
}

//Authentication/page protection middleware
func authenticatePage(f http.HandlerFunc) http.HandlerFunc {
    return func(res http.ResponseWriter, req *http.Request) {
        log.Print("starting authentication")
        session, err := store.Get(req, "session")
        if err != nil {
            log.Fatal("can't open cookie")
        }
        if session.Values["loggedin"] == "false" || session.Values["loggedin"] == nil { //return error/login page
            log.Print("authentication failed")
            log.Print(session.Values["loggedin"])
            failedLoginHandler(res, req)
        } else { //return page
            log.Print("authentication complete")
            f(res, req)
        }
    }
}

//Router
func main() {
    http.HandleFunc("/", rootHandler)
    http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/logout", logoutHandler)
    http.HandleFunc("/protectedpage", authenticatePage(protectedPageHandler))
    http.ListenAndServe(config.Port, nil)
}


//Handlers
func rootHandler(res http.ResponseWriter, req *http.Request) {
    io.WriteString(res, `
<!DOCTYPE html>
<html>
  <head></head>
  <body>
    <p><a href="/login">LOGIN</a></p>
    <p><a href="/logout">LOGOUT</a></p>
    <p><a href="/protectedpage">Test login</a></p>
  </body>
</html>`)
}

func failedLoginHandler(res http.ResponseWriter, req *http.Request) {
    io.WriteString(res, `
<!DOCTYPE html>
<html>
  <head></head>
  <body>
    <p>Failed login, try again.</p>
    <a href="/">Login Page</a>
  </body>
</html>`)
}

func loginHandler(res http.ResponseWriter, req *http.Request) {
    //write return page
    login(res,req)
    io.WriteString(res, `
<!DOCTYPE html>
<html>
  <head></head>
  <body>` +
   `<p>You are now logged in.</p>` +
   `<a href="/protectedpage">Test on this protected page</a>` +
`  </body>
</html>`)
}

func logoutHandler(res http.ResponseWriter, req *http.Request) {
    logout(res, req)
    http.Redirect(res, req, "/", http.StatusSeeOther)
}

func protectedPageHandler(res http.ResponseWriter, req *http.Request) {
    io.WriteString(res, `
<!DOCTYPE html>
<html>
  <head></head>
  <body>
    <p>This page shouldn't be accessible without logging in</p>
    <a href="/logout">Logout</a>
  </body>
</html>`)

}



//Login/logout
func login(res http.ResponseWriter, req *http.Request) {
    session, err := store.Get(req, "session")
    if err != nil {
        log.Fatal("can't open cookie")
    }
    session.Values["loggedin"] = "true"
    session.Save(req, res)
}

func logout(res http.ResponseWriter, req *http.Request) {
    session, err := store.Get(req, "session")
    if err != nil {
        log.Fatal("can't open cookie")
    }
    session.Values["loggedin"] = "false"
    session.Save(req, res)
}
