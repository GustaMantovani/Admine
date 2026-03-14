#!/usr/bin/env bash
set -euo pipefail

# ─── Admine Control Script ──────────────────────────────────────────────────
# Usage: ./admine.sh {start|stop|restart|status|logs}
# ─────────────────────────────────────────────────────────────────────────────

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PID_DIR="/tmp/admine/pids"

SERVICES=("server_handler" "vpn_handler" "bot")

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

log_info()  { echo -e "${CYAN}[INFO]${NC}  $*"; }
log_ok()    { echo -e "${GREEN}[OK]${NC}    $*"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC}  $*"; }
log_err()   { echo -e "${RED}[ERROR]${NC} $*"; }

# ─── Helpers ─────────────────────────────────────────────────────────────────

ensure_dirs() {
    mkdir -p "$PID_DIR" "/tmp/admine/logs"
}

get_pid() {
    local service="$1"
    local pid_file="$PID_DIR/$service.pid"
    if [[ -f "$pid_file" ]]; then
        cat "$pid_file"
    fi
}

is_running() {
    local pid="$1"
    [[ -n "$pid" ]] && kill -0 "$pid" 2>/dev/null
}

# ─── Start ───────────────────────────────────────────────────────────────────

start_redis() {
    log_info "Starting Redis (docker compose)..."
    (cd "$SCRIPT_DIR/pubsub/redis" && docker compose up -d) > /dev/null 2>&1
    sleep 2
    log_ok "Redis started"
}

start_service() {
    local service="$1"
    local pid
    pid=$(get_pid "$service")

    if is_running "$pid"; then
        log_warn "$service is already running (PID $pid)"
        return 0
    fi

    local binary="$SCRIPT_DIR/$service/$service"
    if [[ ! -x "$binary" ]]; then
        log_err "$service binary not found or not executable: $binary"
        return 1
    fi

    log_info "Starting $service..."
    (cd "$SCRIPT_DIR/$service" && exec ./"$service" > /dev/null 2>&1) &
    local new_pid=$!
    echo "$new_pid" > "$PID_DIR/$service.pid"
    
    # Brief check to make sure it didn't die immediately
    sleep 0.5
    if is_running "$new_pid"; then
        log_ok "$service started (PID $new_pid)"
    else
        log_err "$service failed to start — check application logs"
        rm -f "$PID_DIR/$service.pid"
        return 1
    fi
}

do_start() {
    ensure_dirs
    echo -e "${BOLD}━━━ Starting Admine ━━━${NC}"
    start_redis
    for service in "${SERVICES[@]}"; do
        start_service "$service"
    done
    echo -e "${BOLD}━━━ Admine started ━━━${NC}"
}

# ─── Stop ────────────────────────────────────────────────────────────────────

stop_service() {
    local service="$1"
    local pid
    pid=$(get_pid "$service")

    if ! is_running "$pid"; then
        log_warn "$service is not running"
        rm -f "$PID_DIR/$service.pid"
        return 0
    fi

    log_info "Stopping $service (PID $pid)..."
    kill "$pid" 2>/dev/null || true

    # Wait up to 5 seconds for graceful shutdown
    local waited=0
    while is_running "$pid" && (( waited < 5 )); do
        sleep 1
        waited=$((waited + 1))
    done

    if is_running "$pid"; then
        log_warn "$service didn't stop gracefully, sending SIGKILL..."
        kill -9 "$pid" 2>/dev/null || true
    fi

    rm -f "$PID_DIR/$service.pid"
    log_ok "$service stopped"
}

stop_redis() {
    log_info "Stopping Redis (docker compose)..."
    (cd "$SCRIPT_DIR/pubsub/redis" && docker compose down) > /dev/null 2>&1
    log_ok "Redis stopped"
}

