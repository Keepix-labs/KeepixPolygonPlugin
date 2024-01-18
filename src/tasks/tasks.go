package tasks

// TaskFunc defines the signature for task functions
type TaskFunc func(map[string]string) string

// TaskMap maps task names to their corresponding functions
var TaskMap = map[string]TaskFunc{
	"install":      installTask,
	"uninstall":    uninstallTask,
	"installed":    installedTask,
	"status":       statusTask,
	"start":        startTask,
	"stop":         stopTask,
	"sync-state":   syncStateTask,
	"resync":       resyncTask,
	"restart":      restartTask,
	"logs":         logsTask,
	"chain":        getChainTask,
	"wallet-fetch": walletFetchTask,
	"wallet-load":  walletLoadTask,
	"wallet-purge": walletPurgeTask,
	"pools-fetch":  poolsFetchTask,
	"unstake":      unstakeTask,
	"stake":        stakeTask,
	"rewards":      rewardTask,
}

// TaskRequirements maps task names to their required system conditions
var TaskRequirements = map[string][]string{
	"install":      {"docker", "uninstalled", "linux", "cpu4"},
	"uninstall":    {"docker", "stopped"},
	"installed":    {"docker"},
	"status":       {},
	"start":        {"docker", "stopped"},
	"stop":         {"docker", "running"},
	"sync-state":   {"docker", "running"},
	"resync":       {"docker", "installed"},
	"restart":      {"docker", "running"},
	"logs":         {"docker", "running"},
	"chain":        {"docker", "installed"},
	"wallet-fetch": {"installed"},
	"wallet-load":  {"installed"},
	"wallet-purge": {"installed"},
	"pools-fetch":  {"installed"},
	"unstake":      {"installed"},
	"stake":        {"installed"},
	"rewards":      {"installed"},
}

var TarkArgs = map[string][]string{
	"install":      {"ethereumRPC", "testnet", "autostart", "mnemonic"},
	"uninstall":    {},
	"installed":    {},
	"status":       {},
	"start":        {},
	"stop":         {},
	"sync-state":   {},
	"resync":       {"erigon", "heimdall"},
	"restart":      {},
	"logs":         {"erigon", "heimdall", "lines"},
	"chain":        {},
	"wallet-fetch": {},
	"wallet-load":  {"privateKey", "mnemonic"},
	"wallet-purge": {},
	"pools-fetch":  {},
	"unstake":      {"amount", "address"},
	"stake":        {"amount", "address"},
	"rewards":      {"address"},
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
