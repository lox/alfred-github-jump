package main

import (
	"fmt"
	"log"

	"github.com/google/go-github/github"
	"github.com/pascalw/go-alfred"
	"golang.org/x/oauth2"
)

type Repository struct {
	URL, Name, User, Description string
}

func reposCommand(queryTerms []string) {
	alfred.InitTerms(queryTerms)

	response := alfred.NewResponse()
	defer response.Print()

	if !isLoggedIn() {
		response.AddItem(&alfred.AlfredResponseItem{
			Valid: false,
			Uid:   "login",
			Title: "You need to login first with gh-login",
		})
		return
	}

	repos, err := ListRepositories()
	if err != nil {
		response.AddItem(alfredError(err))
		return
	}

	for _, repo := range repos {
		if alfred.MatchesTerms(queryTerms, repo.Name) {
			response.AddItem(&alfred.AlfredResponseItem{
				Valid:    true,
				Uid:      repo.URL,
				Title:    fmt.Sprintf("%s/%s", repo.User, repo.Name),
				Subtitle: repo.Description,
				Arg:      repo.URL,
			})
		}
	}
}

func ListRepositories() ([]Repository, error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT id, url,description, name, user FROM repository")
	if err != nil {
		return nil, err
	}

	repos := []Repository{}

	for rows.Next() {
		var id, url, descr, name, user string
		err = rows.Scan(&id, &url, &descr, &name, &user)
		if err != nil {
			return nil, err
		}

		repos = append(repos, Repository{
			URL:         url,
			Name:        name,
			User:        user,
			Description: descr,
		})
	}

	return repos, nil
}

func nilableString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func UpdateRepositories(token *oauth2.Token) (int64, error) {
	tc := OAuthConf.Client(oauth2.NoContext, token)
	client := github.NewClient(tc)

	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
		Sort:        "pushed",
	}

	db, err := OpenDB()
	if err != nil {
		return 0, err
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	counter := int64(0)
	for {
		result, resp, err := client.Repositories.List("", opt)
		if err != nil {
			return counter, err
		}
		for _, repo := range result {
			log.Printf("Updating %s/%s", *repo.Owner.Login, *repo.Name)

			res, err := db.Exec(
				`INSERT OR REPLACE INTO repository (
					id,
					url,
					description,
					name, user,
					pushed_at,
					updated_at,
					created_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
				fmt.Sprintf("%s/%s", *repo.Owner.Login, *repo.Name),
				nilableString(repo.HTMLURL),
				nilableString(repo.Description),
				*repo.Name,
				*repo.Owner.Login,
				(*repo.PushedAt).Time,
				(*repo.UpdatedAt).Time,
				(*repo.CreatedAt).Time,
			)
			if err != nil {
				return counter, err
			}
			rows, _ := res.RowsAffected()
			counter += rows
		}
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}

	return counter, tx.Commit()
}

func updateCommand() {
	response := alfred.NewResponse()
	defer response.Print()

	token, err := loadToken()
	if err != nil {
		response.AddItem(alfredError(err))
		return
	}

	n, err := UpdateRepositories(token)
	if err != nil {
		response.AddItem(alfredError(err))
		return
	}

	response.AddItem(&alfred.AlfredResponseItem{
		Valid: false,
		Uid:   "updated",
		Title: fmt.Sprintf("Updated %d repositories from github", n),
	})
}
