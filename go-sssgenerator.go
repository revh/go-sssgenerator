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

// check for the commands: init and build
func main() {

	// check if args are enough
	if len(os.Args) < 2 {
		fmt.Println("init or build subcommand is required")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		// exists if there is not projectname
		if len(os.Args) < 3 {
			fmt.Println("a project name is required")
			os.Exit(1)
		}

		initCmd(os.Args[2])
	case "build":
		// if there is no dirnam arg, defaults to current dir
		dirname := "."
		if len(os.Args) >= 3 {
			dirname = os.Args[2]
		}

		buildCmd(dirname)
	default:
		os.Exit(1)
	}
}

// makes all the needed directory
// and creates the .sssogenerator file ready to use
func initCmd(dirname string) {
	fmt.Printf("INIT CMD dirname: %s\n", dirname)

	os.MkdirAll(path.Join(dirname, "src", "posts"), os.ModePerm)
	os.MkdirAll(path.Join(dirname, "src", "templates"), os.ModePerm)

	ioutil.WriteFile(path.Join(dirname, ".sssgenerator"), ssgeneratorTpl, os.ModePerm)
}

// looks for posts and templates, merge all and
// output the results in the output directory
func buildCmd(dirname string) {
	fmt.Printf("BUILD CMD dirname: %s\n", dirname)

	// read configuration
	config := utils.ReadConfig(path.Join(dirname, ".sssgenerator"))

	// collect all templates
	generator.CollectTemplats(config.TemplatesDir)

	// collect all posts
	posts := generator.CollectPosts(config.PostsDir)

	var index *generator.Post
	for i, post := range posts {
		// save the index for a future use
		if post.Index {
			index = post
			continue
		}

		// outputs all the posts in the output directory
		generator.WritePost(post, config.OutputDir)
		fmt.Printf("Generated %d %s\n", i, path.Join(config.OutputDir, post.Status, post.FileName))
	}

	if index != nil {
		// add all the posts to the index's Related prop
		for _, post := range posts {
			if !post.Index {
				index.Related = append(index.Related, post)
			}
		}

		// output the index in the output directory
		generator.WritePost(index, config.OutputDir)
		fmt.Printf("Generated index %s\n", path.Join(config.OutputDir, index.Status, index.FileName))
	}
}
