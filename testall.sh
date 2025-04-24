#!/usr/bin/env bash
set -uo pipefail

fail=0
while IFS= read -r dir; do
  echo "===> go test -count=1 $dir"
  pushd "$dir" > /dev/null || continue
  if ! go test -count=1 ./...; then
    fail=1
  fi
  echo "===> $dir test complete"
  echo
  popd > /dev/null
done < <(go list -m -f '{{.Dir}}')

if [[ "${BASH_SOURCE[0]}" != "${0}" ]]; then
  return $fail
else
  exit $fail
fi

