# go-sssgenerator 
## (go simple static site generator)

# Getting started

installation

`$ go get github.com/revh/go-sssgenerator/...`

To create a new site write the following command

`$ go-sssgenerator init mynewblog`

You should have a new directory tree like this

```
$ tree mynewblog/ -a
mynewblog/
├── .sssgenerator
└── src
    ├── posts
    └── templates
```

Put some frontmatter markdown in the posts directory
Put some html file in the templates directory

Run the buid command

`$ cd mynewblog`

`$ go-sssgenerator build`

If it works you should have the generated files in mynewblog directory

For a working example check https://github.com/revh/revh.github.io

# Warning

This program is only intended as a study project

