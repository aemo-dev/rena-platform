@echo off
REM Run SQL Scripts for Supabase
REM Requires: supabase CLI installed and configured

echo Running SQL scripts...

cd rena\sql-scripts

for %%f in (*.sql) do (
    echo Running %%f...
    supabase db execute --file %%f
    if errorlevel 1 (
        echo Error running %%f
        exit /b 1
    )
)

echo All SQL scripts executed successfully!
