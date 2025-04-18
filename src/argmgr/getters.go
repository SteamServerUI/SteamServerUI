package argmgr

func GetAllArgs() []GameArg {
	if CurrentRunfile == nil {
		return nil
	}

	var allArgs []GameArg
	for _, category := range []string{"basic", "network", "advanced"} {
		if args, exists := CurrentRunfile.Args[category]; exists {
			allArgs = append(allArgs, args...)
		}
	}
	return allArgs
}

func GetUIGroups() []string {
	if CurrentRunfile == nil {
		return nil
	}

	groups := make(map[string]bool)
	for _, arg := range GetAllArgs() {
		groups[arg.UIGroup] = true
	}

	var result []string
	for group := range groups {
		result = append(result, group)
	}
	return result
}

func GetArgsByGroup(group string) []GameArg {
	if CurrentRunfile == nil {
		return nil
	}

	var result []GameArg
	for _, arg := range GetAllArgs() {
		if arg.UIGroup == group {
			result = append(result, arg)
		}
	}
	return result
}
