package generator

import (
	"testing"
)

var testDataOnePath = "testdata/posts/one.md"
var testDataPath = "testdata/posts/"
var testDataTemplatePath = "testdata/templates/"
var testDataOutoutPath = "testdata/output/"

func TestReadPost(t *testing.T) {
	post := ReadPost(testDataOnePath)

	expected := "One"
	got := post.Meta["Title"]

	if got != expected {
		t.Errorf("CollectPosts: expected `%s`, got `%s`", expected, got)
	}

	expected = "<p>FileName: index.html Title: One</p>\n"
	got = string(post.Content)

	if got != expected {
		t.Errorf("CollectPosts: expected `%s`, got `%s`", expected, got)
	}

	expected = "index.html"
	got = post.FileName

	if got != expected {
		t.Errorf("CollectPosts: expected `%s`, got `%s`", expected, got)
	}
}

func TestCollectPosts(t *testing.T) {
	posts := CollectPosts(testDataPath)

	expected := "One"
	got := posts[0].Meta["Title"]

	if got != expected {
		t.Errorf("CollectPosts: expected `%s`, got `%s`", expected, got)
	}
}

func TestCollectTemplats(t *testing.T) {
	tmpl := CollectTemplats(testDataTemplatePath)

	expected := 2
	got := len(tmpl)

	if got != expected {
		t.Errorf("CollectPosts: expected `%d`, got `%d`", expected, got)
	}
}

func TestWritePost(t *testing.T) {
	CollectTemplats(testDataTemplatePath)
	post := ReadPost(testDataOnePath)
	WritePost(post, testDataOutoutPath)
}
