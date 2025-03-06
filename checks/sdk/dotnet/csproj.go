package dotnet

import "encoding/xml"

// CSharpProject represents the root element of a .NET SDK-style project file
type CSharpProject struct {
	XMLName        xml.Name           `xml:"Project"`
	ItemGroups     []CSharpItemGroup  `xml:"ItemGroup"`
	PropertyGroups []CSharpProperties `xml:"PropertyGroup"`
	SDK            string             `xml:"Sdk,attr"`
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
