package tasks

// TaskFunc defines the signature for task functions
type TaskFunc func(map[string]string) string

// TaskMap maps task names to their corresponding functions
var TaskMap = map[string]TaskFunc{
	"install":    installTask,
	"uninstall":  uninstallTask,
	"installed":  installedTask,
	"status":     statusTask,
	"start":      startTask,
	"stop":       stopTask,
	"sync-state": syncStateTask,
	"resync":     resyncTask,
	"restart":    restartTask,
	"logs":       logsTask,
}

// TaskRequirements maps task names to their required system conditions
var TaskRequirements = map[string][]string{
	"install":    {"docker", "uninstalled", "linux", "cpu4"},
	"uninstall":  {"docker", "stopped"},
	"installed":  {"docker"},
	"status":     {},
	"start":      {"docker", "stopped"},
	"stop":       {"docker", "running"},
	"sync-state": {"docker", "running"},
	"resync":     {"docker", "installed"},
	"restart":    {"docker", "running"},
	"logs":       {"docker", "running"},
}

var TarkArgs = map[string][]string{
	"install":    {"ethereumRPC"},
	"uninstall":  {},
	"installed":  {},
	"status":     {},
	"start":      {},
	"stop":       {},
	"sync-state": {},
	"resync":     {"bor", "heimdall"},
	"restart":    {},
	"logs":       {"bor", "heimdall", "lines"},
}

// validateRequirements checks if all requirements for a task are met
func ValidateRequirements(taskName string) (bool, []string) {
	var missingRequirements []string
	requirements, exists := TaskRequirements[taskName]
	if !exists {
		// No specific requirements
		return true, missingRequirements
	}

	for _, requirement := range requirements {
		checkFunc, exists := SystemRequirements[requirement]
		if !exists || !checkFunc() {
			// Requirement not met
			missingRequirements = append(missingRequirements, requirement)
		}
	}
	return len(missingRequirements) == 0, missingRequirements
}

func ValidateArgs(taskName string, args map[string]string) (bool, []string) {
	var missingArgs []string
	requiredArgs, exists := TarkArgs[taskName]
	if !exists {
		// No specific requirements
		return true, missingArgs
	}

	for _, arg := range requiredArgs {
		_, exists := args[arg]
		if !exists {
			// Requirement not met
			missingArgs = append(missingArgs, arg)
		}
	}
	return len(missingArgs) == 0, missingArgs
}
