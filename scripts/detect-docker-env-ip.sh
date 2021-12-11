#!/usr/bin/env sh
set +e

IP="127.0.0.1"

if command -v minikube >/dev/null; then
  dockerEnv=$(minikube status -f '{{- if .DockerEnv }}{{.DockerEnv}}{{- end }}')
  if [ $? -eq 0 ]; then
    if [ "${dockerEnv}" = "in-use" ]; then
      IP=$(minikube ip)
    fi
  fi
fi

echo ${IP}
