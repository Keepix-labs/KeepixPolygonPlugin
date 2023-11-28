package tasks

import (
	"KeepixPlugin/appstate"
	"KeepixPlugin/utils"
)

// RequirementFunc defines the type for requirement validation functions
type RequirementFunc func() bool

// SystemRequirements maps requirement names to their validation functions
var SystemRequirements = map[string]RequirementFunc{
	"docker":      CheckDockerExists,
	"linux":       CheckLinux,
	"windows":     CheckWindows,
	"osx":         CheckOSX,
	"wsl2":        CheckWSL2,
	"installed":   CheckInstalled,
	"uninstalled": CheckUninstalled,
	"cpu4":        CheckHas4CPU,
}

func CheckDockerExists() bool {
	return utils.CheckDockerExists()
}

func CheckLinux() bool {
	return utils.CheckOSType() == "linux"
}

func CheckWindows() bool {
	return utils.CheckOSType() == "windows"
}

func CheckOSX() bool {
	return utils.CheckOSType() == "osx"
}

func CheckWSL2() bool {
	return utils.CheckWSL2()
}

func CheckInstalled() bool {
	return appstate.CurrentState > appstate.ConfiguringNode
}

func CheckUninstalled() bool {
	return appstate.CurrentState == appstate.NoState
}

func CheckHas4CPU() bool {
	return utils.CheckCPUCount() >= 4
}
