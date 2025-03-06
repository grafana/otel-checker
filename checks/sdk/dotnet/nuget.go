package dotnet

type Package struct {
	ID               string `json:"id"`
	RequestedVersion string `json:"requestedVersion,omitempty"`
	ResolvedVersion  string `json:"resolvedVersion"`
}

type Framework struct {
	Framework          string    `json:"framework"`
	TopLevelPackages   []Package `json:"topLevelPackages"`
	TransitivePackages []Package `json:"transitivePackages"`
}

type Project struct {
	Path       string      `json:"path"`
	Frameworks []Framework `json:"frameworks"`
}

type NuGetPackageList struct {
	Version    int       `json:"version"`
	Parameters string    `json:"parameters"`
	Projects   []Project `json:"projects"`
}