do_stop() {
    echo -e "${BOLD}━━━ Stopping Admine ━━━${NC}"
    # Stop in reverse order: bot → vpn_handler → server_handler → redis
    for (( i=${#SERVICES[@]}-1; i>=0; i-- )); do
        stop_service "${SERVICES[$i]}"
    done
    stop_redis
    echo -e "${BOLD}━━━ Admine stopped ━━━${NC}"
}

# ─── Status ──────────────────────────────────────────────────────────────────

do_status() {
    echo -e "${BOLD}━━━ Admine Status ━━━${NC}"
    
    # Redis
    if (cd "$SCRIPT_DIR/pubsub/redis" && docker compose ps --status running 2>/dev/null | grep -q redis); then
        echo -e "  ${GREEN}●${NC} redis          ${GREEN}running${NC}"
    else
        echo -e "  ${RED}●${NC} redis          ${RED}stopped${NC}"
    fi

    # Services
    for service in "${SERVICES[@]}"; do
        local pid
        pid=$(get_pid "$service")
        if is_running "$pid"; then
            echo -e "  ${GREEN}●${NC} ${service}$(printf '%*s' $((15 - ${#service})) '')${GREEN}running${NC}  (PID $pid)"
        else
            echo -e "  ${RED}●${NC} ${service}$(printf '%*s' $((15 - ${#service})) '')${RED}stopped${NC}"
            # Clean stale PID file
            rm -f "$PID_DIR/$service.pid"
        fi
    done
}


# ─── Logs ────────────────────────────────────────────────────────────────────

declare -A SERVICE_LOG=(
    ["server_handler"]="/tmp/admine/logs/server_handler.log"
    ["vpn_handler"]="/tmp/admine/logs/vpn_handler.log"
    ["bot"]="/tmp/admine/logs/bot.log"
)

do_logs() {
    local target="${1:-all}"
    local follow=false
    local lines=50
    shift || true

    while [[ $# -gt 0 ]]; do
        case "$1" in
            -f|--follow) follow=true; shift ;;
            -n|--lines)  lines="${2:?'--lines requires a number'}"; shift 2 ;;
            *) log_err "Unknown option: $1"; usage; return 1 ;;
        esac
    done

    local tail_args=("-n" "$lines")
    $follow && tail_args+=("-f")

    if [[ "$target" == "all" ]]; then
        local existing_logs=()
        for svc in "${SERVICES[@]}"; do
            local log_file="${SERVICE_LOG[$svc]:-}"
            if [[ -n "$log_file" && -f "$log_file" ]]; then
                existing_logs+=("$log_file")
            else
                log_warn "No log file found for $svc (${log_file:-undefined})"
            fi
        done

        if [[ ${#existing_logs[@]} -eq 0 ]]; then
            log_err "No log files found. Have the services been started?"
            return 1
        fi

        tail "${tail_args[@]}" "${existing_logs[@]}"
    else
        local log_file="${SERVICE_LOG[$target]:-}"
        if [[ -z "$log_file" ]]; then
            log_err "Unknown service: $target"
            echo "Available: ${SERVICES[*]} all"
            return 1
        fi
        if [[ ! -f "$log_file" ]]; then
            log_err "Log file not found: $log_file"
            log_warn "Has $target been started at least once?"
            return 1
        fi

        tail "${tail_args[@]}" "$log_file"
    fi
}

# ─── Restart ─────────────────────────────────────────────────────────────────

do_restart() {
    local target="${1:-all}"

    if [[ "$target" == "all" ]]; then
        do_stop
        echo ""
        do_start
    else
        if [[ " ${SERVICES[*]} " =~ " $target " ]]; then
            stop_service "$target"
            start_service "$target"
        else
            log_err "Unknown service: $target"
            echo "Available: ${SERVICES[*]} all"
        fi
    fi
}

# ─── Main ────────────────────────────────────────────────────────────────────

usage() {
    echo -e "${BOLD}Admine Control Script${NC}"
    echo ""
    echo "Usage: $0 <command> [options]"
    echo ""
    echo "Commands:"
    echo "  start                        Start all services"
    echo "  stop                         Stop all services"
    echo "  restart [service|all]        Restart all or a specific service"
    echo "  status                       Show status of all services"
    echo "  logs [service|all] [-f] [-n <lines>]"
    echo "                               Show logs (default: all services, last 50 lines)"
    echo "                               -f / --follow  Follow log output"
    echo "                               -n / --lines   Number of lines to show"
    echo ""
    echo "Services: ${SERVICES[*]}"
}

case "${1:-}" in
    start)   do_start ;;
    stop)    do_stop ;;
    restart) do_restart "${2:-all}" ;;
    status)  do_status ;;
    logs)    do_logs "${2:-all}" "${@:3}" ;;
    *)       usage ;;
esac
