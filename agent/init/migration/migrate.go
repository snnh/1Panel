package migration

import (
	"github.com/1Panel-dev/1Panel/agent/global"
	"github.com/1Panel-dev/1Panel/agent/init/migration/migrations"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Init() {
	InitAgentDB()
	InitTaskDB()
	InitAlertDB()
	global.LOG.Info("Migration run successfully")
}

func InitAgentDB() {
	m := gormigrate.New(global.DB, gormigrate.DefaultOptions, agentDBMigrations())
	if err := m.Migrate(); err != nil {
		global.LOG.Error(err)
		panic(err)
	}
}

func agentDBMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		migrations.AddTable,
		migrations.AddMonitorTable,
		migrations.InitSetting,
		migrations.InitImageRepo,
		migrations.InitDefaultCA,
		migrations.InitPHPExtensions,
		migrations.InitBackup,
		migrations.InitDefault,
		migrations.UpdateWebsiteExpireDate,
		migrations.UpdateRuntime,
		migrations.AddSnapshotRule,
		migrations.UpdatePHPRuntime,
		migrations.AddSnapshotIgnore,
		migrations.InitAppLauncher,
		migrations.AddTableAlert,
		migrations.InitAlertConfig,
		migrations.AddMethodToAlertLog,
		migrations.AddMethodToAlertTask,
		migrations.InitCronjobGroup,
		migrations.AddColumnToAlert,
		migrations.UpdateWebsiteSSL,
		migrations.AddQuickJump,
		migrations.InitLocalSSHConn,
		migrations.InitLocalSSHShow,
		migrations.InitRecordStatus,
		migrations.AddShowNameForQuickJump,
		migrations.AddTimeoutForClam,
		migrations.UpdateCronjobSpec,
		migrations.UpdateWebsiteSSLAddColumn,
		migrations.UpdateMonitorInterval,
		migrations.AddMonitorProcess,
		migrations.UpdateCronJob,
		migrations.AddCommonDescription,
		migrations.UpdateDatabase,
		migrations.UpdateDatabaseMysql,
		migrations.AddDatabaseMongodb,
		migrations.InitIptablesStatus,
		migrations.UpdateWebsite,
		migrations.AddisIPtoWebsiteSSL,
		migrations.InitPingStatus,
		migrations.UpdateApp,
		migrations.AddCronjobArgs,
		migrations.AddWebsiteAcmeAccountColumn,
		migrations.AddAppInstallSortOrder,
		migrations.AddEditionSetting,
		migrations.AddHostTable,
		migrations.AddFileShareTable,
		migrations.AddFileHistoryTable,
		migrations.MigrateLegoV5,
		migrations.InitFirewallPortWhiteList,
		migrations.AddDatabaseUserTable,
		migrations.AddBackupRecordArgs,
		migrations.AddFtpIdentity,
		migrations.AddWebsiteTemplateTable,
		migrations.AddComposePinned,
		migrations.AddFirewallRuleTable,
		migrations.InitDockerPortGuardStatus,
		migrations.NormalizeFirewallBackendSelections,
		migrations.SimplifyFirewallRulePolicy,
		migrations.AddDockerPortGuardReadOnly,
		migrations.MigrateFirewallPortWhitelistSources,
		migrations.RemoveNodeScopeSettings,
		migrations.RemoveAIQuickJump,
	}
}

func InitTaskDB() {
	m := gormigrate.New(global.TaskDB, gormigrate.DefaultOptions, []*gormigrate.Migration{
		migrations.AddTaskTable,
	})
	if err := m.Migrate(); err != nil {
		global.LOG.Error(err)
		panic(err)
	}
}

func InitAlertDB() {
	if err := migrateAlertDB(global.AlertDB); err != nil {
		global.LOG.Error(err)
		panic(err)
	}
}

func migrateAlertDB(db *gorm.DB) error {
	options := *gormigrate.DefaultOptions
	options.UseTransaction = true
	m := gormigrate.New(db, &options, []*gormigrate.Migration{
		migrations.AddAlertConfigUIDAndSecret,
		migrations.MigrateAlertMethodConfigIDs,
		migrations.MigrateAlertLogTaskMethodConfigIDs,
		migrations.AddAlertAuditUser,
		migrations.AddAlertTaskDeliveryLogID,
	})
	return m.Migrate()
}
