# 🏰 Royaka
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)](https://golang.org)
[![Node.js](https://img.shields.io/badge/Node.js-18+-339933?logo=node.js&logoColor=white)](https://nodejs.org)
[![React](https://img.shields.io/badge/React-18+-61DAFB?logo=react&logoColor=black)](https://reactjs.org)

**Royaka** is a turn-based multiplayer tower battle game inspired by *Clash Royale*. Built with Go and React, it features strategic troop deployment, dynamic movement logic, and WebSocket-powered real-time gameplay.

## Academic and Copyright Disclaimer

This repository is a non-commercial student project created solely for academic, educational, and portfolio purposes. It is an unofficial fan-made project and is not affiliated with, endorsed by, sponsored by, or associated with Supercell or the creators and rights holders of *Clash Royale*.

All third-party names, trademarks, characters, artwork, sounds, and other copyrighted materials remain the property of their respective owners. No ownership of those materials is claimed. If you are a rights holder and believe any material in this repository should be removed or credited differently, please open a GitHub issue so it can be reviewed and addressed promptly.

## Features

* **1v1 Multiplayer Matches** via WebSocket
* **Turn-Based Gameplay** with attack, heal, and skip mechanics
* **Two Game Modes**

  * **Simple Mode**: Basic strategic combat
  * **Enhanced Mode**: Adds MANA, EXP, leveling, and critical hits
* **Troop Collection** with tanks, healers, and damage dealers
* **Smart Troop Behavior** (e.g., river crossing only at bridges)
* **User Authentication** (registration, login, and persistent stats)
* **Username Authentication backed by Supabase Postgres**

## Tech Stack

* **Backend**: Go, Gorilla WebSocket, native HTTP server
* **Frontend**: React, Vite, TailwindCSS
* **Storage**: Supabase Postgres (profiles, sessions, troops, and towers)

## Project Structure

```
royaka/
├── client/                 # React frontend (Vite)
│   ├── context/            # Zustand or React Context providers
│   ├── pages/              # Main route pages (Login, Game, etc.)
│   ├── routes/             # Route definitions and utilities
│   ├── App.jsx             # Main application layout
│   ├── main.jsx            # Entry point for React
│   └── index.html          # Vite HTML entry
├── server/                 # Go backend
│   ├── assets/             # JSON data files
│   ├── internal/
│   │ ├── game/             # Game mechanics & logic
│   │ ├── model/            # Data models (players, troops, etc.)
│   │ ├── network/          # WebSocket & HTTP handlers
│   │ └── utils/            # Utility functions
│   └── main.go             # Backend entry point
```

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/trmzaiu/royaka-clash-royale.git
cd royaka-clash-royale
```

### 2. Start Supabase

Install Docker and the [Supabase CLI](https://supabase.com/docs/guides/local-development/cli/getting-started), then run:

```bash
supabase start
supabase db reset
```

The reset applies `supabase/migrations`, then loads development game definitions from `supabase/seed.sql`.

### 3. Configure and Start the Go Backend

```bash
cd server
cp .env.example .env
# Export the values from .env in your shell.
go mod tidy
go run main.go
```

* WebSocket endpoint: `ws://localhost:8080/ws`
* Default port: `8080`

The server requires `DATABASE_URL` and an RFC3339 `LEGACY_SESSION_CUTOFF`. For hosted App Engine deployments, use Supavisor's IPv4 session-mode connection string and provide secrets through the deployment environment; never commit real values.

### 4. Start the React Frontend

```bash
cd client
npm install
npm run dev
```

* Runs on `http://localhost:5173`

## Gameplay Overview

### 1. Turn-Based Mode
- Players take turns (Player 1 → Player 2).
- Each turn, the player gains 3 mana.
- If a player does not have enough mana to attack, they can choose to skip their turn.
- Each player has 30 seconds per turn; if no action is taken, the turn automatically passes to the opponent.
- Victory requires destroying both Guard Towers before accessing the King Tower.

### 2. Timed Match Mode
- Mana increases automatically over time (1 mana every 2 seconds).
- Both players act simultaneously in real-time.
- Towers actively defend by attacking enemy troops within range.
- Matches last 3 minutes with fast-paced, continuous action.
- Victory conditions remain the same: eliminate both Guard Towers before accessing the King Tower.

## Authentication System

* Username/password requests remain on the WebSocket protocol.
* Usernames are unique case-insensitively and may contain arbitrary display characters.
* Passwords are bcrypt hashes in `profiles`; no email address is created or required.
* Random 30-day session IDs are stored in the `sessions` table.
* Every protected game message is bound to the authenticated WebSocket identity.
* Imported JSON sessions remain valid only until `LEGACY_SESSION_CUTOFF`.
* Match stats and game definitions are read from Supabase Postgres.

## Import Existing JSON Data

Create and link the hosted project and apply migrations. Set `LEGACY_SESSION_CUTOFF` to no more than seven days after the import, then run from `server/`:

```bash
go run ./cmd/import-json -data-dir ./assets/data
```

The command requires only `DATABASE_URL` and `LEGACY_SESSION_CUTOFF`. It is idempotent and imports bcrypt hashes without printing credentials. Validate staging first, then run it once against production. Do not use `supabase db push --include-seed` for production; apply schema with `supabase db push` and use the importer for the real game definitions.

After the seven-day rollback period, remove the legacy JSON user/session data and revoke the importer secret from the migration environment.

---

Made with ❤️ by game dev enthusiasts.
