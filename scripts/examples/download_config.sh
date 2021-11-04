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
URL="http://${TREEHUB_SVC}/api/v2/config"
print "url ${URL}"

curl -i -H "x-ats-namespace:default" "${URL}"
