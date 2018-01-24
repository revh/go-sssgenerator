package utils

import (
	"strings"
	"testing"
)

var testDataPath = "../generator/testdata/.sssgenerator"
var testDataCollectPath = "../generator/testdata/"

func TestReadConfig(t *testing.T) {
	config := ReadConfig(testDataPath)
	expected := "generator/testdata/posts"
	got := config.PostsDir

	if !strings.Contains(got, expected) {
		t.Errorf("PostsDir: expected `%s`, got `%s`", expected, got)
	}

	expected = "generator/testdata/templates"
	got = config.TemplatesDir

	if !strings.Contains(got, expected) {
		t.Errorf("TemplatesDir: expected %s, got `%s`", expected, got)
	}

	expected = "generator/testdata/output"
	got = config.OutputDir

	if !strings.Contains(got, expected) {
		t.Errorf("OutputDir: expected `%s`, got `%s`", expected, got)
	}
}

func TestCollectFiles(t *testing.T) {
	files := CollectFiles(testDataCollectPath, ".sssgenerator")

	expected := 1
	got := len(files)

	if got != expected {
		t.Errorf("CollectFiles: expected `%d`, got `%d`", expected, got)
	}
}
