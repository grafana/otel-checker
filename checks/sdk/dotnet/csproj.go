package dotnet

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CSharpProject represents the root element of a .NET SDK-style project file
type CSharpProject struct {
	XMLName        xml.Name           `xml:"Project"`
	ItemGroups     []CSharpItemGroup  `xml:"ItemGroup"`
	PropertyGroups []CSharpProperties `xml:"PropertyGroup"`
	SDK            string             `xml:"Sdk,attr"`
	path           string
}

// CSharpItemGroup represents a group of items in the .NET project file
type CSharpItemGroup struct {
	PackageReferences []CSharpPackageReference `xml:"PackageReference"`
}

// CSharpProperties contains .NET project properties
type CSharpProperties struct {
	TargetFramework string `xml:"TargetFramework"`
}

// CSharpPackageReference represents a NuGet package reference in a .NET project
type CSharpPackageReference struct {
	Include string `xml:"Include,attr"`
	Version string `xml:"Version,attr"`
}

// searches for a .csproj file in the specified directory and returns the path.
// only searches the immediate directory (non-recursive) and expects to find exactly one .csproj file.
func FindCSharpProject(dir string) (string, error) {
	var csprojFiles []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".csproj" {
			csprojFiles = append(csprojFiles, filepath.Join(dir, entry.Name()))
		}
	}

	switch len(csprojFiles) {
	case 0:
		return "", fmt.Errorf("no .csproj files found in directory: %s", dir)
	case 1:
		return csprojFiles[0], nil
	default:
		return "", fmt.Errorf("multiple .csproj files found: %s", strings.Join(csprojFiles, ", "))
	}
}

func LoadCSharpProject(path string) (*CSharpProject, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read project file: %w", err)
	}

	var proj CSharpProject
	if err := xml.Unmarshal(content, &proj); err != nil {
		return nil, fmt.Errorf("failed to parse project file: %w", err)
	}
	proj.path = path

	return &proj, nil
}
