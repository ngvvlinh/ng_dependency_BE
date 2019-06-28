#!/usr/bin/env bash
set -eo pipefail

: ${basic_auth?Missing required variable}
: ${message?Missing required variable}
: ${change_id?Missing required variable}
: ${revision?Missing required variable}

case "$verified" in
  ""|"1"|"-1") ;;
  *) echo "Invalid label"; exit 1 ;;
esac

labels='{}'
if [[ -n "$verified" ]] ; then
  labels=$(jq -n --arg v "$verified" '{ Verified: $v }')
fi

body=$(
  jq -n \
    --arg tag "gitlab-ci" \
    --arg msg "$message" \
    --argjson labels "$labels" \
    '{ tag: $tag, message: $msg, labels: $labels }'
)

curl --user "${basic_auth}" \
  -H "Content-Type: application/json" \
  -d "${body}" \
  "https://g.meta.etop.vn/a/changes/${change_id}/revisions/${revision}/review"
