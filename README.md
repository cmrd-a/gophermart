# gophermart 42

## Что сделано сверх [ТЗ](SPECIFICATION.md)

- Полное описание api в swagger
- Миграции для БД
- golangci-lint + Github Actions

## Руководство

- `cp .env.example .env`
- `docker compose up -d`
- `./cmd/accrual/accrual_darwin_arm64 -a :9090`
- `make run`
- http://localhost:8080/swagger/index.html
- `docker compose down`

---

- `git remote add -m master template https://github.com/yandex-praktikum/go-musthave-diploma-tpl.git`
- `git fetch template && git checkout template/master .github`



