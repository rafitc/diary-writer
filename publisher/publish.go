package publisher

// this is the publisher package
// it clone a given repo
// then create a MDX file in dedicated place
// download and put pictures
// then push the changes to the repo
// Repo is connected to the vercel. so it will trigger build automatically

// with go modules disabled
import (
	"fmt"
	"main/core"
	"main/logger"
	"main/models"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

var log = logger.Logger

func PublishContent(content string, date string, assets []models.ASSET) {
	// Clone the repository
	_, err := cloneAndFetchTheLatestChanges()
	if err != nil {
		log.Errorf("Error in cloning the repository: %v", err)
		return
	}

	// Create a MDX file
	writeMdxFile(content, date)

	// downloadAssets(assets)
	// DownloadAndPutPictures()
	// // Push the changes to the repo
	// PushChanges()
}

func downloadAssets(assets []models.ASSET) {
	for _, asset := range assets {
		log.Info("Downloading asset: ", asset)
	}
}

func writeMdxFile(content string, date string) {
	// Create a MDX file
	log.Debugf("Creating a MDX file")
	file, err := os.Create(core.Config.PUBLISH.CLONE_DIRECTORY + "/content/" + fmt.Sprintf("%s.mdx", date))
	if err != nil {
		log.Errorf("Error in creating a MDX file: %v", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		log.Errorf("Error in writing to the MDX file: %v", err)
		return
	}
	_, err = file.WriteString("\n")
	if err != nil {
		log.Errorf("Error in writing to the MDX file: %v", err)
		return
	}
}

func cloneAndFetchTheLatestChanges() (*git.Repository, error) {

	// Clone the given repository to the given directory
	log.Info("git clone https://github.com/go-git/go-git")

	repo, err := git.PlainClone(core.Config.PUBLISH.CLONE_DIRECTORY, false, &git.CloneOptions{
		URL:      core.Config.PUBLISH.GITHUB_REPO,
		Progress: os.Stdout,
		Auth: &http.BasicAuth{
			Username: core.Config.PUBLISH.GITHUB_USERNAME, // Can be anything for token-based auth (e.g., "x-access-token")
			Password: core.Config.PUBLISH.GITHUB_AUTH_TOKEN,
		},
	})

	// Handle errors during cloning
	if err != nil {
		if err == git.ErrRepositoryAlreadyExists {
			// Open the existing repository
			repo, err = git.PlainOpen(core.Config.PUBLISH.CLONE_DIRECTORY)
			if err != nil {
				log.Fatalf("Failed to open existing repository: %v", err)
				return nil, err
			}
			log.Info("Opened existing repository successfully")
		} else {
			log.Errorf("Error in cloning the repository: %v", err)
			return nil, err
		}
	} else {
		log.Info("Repository cloned successfully")
	}

	log.Debugf("Fetching the latest changes")

	err = repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: core.Config.PUBLISH.GITHUB_USERNAME, // Can be anything for token-based auth (e.g., "x-access-token")
			Password: core.Config.PUBLISH.GITHUB_AUTH_TOKEN,
		},
	})

	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			log.Debug("No new changes to fetch")
		} else {
			log.Errorf("Error in fetching the latest changes: %v", err)
			return nil, err
		}
	} else {
		log.Debug("Fetched the latest changes successfully")
	}

	// pull the latest changes
	log.Debug("Pulling the latest changes")
	worktree, err := repo.Worktree()
	if err != nil {
		log.Errorf("Error getting worktree: %v", err)
		return nil, err
	}
	err = worktree.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: core.Config.PUBLISH.GITHUB_USERNAME, // Can be anything for token-based auth (e.g., "x-access-token")
			Password: core.Config.PUBLISH.GITHUB_AUTH_TOKEN,
		},
	})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			log.Debug("No new changes to pull")
		} else {
			log.Errorf("Error in pulling the latest changes: %v", err)
			return nil, err
		}
	} else {
		log.Debug("Pulled the latest changes successfully")
	}

	return repo, nil

}
