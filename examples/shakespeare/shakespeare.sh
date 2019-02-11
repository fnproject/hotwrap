#!/bin/sh

# Anything sent to stderr will go to function logs
echo running the shakespeare repository 1>&2
read -r cmd
case "$cmd" in
sonnets) cat sonnets.txt ;;
hamlet) cat hamlet.txt ;;
*) echo "pick hamlet or sonnets" ;;
esac