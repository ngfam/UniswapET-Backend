# UniswapET-Backend

## Installation

Run SQL server using docker
```
docker-compose up
```

For database models, look for migrations in `internal/pkg/db/migrations` and apply using golang package `github.com/golang-migrate/migrate/v4/cmd/migrate/`.

## Data

The main project uses 1500 largest crytocurrencies (in terms of marketcap) as Tokens and 10000 pairs of trading markets.

Users is self-creating, initialized with ``1,000,000`` tether. 

## Usage

Main functions can be found in ``graph/schema.graphqls`` (the project uses GraphQL anyway, what other document do you even need)

## Third-party

The project uses coingecko api for Token dashboard.