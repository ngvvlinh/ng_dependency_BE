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

export change_id="$(get_change_id)"
: ${change_id?Can not parse change_id}

export revision="$CI_COMMIT_SHA"
export basic_auth="$(echo "$CI_GERRIT_BASIC_AUTH" | base64 -d)"

pipelines_url="${CI_PROJECT_URL}/commit/${revision}/pipelines (#${CI_PIPELINE_ID})"

if [[ "$1" == "running" ]] ; then
  export verified=""
  export message="Build running. See ${pipelines_url}"
elif [[ ! -f "artifacts/BUILD_SUCCESS" ]] ; then
  export verified="-1"
  export message="Build failure. See ${pipelines_url}"
else
  export verified="1"
  export message="Build successfully. See ${pipelines_url}"
  if [[ -f "artifacts/COVERAGE" ]] ; then
    coverage_url="$(cat artifacts/COVERAGE_URL)"
    export message="Build successfully with coverage $(cat artifacts/COVERAGE) See ${pipelines_url} and ${coverage_url}"
  fi
fi

if [[ $(echo "$CI_COMIT_REF_NAME" | grep "^changes\/") ]]; then
  scripts/ci/report-gerrit-review.sh
else
  scripts/ci/report-gerrit-review.sh || true
fi

echo "$message"
if [[ "$verified" == "-1" ]] ; then exit 255 ; fi
