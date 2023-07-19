#!/usr/bin/env bash

print() {
  BWhite='\033[1;37m'
  Color_Off='\033[0m'
  color=${2:-$BWhite}
  echo -e "${color}$1${Color_Off}"
}

TREEHUB_SVC="localhost:8080"
BRANCH="master"
file="ref_master.txt"
URL="http://${TREEHUB_SVC}/api/v3/refs/heads/${BRANCH}"
print "url ${URL}"

curl -X "POST" \
  -H "Content-Type:application/octet-stream" \
  -H "x-ats-namespace:default" \
  -H "x-ats-ostree-force:true" \
  --data-binary @$file \
  "${URL}" | jq
