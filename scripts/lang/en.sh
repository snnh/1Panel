#!/bin/bash
# Messages for the install script and the 1pctl command line tool (English)
# Provided by scripts/lang/en.sh, copied to /usr/local/bin/lang/en.sh on install

MSG_NEED_ROOT="Please run this command as the root user!"
MSG_NO_SYSTEMD="systemd was not detected. 1Panel v2 requires systemd to manage the core and agent services."
MSG_ARCH_UNSUPPORTED="The current system architecture is not supported, please refer to the documentation for supported systems and architectures."
MSG_OS_UNSUPPORTED="The current system is not supported, 1Panel v2 requires a systemd based Linux distribution."
MSG_PKG_NOT_FOUND="No 1Panel package found. Run 'make package_linux' first, or use quick_start.sh on a machine with internet access."
MSG_INSTALL_START="Installing 1Panel."
MSG_INSTALL_BIN="Installing executables"
MSG_INSTALL_CTL="Installing the 1pctl command line tool"
MSG_INSTALL_UNIT="Installing systemd services"
MSG_INSTALL_LANG="Installing command line language files"
MSG_INSTALL_START_SVC="Starting 1Panel services"
MSG_INSTALL_DONE="1Panel has been installed."
MSG_UNINSTALL_DONE="1Panel has been uninstalled (the data directory is kept, remove it manually if needed)."
MSG_NO_INSTALL="No installed 1Panel was detected."
MSG_1PCTL_USAGE="Usage: 1pctl <command> [options]

Commands:
  version              Show the current version
  user-info            Show the panel login information
  user-list            Show the panel user list
  update username      Change the panel username
  update password      Change the panel password
  update port          Change the panel port
  listen-ip ipv4|ipv6  Show the panel listen address
  reset mfa            Disable MFA two-step verification
  reset https          Disable HTTPS
  reset entrance       Disable the security entrance
  reset ips            Disable the authorized IP restriction
  reset domain         Disable the bound access domain
  reset passkey        Disable Passkey login
  restore             Roll back to the version before the last upgrade
  app init             Initialize the local app repository
  help                 Show this help"
