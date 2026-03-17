@echo off
setlocal enabledelayedexpansion

REM ============================================
REM Supabase SQL Scripts Runner (Migrations Mode)
REM ============================================

set "SOURCE_DIR=rena\sql-scripts"
set "MIGRATIONS_DIR=supabase\migrations"

REM Set to 1 to clear old migrations (WARNING: destructive)
set "CLEAR_MIGRATIONS=1"

REM Validate source directory
if not exist "%SOURCE_DIR%" (
    echo ERROR: Source directory "%SOURCE_DIR%" not found.
    exit /b 1
)

REM Clear migrations if requested
if "%CLEAR_MIGRATIONS%"=="1" (
    echo Clearing existing migrations...
    if exist "%MIGRATIONS_DIR%" (
        rmdir /s /q "%MIGRATIONS_DIR%"
    )
)

REM Ensure migrations directory exists
if not exist "%MIGRATIONS_DIR%" (
    mkdir "%MIGRATIONS_DIR%"
)

echo Preparing SQL scripts...

REM Get timestamp
for /f %%i in ('powershell -command "Get-Date -Format yyyyMMddHHmmss"') do set "timestamp=%%i"
set count=0

REM Loop through SQL files
for %%f in ("%SOURCE_DIR%\*.sql") do (
    set /a count+=1
    set "filename=%%~nf"
    set "newname=!timestamp!!count!_!filename!.sql"

    echo [COPY] %%~nxf -> !MIGRATIONS_DIR!\!newname!
    copy /Y "%%f" "!MIGRATIONS_DIR!\!newname!" >nul

    if errorlevel 1 (
        echo ERROR: Failed to copy %%~nxf
        exit /b 1
    )
)

if %count%==0 (
    echo WARNING: No SQL files found in "%SOURCE_DIR%"
    exit /b 0
)

echo Applying migrations to Supabase...
supabase db push

if errorlevel 1 (
    echo ERROR: Migration failed!
    exit /b 1
)

echo SUCCESS: All SQL scripts applied successfully.
echo Total files processed: %count%

endlocal