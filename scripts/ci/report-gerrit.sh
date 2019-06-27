#!/usr/bin/env bash
set -eo pipefail

get_change_id() {
  echo $CI_COMMIT_REF_NAME | sed "s/changes\/[0-9]\+\/\([0-9]\+\)\/.*/\1/"
}

change_id="$(get_change_id)"
revision="$CI_COMMIT_SHA"
basic_auth="$(echo $CI_GERRIT_BASIC_AUTH | base64 -d)"

if [[ -z "${change_id}" ]] ; then exit 1 ; fi

if [[ -f "artifacts/BUILD_SUCCESS" ]] ; then
  verified="1"
  message="Build successfully. See https://code.eyeteam.vn/etop-backend/backend/commit/${revision}/pipelines"
else
  verified="-1"
  message="Build failure. See https://code.eyeteam.vn/etop-backend/backend/commit/${revision}/pipelines"
fi

body='{
  "tag": "gitlab-ci",
  "message": "'"${message}"'",
  "labels": {"Verified": '$verified'}
}'

curl --user "${basic_auth}" \
  -H "Content-Type: application/json" \
  -d "${body}" \
  "https://g.meta.etop.vn/a/changes/${change_id}/revisions/${revision}/review"

echo "${message}"
if [[ "${verified}" != "1" ]] ; then exit 255 ; fi
