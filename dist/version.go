package dist

// set by ldflags when you "mage build"
var (
	CommitHash string
	Timestamp  string
	GitTag     string
)

func init() {
	if GitTag == "" {
		GitTag = "-not set-"
	}
	if Timestamp == "" {
		Timestamp = "-not set-"
	}
	if CommitHash == "" {
		CommitHash = "-not set-"
	}
}
