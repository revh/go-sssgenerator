package utils

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"path/filepath"
)

// Config is the type for convert .sssgenerator configs
type Config struct {
	BaseDirName  string
	PostsDir     string `json:"posts"`
	TemplatesDir string `json:"templates"`
	OutputDir    string `json:"output"`
}

// ReadConfig reads a given filename and convert to a Config obj
// the paths will be converted to absolute path
func ReadConfig(filename string) *Config {
	config := &Config{}

	// read the file
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// convert to a Config obj
	err = json.Unmarshal(content, &config)
	if err != nil {
		panic(err)
	}

	// convert relative paths to absolute paths
	abs, _ := filepath.Abs(filename)
	base := filepath.Dir(abs)
	config.BaseDirName = base
	config.PostsDir = path.Join(base, config.PostsDir)
	config.TemplatesDir = path.Join(base, config.TemplatesDir)
	config.OutputDir = path.Join(base, config.OutputDir)

	return config
}

// CollectFiles scans a directory and return a list of file of the given extension
func CollectFiles(dirname string, extension string) []string {
	files, err := filepath.Glob(path.Join(dirname, "/*"))
	if err != nil {
		panic(err)
	}

	// iterate over all the file returning only those that match the given extension
	var selectedFiles []string
	for _, file := range files {
		if filepath.Ext(file) == extension {
			selectedFiles = append(selectedFiles, file)
		}
	}
	return selectedFiles
}
