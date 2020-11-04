package main

import (
	"context"
	"log"
	"os"

	"github.com/machinebox/graphql"
)

// Response of API
type Response struct {
	License struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"License"`
}

func main() {
	// Create a client
	client := graphql.NewClient("http://api.github.com/graphql")

	// make a request to github API
	req := graphql.NewRequest(`
		query {
			license(key: "apache-2.0") {
				name
				description
			}
		}
	`)

	var githubToken = os.Getenv("GITHUB_TOKEN")
	req.Header.Add("Authorization", "bearer "+githubToken)

	// define a context for the request
	ctx := context.Background()

	// run it and capture the response
	var respData Response
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}
	log.Println(respData.License.Description)
}
