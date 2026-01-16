#!/bin/bash
# Check test coverage threshold
# Usage: ./scripts/check-coverage.sh <threshold> <coverage_file>
# Example: ./scripts/check-coverage.sh 80.0 coverage.out

THRESHOLD=${1:-90.0}
COVERAGE_FILE=${2:-coverage.out}

if [ ! -f "$COVERAGE_FILE" ]; then
    echo "❌ Coverage file not found: $COVERAGE_FILE"
    exit 1
fi

# Extract total coverage percentage
COVERAGE=$(go tool cover -func="$COVERAGE_FILE" | grep total | awk '{print $3}' | sed 's/%//')

if [ -z "$COVERAGE" ]; then
    echo "❌ Could not extract coverage from $COVERAGE_FILE"
    exit 1
fi

echo "Coverage: ${COVERAGE}%"
echo "Threshold: ${THRESHOLD}%"

# Compare using bc for floating point
RESULT=$(echo "$COVERAGE >= $THRESHOLD" | bc -l)

if [ "$RESULT" -eq 1 ]; then
    echo "✅ Coverage ${COVERAGE}% meets threshold ${THRESHOLD}%"
    exit 0
else
    echo "❌ Coverage ${COVERAGE}% is below threshold ${THRESHOLD}%"
    exit 1
fi
