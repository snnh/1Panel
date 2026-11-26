package service

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/1Panel-dev/1Panel/agent/app/model"
	"github.com/1Panel-dev/1Panel/agent/constant"
	"github.com/1Panel-dev/1Panel/agent/global"
	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func appCacheTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "cache.db")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.App{}, &model.AppDetail{}, &model.AppTag{}, &model.Setting{}); err != nil {
		t.Fatal(err)
	}
	oldDB, oldDir, oldLog := global.DB, global.Dir, global.LOG
	global.DB = db
	global.Dir = global.SystemDir{AppResourceDir: t.TempDir()}
	global.Dir.RemoteAppResourceDir = filepath.Join(global.Dir.AppResourceDir, constant.AppResourceRemote)
	global.Dir.LocalAppResourceDir = filepath.Join(global.Dir.AppResourceDir, constant.AppResourceLocal)
	global.Dir.CustomAppResourceDir = filepath.Join(global.Dir.AppResourceDir, constant.AppResourceCustom)
	global.LOG = logrus.New()
	global.LOG.SetOutput(io.Discard)
	t.Cleanup(func() {
		global.DB, global.Dir, global.LOG = oldDB, oldDir, oldLog
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	})
	if err := db.Create(&model.Setting{Key: "AppStoreLastModified", Value: "42"}).Error; err != nil {
		t.Fatal(err)
	}
	return db
}

func appCacheTestWrite(t *testing.T, name, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(name, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
