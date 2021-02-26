package main

import (
	"os"
	"fmt"
	"strings"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"flag"
	"sort"
	"time"
	"math"
	"html/template"
	"io/ioutil"
	"encoding/json"
)

type ChangeDetail struct {
	CommitId string
	Message  string
	Author   string
	Leadtime float64
}

type ReleaseDetail struct {
	Application string
	TagName 	string
	ReleaseDate time.Time
	Author 		string
	Changes 	[]ChangeDetail
	ChangeVolume int
	LeadTime     float64
}


func main()  {
	var application, gitRepo, gitCred, releaseVersion string
	flag.StringVar(&application, "application", "", " (required) application/service name")
	flag.StringVar(&gitRepo, "gitRepo", "", " (required) git repository url")
	flag.StringVar(&gitCred, "gitCred", "", " (required) credential to access git repository in format (username:password)")
	flag.StringVar(&releaseVersion, "releaseVersion", "", " release version, latest release will be used if not provided")
	flag.Parse()
	if application == "" && gitRepo == "" && gitCred == "" {
		flag.PrintDefaults()
	}
	generateReleaseNotes(application,gitRepo,strings.Split(gitCred, ":")[0],strings.Split(gitCred, ":")[1], releaseVersion)
}

func generateReleaseNotes(application, scm_repo, scm_usr, scm_pwd, releaseVersion string) {
	r, _ := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: scm_repo,
        Auth: &http.BasicAuth{
			Username: scm_usr,
			Password: scm_pwd,
		},
	})
	tagrefs, err := r.TagObjects()
	CheckIfError(err)
	var tags[]*object.Tag 
	err = tagrefs.ForEach(func(t *object.Tag) error {
		tags = append(tags, t)	
		return nil
	})
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Tagger.When.After(tags[j].Tagger.When)
	})
	var changes[]ChangeDetail
	var leadtimeMinutes float64
	for i, t := range tags {
		var releasetagindex int = 0
		if 	releaseVersion != "" && releaseVersion == t.Name {
			releasetagindex = i
		}		
		if i == releasetagindex {
			var breakme bool = false
			cIter, err := r.Log(&git.LogOptions{From: t.Target})
			CheckIfError(err)
			err = cIter.ForEach(func(c *object.Commit) error {
				hash := c.Hash.String()
				if ( ((len(tags)-1)-i != 0) && c.Hash.String() == tags[i+1].Target.String()) {
					breakme = true
				}
				if !breakme {
					line := strings.Split(c.Message, "\n")
					change := ChangeDetail{CommitId : hash, Message: line[0], Author : c.Author.Name, Leadtime : t.Tagger.When.Sub(c.Author.When).Minutes()}
					leadtimeMinutes += t.Tagger.When.Sub(c.Author.When).Minutes()
					changes = append(changes, change)
				}
				
				return nil
			})
		}
		release := ReleaseDetail{Application: application, TagName : t.Name  , ReleaseDate : t.Tagger.When , Author : t.Tagger.Email, Changes : changes, ChangeVolume: len(changes), LeadTime : math.Trunc(leadtimeMinutes/float64(len(changes)))}

		tmpl := template.Must(template.ParseFiles("layout.html"))

		f, err := os.Create("ReleaseNotes.html")
		CheckIfError(err)
		err = tmpl.Execute(f, release)
		CheckIfError(err)
        f.Close()

		metricsJson,_ := json.MarshalIndent(release, ""," ")
		ioutil.WriteFile("ReleaseNotes.json", metricsJson, 0644)
	}
	CheckIfError(err)
}


func CheckIfError(err error) {
	if err == nil {
		return
	}
	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

