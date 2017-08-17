package main
// this requieres db/create_db.sh to be run to create the sqlite database and populate it with tables and some initial data.

import(
"log"
"time"
"net/http"
"encoding/json"
"database/sql"
_ "github.com/mattn/go-sqlite3" //remind me what the underscore does?
"github.com/gorilla/pat" //alternative to net/http
)


type Datum struct {
    Id int64
    First string
    Second time.Time
}



func main() {
    r := pat.New()
    //r.Get("/", rootHandler)
    r.Get("/first_table", getFirstTableHandler)
    r.Post("/first_table", postFirstTableHandler)
    //r.Patch("/first_table", patchFirstTableHandler)
    //r.Put("/first_table", putFirstTableHandler)
    //r.Delete("/first_table", deleteFirstTableHandler)
    http.Handle("/", r)
    log.Print("Listening on localhost:8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func getFirstTableHandler(res http.ResponseWriter, req *http.Request) {
    db, err := sql.Open("sqlite3", "./db/sqlite_example.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    rows, err := db.Query("SELECT * FROM first_table")
    if err != nil {
        panic(err)
    }

    var data []Datum
    for rows.Next() {
        var id int64
        var first_var string
        var second_var time.Time
        err = rows.Scan(&id, &first_var, &second_var)
        if err != nil {
            panic(err)
        }
        data = append(data, Datum{id,first_var,second_var})
    }

    js, err := json.Marshal(data)
    if err != nil {
        panic(err)
    }

    res.Header().Set("Content-Type", "application/json")
    res.Write(js)
}

func postFirstTableHandler(res http.ResponseWriter, req *http.Request) {
    db, err := sql.Open("sqlite3", "./db/sqlite_example.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    stmt, err := db.Prepare("INSERT INTO first_table(first_col, second_col) VALUES(?,?)")
    if err != nil {
        panic(err)
    }


    var m map[string]interface{}
    err = json.NewDecoder(req.Body).Decode(&m)
    req.Body.Close()
    log.Print(err, m)

    //req.ParseForm()
    dbRes, err := stmt.Exec(m["first"], m["second"])
    if err != nil {
        panic(err)
    }

    id, err := dbRes.LastInsertId()
    if err != nil {
        panic(err)
    }

    js, err := json.Marshal(id)
    if err != nil {
        panic(err)
    }

    res.WriteHeader(http.StatusCreated)
    res.Write(js)

}
