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

	"github.com/revh/go-sssgenerator/utils"

	blackfriday "gopkg.in/russross/blackfriday.v2"
)

var templates map[string]*template.Template
var htmlExtension = ".html"
var markdownExtension = ".md"
var breakString = "---"

// Templates() return a list of templates
func Templates() map[string]*template.Template {
	return templates
}

// Post is the type of the post
// it contains some required information, metas and the content of the file
// in the content of the file, special variable can be used to access these properties
// ex: {{.FileName}} output the FileName prop, and Meta.MyProp output a custom prop declared
// directly in the file
type Post struct {
	FilePath string
	FileName string
	Template string
	Content  template.HTML
	Meta     map[string]string

	// Index is a special flag, the field Related
	// will be filled with all other posts
	Index   bool
	Related []*Post

	// Status move the file in a "status" directory
	// ex. outputdir/state/filename
	Status    string
	OutputDir string
}

// ReadPost reads the file and convert it to post object
// a Post file is a file that contains a Header part and a Content part
// for info check "YAML front matter"
func ReadPost(filepath string) *Post {
	p := &Post{}
	p.FilePath = filepath

	// prepare the meta map
	p.Meta = make(map[string]string)

	// open the file
	inFile, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer inFile.Close()

	// create a new text scanner for iterating the content line by line
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	// contain the content of the file
	// the content is all the text after the meta
	var contentBuffer bytes.Buffer

	// make some reflection utils
	elem := reflect.ValueOf(p).Elem()
	st := reflect.TypeOf(*p)

	// we start looking for the Header part
	readingHeader := false
	for scanner.Scan() {
		line := scanner.Text()

		// skip empty lines in the header
		if readingHeader && strings.TrimSpace(line) == "" {
			continue
		}

		// flags the start/end of the header
		if line == breakString {
			if !readingHeader {
				readingHeader = true
				continue
			} else {
				readingHeader = false
				continue
			}
		}

		// for the header part
		if readingHeader {

			// divide the name and value of the field
			pieces := strings.Split(line, ":")
			fieldName, fieldValue := pieces[0], strings.TrimLeft(pieces[1], " ")

			// check if there is a special Index flag
			if fieldName == "Index" && fieldValue == "true" {
				p.Index = true
				continue
			}

			// look if the field in the header exists in the Post struct
			_, ok := st.FieldByName(fieldName)

			// if exists set the value
			// otherwise add the field to Meta map
			if ok {
				elem.FieldByName(fieldName).SetString(fieldValue)
			} else {
				p.Meta[fieldName] = fieldValue
			}
			continue
		}

		// if we are in the Content part, copy the line in the buffer
		contentBuffer.WriteString(line)
		contentBuffer.WriteRune('\n')
	}

	// if there is no custom FileName use the post's filename
	if p.FileName == "" {
		p.FileName = path.Base(filepath)
	}

	// convert the filename from markdown to HTML
	p.FileName = strings.Replace(p.FileName, markdownExtension, htmlExtension, 1)

	// convert the content to a template
	tmpl, err := template.New("content").Parse(contentBuffer.String())
	if err != nil {
		panic(err)
	}

	// replace the content's variables to their values
	var output bytes.Buffer
	err = tmpl.Execute(&output, p)
	if err != nil {
		panic(err)
	}

	// do a first markdown conversion for the Content part
	// and save in the Content field
	html := blackfriday.Run(output.Bytes())
	p.Content = template.HTML(html)

	return p
}

// WritePost generates the html file for the given post
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

// CollectTemplats compiles a list of template searching into the provided dirname
func CollectTemplats(dirname string) []*template.Template {
	// make the template's map
	templates = make(map[string]*template.Template)

	// collect all html files in the templates directory
	files := utils.CollectFiles(dirname, htmlExtension)
	var tmpls []*template.Template

	// loop over all the files
	for _, file := range files {

		//check file in cache
		if _, ok := templates[filepath.Base(file)]; ok {
			continue
		}

		// convert to template
		tmpl, err := template.ParseFiles(file)
		if err != nil {
			panic(err)
		}

		// assign the template to the map
		templates[filepath.Base(file)] = tmpl
	}

	// return the values of the template map
	for _, tmpl := range templates {
		tmpls = append(tmpls, tmpl)
	}
	return tmpls
}

// CollectPosts compile sa list of post searching into a given dirname
func CollectPosts(dirname string) []*Post {

	// collect all the .md files
	files := utils.CollectFiles(dirname, markdownExtension)

	var posts []*Post

	// return a list of Posts obj
	for _, file := range files {
		post := ReadPost(file)
		posts = append(posts, post)
	}
	return posts
}
