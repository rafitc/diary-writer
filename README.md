# Diary Writer

Log your daily activities Quickly.

## Whay i built this

I've been looking for a way to write diary in quick way and publish it on daily. Couldn't find any perfect solution for me. I thought of making one based on my requirement

## WorkFlow

Just text your daily activities to your telegram bot, it fetches all data together and at the end of the day it collect all the data for the day and push into the [Application](https://github.com/rafitc/diary). The Writer application modify the content to fix grammar errors, it generate title and summary as well using LLM (Groq AI)

Backend system for my [Diary](https://github.com/rafitc/diary)

## Tech stacks

- Go lang
- Telegram [Bot API](https://pkg.go.dev/github.com/go-telegram/bot)
- Sqlite3 DB
- [Groq AI](https://groq.com)

## Deployment Guide

1. Clone the repo
2. Get API key from Groq AI
3. Build the docker image
4. Create config.yaml by filling your credentials
5. Create Docker service
   1. Mount sqlite3 db as volume
   2. Create docker config using config/config.yaml
6. Deploy application as service

---

<a href="https://groq.com" target="_blank" rel="noopener noreferrer">
  <img
    src="https://groq.com/wp-content/uploads/2024/03/PBG-mark1-color.svg"
    alt="Powered by Groq for fast inference."
  />
</a>
