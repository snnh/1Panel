package version

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/1Panel-dev/1Panel/agent/app/dto"
	"github.com/1Panel-dev/1Panel/agent/app/model"
	"github.com/1Panel-dev/1Panel/agent/constant"
	"github.com/1Panel-dev/1Panel/agent/global"
	"github.com/1Panel-dev/1Panel/agent/utils/common"
)

func GetUpgradeVersionInfo() (*dto.UpgradeInfo, error) {
	var upgrade dto.UpgradeInfo
	var currentVersion model.Setting
	if err := global.CoreDB.Model(&model.Setting{}).Where("key = ?", "SystemVersion").First(&currentVersion).Error; err != nil {
		global.LOG.Errorf("load %s from db setting failed, err: %v", "SystemVersion", err)
		return nil, err
	}
	var developerMode model.Setting
	if err := global.CoreDB.Model(&model.Setting{}).Where("key = ?", "DeveloperMode").First(&developerMode).Error; err != nil {
		global.LOG.Errorf("load %s from db setting failed, err: %v", "DeveloperMode", err)
		return nil, err
	}

	upgrade.TestVersion, upgrade.NewVersion, upgrade.LatestVersion = loadVersionByMode(developerMode.Value, currentVersion.Value)
	var itemVersion string
	if len(upgrade.NewVersion) != 0 {
		itemVersion = upgrade.NewVersion
	}
	if (global.CONF.Base.Mode == "dev" || developerMode.Value == constant.StatusEnable) && len(upgrade.TestVersion) != 0 {
		itemVersion = upgrade.TestVersion
	}
	if len(upgrade.LatestVersion) != 0 {
		itemVersion = upgrade.LatestVersion
	}
	if len(itemVersion) == 0 {
		return &upgrade, nil
	}
	notes, err := loadReleaseNotes(global.ReleaseAssetURL(itemVersion, fmt.Sprintf("1panel-%s-release-notes", itemVersion)))
	if err != nil {
		return nil, fmt.Errorf("load releases-notes of version %s failed, err: %v", itemVersion, err)
	}
	upgrade.ReleaseNote = notes
	return &upgrade, nil
}

func loadVersionByMode(developer, currentVersion string) (string, string, string) {
	var current, latest string
	if global.CONF.Base.Mode == "dev" {
		betaVersionLatest := loadVersion(true, currentVersion, "beta")
		devVersionLatest := loadVersion(true, currentVersion, "dev")
		if common.ComparePanelVersion(betaVersionLatest, devVersionLatest) {
			return betaVersionLatest, "", ""
		}
		return devVersionLatest, "", ""
	}

	betaVersionLatest := ""
	latest = loadVersion(true, currentVersion, "stable")
	current = loadVersion(false, currentVersion, "stable")
	if developer == constant.StatusEnable {
		betaVersionLatest = loadVersion(true, currentVersion, "beta")
	}
	if current != latest {
		return betaVersionLatest, current, latest
	}

	versionPart := strings.Split(current, ".")
	if len(versionPart) < 3 {
		return betaVersionLatest, current, latest
	}
	num, _ := strconv.Atoi(versionPart[1])
	if num == 0 {
		return betaVersionLatest, current, latest
	}
	if num >= 10 {
		if current[:6] == currentVersion[:6] {
			return betaVersionLatest, current, ""
		}
		return betaVersionLatest, "", latest
	}
	if current[:5] == currentVersion[:5] {
		return betaVersionLatest, current, ""
	}
	return betaVersionLatest, "", latest
}

// loadVersion 从社区仓库的 GitHub Releases 查询最新版本。
// isLatest 参数已无实际意义（社区版只有 stable/beta/dev 三个发布通道），
// 保留是为了兼容原有调用方；GitHub Releases 不存在「不同大版本的 LTS 线」，
// 因此两种查询都返回该通道下的最新版本。
func loadVersion(_ bool, currentVersion, mode string) string {
	version, err := loadLatestVersion(mode)
	if err != nil {
		global.LOG.Errorf("load latest version from github release failed (channel: %s), err: %v", mode, err)
		return ""
	}
	if len(version) == 0 {
		return ""
	}
	return checkVersion(version, currentVersion)
}

// loadLatestVersion 查询指定通道的最新版本号。
// stable 取最新正式版，beta 与 dev 取最新预发布版本。
func loadLatestVersion(mode string) (string, error) {
	isPrerelease := mode == "beta" || mode == "dev"
	_, body, err := HandleRequest(global.LatestReleaseURL(mode), http.MethodGet, constant.TimeOut20s)
	if err != nil {
		return "", err
	}
	if isPrerelease {
		var releases []githubRelease
		if err := json.Unmarshal(body, &releases); err != nil {
			return "", fmt.Errorf("unmarshal github releases failed, err: %v", err)
		}
		for _, item := range releases {
			if item.Prerelease {
				return item.TagName, nil
			}
		}
		return "", nil
	}

	var release githubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return "", fmt.Errorf("unmarshal github release failed, err: %v", err)
	}
	return release.TagName, nil
}

type githubRelease struct {
	TagName    string `json:"tag_name"`
	Prerelease bool   `json:"prerelease"`
}

func checkVersion(v2, v1 string) string {
	addSuffix := false
	if !strings.Contains(v1, "-") {
		v1 = v1 + "-lts"
	}
	if !strings.Contains(v2, "-") {
		addSuffix = true
		v2 = v2 + "-lts"
	}
	if common.ComparePanelVersion(v2, v1) {
		if addSuffix {
			return strings.TrimSuffix(v2, "-lts")
		}
		return v2
	}
	return ""
}

func loadReleaseNotes(path string) (string, error) {
	_, releaseNotes, err := HandleRequest(path, http.MethodGet, constant.TimeOut20s)
	if err != nil {
		return "", err
	}
	return string(releaseNotes), nil
}

func HandleRequest(url, method string, timeout int) (int, []byte, error) {
	defer func() {
		if r := recover(); r != nil {
			global.LOG.Errorf("handle request failed, error message: %v", r)
			return
		}
	}()

	transport := loadRequestTransport()
	client := http.Client{Timeout: time.Duration(timeout) * time.Second, Transport: transport}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return 0, nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, nil, errors.New(resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, body, nil
}

func loadRequestTransport() *http.Transport {
	return &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		IdleConnTimeout:       15 * time.Second,
	}
}
