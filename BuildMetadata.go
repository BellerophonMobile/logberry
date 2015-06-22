package logberry

// BuildMetadata is a simple structured representation of some basic
// properties characterizing the build of a particular binary.
type BuildMetadata struct {
	Host string
	User string
	Date string

	Repositories []RepositoryMetadata
}

// RepositoryMetadata captures basic identifying properties of a
// typical source version control repository instance.
type RepositoryMetadata struct {
	Repository string
	Branch     string
	Commit     string
	Dirty      bool
	Path       string
}
