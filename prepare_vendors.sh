#!/bin/bash
set -e

go_projects=(confirmation_handler distributor notifier resolved_handler scheduler site_checker)
for go_project in ${go_projects[*]}; do
  ( cd "$go_project" && go mod vendor )
done