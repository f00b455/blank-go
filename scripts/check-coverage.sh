#!/bin/bash

THRESHOLD=$1
COVERAGE_FILE=$2

if [ -z "$THRESHOLD" ] || [ -z "$COVERAGE_FILE" ]; then
    echo "Usage: $0 <threshold> <coverage-file>"
    exit 1
fi

if [ ! -f "$COVERAGE_FILE" ]; then
    echo "Coverage file not found: $COVERAGE_FILE"
    exit 1
fi

COVERAGE=$(go tool cover -func="$COVERAGE_FILE" | grep total | awk '{print $3}' | sed 's/%//')

if [ -z "$COVERAGE" ]; then
    echo "Could not parse coverage from $COVERAGE_FILE"
    exit 1
fi

echo "Coverage: ${COVERAGE}%"
echo "Threshold: ${THRESHOLD}%"

if (( $(echo "$COVERAGE < $THRESHOLD" | bc -l) )); then
    echo "Coverage ${COVERAGE}% is below threshold ${THRESHOLD}%"
    exit 1
fi

echo "Coverage check passed!"
