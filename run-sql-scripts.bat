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

copy /Y rena\sql-scripts\*.sql supabase\migrations\

echo Pushing migrations to Supabase...
supabase db push

echo Done!
