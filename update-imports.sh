#!/bin/bash

# Script to update import paths after refactoring

echo "Updating import paths in mobius-server..."

# Update all Go files in mobius-server to use the new module paths
find /Users/awar/Documents/Mobius/mobius-server -name "*.go" -type f -exec sed -i '' \
    -e 's|github\.com/notawar/mobius/internal/server|github.com/notawar/mobius/mobius-server/internal/server|g' \
    -e 's|github\.com/notawar/mobius/pkg/|github.com/notawar/mobius/mobius-server/pkg/|g' \
    {} \;

echo "Updating import paths in mobius-cli..."

# Update all Go files in mobius-cli to use the new module paths
find /Users/awar/Documents/Mobius/mobius-cli -name "*.go" -type f -exec sed -i '' \
    -e 's|github\.com/notawar/mobius/cmd/mobiuscli|github.com/notawar/mobius/mobius-cli/cmd/mobiuscli|g' \
    {} \;

echo "Import path updates complete!"
