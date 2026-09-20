# 1Panel 社区版安装文档

1Panel 社区版是 1Panel v2 的精简版本：移除了 AI、授权/Pro、多节点等商业化内容，界面语言只保留中文与英文。安装包通过社区仓库的 **GitHub Releases** 发布，本文介绍在线安装、离线安装、升级与卸载。

> 仓库：[https://github.com/snnh/1Panel](https://github.com/snnh/1Panel)（分支 `community-dev`）Releases：[https://github.com/snnh/1Panel/releases](https://github.com/snnh/1Panel/releases)

## 1. 环境要求

| 项目 | 要求                                                                            |
| ---- | ------------------------------------------------------------------------------- |
| 系统 | 基于 systemd 的 Linux 发行版（Ubuntu / Debian / CentOS / Rocky / openEuler 等） |
| 架构 | x86_64、arm64、armv7、ppc64le、s390x、riscv64、loongarch64                      |
| 权限 | root 用户，或具备 sudo 的账号                                                   |
| 软件 | `curl`、`tar`（安装脚本依赖）；容器相关功能需要 Docker                          |

安装后的目录结构：

| 路径                                              | 说明                                  |
| ------------------------------------------------- | ------------------------------------- |
| `/usr/local/bin/1panel-core`                      | 面板核心进程                          |
| `/usr/local/bin/1panel-agent`                     | 节点代理进程                          |
| `/usr/local/bin/1pctl`                            | 命令行工具                            |
| `/usr/local/bin/lang/`                            | 命令行语言文件                        |
| `/opt/1panel/`                                    | 数据目录（数据库、日志、备份、GeoIP） |
| `/etc/systemd/system/1panel-{core,agent}.service` | systemd 服务单元                      |

## 2. 在线安装（一键脚本）

```bash
bash -c "$(curl -sSL https://raw.githubusercontent.com/snnh/1Panel/community-dev/quick_start.sh)"
```

脚本会依次完成：识别架构 → 获取最新正式版 → 下载安装包并校验 sha256 → 解压 → 调用包内 `install.sh` → 启动服务，最后打印访问地址、用户名与密码，请务必保存。

可用的环境变量：

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `PANEL_REPO` | `snnh/1Panel` | 发布仓库 |
| `PANEL_VERSION` | 最新正式版 | 指定版本，例如 `v2.0.0` |
| `PANEL_DOWNLOAD_BASE` | `https://github.com` | 下载前缀，可改为代理镜像，例如 `https://ghfast.top/https://github.com` |
| `PANEL_BASE_DIR` | `/opt` | 安装目录 |
| `PANEL_PORT` | `9999` | 面板端口 |
| `PANEL_LANGUAGE` | `zh` | 面板语言，`zh` 或 `en` |

示例（使用代理镜像并指定端口）：

```bash
PANEL_DOWNLOAD_BASE="https://ghfast.top/https://github.com" PANEL_PORT=8888 \
  bash -c "$(curl -sSL https://raw.githubusercontent.com/snnh/1Panel/community-dev/quick_start.sh)"
```

## 3. 离线安装

在没有外网或不允许访问 GitHub 的机器上，先在有网络的机器上下载安装包并拷贝过来：

```bash
# 有网络的机器（version 请替换为实际版本号）
curl -fLO https://github.com/snnh/1Panel/releases/download/v2.0.0/1panel-v2.0.0-linux-amd64.tar.gz
curl -fLO https://github.com/snnh/1Panel/releases/download/v2.0.0/checksums.txt
```

将 `1panel-v2.0.0-linux-amd64.tar.gz` 上传到目标服务器后执行：

```bash
tar zxf 1panel-v2.0.0-linux-amd64.tar.gz
cd 1panel-v2.0.0-linux-amd64
bash install.sh
```

也可以自己在仓库内构建发布包（需要 Go 1.26+ 与 Node.js）：

```bash
make package_linux                     # 生成 build/1panel-<版本>-linux-<架构>.tar.gz
make release_notes                    # 生成随发布使用的 release-notes 文件
```

## 4. 升级

面板内的「系统更新」会自动从社区仓库的 GitHub Releases 查询新版本并完成升级，升级前会自动备份 `1panel-core`、`1panel-agent`、`1pctl`、语言文件到 `/opt/1panel/tmp/upgrade/<版本>/original`，失败时可在「操作记录」中一键回滚。

离线环境可以在目标机器上重新执行第 3 节的 `bash install.sh`，脚本会沿用已安装的账号信息。

## 5. 卸载

```bash
cd 1panel-*-linux-*   # 解压后的发布包目录
bash install.sh uninstall
```

卸载会停止并禁用服务、删除 `/usr/local/bin` 下的程序与 `/etc/systemd/system` 下的单元文件，但**保留** `/opt/1panel` 数据目录；如需彻底删除，请在确认备份后手动执行 `rm -rf /opt/1panel`。

## 6. 命令行工具 1pctl

```bash
1pctl version            # 查看版本
1pctl user-info          # 查看访问地址、用户名与密码
1pctl user-list          # 查看用户列表
1pctl update username    # 修改用户名
1pctl update password    # 修改密码
1pctl update port        # 修改端口
1pctl listen-ip ipv4     # 查看 IPv4 监听地址
1pctl reset mfa          # 关闭 MFA
1pctl reset https        # 关闭 HTTPS
1pctl reset entrance     # 关闭安全入口
1pctl reset ips          # 关闭授权 IP
1pctl reset domain       # 关闭访问域名
1pctl reset passkey      # 关闭 Passkey
1pctl restore            # 回滚到上一次升级前
1pctl app init           # 初始化本地应用仓库
```

## 7. 内容来源说明

| 内容                        | 来源                                                                              |
| --------------------------- | --------------------------------------------------------------------------------- |
| 安装包、升级包、更新日志    | 社区仓库 GitHub Releases                                                          |
| `lang.tar.gz`、`GeoIP.mmdb` | 社区仓库 GitHub Releases（最新版附件）                                            |
| 应用商店、脚本库目录        | 沿用上游已发布的内容目录（`apps-assets.fit2cloud.com`、`resource.fit2cloud.com`） |

如需自建镜像，可覆盖下列环境变量（写入 systemd 单元的 `Environment=` 即可长期生效）：

| 变量                                   | 说明                                                        |
| -------------------------------------- | ----------------------------------------------------------- |
| `PANEL_REPO_OWNER` / `PANEL_REPO_NAME` | 发布仓库                                                    |
| `PANEL_RELEASE_API_BASE`               | 版本查询接口（GitHub API 或镜像）                           |
| `PANEL_RELEASE_DOWNLOAD_BASE`          | 安装包下载前缀                                              |
| `PANEL_MODE`                           | 发布通道，`stable` / `beta` / `dev`，由安装脚本写入单元文件 |

例如把版本查询与下载都指向自建镜像：

```ini
# /etc/systemd/system/1panel-core.service
[Service]
Environment=PANEL_RELEASE_API_BASE=https://mirror.example.com/api/repos/snnh/1Panel
Environment=PANEL_RELEASE_DOWNLOAD_BASE=https://mirror.example.com/snnh/1Panel/releases/download
```

修改后执行 `systemctl daemon-reload && systemctl restart 1panel-core 1panel-agent`。
