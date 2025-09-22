<div align="center" id="top">

# api-watchdog

</div>
<p align="center">
  <img alt="Top language" src="https://img.shields.io/github/languages/top/yourpov/api-watchdog?color=56BEB8">
  <img alt="Language count" src="https://img.shields.io/github/languages/count/yourpov/api-watchdog?color=56BEB8">
  <img alt="Repository size" src="https://img.shields.io/github/repo-size/yourpov/api-watchdog?color=56BEB8">
  <img alt="License" src="https://img.shields.io/github/license/yourpov/api-watchdog?color=56BEB8">
</p>

---

## About

**api-watchdog** is an uptime monitor written in Go
It pings your http endpoints and **alerts when an API goes down** (and when it recovers)

| Feature                | Description                                      |
|------------------------|--------------------------------------------------|
| HTTP checks            | Monitor urls on a set interval                   |
| Latency tracking       | tracks response times in notifications           |
| Status monitoring      | Marks a service as down if it doesn’t return 200 |
| Discord notifications  | Posts clear **DOWN** and **RECOVERED** messages  |

---

## Tech Stack

- [Go](https://go.dev/) (1.22+)

---

## Setup

```bash
# Clone & enter project
git clone https://github.com/yourpov/api-watchdog
cd api-watchdog

# Install deps
go mod tidy

# Build
go build -o api-watchdog .
```

---

## Configuration

```json
{
  "global": {
    "discord_webhook": "",
    "request_timeout_ms": 8000
  },
  "checks": [
    { "name": "API 1", "url": "https://api.com/", "check_interval_seconds": 60 },
    { "name": "Local API", "url": "http://127.0.0.1:8080","check_interval_seconds": 15 }
  ]
}
```

### Notes
- `request_timeout_ms`: how long to wait before marking a request as failed
- Each check needs a `name`, `url`, and `check_interval_seconds`
- All checks default to HTTP GET and expect `200 OK`

---

## Run

```bash
# Run once
./api-watchdog -once

# Run continuously
./api-watchdog -config ./config/config.json
```

---

## Alert Behavior

- Sends **DOWN** when a healthy check becomes unhealthy
- Sends **RECOVERED** when an unhealthy check becomes healthy again
- Each alert includes the endpoint URL and measured latency (if available)

---

## Troubleshooting

- **No alerts sending**  
  Make sure `discord_webhook` is a a valid webhook

- **Always shows “DOWN”**  
  Make sure the API returns `200 OK` (you can adjust the code if you want to allow other statuses)  

- **False alarms from slow APIs**  
  Raise `request_timeout_ms` in `global`

- **Multiple APIs fire at the same time**  
  Checks are staggered, but you can raise `check_interval_seconds` to send them further apart

---

## Showcase
[▶ Watch the showcase](https://www.youtube.com/watch?v=bS96PllP1EcS)
