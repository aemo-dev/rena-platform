@echo off

REM ============================================
REM Supabase SQL Scripts Runner
REM ============================================
REM Setup Instructions:
REM Option A - Using psql (recommended):
REM 1. Install PostgreSQL client (psql)
REM 2. Set DATABASE_URL in .env
REM
REM Option B - Using Supabase CLI:
REM 1. Download Supabase CLI from GitHub:
REM    https://github.com/supabase/cli/releases
REM 2. Extract it and add the executable to PATH
REM 3. Login: supabase login
REM 4. Link: supabase link --project-ref YOUR_PROJECT_ID
REM --------------------------------------------

echo Running SQL scripts...

cd rena\sql-scripts

if defined DATABASE_URL (
    echo Using DATABASE_URL from environment...
    for %%f in (*.sql) do (
        echo Running %%f...
        psql "%DATABASE_URL%" -f %%f
        if errorlevel 1 (
            echo Error running %%f
            exit /b 1
        )
    )
) else (
    echo WARNING: DATABASE_URL not set in environment
    echo Using Supabase CLI...
    for %%f in (*.sql) do (
        echo Running %%f...
        type %%f | supabase db query
        if errorlevel 1 (
            echo Error running %%f
            exit /b 1
        )
    )
)

echo All SQL scripts executed successfully!
