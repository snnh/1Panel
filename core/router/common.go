package router

func commonGroups() []CommonRouter {
	return []CommonRouter{
		&BaseRouter{},
		&BackupRouter{},
		&LogRouter{},
		&SettingRouter{},
		&CommandRouter{},
		&GroupRouter{},
		&ScriptRouter{},
	}
}

// RouterGroupApp is the community router group, registered by core/init/router.
var RouterGroupApp = commonGroups()
