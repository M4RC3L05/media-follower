#!/usr/bin/env sh

set -e

# =====================================
# GLOBALS
# =====================================

SCRIPT_DIR="$(cd "$(dirname -- "$0")" >/dev/null 2>&1 && pwd -P)"
SCRIPT_PATH="$SCRIPT_DIR/$(basename -- "$0")"
PS4='\033[35m$\033[0m '

# =====================================
# UTILS
# =====================================

__utils:print_and_exec() {
  printf "\033[35m$\033[0m %s\n" "$*"
  "$@"
  return $?
}

__utils:list_tasks() {
  while IFS= read -r line; do
    case "$line" in
      task:*)
        final="${line#"task:"}"
        final="${final%"() {"}"
        printf "%s\n" "$final"
      ;;
    esac
  done < $SCRIPT_PATH
}

__utils:task_in_tasks() {
  for task_in_tasks__task in $2; do
    if [ "$task_in_tasks__task" = "$1" ]; then
      return 0
    fi
  done

  return 1
}

# =====================================
# TASKS
# =====================================

__cmd:compile() {
  compile__cmd_name="$1"

  __utils:print_and_exec rm --verbose -f ".bin/$compile__cmd_name"
  __utils:print_and_exec go build -ldflags "-w -s" -o ".bin/$compile__cmd_name" "cmd/$compile__cmd_name/main.go"
}

task:code:test() {
  task:cmd:db_migrate
  __task:cmd:admin_server:assets

  __utils:print_and_exec ginkgo \
    -r \
    -p \
    --randomize-all \
    --randomize-suites \
    --fail-on-pending \
    --fail-on-empty \
    --keep-going \
    --trace
}

task:code:check() {
  __utils:print_and_exec golangci-lint run ./...
  __utils:print_and_exec golangci-lint fmt --diff-colored ./...
}

task:deps:update() {
  __utils:print_and_exec go get -u ./...
  __utils:print_and_exec go mod tidy
}

__task:cmd:admin_server:assets() {
  __utils:print_and_exec rm -rfv cmd/admin_server/.dist
  __utils:print_and_exec go run scripts/bundle_assets/main.go \
    -e cmd/admin_server/frontend/*.css \
    -e cmd/admin_server/frontend/*.ts \
    -o cmd/admin_server/.dist
  __utils:print_and_exec cp --verbose -R cmd/admin_server/frontend/static cmd/admin_server/.dist
  __utils:print_and_exec templ generate
}

task:cmd:admin_server() {
  __task:cmd:admin_server:assets
  __cmd:compile "admin_server"
  __utils:print_and_exec .bin/admin_server "$@" | spretty
}

task:cmd:db_migrate() {
  __cmd:compile "db_migrate"
  __utils:print_and_exec .bin/db_migrate "$@" | spretty

  __utils:print_and_exec go run scripts/gen_db_types/main.go
  set -x; sqlite3 data/app.db ".schema" > internal/store/schema.sql; { set +x; } 2>/dev/null
}

task:cmd:fetch_outputs() {
  __cmd:compile "fetch_outputs_job"
  __utils:print_and_exec .bin/fetch_outputs_job "$@" | spretty
}

main() {
  tasks="$(__utils:list_tasks)"
  task="$1"

  if ! __utils:task_in_tasks "$task" "$tasks"; then
    printf "Task \"%s\" not available\n\n" "$task"
    printf "Possible tasks:\n%s\n" "$tasks"
    exit 1
  fi

  shift

  __utils:print_and_exec task:$task "$@"
}

main "$@"
