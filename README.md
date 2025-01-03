# Diary Writer

Log your daily activities quickly and effortlessly.

## Why I Built This

I wanted a simple and efficient way to maintain a daily diary and publish entries effortlessly. After exploring various tools, I couldn’t find a solution that perfectly suited my needs. So, I decided to build one tailored to my requirements.

## Workflow

1. **Log Activities:** Send your daily activities as messages to a Telegram bot.
2. **Data Collection:** The bot collects and organizes all your messages for the day.
3. **Content Processing:** At the end of the day, the backend system processes the collected data to:
   - Fix grammar errors.
   - Generate a title and summary using a Large Language Model (LLM) powered by [Groq AI](https://groq.com).
4. **Publishing:** The processed content is pushed to the [Diary Application](https://github.com/rafitc/diary) for publishing.

## Tech Stack

- **Language:** Go
- **Bot Integration:** Telegram [Bot API](https://pkg.go.dev/github.com/go-telegram/bot)
- **Database:** SQLite3
- **AI Assistance:** [Groq AI](https://groq.com)
- **Publish the diary** [Go-Git](https://pkg.go.dev/github.com/go-git/go-git/v5)

## Deployment Guide

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/rafitc/diary
   ```
2. **Obtain API Key:** Get an API key from [Groq AI](https://groq.com).
3. **Obtain Telegram Token** Create your own telegram bot [using bot father](https://core.telegram.org/bots/features#creating-a-new-bot) and get the API key
4. **Build the Docker Image:**
   ```bash
   docker build -t diary-writer .
   ```
5. **Configure Application:**
   - Create a `config.yaml` file with your credentials and settings.
6. **DB setup**
   - Install [sqlite3](https://www.sqlite.org/download.html) db in your system
   - Create a folder to create DB. Eg:- `mkdir -p /diary-writer/db`
   - Create DB and tables. Run [this Query](scripts/sqlite-db.sql)
7. **Config setup**
   - Create config.yaml file and update it with your Git credentials, API keys, folder path etc
   - Create a docker config using config.yaml. Eg:- `docker config create diary-writer-config config.yaml`
8. **Set Up Docker Service:**

- Mount the SQLite3 database as a volume.
- Use the `diary-writer-config` config to mount as a config

9.  **Deploy the Service:** Run the application as a Docker service.

- Docker service script
  ```bash
     docker service create \
     --name my-diary-writer \
     --mount type=bind,source=/diary-writer/db,target=/root/sqlite-db \
     --config source=diary-writer-config,target=/root/config/config.yaml \
     diary-writer:latest
  ```

10. **Verify the deployment** You can verify and check the status using `docker logs -f <CONTAINER-ID>`


## Powered by

<a href="https://groq.com" target="_blank" rel="noopener noreferrer">
    <img
        src="https://groq.com/wp-content/uploads/2024/03/PBG-mark1-color.svg"
        alt="Powered by Groq for fast inference."
        width="100"
    />
</a>
