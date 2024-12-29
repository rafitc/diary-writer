package writer

import (
	"main/core"
	"main/db"
	"main/logger"
	"time"

	"github.com/robfig/cron"
)

type DiaryEntry struct {
	Content string `json:"content"`
}

var log = logger.Logger

func StartCronJob() {
	log.Infof("starting %v job")
	// Starting a cronJob to check is there anything in db
	c := cron.New()
	c.AddFunc("0/5 * * * * *", func() { dailyDiaryDataChecker() })
	log.Infof("Started %s cronJob")
	c.Start()

}

func InsertDataIntoDB(content string, asset string) {
	log.Debugf("Inserting data into db")
	db, err := db.NewDatabase(core.Config.DATABASE.NAME)
	if err != nil {
		log.Errorf("Error in connecting to db %v", err)
		return
	}
	defer db.Close()
	query := "INSERT INTO daily_updates (content, asset, creation_date, is_updated) VALUES (?, ?, ?, ?)"
	_, err = db.Insert(query, content, asset, time.Now(), false)
	if err != nil {
		log.Errorf("Error in inserting data into db %v", err)
		return
	}
	log.Debugf("Data inserted into db")
}

func dailyDiaryDataChecker() {
	log.Debugf("Starting daily diary")
	// Check for data in db
	db, err := db.NewDatabase(core.Config.DATABASE.NAME)
	if err != nil {
		log.Errorf("Error in connecting to db %v", err)
		return
	}
	defer db.Close()
	// If data is present, Then compare it with current time, if the data is old then push into github and trigger build
	// get all data from db with is_updated false
	query := `SELECT id, content, asset, creation_date, is_updated FROM daily_updates 
			WHERE is_updated = false
			order by id` // get all data from db with is_updated false
	rows, err := db.Fetch(query)
	if err != nil {
		log.Errorf("Error in fetching data from db %v", err)
		return
	}
	defer rows.Close()

	entries := make(map[string]*DiaryEntry)

	for rows.Next() {
		var id int
		var content string
		var asset string
		var createdAt time.Time
		var isUpdated bool
		err = rows.Scan(&id, &content, &asset, &createdAt, &isUpdated)
		if err != nil {
			log.Errorf("Error in scanning data from db %v", err)
			return
		}
		date := createdAt.Format("2006-01-02")
		if entry, exists := entries[date]; exists {
			entry.Content += " " + content + " \n"
			if len(asset) > 0 {
				entry.Content += "\n (asset)" + asset + " \n"
			}
		} else {
			entries[date] = &DiaryEntry{
				Content: content,
			}
		}
	}

	// If data is not present, then do nothing
	if len(entries) == 0 {
		log.Debugf("No data in db")
		return
	} else {
		log.Debugf("Data present in db")
		for _, entry := range entries {
			log.Debugf("Content:\n %s", entry.Content)
		}
	}

	// If data is present and not old, then do nothing
	// If data is present and old, then push into github and trigger build
	log.Debugf("Completed daily diary data cheker")

}
