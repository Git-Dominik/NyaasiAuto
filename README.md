# NyaasiAuto

A tiny Go project I made to keep track of anime releases on nyaa.si and automatically send new torrents to qBittorrent.

This project uses [qbit-service](https://github.com/Git-Dominik/qbit-service) as the local qBittorrent bridge.

## What it does

- Reads anime targets from `config.json`
- Checks nyaa.si RSS feeds for new episodes
- Matches the configured uploaders
- Sends new magnet links to the local qBittorrent bridge
- Keeps track of the last watched episode so it doesn't keep re-adding old ones

## Main parts

- `cmd/qbit-service` - local qBittorrent bridge
- `cmd/rss-scraper` - checks for new episodes and submits torrents

## Setup

1. Make sure qBittorrent Web UI is running locally
2. Add a `.env` file in `cmd/qbit-service`:

```env
qbUsr=your_username
qbPass=your_password
```

3. Add a `config.json` in the project root with your anime list:

```json
[
  {
    "title": "Anime Title 1",
    "uploaders": ["UploaderA", "UploaderB"],
    "latest_seen": {
      "season": 1,
      "episode": 12
    }
  },
  {
    "title": "Anime Title 2",
    "uploaders": ["UploaderA", "UploaderB"],
    "latest_seen": {
      "season": 2,
      "episode": 5
    }
  }
]
```

You can add as many anime titles as you want.

## Run

Start the qBittorrent bridge:

```bash
go run cmd/qbit-service
```

Then start the scraper:

```bash
go run cmd/rss-scraper
```

## Note

This is basically a personal automation script for keeping up with anime releases without manually checking nyaa.si every few hours.
