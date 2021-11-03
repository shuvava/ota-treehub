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
file="object.txt"
temp_file=$(mktemp)
SHA256=$(< "$file" openssl dgst -sha256)
print "file sha256 $SHA256"
PREFIX="${SHA256[*]:0:2}"
SUFFIX="${SHA256[*]:2}"
URL="http://${TREEHUB_SVC}/api/v3/objects/${PREFIX}/${SUFFIX}.commit"
print "url ${URL}"

curl -H "x-ats-namespace:default" "${URL}" > "$temp_file"
NEW_SHA256=$(< "$temp_file" openssl dgst -sha256)

if [ "$NEW_SHA256" == "$SHA256" ]; then
  print "hashes match" "$BGreen"
else
  print "hashes do not match"  "$BRed"
fi
