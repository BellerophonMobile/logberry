package env

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

func (x *BuildMetadata) Map() map[string]*RepositoryMetadata {

	m := make(map[string]*RepositoryMetadata)

	for _, r := range x.Repositories {
		m[r.Repository] = &r
	}

	return m

}
