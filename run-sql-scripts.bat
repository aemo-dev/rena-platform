@echo off

REM ============================================
REM Supabase SQL Scripts Runner
REM ============================================
REM Prerequisites:
REM - PostgreSQL client (psql) installed and in PATH
REM - DATABASE_URL in .env file
REM --------------------------------------------

echo Loading environment variables...
for /f "usebackq tokens=*" %%a in (".env") do set "%%a"
echo Done loading environment variables.

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
    echo ERROR: DATABASE_URL not set in .env file
    exit /b 1
)

echo All SQL scripts executed successfully!
