package logberry

// BuildMetadata is a simple structured representation of some basic
// properties characterizing the build of a particular binary.
type BuildMetadata struct {
	Host     string
	User     string
	Date     string

  Repositories []RepositoryMetadata
}

type RepositoryMetadata struct {
	Repository string
	Branch string
	Commit string
	Dirty bool
	Path string
}
