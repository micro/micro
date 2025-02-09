package registry

import (
	"path"
	"regexp"
	"strings"
)

var (
	versionRe = regexp.MustCompilePOSIX("^v[0-9]+$")
)

// Translates /foo/bar/zool into api service micro.api.foo method Bar.Zool
// Translates /foo/bar into api service micro.api.foo method Foo.Bar
func apiRoute(p string) (string, string) {
	p = path.Clean(p)
	p = strings.TrimPrefix(p, "/")
	parts := strings.Split(p, "/")

	// if we have 1 part assume name Name.Call
	if len(parts) == 1 && len(parts[0]) > 0 {
		return parts[0], methodName(append(parts, "Call"))
	}

	// If we've got two or less parts
	// Use first part as service
	// Use all parts as method
	if len(parts) <= 2 {
		name := parts[0]
		return name, methodName(parts)
	}

	// Treat /v[0-9]+ as versioning where we have 3 parts
	// /v1/foo/bar => service: v1.foo method: Foo.bar
	if len(parts) == 3 && versionRe.Match([]byte(parts[0])) {
		name := strings.Join(parts[:len(parts)-1], ".")
		return name, methodName(parts[len(parts)-2:])
	}

	// Service is everything minus last two parts
	// Method is the last two parts
	name := strings.Join(parts[:len(parts)-2], ".")
	return name, methodName(parts[len(parts)-2:])
}

func methodName(parts []string) string {
	for i, part := range parts {
		parts[i] = toCamel(part)
	}

	return strings.Join(parts, ".")
}

func toCamel(s string) string {
	words := strings.Split(s, "-")
	var out string
	for _, word := range words {
		out += strings.Title(word)
	}
	return out
}
