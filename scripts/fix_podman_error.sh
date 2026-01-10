#!/bin/bash
set -e

# Define directories
RUNTIME_DIR="/tmp/podman-run-$(id -u)"
DB_FILE="$HOME/.local/share/containers/storage/db.sql"
DB_BACKUP="${DB_FILE}.backup-$(date +%s)"

echo "🔧 Fixing Podman configuration..."

# 1. Create a logical runtime directory in /tmp since /run/user/1000 is missing
if [ ! -d "$RUNTIME_DIR" ]; then
    echo "Creating runtime directory: $RUNTIME_DIR"
    mkdir -p "$RUNTIME_DIR"
    chmod 700 "$RUNTIME_DIR"
else
    echo "Runtime directory exists: $RUNTIME_DIR"
fi

# 2. Fix the corrupted Podman database by resetting it
# The database contains a hardcoded path to the missing /run/user directory.
if [ -f "$DB_FILE" ]; then
    echo "‼  Detected existing Podman database."
    echo "   Moving corrupt database to $DB_BACKUP"
    mv "$DB_FILE" "$DB_BACKUP"
    echo "✅ Database reset. (Images/Containers will need to be recreated)"
else
    echo "No database found (clean state)."
fi

# 3. Export environment variable
echo "--------------------------------------------------------"
echo "✅ Fix applied."
echo ""
echo "👉 IMPORTANT: You must run the following command in your terminal:"
echo "   export XDG_RUNTIME_DIR=$RUNTIME_DIR"
echo ""
echo "Then allow Podman to regenerate its configuration:"
echo "   podman system migrate"
echo "(If that fails, just running 'make docker-up' should work)"
echo "--------------------------------------------------------"
