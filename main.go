package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"

	"google.golang.org/appengine"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	// html templates
	booklistTmpl     = parseTemplate("list.html")
	authordetailTmpl = parseTemplate("detail.html")
	// Google Books API uri format
	googlebooksformat = "https://www.googleapis.com/books/v1/volumes?q=%s&startIndex=%v&maxResults=%v&country=%s&langRestrict=%s&download=epub&printType=books&showPreorders=false&fields=%s"
)

// Book holds the an item (a Volume) from the result of a Google Books query
type Book struct {
	ID         string `json:"id"`
	VolumeInfo struct {
		Title    string   `json:"title"`
		Subtitle string   `json:"subtitle"`
		Pages    int      `json:"pageCount"`
		Authors  []string `json:"authors"`
	} `json:"volumeInfo"`
}

// Books is a collection container that implements sort interface for sorting by pages
type Books []Book

func (slice Books) Len() int {
	return len(slice)
}
func (slice Books) Less(i, j int) bool {
	return slice[i].VolumeInfo.Pages > slice[j].VolumeInfo.Pages
}
func (slice Books) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func main() {
	log.Println("bookshelf")

	r := mux.NewRouter()
	r.Path("/").Methods("GET").Handler(apiHandler(listBooksHandler))
	r.Path("/{author}").Methods("GET").Handler(apiHandler(showBooksByAuthorHandler))

	// AppEngine / Compute Engine Health check endpoint
	r.Methods("GET").Path("/_ah/health").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})

	// Apache format logging
	http.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, r))

	appengine.Main()

}

func listBooksHandler(w http.ResponseWriter, r *http.Request) *apiError {

	authors := struct {
		Authors []string
	}{
		Authors: []string{
			"Miguel de Cervantes",
			"Charles Dickens",
			"Antoine de Saint-Exupéry",
			"J. K. Rowling",
			"J. R. R. Tolkien",
			"Agatha Christie",
			"Lewis Carroll",
			"C. S. Lewis",
			"Dan Brown",
			"Arthur Conan Doyle",
			"Jules Verne",
			"Stephen King",
			"Stieg Larsson",
			"George Orwell",
			"Ian Fleming",
			"James Patterson",
			"Anne Rice",
			"Terry Pratchett",
			"George R. R. Martin",
			"Edgar Rice Burroughs",
			"Michael Connelly",
			"Jo Nesbo",
		},
	}

	return booklistTmpl.Execute(w, r, authors)
}

func showBooksByAuthorHandler(w http.ResponseWriter, r *http.Request) *apiError {
	// grab request variables
	vars := mux.Vars(r)

	books, err := getAuthorBooks(vars["author"])
	if err != nil {
		return apiErrorf(err, "Unable to retrieve Author Info: %v")
	}

	sort.Sort(books)

	authorInfo := struct {
		Author string
		Books  Books
	}{
		Author: vars["author"],
		Books:  books,
	}

	return authordetailTmpl.Execute(w, r, authorInfo)
}

// getAuthorBooks uses the Google Books API to retrieve books by an author
func getAuthorBooks(author string) (Books, error) {

	var books []Book

	googlebooksresponse := struct {
		Items []Book `json:"items"`
	}{}

	// https://www.googleapis.com/books/v1/volumes?q=%s&startIndex=%v&maxResults=%v&country=%s&langRestrict=%s&download=epub&printType=books&showPreorders=false=fields=%s"
	uri := fmt.Sprintf(googlebooksformat,
		url.QueryEscape(author), // query
		0,    // start
		40,   // count
		"US", // country
		"en", // language
		"items(id,accessInfo(epub/isAvailable),volumeInfo(title,subtitle,language,pageCount,authors))", // filter fields
	)
	log.Println(uri)

	resp, err := http.Get(uri)
	if err != nil {
		return books, err
	}
	defer resp.Body.Close()

	bytesdata, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return books, err
	}

	err = json.Unmarshal(bytesdata, &googlebooksresponse)

	books = googlebooksresponse.Items

	return books, nil

}

// http://blog.golang.org/error-handling-and-go
type apiHandler func(http.ResponseWriter, *http.Request) *apiError

type apiError struct {
	Error   error
	Message string
	Code    int
}

func (fn apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *apiError, not os.Error.
		log.Printf(
			"Handler error: status code: %d, message: %s, resulting from err: %#v",
			e.Code, e.Message, e.Error,
		)
		http.Error(w, e.Message, e.Code)
	}
}

func apiErrorf(err error, format string, v ...interface{}) *apiError {
	return &apiError{
		Error:   err,
		Message: fmt.Sprintf(format, v...),
		Code:    500,
	}
}
