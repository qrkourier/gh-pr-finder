package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/shurcooL/graphql"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

func main() {
	app := &cli.App{
		Name:  "gh-pr-finder",
		Usage: "Finds PRs by author across multiple GitHub owners",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "owners",
				Usage:    "Comma-separated list of GitHub repo owners",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "authors",
				Usage:    "Comma-separated list of GitHub usernames",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "branches",
				Usage:    "Limit results to PRs based on a comma-separated list of refs",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			return findPRs(c.String("owners"), c.String("authors"), c.String("branches"))
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func findPRs(owners string, authors string, branches string) error {
	// Extract the GitHub token from the environment variable.
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("GITHUB_TOKEN environment variable is not set")
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := graphql.NewClient("https://api.github.com/graphql", httpClient)

	// Define the GraphQL query.
	var query struct {
		Search struct {
			Edges []struct {
				Node struct {
					PullRequest struct {
						Title graphql.String
						URL   graphql.String
					} `graphql:"... on PullRequest"`
				}
			}
		} `graphql:"search(query: $query, type: ISSUE, first: 10)"`
	}

	// Split the owners string into a slice.
	ownerList := strings.Split(owners, ",")

	// Split the authors string into a slice.
	authorList := strings.Split(authors, ",")

	// Iterate over each owner.
	for _, owner := range ownerList {
		for _, author := range authorList {
			// Construct the search query string.
			searchQuery := fmt.Sprintf("is:pr is:open author:%s user:%s", author, owner)

			if branches != "" {
				refList := strings.Split(branches, ",")
				for _, ref := range refList {
					searchQuery += fmt.Sprintf(" base:%s", ref)
				}
			}

			// Run the query.
			variables := map[string]interface{}{
				"query": graphql.String(searchQuery),
			}
			err := client.Query(context.Background(), &query, variables)
			if err != nil {
				return err
			}

			// Print the URLs of matching PRs.
			for _, edge := range query.Search.Edges {
				fmt.Println(edge.Node.PullRequest.URL)
			}
		}
	}
	return nil
}

