#!/usr/bin/env bash
set -eo pipefail

: ${branch?Missing required variable}
: ${basic_auth?Missing required variable}
: ${message?Missing required variable}
: ${change_id?Missing required variable}
: ${revision?Missing required variable}

request() {
  curl --user "${basic_auth}" -H "Content-Type: application/json" "$@" | tail -n +2
}

case "$verified" in
  ""|"1"|"-1") ;;
  *) echo "Invalid label"; exit 1 ;;
esac

# find the change number
if echo "$branch" | grep -E '^changes/[0-9]{2}/[0-9]+/[0-9]+$' ; then
  # branch /changes/xx/xx/x
  change_number=$(echo "$branch" | sed -E 's/changes\/[0-9]+\///' | sed -E 's/\/[0-9]+//')
  : ${change_number?Invalid change number}

else
  # find the change number using gerrit query
  changes=$(request "https://g.meta.etop.vn/a/changes/?q=change:${change_id}")
  change_number=$(echo "$changes" | jq 'map(select(.branch=="'"$branch"'"))|.[0]._number')
fi

echo "change_number=${change_number}"
if ! echo "$change_number" | grep -E '[0-9]+' >/dev/null ; then
  echo "invalid change number"
  exit 1
fi

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

request -d "${body}" \
  "https://g.meta.etop.vn/a/changes/${change_number}/revisions/${revision}/review"
