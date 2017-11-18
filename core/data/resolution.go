package data

import (
	"fmt"
	"strings"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"os"
)

//// ResolverType is the type of the resolver for the attribute name
//type ResolverType int
//
//const (
//	RES_DEFAULT  ResolverType = iota
//	RES_ENV
//	RES_PROPERTY
//	RES_ACTIVITY
//	RES_TRIGGER
//	RES_SCOPE
//)
//
//var resolvers = make([]Resolver, 6, 6)
//
//type Resolver func(scope Scope, path string) (interface{}, bool)
//
//func GetResolverType(inAttrName string) (ResolverType, error) {
//	if strings.HasPrefix(inAttrName, "${") {
//
//		closeIdx := strings.Index(inAttrName, "}")
//
//		if len(inAttrName) < 4 || closeIdx == -1 {
//			return 0, fmt.Errorf("Invalid resolution expression [%s].", inAttrName)
//		}
//
//		toResolve := inAttrName[2:closeIdx]
//
//		dotIdx := strings.Index(toResolve, ".")
//		if dotIdx != -1 {
//			resType := toResolve[:dotIdx]
//			switch resType {
//			case "property":
//				return RES_PROPERTY, nil
//			case "env":
//				return RES_ENV, nil
//			case "activity":
//				return RES_ACTIVITY, nil
//			case "trigger":
//				return RES_TRIGGER, nil
//			default:
//				return RES_DEFAULT, fmt.Errorf("Unsupported resolver type [%s].", resType)
//			}
//		}
//	}
//	return RES_SCOPE, nil
//}
//
//func GetResolutionInfo(inAttrName string) (ResolverType, string, string, error) {
//	if strings.HasPrefix(inAttrName, "${") {
//
//		closeIdx := strings.Index(inAttrName, "}")
//
//		if len(inAttrName) < 4 || closeIdx == -1 {
//			return 0, "", "", fmt.Errorf("Invalid resolution expression [%s].", inAttrName)
//		}
//
//		toResolve := inAttrName[2:closeIdx]
//
//		var path string
//
//		if closeIdx+1 < len(inAttrName) {
//			path = inAttrName[closeIdx+1:]
//		}
//
//		dotIdx := strings.Index(toResolve, ".")
//		if dotIdx != -1 {
//			resType := toResolve[:dotIdx]
//			switch resType {
//			case "property":
//				return RES_PROPERTY, toResolve[9:], path, nil
//			case "env":
//				return RES_ENV, toResolve[4:], path, nil
//			case "activity":
//				return RES_ACTIVITY, toResolve[9:], path, nil
//			case "trigger":
//				return RES_TRIGGER, toResolve[8:], path, nil
//			default:
//				return 0, "", "", fmt.Errorf("Unsupported resolver type [%s].", resType)
//			}
//		}
//
//		return RES_SCOPE, toResolve, path, nil
//	}
//
//	//todo should we support this?
//
//	idx := strings.IndexFunc(inAttrName, isSep)
//
//	if idx == -1 {
//		return RES_SCOPE, inAttrName, "", nil
//	}
//
//	return RES_SCOPE, inAttrName[:idx], inAttrName[idx:], nil
//}
//
//
//
//func SetResolver(rt ResolverType, resolver Resolver) {
//	resolvers[rt] = resolver
//}
//
//func GetResolver(rt ResolverType) Resolver {
//	return resolvers[rt]
//}

func isSep(r rune) bool {
	return r == '.' || r == '['
}

type Resolver interface {
	Resolve(toResolve string, scope Scope) (value interface{}, err error)
}

type BasicResolver struct {
}

