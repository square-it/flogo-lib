package data

import (
	"fmt"
	"strings"
)

// ResolverType is the type of the resolver for the attribute name
type ResolverType int

const (
	RES_DEFAULT  ResolverType = iota
	RES_ENV
	RES_PROPERTY
	RES_ACTIVITY
	RES_TRIGGER
	RES_SCOPE
)

var resolvers = make([]Resolver, 6, 6)

type Resolver func(scope Scope, path string) (interface{}, bool)

func GetResolverType(inAttrName string) (ResolverType, error) {
	if strings.HasPrefix(inAttrName, "${") {

		closeIdx := strings.Index(inAttrName, "}")

		if len(inAttrName) < 4 || closeIdx == -1 {
			return 0, fmt.Errorf("Invalid resolution expression [%s].", inAttrName)
		}

		toResolve := inAttrName[2:closeIdx]

		dotIdx := strings.Index(toResolve, ".")
		if dotIdx != -1 {
			resType := toResolve[:dotIdx]
			switch resType {
			case "property":
				return RES_PROPERTY, nil
			case "env":
				return RES_ENV, nil
			case "activity":
				return RES_ACTIVITY, nil
			case "trigger":
				return RES_TRIGGER, nil
			default:
				return RES_DEFAULT, fmt.Errorf("Unsupported resolver type [%s].", resType)
			}
		}
	}
	return RES_SCOPE, nil
}

func GetResolutionInfo(inAttrName string) (ResolverType, string, string, error) {
	if strings.HasPrefix(inAttrName, "${") {

		closeIdx := strings.Index(inAttrName, "}")

		if len(inAttrName) < 4 || closeIdx == -1 {
			return 0, "", "", fmt.Errorf("Invalid resolution expression [%s].", inAttrName)
		}

		toResolve := inAttrName[2:closeIdx]

		var path string

		if closeIdx+1 < len(inAttrName) {
			path = inAttrName[closeIdx+1:]
		}

		dotIdx := strings.Index(toResolve, ".")
		if dotIdx != -1 {
			resType := toResolve[:dotIdx]
			switch resType {
			case "property":
				return RES_PROPERTY, toResolve[9:], path, nil
			case "env":
				return RES_ENV, toResolve[4:], path, nil
			case "activity":
				return RES_ACTIVITY, toResolve[9:], path, nil
			case "trigger":
				return RES_TRIGGER, toResolve[8:], path, nil
			default:
				return 0, "", "", fmt.Errorf("Unsupported resolver type [%s].", resType)
			}
		}

		return RES_SCOPE, toResolve, path, nil
	}

	//todo should we support this?

	idx := strings.IndexFunc(inAttrName,isSep)

	if idx == -1 {
		return RES_SCOPE, inAttrName, "", nil
	}

	return RES_SCOPE, inAttrName[:idx],inAttrName[idx:], nil
}

func isSep(r rune) bool  {
	return r == '.' || r == '['
}

func SetResolver(rt ResolverType, resolver Resolver) {
	resolvers[rt] = resolver
}

func GetResolver(rt ResolverType) Resolver {
	return resolvers[rt]
}
