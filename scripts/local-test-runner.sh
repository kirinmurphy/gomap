#!/bin/bash 

output=$(go test -count=1 ./src/... 2>&1)
exit_code=$?

filtered_output=$(echo "$output" | grep -v '\[no test files\]')

if [ $exit_code -ne 0 ]; then
  echo -e "\033[1;91m$filtered_output\033[0m"
  exit $exit_code
fi

echo -e "\033[1;92mTESTS: OK ツ\033[0m"
  