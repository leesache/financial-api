#!/bin/bash

find . -type f \( -name "*.go" -o -name "*.c" -o -name "*.py" -o -name "*.sh" \) -exec wc -l {} + | awk '{total += $1} END {print "Total lines of code:", total}'
