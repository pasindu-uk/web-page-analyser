# Frontend Setup Guide

Step-by-step instructions for setting up and running the React frontend locally.

## Prerequisites

- **Node.js 20 or later** — [Download](https://nodejs.org/) or use [nvm](https://github.com/nvm-sh/nvm)
- **Backend running** — the frontend proxies API calls to `http://localhost:8080`

Verify Node.js is installed:

```bash
node --version
# v20.20.0 or later

npm --version
# 10.x
```

If you use nvm:

```bash
nvm install 20
nvm use 20
```

## Step 1: Clone the repository

```bash
git clone https://github.com/pasindu-uk/web-page-analyser.git
cd web-page-analyser/frontend
```

## Step 2: Install dependencies

```bash
npm install
```

This installs:
- **React 19** + **ReactDOM** — UI library
- **Vite 8** — build tool and dev server
- **TypeScript 5.9** — type checking
- **ESLint** — code linting

## Step 3: Start the backend

The frontend needs the backend API running. In a separate terminal:

```bash
cd ../backend
go run ./cmd/api
```

See the [backend setup guide](../backend/README.md) for full instructions.

## Step 4: Start the dev server

```bash
npm run dev
```

You should see:

```
  VITE v8.0.0  ready in 300 ms

  ➜  Local:   http://localhost:5173/
  ➜  Network: use --host to expose
```

Open `http://localhost:5173` in your browser.

### How the proxy works

In development, Vite proxies all `/api` requests to the Go backend. This is configured in `vite.config.ts`:

```ts
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
  },
},
```

This means:
- `http://localhost:5173/api/analyze` → `http://localhost:8080/api/analyze`
- No CORS issues during development

## Step 5: Using the app

1. Enter a full URL (including `http://` or `https://`) in the input field
2. Click **Analyze**
3. View the results: HTML version, title, heading counts, link breakdown, and login form detection
4. If MySQL is configured on the backend, past analyses appear in the **History** section
5. Click **Clear Cache** next to the History heading to flush the in-memory cache and refresh the list

## Step 6: Lint the code

```bash
npm run lint
```

## Step 7: Build for production

```bash
npm run build
```

The output goes to `frontend/dist/`. This is a static build that can be served by any web server (Nginx, Caddy, etc.) or from the Go backend.

## Project Structure

```
frontend/
├── src/
│   ├── api/
│   │   └── analyzeApi.ts         # API client (POST /api/analyze, GET /api/analyses, DELETE /api/cache)
│   ├── components/
│   │   ├── AnalyzeForm.tsx       # URL input form with validation
│   │   ├── AnalyzeForm.css
│   │   ├── AnalysisResult.tsx    # Full result display card
│   │   ├── AnalysisResult.css
│   │   ├── AnalysisHistory.tsx   # List of past analyses
│   │   ├── AnalysisHistory.css
│   │   ├── ErrorMessage.tsx      # Error display
│   │   ├── ErrorMessage.css
│   │   ├── HeadingSummary.tsx    # H1-H6 count grid
│   │   ├── HeadingSummary.css
│   │   ├── LinkSummary.tsx       # Internal/external/inaccessible link counts
│   │   └── LinkSummary.css
│   ├── types/
│   │   └── analysis.ts          # TypeScript types matching backend models
│   ├── App.tsx                   # Main app component with state management
│   ├── App.css
│   ├── index.css                 # Global styles and reset
│   └── main.tsx                  # React entry point
├── index.html
├── vite.config.ts
├── tsconfig.json
├── tsconfig.app.json
└── package.json
```

## Component Overview

| Component | Purpose |
|---|---|
| `App` | Root component — manages loading, result, error, and history state |
| `AnalyzeForm` | URL input with client-side validation and loading state |
| `AnalysisResult` | Displays all result fields, composes HeadingSummary and LinkSummary |
| `HeadingSummary` | Renders H1–H6 counts in a grid |
| `LinkSummary` | Renders internal/external/inaccessible link counts |
| `AnalysisHistory` | Lists past analyses as clickable items |
| `ErrorMessage` | Displays API errors with status code |

## Troubleshooting

### White screen / blank page

Check the browser console for errors. Common causes:
- Backend not running — start it with `go run ./cmd/api`
- TypeScript import errors — run `npm run build` to see compilation errors

### "Failed to fetch" errors

Make sure the backend is running on port 8080. The Vite dev server proxies `/api` requests there.

### Node version too old

```
Error: Vite requires Node.js version 20.19.0 or higher
```

Upgrade Node.js:
```bash
nvm install 20
nvm use 20
```
