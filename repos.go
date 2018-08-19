package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/go-github/github"
	"github.com/pascalw/go-alfred"
	"github.com/sahilm/fuzzy"

	"golang.org/x/oauth2"
)

type Repository struct {
	URL, Name, User, Description string
	LastUpdated                  time.Time
}

func (r Repository) FullName() string {
	return fmt.Sprintf("%s/%s", r.User, r.Name)
}

func reposCommand(fuzzyFilter string) {
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

	var names []string

	for _, repo := range repos {
		names = append(names, repo.FullName())
	}

	for _, match := range fuzzy.Find(fuzzyFilter, names) {
		response.AddItem(&alfred.AlfredResponseItem{
			Valid:    true,
			Uid:      repos[match.Index].URL,
			Title:    repos[match.Index].FullName(),
			Subtitle: repos[match.Index].Description,
			Arg:      repos[match.Index].URL,
		})
	}
}

func ListRepositories() ([]Repository, error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT id, url,description, name,user,updated_at FROM repository")
	if err != nil {
		return nil, err
	}

	repos := []Repository{}

	for rows.Next() {
		var id, url, descr, name, user string
		var updated time.Time
		err = rows.Scan(&id, &url, &descr, &name, &user, &updated)
		if err != nil {
			return nil, err
		}

		repos = append(repos, Repository{
			URL:         url,
			Name:        name,
			User:        user,
			Description: descr,
			LastUpdated: updated,
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

func githubTime(t *github.Timestamp) *time.Time {
	if t == nil {
		return nil
	}
	return &t.Time
}

func listUserRepositories(client *github.Client) ([]*github.Repository, error) {
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 45},
		Sort:        "pushed",
	}

	repos := []*github.Repository{}

	for {
		ctx := context.Background()
		result, resp, err := client.Repositories.List(ctx, "", opt)
		if err != nil {
			return repos, err
		}
		repos = append(repos, result...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}

	return repos, nil
}

func listStarredRepositories(client *github.Client) ([]*github.Repository, error) {
	opt := &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{PerPage: 45},
		Sort:        "pushed",
	}

	repos := []*github.Repository{}

	for {
		ctx := context.Background()
		result, resp, err := client.Activity.ListStarred(ctx, "", opt)
		if err != nil {
			return repos, err
		}
		for _, starred := range result {
			repos = append(repos, starred.Repository)
		}
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}

	return repos, nil
}

func UpdateRepositories(token *oauth2.Token) (int64, error) {
	tc := OAuthConf.Client(oauth2.NoContext, token)
	client := github.NewClient(tc)

	userRepos, err := listUserRepositories(client)
	if err != nil {
		return 0, err
	}

	starredRepos, err := listStarredRepositories(client)
	if err != nil {
		return 0, err
	}

	db, err := OpenDB()
	if err != nil {
		return 0, err
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	found := map[string]struct{}{}
	counter := int64(0)

	for _, repo := range append(userRepos, starredRepos...) {
		log.Printf("Updating %s/%s", *repo.Owner.Login, *repo.Name)

		name := fmt.Sprintf("%s/%s", *repo.Owner.Login, *repo.Name)
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
			name,
			nilableString(repo.HTMLURL),
			nilableString(repo.Description),
			*repo.Name,
			*repo.Owner.Login,
			githubTime(repo.PushedAt),
			githubTime(repo.UpdatedAt),
			githubTime(repo.CreatedAt),
		)
		if err != nil {
			return counter, err
		}
		found[name] = struct{}{}
		rows, _ := res.RowsAffected()
		counter += rows
	}

	existing, err := ListRepositories()
	if err != nil {
		return 0, err
	}

	// purge repos that don't exit any more
	for _, repo := range existing {
		if _, exists := found[repo.FullName()]; !exists {
			log.Printf("Repo %s doesn't exist, deleting", repo.FullName())

			_, err := db.Exec(
				`DELETE FROM repository WHERE id=?`,
				repo.FullName(),
			)
			if err != nil {
				return 0, err
			}

		}
	}

	return counter, tx.Commit()
}

func updateCommand() {
	token, err := loadToken()
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}

	n, err := UpdateRepositories(token)
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}

	fmt.Printf("Updated %d repositories from github", n)
}
