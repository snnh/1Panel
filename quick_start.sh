#!/bin/bash
# 1Panel 社区版一键安装脚本
#
# 安装包与更新日志均发布在社区仓库的 GitHub Releases 上：
#   https://github.com/snnh/1Panel/releases
#
# 用法：
#   bash -c "$(curl -sSL https://raw.githubusercontent.com/snnh/1Panel/community-dev/quick_start.sh)"
#
# 常用环境变量：
#   PANEL_REPO            发布仓库，默认 snnh/1Panel
#   PANEL_VERSION         指定版本号，如 v2.0.0；默认取最新正式版
#   PANEL_DOWNLOAD_BASE   下载前缀（镜像/代理），默认 https://github.com
#   PANEL_BASE_DIR        安装目录，默认 /opt
#   PANEL_PORT            面板端口，默认 9999
#   PANEL_LANGUAGE        面板语言，zh 或 en，默认 zh

PANEL_REPO="${PANEL_REPO:-snnh/1Panel}"
PANEL_VERSION="${PANEL_VERSION:-}"
PANEL_DOWNLOAD_BASE="${PANEL_DOWNLOAD_BASE:-https://github.com}"
PANEL_BASE_DIR="${PANEL_BASE_DIR:-/opt}"
PANEL_PORT="${PANEL_PORT:-9999}"
PANEL_LANGUAGE="${PANEL_LANGUAGE:-zh}"

osCheck=$(uname -a)
if [[ $osCheck =~ 'x86_64' ]]; then
    architecture="amd64"
elif [[ $osCheck =~ 'arm64' ]] || [[ $osCheck =~ 'aarch64' ]]; then
    architecture="arm64"
elif [[ $osCheck =~ 'armv7l' ]]; then
    architecture="armv7"
elif [[ $osCheck =~ 'ppc64le' ]]; then
    architecture="ppc64le"
elif [[ $osCheck =~ 's390x' ]]; then
    architecture="s390x"
elif [[ $osCheck =~ 'riscv64' ]]; then
    architecture="riscv64"
elif [[ $osCheck =~ 'loongarch64' ]]; then
    architecture="loong64"
else
    echo "当前系统架构暂不支持，请参考文档选择受支持的系统与架构。"
    echo "The current system architecture is not supported."
    exit 1
fi

is_zh() {
    [[ "$PANEL_LANGUAGE" == "zh" ]]
}

die() {
    if is_zh; then
        echo "$1"
    else
        echo "$2"
    fi
    exit 1
}

for cmd in curl tar; do
    if ! command -v "$cmd" >/dev/null 2>&1; then
        if is_zh; then
            die "缺少依赖命令：$cmd，请先安装。" "Missing required command: $cmd"
        fi
    fi
done

if [[ "$(id -u)" != "0" ]]; then
    die "请使用 root 用户执行安装。" "Please run the installation as the root user."
fi

# 1. 解析版本号
if [[ -z "$PANEL_VERSION" ]]; then
    PANEL_VERSION=$(curl -sfL --max-time 20 \
        "https://api.github.com/repos/${PANEL_REPO}/releases/latest" |
        grep -m1 '"tag_name"' | cut -d '"' -f 4)
fi

PACKAGE_FILE_NAME="1panel-${PANEL_VERSION}-linux-${architecture}.tar.gz"
PACKAGE_URL="${PANEL_DOWNLOAD_BASE}/${PANEL_REPO}/releases/download/${PANEL_VERSION}/${PACKAGE_FILE_NAME}"
CHECKSUM_URL="${PANEL_DOWNLOAD_BASE}/${PANEL_REPO}/releases/download/${PANEL_VERSION}/checksums.txt"

# 2. 本地已有安装包时优先复用
if [[ -f "$PACKAGE_FILE_NAME" ]]; then
    if [[ -n "$PANEL_VERSION" ]]; then
        EXPECTED_HASH=$(curl -sfL --max-time 20 "$CHECKSUM_URL" | grep "$PACKAGE_FILE_NAME" | awk '{print $1}')
        ACTUAL_HASH=$(sha256sum "$PACKAGE_FILE_NAME" | awk '{print $1}')
        if [[ "$EXPECTED_HASH" == "$ACTUAL_HASH" ]]; then
            if is_zh; then
                echo "检测到本地安装包且校验通过，跳过下载。"
            else
                echo "Local package found and verified, skipping download."
            fi
        else
            rm -f "$PACKAGE_FILE_NAME"
        fi
    else
        rm -f "$PACKAGE_FILE_NAME"
    fi
fi

# 3. 下载
if [[ ! -f "$PACKAGE_FILE_NAME" ]]; then
    if [[ -z "$PANEL_VERSION" ]]; then
        die "获取最新版本失败，请稍后重试；也可通过 PANEL_VERSION 指定版本后重试。" \
            "Failed to resolve the latest version, please retry later or set PANEL_VERSION."
    fi
    if is_zh; then
        echo "准备下载 1Panel ${PANEL_VERSION}（架构：${architecture}）。"
        echo "下载地址：${PACKAGE_URL}"
    else
        echo "Downloading 1Panel ${PANEL_VERSION} (arch: ${architecture})."
        echo "Download URL: ${PACKAGE_URL}"
    fi
    if ! curl -fL --retry 3 --connect-timeout 20 -o "$PACKAGE_FILE_NAME" "$PACKAGE_URL"; then
        rm -f "$PACKAGE_FILE_NAME"
        die "下载安装包失败，请检查网络连接后重试；国内网络可设置 PANEL_DOWNLOAD_BASE 使用代理镜像。" \
            "Failed to download the package. Set PANEL_DOWNLOAD_BASE to use a mirror."
    fi
fi

# 4. 校验
if [[ -n "$PANEL_VERSION" ]]; then
    EXPECTED_HASH=$(curl -sfL --max-time 20 "$CHECKSUM_URL" | grep "$PACKAGE_FILE_NAME" | awk '{print $1}')
    if [[ -n "$EXPECTED_HASH" ]]; then
        ACTUAL_HASH=$(sha256sum "$PACKAGE_FILE_NAME" | awk '{print $1}')
        if [[ "$EXPECTED_HASH" != "$ACTUAL_HASH" ]]; then
            rm -f "$PACKAGE_FILE_NAME"
            die "安装包校验失败，下载文件可能不完整或已损坏。" "Package checksum verification failed."
        fi
        if is_zh; then
            echo "安装包校验通过。"
        else
            echo "Package checksum verified."
        fi
    else
        if is_zh; then
            echo "未获取到校验文件，跳过校验。"
        else
            echo "No checksum file found, skipping verification."
        fi
    fi
fi

# 5. 解压并安装
PACKAGE_DIR="1panel-${PANEL_VERSION}-linux-${architecture}"
rm -rf "$PACKAGE_DIR"
if ! tar zxf "$PACKAGE_FILE_NAME"; then
    rm -rf "$PACKAGE_DIR" "$PACKAGE_FILE_NAME"
    die "解压安装包失败，下载文件可能不完整或已损坏。" "Failed to extract the package."
fi
if [[ ! -f "${PACKAGE_DIR}/install.sh" ]]; then
    rm -rf "$PACKAGE_DIR" "$PACKAGE_FILE_NAME"
    die "安装包内容不完整，缺少 install.sh。" "The package is incomplete, install.sh is missing."
fi

cd "$PACKAGE_DIR" || exit 1
PANEL_BASE_DIR="$PANEL_BASE_DIR" PANEL_PORT="$PANEL_PORT" PANEL_LANGUAGE="$PANEL_LANGUAGE" \
    /bin/bash install.sh "$@"
