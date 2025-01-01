package publisher

// this is the publisher package
// it clone a given repo
// then create a MDX file in dedicated place
// download and put pictures
// then push the changes to the repo
// Repo is connected to the vercel. so it will trigger build automatically

// with go modules disabled
import (
	"encoding/json"
	"fmt"
	"io"
	"main/core"
	"main/logger"
	"main/models"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

var log = logger.Logger

func PublishContent(content string, date string, assets []models.ASSET, titleForJsonLog string, summaryForJsonLog string) bool {
	var status bool
	// Clone the repository
	_, err := cloneAndFetchTheLatestChanges()
	if err != nil {
		log.Errorf("Error in cloning the repository: %v", err)
		return false
	}

	// Create a MDX file
	status = writeMdxFile(content, date)
	if !status {
		log.Errorf("Error in writing the MDX file")
		return false
	}

	// download the assets
	status, dailyLogIcon := downloadAssets(assets, date)
	if !status {
		log.Errorf("Error in downloading the assets")
		return false
	}

	// update the log json
	status = updateLogJSON(titleForJsonLog, summaryForJsonLog, date, dailyLogIcon)
	if !status {
		log.Errorf("Error in updating the log json")
		return false
	}

	// Push the changes to the repo
	status = pushChanges(date)
	if !status {
		log.Errorf("Error in pushing the changes")
		return false
	}

	// Hope this will trigger the build automatically
	return true
}
func updateLogJSON(title string, summary string, date string, dailyLogIcon string) bool {

	// build the json
	var newEntry models.LogEntry

	// change the date format
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Errorf("Error in parsing the date: %v", err)
		return false
	}
	formattedDate := parsedDate.Format("02 January (Monday) - 2006")

	newEntry.Title = title
	newEntry.Dates = formattedDate
	newEntry.Description = summary
	newEntry.Image = dailyLogIcon
	newEntry.Tags = ""
	newEntry.Links = append(newEntry.Links, models.Link{
		Href: fmt.Sprintf("/diary/%s", date),
	})

	// read it from the source add it in the top
	// write it back to the source
	// Open the file
	// Open the JSON file
	file, err := os.Open(core.Config.PUBLISH.CLONE_DIRECTORY + "public/daily_updates.json")
	if err != nil {
		log.Fatalf("Failed to open JSON file: %v", err)
		return false
	}
	defer file.Close()

	// Read the file's content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
		return false
	}

	// Unmarshal the JSON data into a variable
	var diaryEntries []models.LogEntry
	err = json.Unmarshal(fileContent, &diaryEntries)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
		return false
	}

	// Add the new entry to the top
	diaryEntries = append([]models.LogEntry{newEntry}, diaryEntries...)
	//write it back to the source
	// Marshal the JSON data
	newFileContent, err := json.MarshalIndent(diaryEntries, "", "    ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
		return false
	}
	// writ it back to the source
	err = os.WriteFile(core.Config.PUBLISH.CLONE_DIRECTORY+"public/daily_updates.json", newFileContent, 0644)
	if err != nil {
		log.Fatalf("Failed to write JSON: %v", err)
		return false
	}

	return true
}

func pushChanges(date string) bool {
	// Open the repository
	repo, err := git.PlainOpen(core.Config.PUBLISH.CLONE_DIRECTORY)
	if err != nil {
		log.Errorf("Error in opening the repository: %v", err)
		return false
	}

	// Create a worktree
	worktree, err := repo.Worktree()
	if err != nil {
		log.Errorf("Error in getting the worktree: %v", err)
		return false
	}

	// Add the files
	_, err = worktree.Add(".")
	if err != nil {
		log.Errorf("Error in adding the files: %v", err)
		return false
	}
	commitMsg := fmt.Sprintf("Diary completed for %s by writerBot", date)
	// Commit the changes
	commit, err := worktree.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  core.Config.PUBLISH.GITHUB_COMMIT_USER,
			Email: core.Config.PUBLISH.GITHUB_COMMIT_EMAIL,
			When:  time.Now(),
		},
	})
	if err != nil {
		log.Errorf("Error in committing the changes: %v", err)
		return false
	}
	log.Infof("Committed the changes: %s", commit.String())

	// Push the changes
	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: core.Config.PUBLISH.GITHUB_USERNAME, // Can be anything for token-based auth (e.g., "x-access-token")
			Password: core.Config.PUBLISH.GITHUB_AUTH_TOKEN,
		},
	})
	if err != nil {
		log.Errorf("Error in pushing the changes: %v", err)
		return false
	}
	log.Info("Pushed the changes successfully")
	return true
}

func downloadAssets(assets []models.ASSET, date string) (bool, string) {
	dailyLogIcon := ""
	for _, asset := range assets {

		// Determine the file extension and set the file path
		// Create the directory if it doesn't exist
		dirPath := fmt.Sprintf("%spublic/images/%s/", core.Config.PUBLISH.CLONE_DIRECTORY, date)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			log.Errorf("Error in creating directory %s: %v", dirPath, err)
			return false, dailyLogIcon
		}

		filePath := fmt.Sprintf("%spublic/images/%s/%s%s", core.Config.PUBLISH.CLONE_DIRECTORY, date, asset.Asset, asset.Extension)
		// store first phots link as daily log icon
		// if no photo present keep it empty. so it will pick default image
		if dailyLogIcon == "" {
			dailyLogIcon = fmt.Sprintf("../images/%s/%s%s", date, asset.Asset, asset.Extension)
		}
		// Create the file
		file, err := os.Create(filePath)
		if err != nil {
			log.Errorf("Error in creating file %s: %v", filePath, err)
			return false, dailyLogIcon
		}
		defer file.Close()

		// Write the blob data to the file
		_, err = file.Write(asset.Blob)
		if err != nil {
			log.Errorf("Error in writing to file %s: %v", filePath, err)
			return false, dailyLogIcon
		}
		log.Infof("Successfully saved asset %s", filePath)
	}
	return true, dailyLogIcon
}

func writeMdxFile(content string, date string) bool {
	// Create a MDX file
	log.Debugf("Creating a MDX file")
	file, err := os.Create(core.Config.PUBLISH.CLONE_DIRECTORY + "/content/" + fmt.Sprintf("%s.mdx", date))
	if err != nil {
		log.Errorf("Error in creating a MDX file: %v", err)
		return false
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		log.Errorf("Error in writing to the MDX file: %v", err)
		return false
	}
	_, err = file.WriteString("\n")
	if err != nil {
		log.Errorf("Error in writing to the MDX file: %v", err)
		return false
	}
	return true
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
