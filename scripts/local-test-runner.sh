#!/bin/bash 

output=$(go test -count=1 -v ./src/... 2>&1)
exit_code=$?

filtered_output=$(echo "$output" | grep -v '\[no test files\]' | grep -v -- "=== RUN" | grep -v -- "--- PASS:")

fail_count=$(echo "$filtered_output" | grep -c -- "--- FAIL:")

fail_lines=$(echo "$filtered_output" | grep "FAIL:")

if [ $exit_code -ne 0 ]; then
  filtered_output=$(echo "$filtered_output" | grep -v -- "FAIL:")
  echo -e "\033[1;91m===========================\033[0m"
  echo "$filtered_output"
  echo -e "\033[1;91mFailed Tests: $fail_count\033[0m"
  echo -e "\033[1;91m$fail_lines\033[0m"
  exit $exit_code
else
  echo "$filtered_output"
  echo -e "\033[1;92mALL TESTS OK ツ\033[0m"
fi
