package github

import (
	"log"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestGithubFinder(t *testing.T) {
	finder := &GithubFinder{
		Token: os.Getenv("GITHUB_TOKEN"),
	}

	if err := finder.Init(); err != nil {
		log.Fatal(err)
	}

	query := "org:luxtagofficial apostille"
	result, err := finder.Find(query)
	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(*result)
}
