package projectscan

type Candidate struct {
	DBType     string
	Host       string
	Port       string
	Database   string
	User       string
	Password   string
	URL        string
	SQLitePath string
	SourceFile string
	Confidence int
	Parser     string
	Evidence   []string
}

type Result struct {
	Candidates []Candidate
}
