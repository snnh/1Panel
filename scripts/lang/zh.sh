#!/bin/bash
# 1Panel 安装脚本与 1pctl 命令行文案（简体中文）
# 该文件由 scripts/lang/zh.sh 提供，安装时复制到 /usr/local/bin/lang/zh.sh

MSG_NEED_ROOT="请使用 root 用户执行该命令！"
MSG_NO_SYSTEMD="未检测到 systemd，1Panel v2 需要 systemd 管理 core 与 agent 服务。"
MSG_ARCH_UNSUPPORTED="当前系统架构暂不支持，请参考文档选择受支持的系统与架构。"
MSG_OS_UNSUPPORTED="当前系统暂不支持，1Panel v2 需要基于 systemd 的 Linux 发行版。"
MSG_PKG_NOT_FOUND="未找到 1Panel 安装包，请先执行 make package_linux 打包，或使用 quick_start.sh 在联网环境安装。"
MSG_INSTALL_START="开始安装 1Panel。"
MSG_INSTALL_BIN="安装可执行文件"
MSG_INSTALL_CTL="安装 1pctl 命令行工具"
MSG_INSTALL_UNIT="安装 systemd 服务"
MSG_INSTALL_LANG="安装命令行语言文件"
MSG_INSTALL_START_SVC="启动 1Panel 服务"
MSG_INSTALL_DONE="1Panel 安装完成。"
MSG_UNINSTALL_DONE="1Panel 已卸载（数据目录保留，如需彻底删除请手动清理）。"
MSG_NO_INSTALL="未检测到已安装的 1Panel。"
MSG_1PCTL_USAGE="用法: 1pctl <command> [options]

命令:
  version              显示当前版本
  user-info            显示面板登录信息
  user-list            显示面板用户列表
  update username      修改面板用户名
  update password      修改面板密码
  update port          修改面板端口
  listen-ip ipv4|ipv6  显示面板监听地址
  reset mfa            关闭 MFA 二次验证
  reset https          关闭 HTTPS
  reset entrance       关闭安全入口
  reset ips            关闭授权 IP 限制
  reset domain         关闭访问域名绑定
  reset passkey        关闭 Passkey 登录
  restore              回滚到上一次升级前的版本
  app init             初始化本地应用仓库
  help                 显示该帮助"
