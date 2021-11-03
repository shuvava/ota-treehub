#!/usr/bin/env bash

BRed='\033[1;31m'
BGreen='\033[1;32m'
BWhite='\033[1;37m'
Color_Off='\033[0m'
print() {
  color=${2:-$BWhite}
  echo -e "${color}$1${Color_Off}"
}

TREEHUB_SVC="localhost:8080"
file="ref_master.txt"
temp_file=$(mktemp)
SHA256=$(cat "$file")
print "file sha256 $SHA256"
URL="http://${TREEHUB_SVC}/api/v3/refs/heads/master"
print "url ${URL}"

curl -H "x-ats-namespace:default" "${URL}" > "$temp_file"
NEW_SHA256=$(cat "$temp_file")

if [ "$NEW_SHA256" == "$SHA256" ]; then
  print "hashes match" "$BGreen"
else
  print "hashes do not match"  "$BRed"
fi
