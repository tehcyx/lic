package report

//License represents a license, e.g. Apache 2, GNU GPL v2
type License struct {
	Name      string
	AltName   string
	ShortName string
	Text      string
	Link      string
}

//Licenses map of all licenses supported by this library
var Licenses map[string]License

func init() {
	Licenses = make(map[string]License)
	Licenses = map[string]License{
		"aflv3":      License{Name: "Academic Free License v3.0", ShortName: "afl-3.0"},
		"apache20":   License{Name: "Apache license 2.0", ShortName: "apache-2.0"},
		"artistic20": License{Name: "Artistic license 2.0", ShortName: "artistic-2.0"},
		"bsl10":      License{Name: "Boost Software License 1.0", ShortName: "bsl-1.0"},
		"bsd2":       License{Name: "BSD 2-clause \"Simplified\" license", ShortName: "bsd-2-clause"},
		"bsd3":       License{Name: "BSD 3-clause \"New\" or \"Revised\" license", ShortName: "bsd-3-clause"},
		"bsd3clear":  License{Name: "BSD 3-clause Clear license", ShortName: "bsd-3-clause-clear"},
		"cc":         License{Name: "Creative Commons license family", ShortName: "cc"},
		"cc0":        License{Name: "Creative Commons Zero v1.0 Universal", ShortName: "cc0-1.0"},
		"ccby3":      License{Name: "Creative Commons Attribution 3.0", ShortName: "cc-by-3.0"},
		"ccby4":      License{Name: "Creative Commons Attribution 4.0", ShortName: "cc-by-4.0"},
		"ccbysa":     License{Name: "Creative Commons Attribution Share Alike 4.0", ShortName: "cc-by-sa-4.0"},
		"wtfpl":      License{Name: "Do What The F*ck You Want To Public License", ShortName: "wtfpl"},
		"ecl20":      License{Name: "Educational Community License v2.0", ShortName: "ecl-2.0"},
		"epl10":      License{Name: "Eclipse Public License 1.0", ShortName: "epl-1.0"},
		"eupl11":     License{Name: "European Union Public License 1.1", ShortName: "eupl-1.1"},
		"agpl30":     License{Name: "GNU Affero General Public License v3.0", ShortName: "agpl-3.0"},
		"gpl":        License{Name: "GNU General Public License family", ShortName: "gpl"},
		"gpl20":      License{Name: "GNU General Public License v2.0", ShortName: "gpl-2.0"},
		"gpl30":      License{Name: "GNU General Public License v3.0", ShortName: "gpl-3.0"},
		"lgpl":       License{Name: "GNU Lesser General Public License family", ShortName: "lgpl"},
		"lgpl21":     License{Name: "GNU Lesser General Public License v2.1", ShortName: "lgpl-2.1"},
		"lgpl30":     License{Name: "GNU Lesser General Public License v3.0", ShortName: "lgpl-3.0"},
		"isc":        License{Name: "ISC", ShortName: "isc"},
		"lppl13c":    License{Name: "LaTeX Project Public License v1.3c", ShortName: "lppl-1.3c"},
		"mspl":       License{Name: "Microsoft Public License", ShortName: "ms-pl"},
		"mit":        License{Name: "MIT", ShortName: "mit"},
		"mpl20":      License{Name: "Mozilla Public License 2.0", ShortName: "mpl-2.0"},
		"osl30":      License{Name: "Open Software License 3.0", ShortName: "osl-3.0"},
		"postgresql": License{Name: "PostgreSQL License", ShortName: "postgresql"},
		"ofl11":      License{Name: "SIL Open Font License 1.1", ShortName: "ofl-1.1"},
		"ncsa":       License{Name: "University of Illinois/NCSA Open Source License", ShortName: "ncsa"},
		"unlicense":  License{Name: "The Unlicense", ShortName: "unlicense"},
		"zlib":       License{Name: "zLib License", ShortName: "zlib"},

		//Unknown license placeholder
		"na": License{Name: "Not Available", ShortName: "na", Text: "Placeholder for unknown license", AltName: "N/A"},
	}
}
