#!/usr/bin/env bash
set -euo pipefail

APP_PORT="${APP_PORT:-8081}"
PING_NAME="${PING_NAME:-loadgen}"
RPS_BURST="${RPS_BURST:-5}"
INTERVAL_SEC="${INTERVAL_SEC:-1}"
URL="${URL:-http://localhost:${APP_PORT}/api/v1/ping?name=${PING_NAME}}"

STATE_DIR="${STATE_DIR:-.tmp}"
PID_FILE="${STATE_DIR}/ping_load.pid"
LOG_FILE="${STATE_DIR}/ping_load.log"

mkdir -p "${STATE_DIR}"

is_running() {
  [[ -f "${PID_FILE}" ]] || return 1
  local pid
  pid="$(cat "${PID_FILE}")"
  kill -0 "${pid}" 2>/dev/null
}

run_loop() {
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] started: url=${URL} burst=${RPS_BURST}/sec" >>"${LOG_FILE}"

  while true; do
    local i
    for ((i = 1; i <= RPS_BURST; i++)); do
      (
        curl -sS -o /dev/null "${URL}" || true
      ) &
    done
    wait
    sleep "${INTERVAL_SEC}"
  done
}

start() {
  if is_running; then
    echo "already running, pid=$(cat "${PID_FILE}")"
    exit 0
  fi

  nohup bash -c "$(declare -f run_loop); run_loop" >>"${LOG_FILE}" 2>&1 &
  echo $! >"${PID_FILE}"
  echo "started, pid=$(cat "${PID_FILE}")"
  echo "logs: ${LOG_FILE}"
}

stop() {
  if ! is_running; then
    echo "not running"
    rm -f "${PID_FILE}"
    exit 0
  fi

  local pid
  pid="$(cat "${PID_FILE}")"
  kill "${pid}" 2>/dev/null || true
  rm -f "${PID_FILE}"
  echo "stopped"
}

status() {
  if is_running; then
    echo "running, pid=$(cat "${PID_FILE}")"
  else
    echo "not running"
  fi
  echo "log: ${LOG_FILE}"
}

logs() {
  tail -n 50 -f "${LOG_FILE}"
}

usage() {
  cat <<USAGE
Usage: $0 <start|stop|status|logs>

Env vars:
  APP_PORT     default: 8081
  PING_NAME    default: loadgen
  RPS_BURST    default: 5
  INTERVAL_SEC default: 1
  URL          default: http://localhost:
USAGE
}

cmd="${1:-}"
case "${cmd}" in
  start) start ;;
  stop) stop ;;
  status) status ;;
  logs) logs ;;
  *) usage; exit 1 ;;
esac
