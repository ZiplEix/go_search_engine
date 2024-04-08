# Go Search Engine

## Features

- [ ] Web crawler
- [ ] Custom Full Text Search Engine
- [ ] Cronjobs
- [ ] Logged in section to control the crawler
- [ ] JSON API
- [ ] Rate limiting
- [ ] Docker ??

## Tech Stack

- GO
  - Fiber (not needed but i'm confortable with it)
  - GORM
  - Go HTML ot Templ
- Postgres
- Custom JWT Auth
- Docker
- Tailwind & daisyUI 😒

## Database Start

```bash
docker build -t go_search_engine_database .
```

```bash
docker run -d --name go_search_engine_database -p 5432:5432 go_search_engine_database
```
