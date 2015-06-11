package logberry

// BuildMetadata is a simple structured representation of some basic
// properties characterizing the build of a particular binary.
type BuildMetadata struct {
	RepoRoot string
	Branch   string
	Commit   string
	Host     string
	User     string
	Date     string
}
