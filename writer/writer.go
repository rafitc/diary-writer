package writer

import (
	"fmt"
	"io"
	"io/ioutil"
	"main/core"
	"main/db"
	"main/editor"
	"main/logger"
	"main/models"
	"main/publisher"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/robfig/cron"
)

var log = logger.Logger

func StartCronJob() {
	log.Debug("Scheduling Cron Job")
	// Starting a cronJob to check is there anything in db
	c := cron.New()
	c.AddFunc(core.Config.PUBLISH.PUBLISH_JOB_CRON, func() {
		dailyDiaryDataChecker()
	})
	log.Info("Configured cron job to check for data in db")
	c.Start()
}

func InsertDataIntoDB(content string, asset string, extension string, download_link string, creation_date time.Time) {
	log.Debug("Inserting data into db")

	db, err := db.NewDatabase(core.Config.DATABASE.NAME)
	if err != nil {
		log.Errorf("Error in connecting to db %v", err)
		return
	}
	defer db.Close()

	if len(download_link) > 0 {
		log.Infof("Downloading file from %s", download_link)
		resp, err := http.Get(download_link)
		if err != nil {
			log.Errorf("Error downloading file: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Errorf("Error downloading file: received status code %d", resp.StatusCode)
			return
		}

		tempFile, err := os.CreateTemp("", "downloaded-*"+extension)
		if err != nil {
			log.Errorf("Error creating temp file: %v", err)
			return
		}
		defer tempFile.Close()

		_, err = io.Copy(tempFile, resp.Body)
		if err != nil {
			log.Errorf("Error saving file to temp location: %v", err)
			return
		}

		tempFilePath := tempFile.Name()
		log.Infof("File downloaded to %s", tempFilePath)

		fileData, err := ioutil.ReadFile(tempFilePath)
		if err != nil {
			log.Errorf("Error reading temp file: %v", err)
			return
		}

		query := "INSERT INTO daily_updates (content, asset, asset_extension, asset_download_link, creation_date, is_updated, asset_blob) VALUES (?, ?, ?, ?, ?, ?, ?)"
		_, err = db.Insert(query, content, asset, extension, download_link, creation_date, false, fileData)
		if err != nil {
			log.Errorf("Error in inserting data into db %v", err)
			return
		}
	} else {
		query := "INSERT INTO daily_updates (content, asset, asset_extension, asset_download_link, creation_date, is_updated) VALUES (?, ?, ?, ?, ?, ?)"
		_, err = db.Insert(query, content, asset, extension, download_link, creation_date, false)
		if err != nil {
			log.Errorf("Error in inserting data into db %v", err)
			return
		}
	}

	log.Info("Data inserted into sqlite-db")
}

func dailyDiaryDataChecker() {
	log.Info("Starting daily diary data checker")
	// Check for data in db
	db, err := db.NewDatabase(core.Config.DATABASE.NAME)
	if err != nil {
		log.Errorf("Error in connecting to db %v", err)
		return
	}
	defer db.Close()
	// If data is present, Then compare it with current time, if the data is old then push into github and trigger build
	// get all previous days data from db where is_updated false
	query := `SELECT id, content, asset, asset_extension, creation_date, asset_blob FROM daily_updates 
			WHERE is_updated = false AND creation_date < DATE('now', 'localtime')
			order by id` // get all data from db with is_updated false
	rows, err := db.Fetch(query)
	if err != nil {
		log.Errorf("Error in fetching data from db %v", err)
		return
	}
	defer rows.Close()

	entries := make(map[string]*models.DiaryEntry)

	for rows.Next() {
		var id int
		var content string
		var asset string
		var asset_extension string
		var createdAt time.Time
		var fileblob []byte

		err = rows.Scan(&id, &content, &asset, &asset_extension, &createdAt, &fileblob)
		if err != nil {
			log.Errorf("Error in scanning data from db %v", err)
			return
		}
		date := createdAt.Format("2006-01-02")
		if entry, exists := entries[date]; exists {
			entry.Ids = append(entry.Ids, id)
			entry.Content += "" + content + "\n"
			if len(asset) > 0 {
				entry.Content += fmt.Sprintf("\n![Alt Text](../images/%s/%s%s)\n\n", date, asset, asset_extension)
				assetEntry := models.ASSET{
					Asset:     asset,
					Extension: asset_extension,
					Blob:      fileblob,
				}
				entry.Asset = append(entry.Asset, assetEntry)
			}
		} else {
			assetEntries := []models.ASSET{}
			if len(asset) > 0 {
				entry.Content += fmt.Sprintf("\n![Alt Text](../images/%s/%s%s)\n\n", date, asset, asset_extension)
				assetEntries = append(assetEntries, models.ASSET{
					Asset:     asset,
					Extension: asset_extension,
					Blob:      fileblob,
				})
			}
			entries[date] = &models.DiaryEntry{
				Content: content,
				Date:    date,
				Asset:   assetEntries,
				Ids:     []int{id},
			}
		}
	}

	// If data is not present, then do nothing
	if len(entries) == 0 {
		log.Debug("No data in db")
		return
	} else {
		log.Debug("Data present in db")
		for _, entry := range entries {

			// the expecting content is fully qualified mdx format
			finalContent, titleForJsonLog, summaryForJsonLog := editor.EditContent(entry.Content, entry.Date)

			// Call the publisher to push the data into github
			status := publisher.PublishContent(finalContent, entry.Date, entry.Asset, titleForJsonLog, summaryForJsonLog)
			if !status {
				log.Errorf("Error in Publishing data to github")
				return
			}

			// Update the db with is_updated true
			// Later will add an script to delete the is_updated true data
			// Build the query with placeholders
			placeholders := make([]string, len(entry.Ids))
			args := make([]interface{}, len(entry.Ids))

			for i, id := range entry.Ids {
				placeholders[i] = "?"
				args[i] = id
			}

			query := fmt.Sprintf("UPDATE daily_updates SET is_updated = true WHERE id IN (%s)", strings.Join(placeholders, ","))

			_, err = db.Update(query, args...)
			if err != nil {
				log.Errorf("Error in updating data in db %v", err)
				return
			}
		}
	}

	// If data is present and not old, then do nothing
	// If data is present and old, then push into github and trigger build
	log.Debug("Completed daily diary data checker")

}
