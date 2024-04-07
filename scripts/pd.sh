#!/bin/sh

PD_EMAIL=EMAIL
PD_TOKEN=PDTOKEN

curl -vvv -H "From: ${PD_EMAIL}" -H "Authorization: Token token=${TOKEN}" \
   -H 'Content-Type: application/json' -H 'Accept: application/vnd.pagerduty+json;version=2' \
   -d '{ "incident": { "type": "incident", "title": "fire, fire", "service": { "id": "PFCSJTG", "type": "service_reference" } } }'\
   https://api.pagerduty.com/incidents
