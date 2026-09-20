#!/bin/bash
# 1Panel 社区版安装脚本
#
# 该脚本位于发布包内，由 quick_start.sh 自动调用，也可在解压后的目录中手动执行：
#   bash install.sh                # 安装
#   bash install.sh uninstall      # 卸载（保留数据目录）
#
# 环境变量：
#   PANEL_BASE_DIR   安装目录，默认 /opt
#   PANEL_PORT       面板端口，默认 9999
#   PANEL_LANGUAGE   面板语言，zh 或 en，默认 zh
#   PANEL_SKIP_START 置为 1 时只安装文件不启动服务

PANEL_BASE_DIR="${PANEL_BASE_DIR:-/opt}"
PANEL_PORT="${PANEL_PORT:-9999}"
PANEL_LANGUAGE="${PANEL_LANGUAGE:-zh}"
PANEL_SKIP_START="${PANEL_SKIP_START:-0}"

BIN_DIR="/usr/local/bin"
LANG_DIR="${BIN_DIR}/lang"
UNIT_DIR="/etc/systemd/system"
PKG_DIR=$(cd "$(dirname "$0")" && pwd)

# 从发布包目录名（1panel-v2.0.0-linux-amd64）中解析版本号
PKG_DIR_NAME=$(basename "$PKG_DIR")
VERSION="${PANEL_VERSION:-}"
if [[ "$PKG_DIR_NAME" =~ ^1panel-(.+)-linux-[a-z0-9]+$ ]]; then
    VERSION="${BASH_REMATCH[1]}"
fi
[[ -n "$VERSION" ]] || VERSION="v2.0.0"

load_lang() {
    if [[ -f "${PKG_DIR}/lang/${PANEL_LANGUAGE}.sh" ]]; then
        # shellcheck disable=SC1090
        . "${PKG_DIR}/lang/${PANEL_LANGUAGE}.sh"
    elif [[ -f "${PKG_DIR}/lang/zh.sh" ]]; then
        # shellcheck disable=SC1090
        . "${PKG_DIR}/lang/zh.sh"
    fi
}

msg() {
    # msg <zh> <en>
    if [[ "$PANEL_LANGUAGE" == "en" ]]; then
        echo "$2"
    else
        echo "$1"
    fi
}

info() {
    msg "$1" "$2"
}

load_lang

fail() {
    msg "$1" "$2"
    exit 1
}

is_root() {
    [[ "$(id -u)" == "0" ]]
}

detect_arch() {
    local osCheck
    osCheck=$(uname -a)
    if [[ $osCheck =~ 'x86_64' ]]; then
        echo "amd64"
    elif [[ $osCheck =~ 'arm64' ]] || [[ $osCheck =~ 'aarch64' ]]; then
        echo "arm64"
    elif [[ $osCheck =~ 'armv7l' ]]; then
        echo "armv7"
    elif [[ $osCheck =~ 'ppc64le' ]]; then
        echo "ppc64le"
    elif [[ $osCheck =~ 's390x' ]]; then
        echo "s390x"
    elif [[ $osCheck =~ 'riscv64' ]]; then
        echo "riscv64"
    elif [[ $osCheck =~ 'loongarch64' ]]; then
        echo "loong64"
    else
        echo ""
    fi
}

