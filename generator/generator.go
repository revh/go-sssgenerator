package generator

import (
	"bufio"
	"bytes"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	blackfriday "gopkg.in/russross/blackfriday.v2"
)

import "../utils"

var templates map[string]*template.Template
var htmlExtension = ".html"
var markdownExtension = ".md"

func Templates() map[string]*template.Template {
	return templates
}

//Post is the type of the post
//it contains some required information, metas and the content of the file
//in the content of the file, special variable can be used to access these properties
//ex: {{.FileName}} output the FileName prop, and Meta.MyProp output a custom prop declared
//directly in the file
type Post struct {
	Index     bool
	FilePath  string
	FileName  string
	Template  string
	Meta      map[string]string
	Content   template.HTML
	Related   []*Post
	Status    string
	OutputDir string
}

//ReadPost read the file and convert it to post
func ReadPost(filepath string) *Post {
	p := &Post{}
	p.FilePath = filepath

	p.Meta = make(map[string]string)

	inFile, _ := os.Open(filepath)
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	readingMetas := false
	var contentBuffer bytes.Buffer

	elem := reflect.ValueOf(p).Elem()
	st := reflect.TypeOf(*p)

	for scanner.Scan() {
		line := scanner.Text()

		if readingMetas && strings.TrimSpace(line) == "" {
			continue
		}

		if line == "---" {
			if !readingMetas {
				readingMetas = true
				continue
			} else {
				readingMetas = false
				continue
			}
		}

		if readingMetas {
			pieces := strings.Split(line, ":")
			fieldName, fieldValue := pieces[0], strings.TrimLeft(pieces[1], " ")

			if fieldName == "Index" && fieldValue == "true" {
				p.Index = true
				continue
			}

			_, ok := st.FieldByName(fieldName)
			if ok {
				elem.FieldByName(fieldName).SetString(fieldValue)
			} else {
				p.Meta[fieldName] = fieldValue
			}
		} else {
			contentBuffer.WriteString(line)
			contentBuffer.WriteRune('\n')
		}
	}

	if p.FileName == "" {
		p.FileName = path.Base(filepath)
	}

	p.FileName = strings.Replace(p.FileName, markdownExtension, htmlExtension, 1)

	tmpl, err := template.New("content").Parse(contentBuffer.String())
	if err != nil {
		panic(err)
	}

	var output bytes.Buffer
	err = tmpl.Execute(&output, p)
	if err != nil {
		panic(err)
	}

	html := blackfriday.Run(output.Bytes())

	p.Content = template.HTML(html)
	return p
}

//WritePost generate the html file for the given post
func WritePost(post *Post, dirname string) {
	tmpl := templates[post.Template]
	outPath := path.Join(dirname, post.Status, post.FileName)

	os.MkdirAll(filepath.Dir(outPath), os.ModePerm)

	f, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(f)
	err = tmpl.Execute(w, post)
	if err != nil {
		panic(err)
	}

	w.Flush()
	f.Sync()
}

//CollectTemplats compile a list of template searching into the provided dirname
func CollectTemplats(dirname string) []*template.Template {
	templates = make(map[string]*template.Template)

	files := utils.CollectFiles(dirname, htmlExtension)
	var tmpls []*template.Template

	for _, file := range files {

		//check file in cache
		if _, ok := templates[filepath.Base(file)]; ok {
			continue
		}

		tmpl, err := template.ParseFiles(file)
		if err != nil {
			panic(err)
		}

		templates[filepath.Base(file)] = tmpl
	}

	for _, tmpl := range templates {
		tmpls = append(tmpls, tmpl)
	}

	return tmpls
}

//CollectPosts compile a list of post searching into the provided dirname
func CollectPosts(dirname string) []*Post {
	files := utils.CollectFiles(dirname, markdownExtension)

	var posts []*Post

	for _, file := range files {
		post := ReadPost(file)
		posts = append(posts, post)
	}
	return posts
}
