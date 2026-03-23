# 🌿 Botanica API

A herbal remedies and foraging API built in Go — part of the **Earthy Mujer** wellness brand by Cynthia Rico Cook.

**Live API:** https://botanica-production.up.railway.app

---

## About

Botanica is the backend that powers my wellness apps. It serves two datasets I know deeply — herbal remedy recipes from my personal practice, and foraging data from *Nature's Cookbook*, a book I co-authored as part of CEED's educational curriculum.

I'm a certified Integrative Nutrition Health Coach, a forager, and a senior software engineer. This project is where those worlds meet.

---

## Architecture

Botanica uses the same pattern I work with daily at Alaska Airlines:

```
React Frontend → HTTP (BFF layer) → gRPC Server
```

- **gRPC server** — handles all business logic and data
- **HTTP/REST wrapper** — translates HTTP requests into gRPC calls, making the API accessible from any browser or client
- **Proto file** — defines the service contracts between layers

This is the BFF (Backend for Frontend) pattern — the HTTP layer exists purely to serve the frontend's needs while the gRPC layer handles the core logic.

---

## Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | API info and available endpoints |
| GET | `/remedies` | List all remedies |
| GET | `/remedies?type=tea` | Filter by type (tea, bath, tincture) |
| GET | `/remedies?property=sleep` | Filter by healing property |
| GET | `/remedies/{id}` | Get a remedy by ID |
| GET | `/ingredients` | List all unique remedy ingredients |
| GET | `/forageables` | List all plants and mushrooms |
| GET | `/forageables?category=Plant` | Filter by category (Plant, Mushroom) |
| GET | `/forageables?property=immune` | Filter by healing property |
| GET | `/forageables/{id}` | Get a forageable by ID |

---

## Example Requests

```bash
# All tea remedies
curl https://botanica-production.up.railway.app/remedies?type=tea

# Plants good for sleep
curl https://botanica-production.up.railway.app/forageables?property=sleep

# Get a specific remedy
curl https://botanica-production.up.railway.app/remedies/1

# All mushrooms
curl https://botanica-production.up.railway.app/forageables?category=Mushroom
```

---

## Running Locally

```bash
git clone https://github.com/ricocynthia/botanica.git
cd botanica
go mod tidy
go run main.go
```

The API will be available at `http://localhost:8080`

To regenerate proto files after changes:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/remedies.proto
```

---

## Tech Stack

- **Go** — standard library (`net/http`, `encoding/json`)
- **gRPC** — service layer
- **Protocol Buffers** — service contract definitions
- **Railway** — deployment

---

## Data Sources

**Remedies** — personal recipes from my wellness practice:
- Sick Prevention Tea *(from my doula Erica)*
- Tea Para La Tos *(family cough remedy)*
- Chamomile y Canela *(my nightly wind-down blend)*
- Detox Bath *(a winter staple)*

**Forageables** — 10 plants and 10 mushrooms from:

> *Nature's Cookbook — Acknowledging The Gifts of Nature Through Cooking, Gardening & Healing*
> By Cynthia Rico Cook, Nateya Arrieta, and Eva Ryrie Barnett
> Part of CEED's educational curriculum

[Read the book →](https://ceed.org/resource/natures-cookbook/#_)

---

## Related Projects

- [Forage & Heal](https://ricocynthia.github.io/forage-and-heal/) — React field guide app powered by this API
- Earthy Mujer Remedies App — coming soon 🌿

---

## Disclaimer

Not medical advice. Always consult a qualified healthcare professional before using any plant or mushroom medicinally.

---

Built with 🌿 by [Cynthia Rico Cook](https://ricocynthia.github.io) — Earthy Mujer
