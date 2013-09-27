package hello

import (
    "html/template"
    "net/http"
    "time"

    "appengine"
    "appengine/datastore"
    "appengine/user"
)

type Greeting struct {
    Author  string
    Content string
    Date    time.Time
}

type Vegetable struct {
	Name string
}

type Cultivation struct {
	Name string
    Date    time.Time
	Veggy Vegetable
	
}

type Garden struct {
	Name string
    Date    time.Time
	Cultivations []Cultivation
	
}

func init() {
    http.HandleFunc("/", root)
    http.HandleFunc("/sign", sign)
	
	http.HandleFunc("/listGarden", listGarden)
//	http.HandleFunc("/viewGarden", viewGarden)
//	http.HandleFunc("/editGarden", editGarden)
	http.HandleFunc("/newGarden", createGarden)
}

// guestbookKey returns the key used for all guestbook entries.
func guestbookKey(c appengine.Context) *datastore.Key {
    // The string "default_guestbook" here could be varied to have multiple guestbooks.
    return datastore.NewKey(c, "Guestbook", "default_guestbook", 0, nil)
}

func gardenContainerKey(c appengine.Context) *datastore.Key {
    // The string "default_guestbook" here could be varied to have multiple guestbooks.
    return datastore.NewKey(c, "GardenContainer", "default_gardenContainer", 0, nil)
}

func root(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    // Ancestor queries, as shown here, are strongly consistent with the High
    // Replication Datastore. Queries that span entity groups are eventually
    // consistent. If we omitted the .Ancestor from this query there would be
    // a slight chance that Greeting that had just been written would not
    // show up in a query.
    q := datastore.NewQuery("Greeting").Ancestor(guestbookKey(c)).Order("-Date").Limit(10)
    greetings := make([]Greeting, 0, 10)
    if _, err := q.GetAll(c, &greetings); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if err := guestbookTemplate.Execute(w, greetings); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

var guestbookTemplate = template.Must(template.ParseFiles("guestbook.html"))
var gardenListTemplate = template.Must(template.ParseFiles("gardenList.html"))
var newGardenTemplate = template.Must(template.ParseFiles("newGarden.html"))


func sign(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    g := Greeting{
        Content: r.FormValue("content"),
        Date:    time.Now(),
    }
    if u := user.Current(c); u != nil {
        g.Author = u.String()
    }
    // We set the same parent key on every Greeting entity to ensure each Greeting
    // is in the same entity group. Queries across the single entity group
    // will be consistent. However, the write rate to a single entity group
    // should be limited to ~1/second.
    key := datastore.NewIncompleteKey(c, "Greeting", guestbookKey(c))
    _, err := datastore.Put(c, key, &g)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/", http.StatusFound)
}


func listGarden(w http.ResponseWriter, r *http.Request) {
	
    c := appengine.NewContext(r)
    // Ancestor queries, as shown here, are strongly consistent with the High
    // Replication Datastore. Queries that span entity groups are eventually
    // consistent. If we omitted the .Ancestor from this query there would be
    // a slight chance that Greeting that had just been written would not
    // show up in a query.
    q := datastore.NewQuery("Garden").Ancestor(gardenContainerKey(c)).Order("-Date").Limit(10)
    gardens := make([]Garden, 0, 10)
    if _, err := q.GetAll(c, &gardens); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if err := gardenListTemplate.Execute(w, gardens); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
	
	
}

func createGarden(w http.ResponseWriter, r *http.Request) {

    if err := newGardenTemplate.Execute(w, 1); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}