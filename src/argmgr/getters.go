package argmgr

func GetAllArgs(template *GameTemplate) []GameArg {
	var allArgs []GameArg
	for _, category := range []string{"basic", "network", "advanced"} {
		if args, exists := template.Args[category]; exists {
			allArgs = append(allArgs, args...)
		}
	}
	return allArgs
}

func GetUIGroups(template *GameTemplate) []string {
	groups := make(map[string]bool)
	for _, arg := range GetAllArgs(template) {
		groups[arg.UIGroup] = true
	}

	var result []string
	for group := range groups {
		result = append(result, group)
	}
	return result
}

func GetArgsByGroup(template *GameTemplate, group string) []GameArg {
	var result []GameArg
	for _, arg := range GetAllArgs(template) {
		if arg.UIGroup == group {
			result = append(result, arg)
		}
	}
	return result
}
