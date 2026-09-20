# 1Panel Community Edition — Installation Guide

The 1Panel community edition is a trimmed build of 1Panel v2: the AI, license/Pro and multi-node features are removed, and the UI keeps Chinese and English only. Installation packages are published through the **GitHub Releases** of the community repository. This document covers online installation, offline installation, upgrading and uninstalling.

> Repository: [https://github.com/snnh/1Panel](https://github.com/snnh/1Panel) (branch `community-dev`) Releases: [https://github.com/snnh/1Panel/releases](https://github.com/snnh/1Panel/releases)

## 1. Requirements

| Item      | Requirement                                                                           |
| --------- | ------------------------------------------------------------------------------------- |
| OS        | systemd based Linux distributions (Ubuntu / Debian / CentOS / Rocky / openEuler, ...) |
| Arch      | x86_64, arm64, armv7, ppc64le, s390x, riscv64, loongarch64                            |
| Privilege | root, or an account with sudo                                                         |
| Tools     | `curl` and `tar` (required by the install scripts); Docker for container features     |

Layout after installation:

| Path                                              | Description                                      |
| ------------------------------------------------- | ------------------------------------------------ |
| `/usr/local/bin/1panel-core`                      | panel core process                               |
| `/usr/local/bin/1panel-agent`                     | node agent process                               |
| `/usr/local/bin/1pctl`                            | command line tool                                |
| `/usr/local/bin/lang/`                            | command line language files                      |
| `/opt/1panel/`                                    | data directory (databases, logs, backups, GeoIP) |
| `/etc/systemd/system/1panel-{core,agent}.service` | systemd unit files                               |

## 2. Online installation (one-liner)

```bash
bash -c "$(curl -sSL https://raw.githubusercontent.com/snnh/1Panel/community-dev/quick_start.sh)"
```

The script detects the architecture, resolves the latest stable release, downloads and verifies the package with sha256, extracts it, runs the bundled `install.sh` and finally starts the services and prints the access URL, username and password. Keep them safe.

Environment variables:

| Variable              | Default              | Description                                                   |
| --------------------- | -------------------- | ------------------------------------------------------------- |
| `PANEL_REPO`          | `snnh/1Panel`        | release repository                                            |
| `PANEL_VERSION`       | latest stable        | pin a version, e.g. `v2.0.0`                                  |
| `PANEL_DOWNLOAD_BASE` | `https://github.com` | download prefix, e.g. `https://ghfast.top/https://github.com` |
| `PANEL_BASE_DIR`      | `/opt`               | install directory                                             |
| `PANEL_PORT`          | `9999`               | panel port                                                    |
| `PANEL_LANGUAGE`      | `zh`                 | panel language, `zh` or `en`                                  |

Example (mirror and custom port):

```bash
PANEL_DOWNLOAD_BASE="https://ghfast.top/https://github.com" PANEL_PORT=8888 \
  bash -c "$(curl -sSL https://raw.githubusercontent.com/snnh/1Panel/community-dev/quick_start.sh)"
```

## 3. Offline installation

Download the package on a machine with internet access and copy it to the target server:

```bash
# on a connected machine (replace the version accordingly)
curl -fLO https://github.com/snnh/1Panel/releases/download/v2.0.0/1panel-v2.0.0-linux-amd64.tar.gz
curl -fLO https://github.com/snnh/1Panel/releases/download/v2.0.0/checksums.txt
```

Upload `1panel-v2.0.0-linux-amd64.tar.gz` to the server and run:

```bash
tar zxf 1panel-v2.0.0-linux-amd64.tar.gz
cd 1panel-v2.0.0-linux-amd64
bash install.sh
```

You can also build the package locally (Go 1.26+ and Node.js required):

```bash
make package_linux                     # build/1panel-<version>-linux-<arch>.tar.gz
make release_notes                     # release notes shipped with a release
```

## 4. Upgrading

The built-in updater queries the community GitHub Releases for new versions and performs the upgrade. The previous `1panel-core`, `1panel-agent`, `1pctl` and language files are backed up to `/opt/1panel/tmp/upgrade/<version>/original`, and a failed upgrade can be rolled back from the operation log.

On offline machines simply re-run `bash install.sh` from a newer package; the existing credentials are reused.

## 5. Uninstalling

```bash
cd 1panel-*-linux-*   # the extracted package directory
bash install.sh uninstall
```

This stops and disables the services and removes the files under `/usr/local/bin` and the systemd units, but **keeps** the `/opt/1panel` data directory. Run `rm -rf /opt/1panel` manually only after you have backed up your data.

## 6. Command line tool 1pctl

```bash
1pctl version            # show the version
1pctl user-info          # show the access URL, username and password
1pctl user-list          # list the panel users
1pctl update username    # change the username
1pctl update password    # change the password
1pctl update port        # change the port
1pctl listen-ip ipv4     # show the IPv4 listen address
1pctl reset mfa          # disable MFA
1pctl reset https        # disable HTTPS
1pctl reset entrance     # disable the security entrance
1pctl reset ips          # disable the authorized IP restriction
1pctl reset domain       # disable the bound access domain
1pctl reset passkey      # disable Passkey
1pctl restore            # roll back to the version before the last upgrade
1pctl app init           # initialize the local app repository
```

## 7. Content sources

| Content | Source |
| --- | --- |
| Installation/upgrade packages, release notes | community GitHub Releases |
| `lang.tar.gz`, `GeoIP.mmdb` | latest assets of the community GitHub Releases |
| App store and script library catalogs | the upstream published catalogs (`apps-assets.fit2cloud.com`, `resource.fit2cloud.com`) |

To use your own mirror, override the following environment variables (persist them with `Environment=` entries in the systemd units):

| Variable                               | Description                                                          |
| -------------------------------------- | -------------------------------------------------------------------- |
| `PANEL_REPO_OWNER` / `PANEL_REPO_NAME` | release repository                                                   |
| `PANEL_RELEASE_API_BASE`               | version query API (GitHub API or a mirror)                           |
| `PANEL_RELEASE_DOWNLOAD_BASE`          | package download prefix                                              |
| `PANEL_MODE`                           | release channel: `stable` / `beta` / `dev`, written by the installer |

Example:

```ini
# /etc/systemd/system/1panel-core.service
[Service]
Environment=PANEL_RELEASE_API_BASE=https://mirror.example.com/api/repos/snnh/1Panel
Environment=PANEL_RELEASE_DOWNLOAD_BASE=https://mirror.example.com/snnh/1Panel/releases/download
```

Then run `systemctl daemon-reload && systemctl restart 1panel-core 1panel-agent`.
