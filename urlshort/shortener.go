package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/go-yaml/yaml"
)

const (
	yamlDefault = `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	jsonDefault = `[ 
{"path": "/urlshort", "url": "https://github.com/gophercises/urlshort"},
{"path": "/urlshort-final",	"url": "https://github.com/gophercises/urlshort/tree/solution"}
]`
)

var (
	dbBucket = []byte("PahtsToUrlsBucket")
)

/*
 * create an http.Handler that will look at the path of every incoming
 * web request and determine if it will redirect the user or not
 * 	e.g.: if we had a redirect setup for /dogs to a specific site like this
 *		  /dogs -> https://www.somesite.com/a-story-about-dogs
 *  	  every request containing /dogs would be substituted accoridingly
 */
func main() {
	var (
		yamlFilePath string // path to yaml mapping file
		jsonFilePath string // path to json mapping file
	)

	// parse flags
	flag.StringVar(&yamlFilePath, "yaml_path", "", "YAML config file to use for mapping URLs to their shortened paths")
	flag.StringVar(&jsonFilePath, "json_path", "", "JSON config file to use for mapping URLs to their shortened paths")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := MapHandler(pathsToUrls, mux)

	// If no yaml file was given by the user via flag, use default value
	var yaml string
	if yamlFilePath == "" {
		yaml = yamlDefault
	} else {
		var err error
		if yaml, err = readYamlFile(yamlFilePath); err != nil {
			log.Fatalf("error parsing yaml file: %s", err)
		}
	}

	// declare yaml handler with mapHandler as a fallback
	yamlHandler, err := YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		log.Fatalf("error loading yaml handler: %s", err)
	}

	// If no json file was given by the user via flag, use default value
	var json string
	if jsonFilePath == "" {
		json = jsonDefault
	} else {
		var err error
		if json, err = readJSONFile(jsonFilePath); err != nil {
			log.Fatalf("error parsing json file: %s", err)
		}
	}

	// declare JSON handler with yamlHandler as a fallback
	jsonHandler, err := JSONHandler([]byte(json), yamlHandler)
	if err != nil {
		log.Fatalf("error loading json handler: %s", err)
	}

	// allocate data base
	fillDataBase()

	// declare DB handler with jsonHandler as a fallback
	dbHandler, err := DBHandler(jsonHandler)
	if err != nil {
		log.Fatalf("error loading db handler: %s", err)
	}

	fmt.Println("Starting the server on :8080")
	// default with all cascading handlers
	http.ListenAndServe(":8080", dbHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v, ok := pathsToUrls[r.URL.String()]
		if !ok {
			// fallback
			fallback.ServeHTTP(w, r)
		}
		// redirect
		fmt.Fprintln(w, v)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// parse yaml
	var pathsToUrls []Shorts
	err := yaml.Unmarshal(yml, &pathsToUrls)
	if err != nil {
		log.Fatalf("error parsing yaml data: %s", err)
		return nil, err
	}

	// no errors, thus returning http.HanderFunc
	return func(w http.ResponseWriter, r *http.Request) {
		req := r.URL.String()
		// check if requested short has a url entry
		for _, short := range pathsToUrls {
			if short.Path == req {
				fmt.Fprintln(w, short.URL)
				return
			}
		}
		// if requested path is not in yaml, fall back to MapHandler
		fallback.ServeHTTP(w, r)
	}, nil
}

// JSONHandler will parse provided JSON and return an http.HandlerFunc
// exact same procedure as in YAMLHandler
func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathsToUrls []Shorts

	err := json.Unmarshal(jsn, &pathsToUrls)
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, r *http.Request) {
		for _, short := range pathsToUrls {
			if short.Path == r.URL.String() {
				fmt.Fprintln(w, short.URL)
				return
			}
		}
		fallback.ServeHTTP(w, r)
	}, nil
}

// Shorts structure storing the config for shortened urls where Path
// is the provided shortcut and URL the actual URL
type Shorts struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

func readYamlFile(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func readJSONFile(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// DBHandler returns a handler with access to data base (here: BoltDB)
func DBHandler(fallback http.Handler) (http.HandlerFunc, error) {
	return func(w http.ResponseWriter, r *http.Request) {
		// connect to db
		db, err := bolt.Open("my.db", 0644, nil)
		if err != nil {
			log.Fatalf("error opening data base: %s", err)
		}
		defer db.Close()

		// retreive data from db
		// read-only transaction
		var url []byte
		req := r.URL.String()
		if err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(dbBucket)
			if b == nil {
				return err
			}

			url = b.Get([]byte(req))

			return nil
		}); err != nil {
			log.Fatalf("error reading data base: %s", err)
		}

		// if result is empty, no entry found with given key -> fallback
		if string(url) == "" {
			fallback.ServeHTTP(w, r)
			return
		}
		fmt.Fprintf(w, "DBHandler: %s\n", string(url))
	}, nil
}

func fillDataBase() error {
	// connect to db
	db, err := bolt.Open("my.db", 0644, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	// fill database
	// read-write transaction
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(dbBucket)
		if err != nil {
			return err
		}

		if err := b.Put([]byte("/urlshort-godoc"), []byte("https://godoc.org/github.com/gophercises/urlshort")); err != nil {
			return err
		}
		if err := b.Put([]byte("/yaml-godoc"), []byte("https://godoc.org/gopkg.in/yaml.v2")); err != nil {
			return err
		}

		return nil
	})
}
