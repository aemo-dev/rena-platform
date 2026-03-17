@echo off

REM ============================================
REM Supabase SQL Scripts Runner
REM ============================================
REM Prerequisites:
REM - Supabase CLI installed
REM - Run from project root directory
REM --------------------------------------------

echo Copying SQL scripts to supabase/migrations...

if not exist supabase/migrations mkdir supabase/migrations

REM Get current timestamp
for /f "tokens=2 delims==" %%a in ('wmic OS Get localdatetime /value') do set "dt=%%a"
set timestamp=%dt:~0,8%000000

REM Copy and rename files with timestamp
copy /Y rena\sql-scripts\keystores.sql supabase\migrations\%timestamp%_keystores.sql
copy /Y rena\sql-scripts\projects.sql supabase\migrations\%timestamp%_projects.sql

echo Pushing migrations to Supabase...
supabase db push

echo Done!
