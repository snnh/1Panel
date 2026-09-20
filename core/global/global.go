package global

import (
	"fmt"
	"os"
	"strings"

	"github.com/1Panel-dev/1Panel/core/init/auth"
	"github.com/1Panel-dev/1Panel/core/init/session/psession"
	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	DB      *gorm.DB
	AlertDB *gorm.DB
	TaskDB  *gorm.DB
	AgentDB *gorm.DB
	LOG     *logrus.Logger
	CONF    ServerConfig
	VALID   *validator.Validate
	SESSION *psession.PSession
	Viper   *viper.Viper

	I18n       *i18n.Localizer
	I18nForCmd *i18n.Localizer

	Cron *cron.Cron

	ScriptSyncJobID cron.EntryID

	IPTracker *auth.IPTracker
)

type DBOption func(*gorm.DB) *gorm.DB

const (
	// DefaultRepoOwner 与 DefaultRepoName 是社区版默认的发布仓库，
	// 安装、升级与运行期资源均从该仓库的 GitHub Releases 获取。
	// 自建镜像或 fork 可通过环境变量覆盖：
	//   PANEL_REPO_OWNER / PANEL_REPO_NAME            发布仓库
	//   PANEL_RELEASE_API_BASE                        版本查询接口（可指向镜像）
	//   PANEL_RELEASE_DOWNLOAD_BASE                   安装包下载前缀（可指向镜像）
	DefaultRepoOwner = "snnh"
	DefaultRepoName  = "1Panel"
)

func repoOwner() string {
	if owner := strings.TrimSpace(os.Getenv("PANEL_REPO_OWNER")); owner != "" {
		return owner
	}
	return DefaultRepoOwner
}

func repoName() string {
	if name := strings.TrimSpace(os.Getenv("PANEL_REPO_NAME")); name != "" {
		return name
	}
	return DefaultRepoName
}

// ReleaseAPIBase 返回 GitHub Releases API 地址，用于查询最新版本。
func ReleaseAPIBase() string {
	if base := strings.TrimSpace(os.Getenv("PANEL_RELEASE_API_BASE")); base != "" {
		return strings.TrimRight(base, "/")
	}
	return "https://api.github.com/repos/" + repoOwner() + "/" + repoName()
}

// ReleaseDownloadBase 返回安装包的下载前缀。
func ReleaseDownloadBase() string {
	if base := strings.TrimSpace(os.Getenv("PANEL_RELEASE_DOWNLOAD_BASE")); base != "" {
		return strings.TrimRight(base, "/")
	}
	return "https://github.com/" + repoOwner() + "/" + repoName() + "/releases/download"
}

// ReleaseAssetURL 返回某个版本的附件地址，
// 例如 https://github.com/1Panel-dev/1Panel/releases/download/v2.0.0/1panel-v2.0.0-linux-amd64.tar.gz
func ReleaseAssetURL(version, fileName string) string {
	return fmt.Sprintf("%s/%s/%s", ReleaseDownloadBase(), version, fileName)
}

// LatestReleaseURL 返回查询最新版本的地址。
// stable 只取最新正式版；beta 与 dev 取最新预发布版本。
func LatestReleaseURL(mode string) string {
	switch mode {
	case "beta", "dev":
		return ReleaseAPIBase() + "/releases?per_page=20"
	default:
		return ReleaseAPIBase() + "/releases/latest"
	}
}

// LatestReleaseAssetURL 返回最新版本附件的地址，
// 用于 lang.tar.gz、GeoIP.mmdb 等跟随发布走、但与具体版本号无关的资源。
func LatestReleaseAssetURL(fileName string) string {
	if base := strings.TrimSpace(os.Getenv("PANEL_RELEASE_DOWNLOAD_BASE")); base != "" {
		return fmt.Sprintf("%s/latest/download/%s", strings.TrimRight(base, "/"), fileName)
	}
	return fmt.Sprintf("https://github.com/%s/%s/releases/latest/download/%s", repoOwner(), repoName(), fileName)
}

// ResourceURL 返回脚本库等内容目录的地址（社区版沿用上游已发布的内容目录）。
func ResourceURL() string {
	return "https://resource.fit2cloud.com/1panel/resource/v2"
}

// AppRepoURL 返回应用商店内容目录的地址（社区版沿用上游已发布的应用内容）。
func AppRepoURL() string {
	return "https://apps-assets.fit2cloud.com"
}
