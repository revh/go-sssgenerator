package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/revh/go-sssgenerator/generator"
	"github.com/revh/go-sssgenerator/utils"
)

var ssgeneratorTpl = []byte(`{
    "posts": "src/posts/",
    "templates":"src/templates/",
    "output":"."
}`)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("init or build subcommand is required")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		if len(os.Args) < 3 {
			fmt.Println("a project name is required")
			os.Exit(1)
		}
		initCmd(os.Args[2])
	case "build":
		dirname := "."
		if len(os.Args) >= 3 {
			dirname = os.Args[2]
		}
		buildCmd(dirname)
	default:
		os.Exit(1)
	}
}

func initCmd(dirname string) {
	fmt.Printf("INIT CMD dirname: %s\n", dirname)

	os.MkdirAll(path.Join(dirname, "src", "posts"), os.ModePerm)
	os.MkdirAll(path.Join(dirname, "src", "templates"), os.ModePerm)
	ioutil.WriteFile(path.Join(dirname, ".sssgenerator"), ssgeneratorTpl, os.ModePerm)
}

func buildCmd(dirname string) {
	fmt.Printf("BUILD CMD dirname: %s\n", dirname)

	config := utils.ReadConfig(path.Join(dirname, ".sssgenerator"))

	generator.CollectTemplats(config.TemplatesDir)
	posts := generator.CollectPosts(config.PostsDir)

	var index *generator.Post
	for i, post := range posts {
		if post.Index {
			index = post
		} else {
			generator.WritePost(post, config.OutputDir)
			fmt.Printf("Generated %d %s\n", i, path.Join(config.OutputDir, post.Status, post.FileName))
		}
	}

	if index != nil {
		for _, post := range posts {
			if !post.Index {
				index.Related = append(index.Related, post)
			}
		}
		generator.WritePost(index, config.OutputDir)
		fmt.Printf("Generated index %s\n", path.Join(config.OutputDir, index.Status, index.FileName))
	}
}
