#!/usr/bin/env bash

print() {
  BWhite='\033[1;37m'
  Color_Off='\033[0m'
  color=${2:-$BWhite}
  echo -e "${color}$1${Color_Off}"
}

TREEHUB_SVC="localhost:8080"
file="object.txt"
SHA256=$(< "$file" openssl dgst -sha256)
print "file sha256 $SHA256"
PREFIX="${SHA256[*]:0:2}"
SUFFIX="${SHA256[*]:2}"
URL="http://${TREEHUB_SVC}/api/v3/objects/${PREFIX}/${SUFFIX}.commit"
print "url ${URL}"

curl -X "POST" \
  -H "Content-Type:application/octet-stream" \
  -H "x-ats-namespace:default" \
  --data-binary @$file \
  "${URL}" | jq
