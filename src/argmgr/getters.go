package argmgr

func GetAllArgs(runfile *RunFile) []GameArg {
	var allArgs []GameArg
	for _, category := range []string{"basic", "network", "advanced"} {
		if args, exists := runfile.Args[category]; exists {
			allArgs = append(allArgs, args...)
		}
	}
	return allArgs
}

func GetUIGroups(runfile *RunFile) []string {
	groups := make(map[string]bool)
	for _, arg := range GetAllArgs(runfile) {
		groups[arg.UIGroup] = true
	}

	var result []string
	for group := range groups {
		result = append(result, group)
	}
	return result
}

func GetArgsByGroup(runfile *RunFile, group string) []GameArg {
	var result []GameArg
	for _, arg := range GetAllArgs(runfile) {
		if arg.UIGroup == group {
			result = append(result, arg)
		}
	}
	return result
}
