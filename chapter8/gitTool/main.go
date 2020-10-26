package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/levigross/grequests"
	"github.com/urfave/cli"
)

// GITHUB_TOKEN holds the value of the github token
var githubToken = os.Getenv("GITHUB_TOKEN")
var requestOptions = &grequests.RequestOptions{Auth: []string{githubToken, "x-oauth-basic"}}

// Repo holds github repo data
type Repo struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Forks    int    `json:"forks"`
	Private  bool   `json:"private"`
}

// File Structs for modelling JSON body in create Gist
type File struct {
	Content string `json:"content"`
}

// Gist holds the Gist data structure
type Gist struct {
	Description string          `json:"description"`
	Public      bool            `json:"public"`
	Files       map[string]File `json:"files"`
}

func getStats(url string) *grequests.Response {
	resp, err := grequests.Get(url, requestOptions)
	if err != nil {
		log.Fatalln("Unable to make request : ", err)
	}
	return resp
}

// Reads the files provided and creates Gist on github
func createGist(url string, args []string) *grequests.Response {
	// get first two arguments
	description := args[0]
	// remaining arguments are file names with path
	var fileContents = make(map[string]File)
	for i := 1; i < len(args); i++ {
		dat, err := ioutil.ReadFile(args[i])
		if err != nil {
			log.Println("Please check the filenames. Absolute path (or) same directory are allowed")
			return nil
		}
		var file File
		file.Content = string(dat)
		fileContents[args[i]] = file
	}
	var gist = Gist{Description: description, Public: true, Files: fileContents}
	var postBody, _ = json.Marshal(gist)
	var requestOptionsCopy = requestOptions
	// Add data to JSON field
	requestOptionsCopy.JSON = string(postBody)
	// make a Post request to Github
	resp, err := grequests.Post(url, requestOptionsCopy)
	if err != nil {
		log.Println("Create request failed for Github API")
	}
	return resp
}

func main() {
	app := cli.NewApp()
	// Define command for our client
	app.Commands = []cli.Command{
		{
			Name:    "fetch",
			Aliases: []string{"f"},
			Usage:   "Fetch the repo details with user. [Usage]: githubAPI fetch user",
			Action: func(c *cli.Context) error {
				if c.NArg() > 0 {
					// Github API logic
					var repos []Repo
					user := c.Args()[0]
					var repoURL = fmt.Sprintf("https://api.github.com/users/%s/repos", user)
					resp := getStats(repoURL)
					resp.JSON(&repos)
					log.Println(repos)
				} else {
					log.Println("Please give a username. See -h to see help")
				}
				return nil
			},
		},
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "Creates a gist from the given text. [Usage]: githubAPI name 'description' sample.txt",
			Action: func(c *cli.Context) error {
				if c.NArg() > 1 {
					// Github API Logic
					args := c.Args()
					var postURL = "https://api.github.com/gists"
					resp := createGist(postURL, args)
					log.Println(resp.String())
				} else {
					log.Println("Please give sufficient arguments. See -h to see help")
				}
				return nil
			},
		},
	}

	app.Version = "1.0"
	app.Run(os.Args)
}
