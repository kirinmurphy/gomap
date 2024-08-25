#!/bin/bash 

output=$(go test -count=1 -v ./src/... 2>&1)
exit_code=$?

filtered_output=$(echo "$output" | grep -v '\[no test files\]')

fail_lines=$(echo "$filtered_output" | grep "FAIL:")

if [ $exit_code -ne 0 ]; then
  filtereder_ouput=$(echo "$filtered_output" | grep -v "FAIL:")
  echo -e "\033[1;91m$fail_lines\033[0m"
  echo "$filtereder_ouput"
  echo -e "\033[1;91m===============================\033[0m"
  exit $exit_code
else
  echo "$filtered_output"
  echo -e "\033[1;92mTESTS: OK ツ\033[0m"
fi
