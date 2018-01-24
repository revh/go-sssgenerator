# go-sssgenerator 
## (go simple static site generator)

# Getting started

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

run the buid command

`$ cd mynewblog`

`$ go-sssgenerator init mynewblog`

if it works you should have the generated files in mynewblog directory

for a working example check https://github.com/revh/revh.github.io

# Warning
This program is only intened as a study project

