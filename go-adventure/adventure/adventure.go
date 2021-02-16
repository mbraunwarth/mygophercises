package adventure

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

var (
	// adventures Stores all loadable adventures.
	adventures []*Adventure
)

// Adventure structure represents a whole adventure which is intended to be
// stored in a single JSON file. An adventure is made up of one or more arcs.
type Adventure struct {
	Name     string
	FilePath string
	Arcs     []Arc
}

// New returns the adventure indicated by its name.
func New(name, path string) *Adventure {
	arcs, err := parseArcs(path)
	if err != nil {
		log.Fatalf("error parsing arcs: %s", err)
	}
	a := &Adventure{Name: name, FilePath: path, Arcs: arcs}
	return a
}

func RTFV(name string) *Adventure {
	var result *Adventure
	for _, a := range adventures {
		if a.Name == name {
			result = a
		}
	}
	return result
}

// The Loader interface.
type Loader interface {
	Load(paths ...string) error
}

// Service serves adventures.
type Service interface {
	Load(paths ...string) error
	ListAdventureNames() []string
}

// Load loads the directory intended to hold the adventures in JSON format and fill
// adventures list with their names and adventureToFilePath with corresponding data.
// If paths is specified it will read from those instead.
func Load(paths ...string) error {
	var cwd string
	if paths == nil {
		// get directory where adventure files are located
		// normally <programs root dir>/adventuers, maybe later changeable via flag

		// obtain file path to current working directory of the program (root dir)
		var err error
		if cwd, err = os.Getwd(); err != nil {
			return err
		}
	}

	// make directory out of file path with actual adventures directory appended
	adventuresDir := path.Join(cwd, "adventures")
	dir, err := os.Open(adventuresDir)
	if err != nil {
		return err
	}

	// get content of directory
	dirContent, err := dir.Readdirnames(0)
	if err != nil {
		return err
	}

	// read file names, only include json files
	for _, f := range dirContent {
		// treat each file name as the adventures name,
		// create new Adventure object and write in adventures list
		if path.Ext(f) == ".json" {
			// format name a bit
			n, p := f[:len(f)-len(".json")], path.Join(adventuresDir, f)
			n = strings.ReplaceAll(n, "-", " ")
			// fill data variables
			a := New(n, p)
			adventures = append(adventures, a)
		}
	}

	return nil
}

// ListAdventureNames returns a list of the names of all possible adventures.
func ListAdventureNames() []string {
	var names []string
	for _, a := range adventures {
		names = append(names, a.Name)
	}
	return names
}

// Arc structure represents a single arc of an adventure.
type Arc struct {
	ID      string
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	NextArcText string `json:"text"`
	NextArcName string `json:"arc"`
}

// ------------- Unexported Stuff -------------

// parseArcs parses the JSON data for an adventure and returns
// a list of arcs.
func parseArcs(adventurePath string) ([]Arc, error) {
	// preparing JSON file storing adventure
	content, err := ioutil.ReadFile(adventurePath)
	if err != nil {
		log.Fatalf("error reading file: %s", err)
	}

	// As the JSON is in an unconventional format
	// where arc IDs itself represent the key to the arc object,
	// the parsing process is a bit more involving and uses an empty
	// interface as a temporary main container for the arcs.
	var data interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		log.Fatalf("error unmarshalling json: %s", err)
	}

	// temporary container for data castet to a map
	container := data.(map[string]interface{})

	// manually parse the data by casting
	var arcs []Arc = make([]Arc, 0)
	for k, v := range container {
		values := v.(map[string]interface{})

		// story case
		sTmp := values["story"].([]interface{})
		story := make([]string, 0)
		for _, s := range sTmp {
			story = append(story, s.(string))
		}

		// options case
		oTmp := values["options"].([]interface{})
		options := make([]Option, 0)
		for _, tmp := range oTmp {
			o := tmp.(map[string]interface{})

			option := Option{
				NextArcText: o["text"].(string),
				NextArcName: o["arc"].(string),
			}

			options = append(options, option)
		}

		a := &Arc{
			ID:      k,
			Title:   values["title"].(string),
			Story:   story,
			Options: options,
		}

		arcs = append(arcs, *a)
	}

	return arcs, nil
}
