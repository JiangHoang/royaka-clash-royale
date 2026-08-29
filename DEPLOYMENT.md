# Deploy Royaka to Vercel and Render

The frontend is deployed as a Vite SPA on Vercel. The Go WebSocket backend runs
as a single Render service and stores persistent data in Supabase Postgres.

## 1. Create and migrate Supabase

1. Create a hosted Supabase project in a region close to the Render service.
2. Copy the IPv4 Supavisor session-mode connection string and append the
   required SSL parameters. Keep this value as `DATABASE_URL`; never commit it.
3. Link the local Supabase CLI to the hosted project, then apply migrations:

   ```bash
   supabase link --project-ref PROJECT_REF
   supabase db push
   ```

4. Choose an RFC3339 UTC cutoff no more than seven days in the future. Use the
   exact same value for the import and the Render environment:

   ```powershell
   $env:DATABASE_URL = "postgresql://..."
   $env:LEGACY_SESSION_CUTOFF = "2026-09-05T00:00:00Z"
   Set-Location server
   go run ./cmd/import-json -data-dir ./assets/data
   ```

The importer is idempotent and migrates the users, sessions, troops, and towers
from `server/assets/data`. Do not use the local development seed for production.

## 2. Create the Vercel frontend

Import this repository as a Vercel project and set:

- Root Directory: `client`
- Framework Preset: Vite
- Build Command: `npm run build`
- Output Directory: `dist`

Deploy once to reserve the production `https://PROJECT.vercel.app` hostname.
The first deployment cannot connect to the backend until `VITE_WS_URL` is set.

## 3. Create the Render backend

Create a Render Blueprint from the repository's `render.yaml`. During the
initial Blueprint setup, provide:

- `DATABASE_URL`: the Supabase connection string from step 1
- `LEGACY_SESSION_CUTOFF`: the same timestamp used by the importer
- `ALLOWED_ORIGINS`: the exact Vercel origin, for example
  `https://PROJECT.vercel.app`

Separate multiple allowed origins with commas. Add an exact preview origin only
when that preview should be able to use the production game server. Do not use a
wildcard because browsers send an Origin header during the WebSocket upgrade.

Wait for `https://SERVICE.onrender.com/healthz` to return HTTP 200.

## 4. Connect and redeploy Vercel

Add this Vercel environment variable to Production and any enabled Preview
environment, then redeploy:

```text
VITE_WS_URL=wss://SERVICE.onrender.com/ws
```

Vite embeds this value at build time, so changing it does not affect an existing
deployment until a new frontend build is created.

## 5. Smoke test

1. Open and refresh `/auth`, `/lobby`, `/game-simple`, `/game-enhanced`, and
   `/card-desk`; none should return a Vercel 404.
2. Log in with an imported account, refresh, and verify session restoration.
3. Register and log in with a new account.
4. Open two browsers with different users and complete matchmaking in both game
   modes.
5. Restart the Render service and verify that the client reconnects. The current
   in-memory match and room state is intentionally lost on restart.

Keep Render at one instance. Matchmaking, active rooms, connections, and game
state are process-local and are not safe for horizontal scaling yet. The free
Render plan can also cold-start after an idle period, so it is suitable for a
demo rather than an always-on production service.
