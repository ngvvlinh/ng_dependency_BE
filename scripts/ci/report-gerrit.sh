#!/usr/bin/env bash
set -eo pipefail

: ${CI_COMMIT_REF_NAME?Missing required variable}
: ${CI_COMMIT_SHA?Missing required variable}
: ${CI_GERRIT_BASIC_AUTH?Missing required variable}
: ${CI_PIPELINE_URL?Missing required variable}
: ${CI_PROJECT_URL?Missing required variable}

get_change_id() {
  echo $CI_COMMIT_REF_NAME | sed "s/changes\/[0-9]\+\/\([0-9]\+\)\/.*/\1/"
}

change_id="$(get_change_id)"
: ${change_id?Can not parse change_id}

revision="$CI_COMMIT_SHA"
pipelines_url="${CI_PIPELINE_URL}"
basic_auth="$(echo $CI_GERRIT_BASIC_AUTH | base64 -d)"

if [[ "$1" == "running" ]] ; then
  verified=""
  message="Build running. See ${pipelines_url}"
elif [[ -f "artifacts/BUILD_SUCCESS" ]] ; then
  verified="1"
  message="Build successfully. See ${pipelines_url}"
else
  verified="-1"
  message="Build failure. See ${pipelines_url}"
fi

if [[ -z "${verified}" ]] ; then
  body='{"tag": "gitlab-ci", "message": "'"${message}"'"}'
else
  body='{"tag": "gitlab-ci", "message": "'"${message}"'", "labels": {"Verified": '${verified}'}}'
fi

curl --user "${basic_auth}" \
  -H "Content-Type: application/json" \
  -d "${body}" \
  "https://g.meta.etop.vn/a/changes/${change_id}/revisions/${revision}/review"

echo "${message}"
if [[ "${verified}" == "-1" ]] ; then exit 255 ; fi
