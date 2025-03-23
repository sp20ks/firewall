set -e

echo "Resolving modules in $(pwd)..."

find . -mindepth 2 -type f -name go.mod | while read -r modfile; do
    dir=$(dirname "$modfile")
    echo "Running linter in $dir ..."
    (cd "$dir" && golangci-lint run) || exit 1
done