rand_str() {
    local length="${1:-10}"
    local result=""
    while [[ ${#result} -lt $length ]]; do
        if [[ -r /dev/urandom ]]; then
            result="${result}$(head -c 128 /dev/urandom | base64 2>/dev/null | tr -dc 'a-zA-Z0-9')"
        else
            result="${result}$(openssl rand -base64 128 2>/dev/null | tr -dc 'a-zA-Z0-9')"
        fi
    done
    echo "${result:0:$length}"
}

port_in_use() {
    local port="$1"
    if command -v ss >/dev/null 2>&1; then
        ss -ltn 2>/dev/null | awk '{print $4}' | grep -Eq ":${port}\$"
        return $?
    fi
    (echo >"/dev/tcp/127.0.0.1/${port}") >/dev/null 2>&1 && return 0
    return 1
}

# 读取已有 1pctl 中的值，兼容 KEY=value 与历史版本的 KEY="value" 两种写法
read_1pctl_value() {
    local key="$1"
    sed -n "s/^${key}=//p" "${BIN_DIR}/1pctl" 2>/dev/null |
        head -n 1 |
        sed -e 's/^"//' -e 's/"$//'
}

write_1pctl() {
    local username="$1" password="$2" entrance="$3"
    # 注意：core/agent 会逐行读取这些值，必须写成 KEY=value（值两边不能加引号），
    # 加引号会让面板读到的端口、安装目录带上引号而无法启动。
    cat >"${BIN_DIR}/1pctl" <<EOF
#!/bin/bash
# 1Panel 命令行工具（社区版），由 scripts/1pctl 安装生成
# 注意：BASE_DIR、LANGUAGE、ORIGINAL_* 会被 core/agent 逐行解析，
# 必须写成 KEY=value（值两边不能加引号），也不要调整行结构。

BASE_DIR=${PANEL_BASE_DIR}
LANGUAGE=${PANEL_LANGUAGE}
ORIGINAL_PORT=${PANEL_PORT}
ORIGINAL_VERSION=${VERSION}
ORIGINAL_USERNAME=${username}
ORIGINAL_PASSWORD=${password}
ORIGINAL_ENTRANCE=${entrance}
EOF
    sed -n '/^CORE_BIN=/,$p' "${PKG_DIR}/1pctl" >>"${BIN_DIR}/1pctl"
    chmod 755 "${BIN_DIR}/1pctl"
}

uninstall() {
    if [[ ! -f "${BIN_DIR}/1panel-core" ]]; then
        msg "$MSG_NO_INSTALL" "No installed 1Panel was detected."
        exit 0
    fi
    for svc in 1panel-core 1panel-agent; do
        systemctl disable --now "$svc" >/dev/null 2>&1 || true
        rm -f "${UNIT_DIR}/${svc}.service"
    done
    systemctl daemon-reload >/dev/null 2>&1 || true
    rm -f "${BIN_DIR}/1panel-core" "${BIN_DIR}/1panel-agent" "${BIN_DIR}/1pctl"
    rm -rf "${LANG_DIR}"
    info "$MSG_UNINSTALL_DONE" "1Panel has been uninstalled."
    exit 0
}

case "${1:-}" in
    uninstall) uninstall ;;
esac

is_root || fail "$MSG_NEED_ROOT" "Please run this command as the root user!"
[[ -d /run/systemd/system ]] || fail "$MSG_NO_SYSTEMD" "systemd is required by 1Panel v2."

arch=$(detect_arch)
[[ -n "$arch" ]] || fail "$MSG_ARCH_UNSUPPORTED" "Unsupported architecture."

for file in "1panel-core" "1panel-agent" "1pctl" "install.sh"; do
    [[ -f "${PKG_DIR}/${file}" ]] || fail "$MSG_PKG_NOT_FOUND (missing ${file})" "Package is incomplete (missing ${file})."
done

if port_in_use "$PANEL_PORT"; then
    fail "端口 ${PANEL_PORT} 已被占用，请通过 PANEL_PORT 指定其他端口。" "Port ${PANEL_PORT} is already in use, set PANEL_PORT to another one."
fi

info "$MSG_INSTALL_START" "Installing 1Panel."

mkdir -p "$BIN_DIR" "$LANG_DIR" "${PANEL_BASE_DIR}/1panel/db" "${PANEL_BASE_DIR}/1panel/log" "${PANEL_BASE_DIR}/1panel/geo"

info "$MSG_INSTALL_BIN ..." "Installing executables ..."
install -m 755 "${PKG_DIR}/1panel-core" "${BIN_DIR}/1panel-core"
install -m 755 "${PKG_DIR}/1panel-agent" "${BIN_DIR}/1panel-agent"

# 已安装过时沿用原有账号信息，保证升级后仍能登录；
# 面板首次启动后会把 1pctl 中的密码改写为 **********，此时需要重新生成
if [[ -f "${BIN_DIR}/1pctl" ]]; then
    old_username=$(read_1pctl_value ORIGINAL_USERNAME)
    old_password=$(read_1pctl_value ORIGINAL_PASSWORD)
    old_entrance=$(read_1pctl_value ORIGINAL_ENTRANCE)
fi
[[ -n "$old_password" && "$old_password" != *"*"* ]] || old_password=""
[[ -n "$old_username" && "$old_username" != *"*"* ]] || old_username=""
[[ -n "$old_entrance" && "$old_entrance" != *"*"* ]] || old_entrance=""
username="${old_username:-$(rand_str 10)}"
password="${old_password:-$(rand_str 10)}"
entrance="${old_entrance:-$(rand_str 10)}"



info "$MSG_INSTALL_CTL ..." "Installing the 1pctl command line tool ..."
write_1pctl "$username" "$password" "$entrance"

info "$MSG_INSTALL_LANG ..." "Installing command line language files ..."
for lang_file in "${PKG_DIR}"/lang/*.sh; do
    [[ -f "$lang_file" ]] || continue
    install -m 644 "$lang_file" "${LANG_DIR}/$(basename "$lang_file")"
done

if [[ -f "${PKG_DIR}/GeoIP.mmdb" ]]; then
    install -m 644 "${PKG_DIR}/GeoIP.mmdb" "${PANEL_BASE_DIR}/1panel/geo/GeoIP.mmdb"
fi

info "$MSG_INSTALL_UNIT ..." "Installing systemd services ..."
for unit in "${PKG_DIR}"/initscript/*.service; do
    [[ -f "$unit" ]] || continue
    install -m 644 "$unit" "${UNIT_DIR}/$(basename "$unit")"
done
systemctl daemon-reload

if [[ "$PANEL_SKIP_START" == "1" ]]; then
    info "$MSG_INSTALL_DONE" "1Panel has been installed."
    exit 0
fi

info "$MSG_INSTALL_START_SVC ..." "Starting 1Panel services ..."
systemctl enable 1panel-core 1panel-agent >/dev/null 2>&1
systemctl restart 1panel-core >/dev/null 2>&1
systemctl restart 1panel-agent >/dev/null 2>&1

local_ip=$(hostname -I 2>/dev/null | awk '{print $1}')
if [[ -z "$local_ip" ]]; then
    local_ip=$(hostname -i 2>/dev/null | awk '{print $1}')
fi
[[ -n "$local_ip" ]] || local_ip="127.0.0.1"

if [[ "$PANEL_LANGUAGE" == "en" ]]; then
    echo "=========================================================="
    echo " 1Panel ${VERSION} installed."
    echo " URL:      http://${local_ip}:${PANEL_PORT}/${entrance}"
    echo " Username: ${username}"
    echo " Password: ${password}"
    echo " Command:  1pctl help"
    echo "=========================================================="
else
    echo "=========================================================="
    echo " 1Panel ${VERSION} 安装完成。"
    echo " 访问地址：http://${local_ip}:${PANEL_PORT}/${entrance}"
    echo " 用户名：${username}"
    echo " 密码：${password}"
    echo " 命令行帮助：1pctl help"
    echo "=========================================================="
fi
