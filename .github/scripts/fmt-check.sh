#!/bin/sh
gofiles=$(find . -type f -name '*.go')
[ -z "$gofiles" ] && exit 0

unformatted=$(gofmt -l $gofiles)
[ -z "$unformatted" ] && exit 0

echo >&2 "Go files must be formatted with go fmt. Please format these files:"
for fn in $unformatted; do
    echo >&2 "  $fn"
done

exit 1