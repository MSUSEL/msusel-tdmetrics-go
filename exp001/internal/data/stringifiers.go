package data

import (
	"fmt"
	"sort"
	"strings"
)

// FormatParticipation formats the participantion results for a package.
func (p *Package) FormatParticipation() string {
	participants := p.Participation()

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
		fmt.Sprintf("Package: %s", p.Package.Path()),
		fmt.Sprintf("Name:    %s", p.Package.Name()),
	}, parts...)
	return strings.Join(parts, "\n")
}

// FormatParticipation formats the participantion results for a project.
func (p *Project) FormatParticipation() string {
	packageStrings := make([]string, len(p.Packages))
	for i, pkg := range p.Packages {
		packageStrings[i] = pkg.FormatParticipation()
	}
	sort.Strings(packageStrings)
	return strings.Join(packageStrings, "\n\n")
}

// PrintParticipation formats and prints the participantion results.
func (p *Project) PrintParticipation() {
	fmt.Println(p.FormatParticipation())
}

// PrintFuncs formats the functions for a package.
func (p *Package) FormatFuncs() string {
	path := p.Package.Path()
	defFuncs := p.DefinedFuncs()

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
		fmt.Sprintf("Name:    %s", p.Package.Name()),
	}, parts...)
	return strings.Join(parts, "\n")
}

// PrintFuncs formats the functions for a project
func (p *Project) FormatFuncs() string {
	packageStrings := make([]string, len(p.Packages))
	for i, pkg := range p.Packages {
		packageStrings[i] = pkg.FormatFuncs()
	}
	sort.Strings(packageStrings)
	return strings.Join(packageStrings, "\n\n")
}

// PrintFuncs formats and prints the functions.
func (p *Project) PrintFuncs() {
	fmt.Println(p.FormatFuncs())
}
