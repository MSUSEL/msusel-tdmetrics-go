package reader

import (
	"fmt"
	"go/token"
	"sort"
	"strings"
)

// Project is the collection of compiled data for the project.
type Project struct {
	BasePath string
	FileSet  *token.FileSet
	Packages []*Package
}

// PrintParticipation formats and prints the participantion results.
func (p *Project) PrintParticipation() {
	packageStrings := make([]string, len(p.Packages))
	for i, pkg := range p.Packages {
		participants := pkg.Participation()

		maxWidth := 0
		for id := range participants {
			if width := len(id.String()); maxWidth < width {
				maxWidth = width
			}
		}

		parts := []string{}
		for id, f := range participants {
			subs := make([]string, len(f))
			for j, sub := range f {
				subs[j] = sub.String()
			}
			sort.Strings(subs)

			part := fmt.Sprintf(`   %-*s => [%s]`,
				maxWidth, id.String(), strings.Join(subs, `, `))
			parts = append(parts, part)
		}

		sort.Strings(parts)
		parts = append([]string{
			fmt.Sprintf("Package: %s", pkg.Package.Path()),
			fmt.Sprintf("Name:    %s", pkg.Package.Name()),
		}, parts...)
		packageStrings[i] = strings.Join(parts, "\n")
	}
	sort.Strings(packageStrings)
	result := strings.Join(packageStrings, "\n\n")
	fmt.Println(result)
}

// PrintFuncs formats and prints the functions.
func (p *Project) PrintFuncs() {
	packageStrings := make([]string, len(p.Packages))
	for i, pkg := range p.Packages {
		path := pkg.Package.Path()
		defFuncs := pkg.DefinedFuncs()

		maxWidth := 0
		for id := range defFuncs {
			if width := len(id.String()); maxWidth < width {
				maxWidth = width
			}
		}

		parts := []string{}
		for id, f := range defFuncs {
			funcStr := strings.ReplaceAll(f.String(), path+`.`, ``)
			part := fmt.Sprintf(`   %-*s => [%s]`, maxWidth, id.String(), funcStr)
			parts = append(parts, part)
		}

		sort.Strings(parts)
		parts = append([]string{
			fmt.Sprintf("Package: %s", path),
			fmt.Sprintf("Name:    %s", pkg.Package.Name()),
		}, parts...)
		packageStrings[i] = strings.Join(parts, "\n")
	}
	sort.Strings(packageStrings)
	result := strings.Join(packageStrings, "\n\n")
	fmt.Println(result)
}
