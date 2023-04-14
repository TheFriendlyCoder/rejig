#!/usr/bin/env bash
# Stop on first error
set -e

# -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
#                                                         CONSTANTS
# minimum code coverage we expect from our code project
TESTCOVERAGE_THRESHOLD=80
# Red ASCII color
RED='\033[0;31m'
GREEN='\033[0;32m'
# NoColor reset code
NC='\033[0m'
# Output file containing coverage metrics
COVERAGE_FILE=coverage.out
# -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=

# Run all tests recursively
#go test ./...

# With shell coverage report
#go test --cover ./...

# Run tests and store coverage report to disk
go test ./... -race -coverprofile "${COVERAGE_FILE}" -covermode atomic

# Remove ignored files from coverage report
#declare -a ignoreFiles=("lib/templateManager/aferoWrapper.go")
#for curFile in "${ignoreFiles[@]}"; do
#  # NOTE: GREP needs to store output to a temp file otherwise the bash
#  #       file redirect will overwrite the original file with an empty
#  #       one before GREP even runs
#  grep -F -v "${curFile}" "${COVERAGE_FILE}" > "${COVERAGE_FILE}.tmp" && mv "${COVERAGE_FILE}.tmp" "${COVERAGE_FILE}"
#  # NOTE: Would be easier to do with SED but this doesn't work on Github
#  #       apparently the -i option is not supported cross platform
#  #  sed -i '' "/${curFile//\//\\/}/d" "${COVERAGE_FILE}"
#done
# Display ASCII report to console
go tool cover -func="${COVERAGE_FILE}"
# Display HTML report in default browser
go tool cover -html="${COVERAGE_FILE}"

# Check coverage threshold
# Extract the coverage total from the coverage report. The line being parsed should look
# something like this:
#       total:								(statements)		12.3%
# and we extract the "12.3" decimal number from the line
totalCoverage=$(go tool cover -func="${COVERAGE_FILE}" | grep "^total" | grep -Eo '[0-9]+\.[0-9]+')
# Here we check 2 floating point values for percentage coverage using awk
if (( $(echo "$totalCoverage $TESTCOVERAGE_THRESHOLD" | awk '{print ($1 > $2)}') )); then
    echo -e "${GREEN}Coverage of ${totalCoverage}% meets minimum ${TESTCOVERAGE_THRESHOLD}%.${NC}"
else
    echo -e "${RED}Coverage of ${totalCoverage}% doesn't meet minimum coverage of ${TESTCOVERAGE_THRESHOLD}%.${NC}"
    exit 1
fi
