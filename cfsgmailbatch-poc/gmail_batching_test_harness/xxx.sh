#!/bin/sh

echo "Making POST"

#
#
# You need to do this first:
#   1. find your gmail access token, and ensure it's working
#   2. export GMAIL_ACCESS_TOKEN="ya29.a0_____blahblah_______"
#   3. Run this script with:
#          ./xxx.sh
#
# Gmail batching docs:
#   https://developers.google.com/gmail/api/guides/batch#:~:text=The%20Gmail%20API%20supports%20batching,lot%20of%20data%20to%20upload.
#
# Notes:
# - DONT use --data as it removes all newlines!!!!!!!!!!!
#

if [[ "$GMAIL_ACCESS_TOKEN" == "" ]]; then
	echo "Error - GMAIL_ACCESS_TOKEN is not set"
	exit 999
fi

curl -v \
	--http1.1 \
	-X POST \
	--data-binary "@xxx.txt" \
	-H "Content-Type: multipart/mixed; boundary=\"foo_bar\"" \
	-H "Authorization: Bearer $GMAIL_ACCESS_TOKEN" \
	https://www.googleapis.com/batch/gmail/v1


# actual url: 
# dummy url: https://webhook.site/5082de3f-6ddc-456a-b00b-af7943fd90a0

# todo: not sure if we should be hitting "https://gmail.googleapis.com/batch/gmail/v1" or "https://www.googleapis.com/batch/gmail/v1"

#  https://gmail.googleapis.com/gmail/v1/users/me/messages/15a4a1151a83d7ab?alt=json&format=metadata&metadataHeaders=Subject&metadataHeaders=Date&metadataHeaders=From&prettyPrint=false

