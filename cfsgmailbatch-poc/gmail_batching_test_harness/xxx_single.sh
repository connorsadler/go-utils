#!/bin/sh

echo "Making single request"

#
# You need to do this first:
#   1. find your gmail access token
#   2. export GMAIL_ACCESS_TOKEN="ya29.a0_____blahblah_______"
#   3. Run this script with:
#          ./xxx_single.sh
#
#
# Gmail batching docs:
#   https://developers.google.com/gmail/api/guides/batch#:~:text=The%20Gmail%20API%20supports%20batching,lot%20of%20data%20to%20upload.
#

if [[ "$GMAIL_ACCESS_TOKEN" == "" ]]; then
	echo "Error - GMAIL_ACCESS_TOKEN is not set"
	exit 999
fi

curl -v \
	-X GET \
	-H "Authorization: Bearer $GMAIL_ACCESS_TOKEN" \
	https://gmail.googleapis.com/gmail/v1/users/me/messages/15a4a1151a83d7ab?alt=json&format=metadata&metadataHeaders=Subject&metadataHeaders=Date&metadataHeaders=From&prettyPrint=false