func (r *BasicResolver) Resolve(toResolve string, scope Scope) (value interface{}, err error) {

	var details *ResolutionDetails

	if strings.HasPrefix(toResolve, "${") {
		details, err = GetResolutionDetailsOld(toResolve)
	} else if strings.HasPrefix(toResolve, "$") {
		details, err = GetResolutionDetails(toResolve[1:])
	} else {
		return SimpleScopeResolve(toResolve, scope)
	}

	if err != nil {
		return nil, err
	}

	if details == nil {
		return nil, fmt.Errorf("unable to resolve '%s'", toResolve)
	}

	var exists bool

	switch details.ResolverName {
	case "property":
		// Property resolution
		provider := GetPropertyProvider()
		value, exists = provider.GetProperty(details.Property + details.Path) //should we add the path and reset it to ""
		if !exists {
			err := fmt.Errorf("failed to resolve Property: '%s', ensure that property is configured in the application", details.Property)
			logger.Error(err.Error())
			return nil, err
		}
	case "env":
		// Environment resolution
		value, exists = os.LookupEnv(details.Property + details.Path)
		if !exists {
			err := fmt.Errorf("failed to resolve Environment Variable: '%s', ensure that variable is configured", details.Property)
			logger.Error(err.Error())
			return "", err
		}
	default:
		return nil, fmt.Errorf("unsupported resolver: %s", details.ResolverName)
	}

	return value, nil
}

func SimpleScopeResolve(toResolve string, scope Scope) (value interface{}, err error) {
	idx := strings.Index(toResolve, ".")

	if idx != -1 {
		attr, found := scope.GetAttr(toResolve[:idx])
		if !found {
			return nil, fmt.Errorf("could not resolve '%s'", toResolve)
		}
		value, err := PathGetValue(attr.Value, toResolve[idx:])
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		return value, nil

	} else {
		attr, found := scope.GetAttr(toResolve)
		if !found {
			return nil, fmt.Errorf("could not resolve '%s'", toResolve)
		}

		return attr.Value, nil
	}
}

type ResolutionDetails struct {
	ResolverName string
	Item         string
	Property     string
	Path         string
}

func GetResolutionDetails(toResolve string) (*ResolutionDetails, error) {

	//todo optimize, maybe tokenize first

	dotIdx := strings.Index(toResolve, ".")

	if dotIdx == -1 {
		return nil, fmt.Errorf("invalid resolution expression [%s]", toResolve)
	}

	details := &ResolutionDetails{}
	itemIdx := strings.Index(toResolve[:dotIdx], "[")

	if itemIdx != -1 {
		details.Item = toResolve[itemIdx+1:dotIdx-1]
		details.ResolverName = toResolve[:itemIdx]
	} else {
		details.ResolverName = toResolve[:dotIdx]

		//special case for activity without brackets
		if strings.HasPrefix(toResolve, "activity") {
			nextDot := strings.Index(toResolve[dotIdx+1:], ".") + dotIdx + 1
			details.Item = toResolve[dotIdx+1:nextDot]
			dotIdx = nextDot
		}
	}

	pathIdx := strings.IndexFunc(toResolve[dotIdx+1:], isSep)

	if pathIdx != -1 {
		pathStart := pathIdx + dotIdx + 1
		details.Path = toResolve[pathStart:]
		details.Property = toResolve[dotIdx+1:pathStart]
	} else {
		details.Property = toResolve[dotIdx+1:]
	}

	return details, nil
}

func GetResolutionDetailsOld(toResolve string) (*ResolutionDetails, error) {

	//todo optimize, maybe tokenize first

	closeIdx := strings.Index(toResolve, "}")

	if len(toResolve) < 4 || closeIdx == -1 {
		return nil, fmt.Errorf("invalid resolution expression [%s]", toResolve)
	}

	details := &ResolutionDetails{}

	dotIdx := strings.Index(toResolve, ".")

	if dotIdx == -1 {
		return nil, fmt.Errorf("invalid resolution expression [%s]", toResolve)
	}

	details.ResolverName = toResolve[2:dotIdx]

	if details.ResolverName == "activity" {
		nextDot := strings.Index(toResolve[dotIdx+1:], ".") + dotIdx + 1
		details.Item = toResolve[dotIdx+1:nextDot]
		dotIdx = nextDot
	}
	details.Property = toResolve[dotIdx+1:closeIdx]

	if closeIdx+1 < len(toResolve) {
		details.Path = toResolve[closeIdx+1:]
	}

	return details, nil
}
