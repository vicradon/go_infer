#!/bin/bash

# Install git hooks from the hooks directory
set -e

HOOKS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/hooks"
GIT_DIR="$(git rev-parse --git-common-dir 2>/dev/null || echo .git)"
HOOKS_TARGET="$GIT_DIR/hooks"

echo "Installing git hooks from $HOOKS_DIR to $HOOKS_TARGET"

# Create hooks directory if it doesn't exist
mkdir -p "$HOOKS_TARGET"

# Copy all hooks from hooks/ to .git/hooks/
for hook in "$HOOKS_DIR"/*; do
    if [ -f "$hook" ]; then
        hook_name=$(basename "$hook")
        echo "  Installing $hook_name"
        cp "$hook" "$HOOKS_TARGET/$hook_name"
        chmod +x "$HOOKS_TARGET/$hook_name"
    fi
done

echo "✅ Git hooks installed successfully"
