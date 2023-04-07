#!/usr/bin/env sh -e

# minimum code coverage we expect from our code project
TESTCOVERAGE_THRESHOLD=80
# Red ASCII color
RED='\033[0;31m'
GREEN='\033[0;32m'
# NoColor reset code
NC='\033[0m'

# Run all tests recursively
#go test ./...

# With shell coverage report
#go test --cover ./...

# Run tests and store coverage report to disk
go test ./... -coverprofile coverage.out -covermode set

# Remove ignored files from coverage report
declare -a ignoreFiles=("lib/aferoWrapper.go")
for curFile in "${ignoreFiles[@]}"; do
  sed -i '' "/${curFile//\//\\/}/d" ./coverage.out
done

# Display ASCII report to console
go tool cover -func coverage.out
# Display HTML report in default browser
#go tool cover -html=coverage.out

# Check coverage threshold
# Extract the coverage total from the coverage report. The line being parsed should look
# something like this:
#       total:								(statements)		12.3%
# and we extract the "12.3" decimal number from the line
totalCoverage=$(go tool cover -func=coverage.out | grep "^total" | grep -Eo '[0-9]+\.[0-9]+')
# Here we check 2 floating point values for percentage coverage using awk
if (( $(echo "$totalCoverage $TESTCOVERAGE_THRESHOLD" | awk '{print ($1 > $2)}') )); then
    echo -e "${GREEN}Coverage of ${totalCoverage}% meets minimum ${TESTCOVERAGE_THRESHOLD}%.${NC}"
else
    echo -e "${RED}Coverage of ${totalCoverage}% doesn't meet minimum coverage of ${TESTCOVERAGE_THRESHOLD}%.${NC}"
    exit 1
fi
