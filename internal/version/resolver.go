package version

import (
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
)

func FindBestVersion(available []string, constraint string) (string, error) {
	constraint = strings.TrimSpace(constraint)

	switch {
	case constraint == "*": // Latest
		matching := filterVersions(available, nil)
		if len(matching) == 0 {
			return "", nil
		}
		sort.Sort(sort.Reverse(matching))
		return matching[0].String(), nil

	case strings.Contains(constraint, "||"): // OR
		for _, part := range strings.Split(constraint, "||") {
			if ver, err := FindBestVersion(available, strings.TrimSpace(part)); err == nil && ver != "" {
				return ver, nil
			}
		}
		return "", nil

	default:
		c, err := semver.NewConstraint(constraint)
		if err != nil {
			return "", err
		}

		matching := filterVersions(available, c)
		if len(matching) == 0 {
			return "", nil
		}

		sort.Sort(sort.Reverse(matching))
		return matching[0].String(), nil
	}
}

func filterVersions(versions []string, c *semver.Constraints) semver.Collection {
	var matching semver.Collection
	for _, v := range versions {
		if ver, err := semver.NewVersion(v); err == nil {
			if c == nil || c.Check(ver) {
				matching = append(matching, ver)
			}
		}
	}
	return matching
}
