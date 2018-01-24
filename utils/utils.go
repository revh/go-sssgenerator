package utils

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"path/filepath"
)

type Config struct {
	BaseDirName  string
	PostsDir     string `json:"posts"`
	TemplatesDir string `json:"templates"`
	OutputDir    string `json:"output"`
}

func ReadConfig(filename string) *Config {
	config := &Config{}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		panic(err)
	}

	abs, _ := filepath.Abs(filename)
	base := filepath.Dir(abs)
	config.BaseDirName = base
	config.PostsDir = path.Join(base, config.PostsDir)
	config.TemplatesDir = path.Join(base, config.TemplatesDir)
	config.OutputDir = path.Join(base, config.OutputDir)

	return config
}

func CollectFiles(dirname string, extension string) []string {
	files, err := filepath.Glob(path.Join(dirname, "/*"))
	if err != nil {
		panic(err)
	}

	var selectedFiles []string
	for _, file := range files {
		if filepath.Ext(file) == extension {
			selectedFiles = append(selectedFiles, file)
		}
	}
	return selectedFiles
}
