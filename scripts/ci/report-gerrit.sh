#!/usr/bin/env bash
set -eo pipefail

: ${CI_COMMIT_REF_NAME?Missing required variable}
: ${CI_COMMIT_SHA?Missing required variable}
: ${CI_GERRIT_BASIC_AUTH?Missing required variable}
: ${CI_PIPELINE_URL?Missing required variable}
: ${CI_PROJECT_URL?Missing required variable}
: ${CI_JOB_URL?Missing required variable}

get_numeric_change_id() {
  echo $CI_COMMIT_REF_NAME | sed "s/changes\/[0-9]\+\/\([0-9]\+\)\/.*/\1/"
}

get_change_id() {
  git log -1 | grep "Change-Id" | tail -1 | grep -o "[A-z0-9]\+$"
}

export branch="$CI_COMMIT_REF_NAME"

change_id="$(get_change_id)"
export change_id
: ${change_id?Can not parse change_id}

export revision="$CI_COMMIT_SHA"
basic_auth="$(echo "$CI_GERRIT_BASIC_AUTH" | base64 -d)"
export basic_auth

pipelines_url="${CI_PROJECT_URL}/commit/${revision}/pipelines (#${CI_PIPELINE_ID})"

if [[ "$1" == "running" ]] ; then
  export verified=""
  export message="Build running on ${CI_COMMIT_REF_NAME}. See ${pipelines_url}"

elif [[ ! -f "artifacts/BUILD_SUCCESS" ]] ; then
  export verified="-1"
  export message="Build failure on ${CI_COMMIT_REF_NAME}. See ${pipelines_url}"

else
  export verified="1"
  export message="Build successfully on ${CI_COMMIT_REF_NAME}. See ${pipelines_url}"
  if [[ -f "artifacts/COVERAGE%" ]] ; then
    coverage_url="$(cat artifacts/ARTIFACTS_URL)/browse/artifacts/coverage/"
    message="Build successfully on ${CI_COMMIT_REF_NAME} with coverage $(cat artifacts/COVERAGE%) See ${pipelines_url} and ${coverage_url}"
    export message
  fi
fi

if [[ $(echo "$CI_COMMIT_REF_NAME" | grep "^changes\/") ]]; then
  scripts/ci/report-gerrit-review.sh
else
  scripts/ci/report-gerrit-review.sh || true
fi

echo "$message"
if [[ "$verified" == "-1" ]] ; then exit 255 ; fi
