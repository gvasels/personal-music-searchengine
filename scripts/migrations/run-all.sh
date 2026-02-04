#!/bin/bash
# run-all.sh - Run all migrations for global-user-type feature
#
# Usage: ./run-all.sh [--dry-run]

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Pass through arguments
ARGS="$@"

echo "======================================"
echo "Running Global User Type Migrations"
echo "======================================"
echo ""

# Migration 1: User Roles
echo ">>> Migration 1: User Roles"
"$SCRIPT_DIR/migrate-user-roles.sh" $ARGS

echo ""
echo "--------------------------------------"
echo ""

# Migration 2: Playlist Visibility
echo ">>> Migration 2: Playlist Visibility"
"$SCRIPT_DIR/migrate-playlist-visibility.sh" $ARGS

echo ""
echo "======================================"
echo "All migrations completed!"
echo "======================================"
