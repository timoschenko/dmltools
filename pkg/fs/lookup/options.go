// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lookup

import "time"

type config struct {
	cache         Cache
	filter        Filter
	ignoreMatcher IgnoreMatcher

	includePattern       string
	excludePattern       string
	extendExcludePattern string
	forceExcludePattern  string
	ignoreFiles          []string

	vcsMarkers       []string
	stopAtVCSRoot    bool
	stopAtMountPoint bool
	verifyOwnership  bool
	allowSymlinks    bool
	allowSymlinkEsc  bool
	ceilingDirs      []string
	maxDepth         int
	maxFileSize      int64
	cacheTTL         time.Duration
	strictSecurity   bool
	onSkip           SkipCallback
	onIgnore         IgnoreCallback
	fsys             FileSystem
	clock            Clock
}

func defaultConfig() config {
	return config{
		vcsMarkers:       []string{VCSMarkerGit, VCSMarkerHg},
		stopAtVCSRoot:    true,
		stopAtMountPoint: true,
		verifyOwnership:  true,
		allowSymlinks:    true,
		allowSymlinkEsc:  false,
		maxDepth:         DefaultMaxDepth,
		maxFileSize:      DefaultMaxFileSize,
		fsys:             NewOSFileSystem(),
		clock:            realClock{},
	}
}

// Option configures a [Cache], [Filter], or [Finder] instance.
type Option func(*config)

// WithCache injects an existing [Cache] into [NewFinder].
func WithCache(cache Cache) Option {
	return func(c *config) {
		c.cache = cache
	}
}

// WithFilter injects an existing [Filter] into [NewFinder].
func WithFilter(filter Filter) Option {
	return func(c *config) {
		c.filter = filter
	}
}

// WithIgnoreMatcher injects an existing [IgnoreMatcher] into [NewFinder].
func WithIgnoreMatcher(matcher IgnoreMatcher) Option {
	return func(c *config) {
		c.ignoreMatcher = matcher
	}
}

// WithIgnoreFiles configures hierarchical ignore filenames (such as
// [FileGitIgnore] and [FileIgnore]) evaluated during [Finder.Collect].
// Within each directory, files are evaluated in the order specified.
func WithIgnoreFiles(names ...string) Option {
	return func(c *config) {
		c.ignoreFiles = normalizeIgnoreFileNames(names)
	}
}

// WithInclude sets the regular expression matching paths to include.
func WithInclude(pattern string) Option {
	return func(c *config) {
		c.includePattern = pattern
	}
}

// WithExclude sets the regular expression matching paths to exclude.
func WithExclude(pattern string) Option {
	return func(c *config) {
		c.excludePattern = pattern
	}
}

// WithExtendExclude sets an additional regular expression matching paths to exclude.
func WithExtendExclude(pattern string) Option {
	return func(c *config) {
		c.extendExcludePattern = pattern
	}
}

// WithForceExclude sets the regular expression matching paths to forcibly exclude.
func WithForceExclude(pattern string) Option {
	return func(c *config) {
		c.forceExcludePattern = pattern
	}
}

// WithVCSMarkers overrides the VCS root marker names (default: .git, .hg).
func WithVCSMarkers(markers ...string) Option {
	return func(c *config) {
		filtered := make([]string, 0, len(markers))
		for _, m := range markers {
			if m != "" {
				filtered = append(filtered, m)
			}
		}
		c.vcsMarkers = filtered
	}
}

// WithStopAtVCSRoot controls whether upward traversal stops at VCS roots.
func WithStopAtVCSRoot(enabled bool) Option {
	return func(c *config) {
		c.stopAtVCSRoot = enabled
	}
}

// WithStopAtMountPoint controls whether upward traversal stops at mount boundaries.
func WithStopAtMountPoint(enabled bool) Option {
	return func(c *config) {
		c.stopAtMountPoint = enabled
	}
}

// WithVerifyOwnership controls whether upward traversal verifies directory UIDs.
func WithVerifyOwnership(enabled bool) Option {
	return func(c *config) {
		c.verifyOwnership = enabled
	}
}

// WithAllowSymlinks controls whether symlinked utility files are followed.
func WithAllowSymlinks(allow bool) Option {
	return func(c *config) {
		c.allowSymlinks = allow
	}
}

// WithAllowSymlinkEscape controls whether symlinks may resolve outside the boundary root.
func WithAllowSymlinkEscape(allow bool) Option {
	return func(c *config) {
		c.allowSymlinkEsc = allow
	}
}

// WithCeilingDirs configures directory paths where upward traversal must stop.
func WithCeilingDirs(dirs ...string) Option {
	return func(c *config) {
		c.ceilingDirs = append([]string(nil), dirs...)
	}
}

// WithMaxDepth sets the maximum number of parent directories to ascend.
func WithMaxDepth(depth int) Option {
	return func(c *config) {
		if depth >= 0 {
			c.maxDepth = depth
		}
	}
}

// WithMaxFileSize sets the maximum permitted size in bytes for candidate files.
// Non-positive values disable the size limit.
func WithMaxFileSize(maxBytes int64) Option {
	return func(c *config) {
		c.maxFileSize = maxBytes
	}
}

// WithCacheTTL sets a cache validity window during which repeated probes skip
// re-statting the filesystem. Zero forces stat-fingerprint validation on every call.
func WithCacheTTL(ttl time.Duration) Option {
	return func(c *config) {
		if ttl >= 0 {
			c.cacheTTL = ttl
		}
	}
}

// WithStrictSecurity causes security check failures on candidate files to
// return a [*SecurityError] instead of skipping the file.
func WithStrictSecurity(strict bool) Option {
	return func(c *config) {
		c.strictSecurity = strict
	}
}

// WithOnSkip registers a callback notified whenever a candidate is skipped.
func WithOnSkip(cb SkipCallback) Option {
	return func(c *config) {
		c.onSkip = cb
	}
}

// WithOnIgnore registers a callback notified when [Finder.Collect] ignores a path.
func WithOnIgnore(cb IgnoreCallback) Option {
	return func(c *config) {
		c.onIgnore = cb
	}
}

// WithFileSystem overrides the backing [FileSystem] (useful for tests).
func WithFileSystem(fsys FileSystem) Option {
	return func(c *config) {
		if fsys != nil {
			c.fsys = fsys
		}
	}
}

// WithClock overrides the [Clock] used for cache TTL evaluation.
func WithClock(clock Clock) Option {
	return func(c *config) {
		if clock != nil {
			c.clock = clock
		}
	}
}
