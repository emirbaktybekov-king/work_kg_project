# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

WorkKG is a full-stack job portal platform for Kyrgyzstan with a Telegram bot integration and Admin CRM dashboard.

- **Backend:** Go 1.21 with Gorilla Mux, PostgreSQL, Telegram Bot API
- **Frontend:** Next.js 16 with React 19, TypeScript, Tailwind CSS 4, Radix UI

## Project Structure

```
work_kg_project/
├── work_kg_backend/        # Go backend server
│   ├── cmd/seed/           # Database seed scripts
│   ├── internal/
│   │   ├── bot/            # Telegram bot logic & handlers
│   │   ├── config/         # Environment configuration
│   │   ├── database/       # PostgreSQL operations & schema
│   │   ├── handlers/       # HTTP request handlers
│   │   └── models/         # Data models & types
│   └── main.go             # Entry point (starts HTTP server + bot)
│
└── work_kg_frontend/       # Next.js admin CRM
    ├── src/
    │   ├── app/            # Next.js app router pages
    │   │   ├── auth/login/ # Admin login
    │   │   └── crm/        # Main CRM dashboard
    │   ├── components/ui/  # Radix UI components
    │   └── lib/
    │       ├── api.ts      # API client & TypeScript interfaces
    │       └── utils.ts    # Utilities
    └── package.json
```

## Development Commands

### Frontend (work_kg_frontend/)

```bash
npm run dev      # Start dev server on port 3000
npm run build    # Production build
npm run lint     # ESLint
```

### Backend (work_kg_backend/)

```bash
go run main.go           # Run development server
go build -o server .     # Build binary
go mod tidy              # Update dependencies
```

## Environment Variables

Backend `.env`:
```
TELEGRAM_TOKEN=<bot_token>
DATABASE_URL=postgresql://user:password@host:5432/database
SERVER_PORT=8080
```

Frontend `.env`:
```
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

## Architecture Notes

### Backend

- HTTP server (Gorilla Mux) and Telegram bot run concurrently via goroutine in main.go
- Auth middleware at `internal/handlers/auth.go` uses simple email-as-token pattern
- Bot state management uses global `userStates` map in `internal/bot/bot.go`
- Database schema defined in `internal/database/db.go` with auto-migration on startup
- Default admin: `admin@workkg.com` / `admin123`

### Frontend

- Single large CRM page at `src/app/crm/page.tsx` (~1160 lines) handles dashboard, jobs, resumes, users tabs
- API client in `src/lib/api.ts` defines all interfaces and HTTP methods
- Token stored in localStorage, passed via Authorization header

### API Endpoints

- `POST /api/auth/login` - Admin login
- `GET /api/auth/me` - Current user
- `GET|POST /api/jobs`, `PUT|DELETE /api/jobs/{id}` - Job CRUD
- `GET /api/users` - List Telegram users
- `GET /api/resumes` - List resumes
- `GET /api/stats` - Dashboard statistics

### Database Tables

- `users` - Telegram bot users (telegram_id, city, specialty, experience)
- `admin_users` - CRM admins (email, bcrypt password)
- `jobs` - Job listings (title, category, subcategory, city, salary)
- `resumes` - User applications

### Telegram Bot

- Entry point: `internal/bot/bot.go`
- Handles user registration flow with state-based forms
- Main menu: profile, search jobs, search employees, entertainment, subscriptions
- Cities: Bishkek, Osh, Talas, Naryn, Karakol, Jalal-Abad, Cholpon-Ata
- Job categories: Construction, Food Service, Sewing, IT, Sales, Transport
