#!/bin/bash -e

# print timestamp and command then execute the command
echo_and_run() { echo "$(date +%T) $*"; $*; }

for pod_file in $(find pods -type f -name *.y*ml); do
  cmd="go run kube_obj_lint.go $pod_file"

  if [ "$pod_file" = "pods/undeclared-volumes-mounted.yml" ]; then
    # object files that are expected to fail should have non-zero exit code
    if echo_and_run $cmd; then
      echo "'$cmd' expected to fail but did not"
      exit 1
    fi
  else
    # valid object files should return zero exit code
    echo_and_run $cmd
  fi
done
