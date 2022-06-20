#!/bin/bash
set -e

for go_project in */; do
  ( cd "$go_project" && go test -v)
done