# Rena Builder Platform

Rena Builder is a professional-grade visual development platform for creating modern, scalable applications without writing code.

## Project Structure

- `rena/frontend`: Vue 3 + Vite + TypeScript application for the landing page and builder dashboard.
- `rena/backend`: Go (Gin) backend for project management and builder services.
- `rena/languages`: Multi-language support (English and Arabic).
- `rena/lib`: Support libraries and binaries (Android SDK tools, etc.).

## Prerequisites

- Node.js (v18+)
- Go (v1.20+)
- Android SDK (for app generation)

## Getting Started

1.  Clone the repository.
2.  Configure the `.env` file in the root directory.
3.  **Run SQL Scripts** (first time setup only):
    
    **Option A - Manual via Supabase Dashboard (Recommended):**
    - Go to https://app.supabase.com and select your project
    - Open SQL Editor
    - Copy and paste content from `rena/sql-scripts/keystores.sql` and run it
    - Copy and paste content from `rena/sql-scripts/projects.sql` and run it
    - Verify tables are created in Table Editor
    
    **Option B - Using PowerShell script (requires supabase CLI or psql):**
    1. Add your database connection to `.env`:
    ```bash
    SUPABASE_DB_URL=postgresql://postgres:<password>@<db-host>:5432/postgres
    ```
    2. Then run:
    ```powershell
    .\run-sql-scripts.ps1
    ```
    
    Or use `run.bat` and select 'y' when prompted.
    
    **Note:** `run.bat` uses `psql` to execute SQL files. Install PostgreSQL client and ensure `psql` is on PATH.
    
4.  Run the application:
    ```bash
    ./run.bat
    ```

## Features

- Modern landing page with GSAP animations.
- Supabase-powered authentication (Google/GitHub).
- Project management dashboard.
- Automatic RTL support for Arabic.
- User-specific Android keystore generation.
