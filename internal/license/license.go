package license

import (
	"context"
	"log"

	"github.com/tehcyx/lic/internal/license/github"
)

// License represents a license, e.g. Apache 2, GNU GPL v2
type License struct {
	Name      string
	AltName   string
	ShortName string
	Text      string
	Link      string
}

// Licenses map of all licenses supported by this library
var Licenses map[string]License

const (
	licenseUnknownKey = "na"
)

func init() {
	Licenses = make(map[string]License)
	// Map using SPDX license identifiers as keys (what GitHub API returns)
	// Organized by license type for clarity
	Licenses = map[string]License{
		// Permissive Licenses - MIT/BSD Style
		"0bsd":              License{Name: "BSD Zero Clause License", ShortName: "0bsd"},
		"mit":               License{Name: "MIT License", ShortName: "mit"},
		"mit-0":             License{Name: "MIT No Attribution", ShortName: "mit-0"},
		"bsd-1-clause":      License{Name: "BSD 1-Clause License", ShortName: "bsd-1-clause"},
		"bsd-2-clause":      License{Name: "BSD 2-Clause \"Simplified\" License", ShortName: "bsd-2-clause"},
		"bsd-2-clause-patent": License{Name: "BSD 2-Clause Plus Patent License", ShortName: "bsd-2-clause-patent"},
		"bsd-3-clause":      License{Name: "BSD 3-Clause \"New\" or \"Revised\" License", ShortName: "bsd-3-clause"},
		"bsd-3-clause-clear": License{Name: "BSD 3-Clause Clear License", ShortName: "bsd-3-clause-clear"},
		"bsd-4-clause":      License{Name: "BSD 4-Clause \"Original\" or \"Old\" License", ShortName: "bsd-4-clause"},
		"isc":               License{Name: "ISC License", ShortName: "isc"},
		"ncsa":              License{Name: "University of Illinois/NCSA Open Source License", ShortName: "ncsa"},

		// Permissive Licenses - Apache Style
		"apache-1.0":        License{Name: "Apache License 1.0", ShortName: "apache-1.0"},
		"apache-1.1":        License{Name: "Apache License 1.1", ShortName: "apache-1.1"},
		"apache-2.0":        License{Name: "Apache License 2.0", ShortName: "apache-2.0"},

		// Permissive Licenses - Academic/Educational
		"afl-1.1":           License{Name: "Academic Free License v1.1", ShortName: "afl-1.1"},
		"afl-1.2":           License{Name: "Academic Free License v1.2", ShortName: "afl-1.2"},
		"afl-2.0":           License{Name: "Academic Free License v2.0", ShortName: "afl-2.0"},
		"afl-2.1":           License{Name: "Academic Free License v2.1", ShortName: "afl-2.1"},
		"afl-3.0":           License{Name: "Academic Free License v3.0", ShortName: "afl-3.0"},
		"ecl-1.0":           License{Name: "Educational Community License v1.0", ShortName: "ecl-1.0"},
		"ecl-2.0":           License{Name: "Educational Community License v2.0", ShortName: "ecl-2.0"},

		// Permissive Licenses - Other
		"bsl-1.0":           License{Name: "Boost Software License 1.0", ShortName: "bsl-1.0"},
		"unlicense":         License{Name: "The Unlicense", ShortName: "unlicense"},
		"zlib":              License{Name: "zLib License", ShortName: "zlib"},
		"postgresql":        License{Name: "PostgreSQL License", ShortName: "postgresql"},
		"wtfpl":             License{Name: "Do What The F*ck You Want To Public License", ShortName: "wtfpl"},
		"artistic-1.0":      License{Name: "Artistic License 1.0", ShortName: "artistic-1.0"},
		"artistic-2.0":      License{Name: "Artistic License 2.0", ShortName: "artistic-2.0"},
		"python-2.0":        License{Name: "Python License 2.0", ShortName: "python-2.0"},

		// Copyleft Licenses - Strong (GPL)
		"gpl":               License{Name: "GNU General Public License Family", ShortName: "gpl"},
		"gpl-1.0":           License{Name: "GNU General Public License v1.0", ShortName: "gpl-1.0"},
		"gpl-1.0-only":      License{Name: "GNU General Public License v1.0 only", ShortName: "gpl-1.0-only"},
		"gpl-1.0-or-later":  License{Name: "GNU General Public License v1.0 or later", ShortName: "gpl-1.0-or-later"},
		"gpl-2.0":           License{Name: "GNU General Public License v2.0", ShortName: "gpl-2.0"},
		"gpl-2.0-only":      License{Name: "GNU General Public License v2.0 only", ShortName: "gpl-2.0-only"},
		"gpl-2.0-or-later":  License{Name: "GNU General Public License v2.0 or later", ShortName: "gpl-2.0-or-later"},
		"gpl-3.0":           License{Name: "GNU General Public License v3.0", ShortName: "gpl-3.0"},
		"gpl-3.0-only":      License{Name: "GNU General Public License v3.0 only", ShortName: "gpl-3.0-only"},
		"gpl-3.0-or-later":  License{Name: "GNU General Public License v3.0 or later", ShortName: "gpl-3.0-or-later"},
		"agpl-1.0":          License{Name: "Affero General Public License v1.0", ShortName: "agpl-1.0"},
		"agpl-3.0":          License{Name: "GNU Affero General Public License v3.0", ShortName: "agpl-3.0"},
		"agpl-3.0-only":     License{Name: "GNU Affero General Public License v3.0 only", ShortName: "agpl-3.0-only"},
		"agpl-3.0-or-later": License{Name: "GNU Affero General Public License v3.0 or later", ShortName: "agpl-3.0-or-later"},

		// Copyleft Licenses - Weak (LGPL)
		"lgpl":              License{Name: "GNU Lesser General Public License Family", ShortName: "lgpl"},
		"lgpl-2.0":          License{Name: "GNU Lesser General Public License v2.0", ShortName: "lgpl-2.0"},
		"lgpl-2.0-only":     License{Name: "GNU Lesser General Public License v2.0 only", ShortName: "lgpl-2.0-only"},
		"lgpl-2.0-or-later": License{Name: "GNU Lesser General Public License v2.0 or later", ShortName: "lgpl-2.0-or-later"},
		"lgpl-2.1":          License{Name: "GNU Lesser General Public License v2.1", ShortName: "lgpl-2.1"},
		"lgpl-2.1-only":     License{Name: "GNU Lesser General Public License v2.1 only", ShortName: "lgpl-2.1-only"},
		"lgpl-2.1-or-later": License{Name: "GNU Lesser General Public License v2.1 or later", ShortName: "lgpl-2.1-or-later"},
		"lgpl-3.0":          License{Name: "GNU Lesser General Public License v3.0", ShortName: "lgpl-3.0"},
		"lgpl-3.0-only":     License{Name: "GNU Lesser General Public License v3.0 only", ShortName: "lgpl-3.0-only"},
		"lgpl-3.0-or-later": License{Name: "GNU Lesser General Public License v3.0 or later", ShortName: "lgpl-3.0-or-later"},

		// Copyleft Licenses - MPL/EPL Style (File-Level)
		"mpl-1.0":           License{Name: "Mozilla Public License 1.0", ShortName: "mpl-1.0"},
		"mpl-1.1":           License{Name: "Mozilla Public License 1.1", ShortName: "mpl-1.1"},
		"mpl-2.0":           License{Name: "Mozilla Public License 2.0", ShortName: "mpl-2.0"},
		"mpl-2.0-no-copyleft-exception": License{Name: "Mozilla Public License 2.0 (no copyleft exception)", ShortName: "mpl-2.0-no-copyleft-exception"},
		"epl-1.0":           License{Name: "Eclipse Public License 1.0", ShortName: "epl-1.0"},
		"epl-2.0":           License{Name: "Eclipse Public License 2.0", ShortName: "epl-2.0"},
		"eupl-1.0":          License{Name: "European Union Public License 1.0", ShortName: "eupl-1.0"},
		"eupl-1.1":          License{Name: "European Union Public License 1.1", ShortName: "eupl-1.1"},
		"eupl-1.2":          License{Name: "European Union Public License 1.2", ShortName: "eupl-1.2"},
		"cddl-1.0":          License{Name: "Common Development and Distribution License 1.0", ShortName: "cddl-1.0"},
		"cddl-1.1":          License{Name: "Common Development and Distribution License 1.1", ShortName: "cddl-1.1"},

		// Copyleft Licenses - Other
		"cpl-1.0":           License{Name: "Common Public License 1.0", ShortName: "cpl-1.0"},
		"osl-1.0":           License{Name: "Open Software License 1.0", ShortName: "osl-1.0"},
		"osl-1.1":           License{Name: "Open Software License 1.1", ShortName: "osl-1.1"},
		"osl-2.0":           License{Name: "Open Software License 2.0", ShortName: "osl-2.0"},
		"osl-2.1":           License{Name: "Open Software License 2.1", ShortName: "osl-2.1"},
		"osl-3.0":           License{Name: "Open Software License 3.0", ShortName: "osl-3.0"},

		// Creative Commons Licenses
		"cc":                License{Name: "Creative Commons License Family", ShortName: "cc"},
		"cc0-1.0":           License{Name: "Creative Commons Zero v1.0 Universal", ShortName: "cc0-1.0"},
		"cc-by-1.0":         License{Name: "Creative Commons Attribution 1.0 Generic", ShortName: "cc-by-1.0"},
		"cc-by-2.0":         License{Name: "Creative Commons Attribution 2.0 Generic", ShortName: "cc-by-2.0"},
		"cc-by-2.5":         License{Name: "Creative Commons Attribution 2.5 Generic", ShortName: "cc-by-2.5"},
		"cc-by-3.0":         License{Name: "Creative Commons Attribution 3.0 Unported", ShortName: "cc-by-3.0"},
		"cc-by-4.0":         License{Name: "Creative Commons Attribution 4.0 International", ShortName: "cc-by-4.0"},
		"cc-by-sa-1.0":      License{Name: "Creative Commons Attribution ShareAlike 1.0 Generic", ShortName: "cc-by-sa-1.0"},
		"cc-by-sa-2.0":      License{Name: "Creative Commons Attribution ShareAlike 2.0 Generic", ShortName: "cc-by-sa-2.0"},
		"cc-by-sa-2.5":      License{Name: "Creative Commons Attribution ShareAlike 2.5 Generic", ShortName: "cc-by-sa-2.5"},
		"cc-by-sa-3.0":      License{Name: "Creative Commons Attribution ShareAlike 3.0 Unported", ShortName: "cc-by-sa-3.0"},
		"cc-by-sa-4.0":      License{Name: "Creative Commons Attribution ShareAlike 4.0 International", ShortName: "cc-by-sa-4.0"},
		"cc-by-nc-1.0":      License{Name: "Creative Commons Attribution Non Commercial 1.0 Generic", ShortName: "cc-by-nc-1.0"},
		"cc-by-nc-2.0":      License{Name: "Creative Commons Attribution Non Commercial 2.0 Generic", ShortName: "cc-by-nc-2.0"},
		"cc-by-nc-2.5":      License{Name: "Creative Commons Attribution Non Commercial 2.5 Generic", ShortName: "cc-by-nc-2.5"},
		"cc-by-nc-3.0":      License{Name: "Creative Commons Attribution Non Commercial 3.0 Unported", ShortName: "cc-by-nc-3.0"},
		"cc-by-nc-4.0":      License{Name: "Creative Commons Attribution Non Commercial 4.0 International", ShortName: "cc-by-nc-4.0"},
		"cc-by-nc-nd-1.0":   License{Name: "Creative Commons Attribution Non Commercial No Derivatives 1.0 Generic", ShortName: "cc-by-nc-nd-1.0"},
		"cc-by-nc-nd-2.0":   License{Name: "Creative Commons Attribution Non Commercial No Derivatives 2.0 Generic", ShortName: "cc-by-nc-nd-2.0"},
		"cc-by-nc-nd-2.5":   License{Name: "Creative Commons Attribution Non Commercial No Derivatives 2.5 Generic", ShortName: "cc-by-nc-nd-2.5"},
		"cc-by-nc-nd-3.0":   License{Name: "Creative Commons Attribution Non Commercial No Derivatives 3.0 Unported", ShortName: "cc-by-nc-nd-3.0"},
		"cc-by-nc-nd-4.0":   License{Name: "Creative Commons Attribution Non Commercial No Derivatives 4.0 International", ShortName: "cc-by-nc-nd-4.0"},
		"cc-by-nc-sa-1.0":   License{Name: "Creative Commons Attribution Non Commercial ShareAlike 1.0 Generic", ShortName: "cc-by-nc-sa-1.0"},
		"cc-by-nc-sa-2.0":   License{Name: "Creative Commons Attribution Non Commercial ShareAlike 2.0 Generic", ShortName: "cc-by-nc-sa-2.0"},
		"cc-by-nc-sa-2.5":   License{Name: "Creative Commons Attribution Non Commercial ShareAlike 2.5 Generic", ShortName: "cc-by-nc-sa-2.5"},
		"cc-by-nc-sa-3.0":   License{Name: "Creative Commons Attribution Non Commercial ShareAlike 3.0 Unported", ShortName: "cc-by-nc-sa-3.0"},
		"cc-by-nc-sa-4.0":   License{Name: "Creative Commons Attribution Non Commercial ShareAlike 4.0 International", ShortName: "cc-by-nc-sa-4.0"},
		"cc-by-nd-1.0":      License{Name: "Creative Commons Attribution No Derivatives 1.0 Generic", ShortName: "cc-by-nd-1.0"},
		"cc-by-nd-2.0":      License{Name: "Creative Commons Attribution No Derivatives 2.0 Generic", ShortName: "cc-by-nd-2.0"},
		"cc-by-nd-2.5":      License{Name: "Creative Commons Attribution No Derivatives 2.5 Generic", ShortName: "cc-by-nd-2.5"},
		"cc-by-nd-3.0":      License{Name: "Creative Commons Attribution No Derivatives 3.0 Unported", ShortName: "cc-by-nd-3.0"},
		"cc-by-nd-4.0":      License{Name: "Creative Commons Attribution No Derivatives 4.0 International", ShortName: "cc-by-nd-4.0"},

		// Microsoft Licenses
		"ms-pl":             License{Name: "Microsoft Public License", ShortName: "ms-pl"},
		"ms-rl":             License{Name: "Microsoft Reciprocal License", ShortName: "ms-rl"},

		// Specialized Licenses
		"lppl-1.0":          License{Name: "LaTeX Project Public License v1.0", ShortName: "lppl-1.0"},
		"lppl-1.1":          License{Name: "LaTeX Project Public License v1.1", ShortName: "lppl-1.1"},
		"lppl-1.2":          License{Name: "LaTeX Project Public License v1.2", ShortName: "lppl-1.2"},
		"lppl-1.3a":         License{Name: "LaTeX Project Public License v1.3a", ShortName: "lppl-1.3a"},
		"lppl-1.3c":         License{Name: "LaTeX Project Public License v1.3c", ShortName: "lppl-1.3c"},
		"ofl-1.0":           License{Name: "SIL Open Font License 1.0", ShortName: "ofl-1.0"},
		"ofl-1.1":           License{Name: "SIL Open Font License 1.1", ShortName: "ofl-1.1"},
		"ofl-1.1-rfn":       License{Name: "SIL Open Font License 1.1 with Reserved Font Name", ShortName: "ofl-1.1-rfn"},
		"ofl-1.1-no-rfn":    License{Name: "SIL Open Font License 1.1 with no Reserved Font Name", ShortName: "ofl-1.1-no-rfn"},

		// Other/Proprietary/Unknown
		"other":             License{Name: "Other", ShortName: "other", Text: "Other license type"},
		"proprietary":       License{Name: "Proprietary", ShortName: "proprietary", Text: "Proprietary license"},

		// Unknown license placeholder
		"na": License{Name: "Not Available", ShortName: "na", Text: "Placeholder for unknown license", AltName: "N/A"},
	}
}

// getProviders returns the list of license providers in priority order
func getProviders() []Provider {
	return []Provider{
		github.NewProvider(), // GitHub repositories
		// Future providers can be added here (k8s.io, golang.org, gopkg.in, etc.)
	}
}

// Get retrieves license information using the first provider that supports the import path
func Get(name, version, branch, url string) License {
	return GetWithContext(context.Background(), name, version, branch, url)
}

// GetWithContext retrieves license information using the first provider that supports the import path
func GetWithContext(ctx context.Context, name, version, branch, url string) License {
	providers := getProviders()

	// Try each provider in order
	for _, provider := range providers {
		// Check for cancellation
		select {
		case <-ctx.Done():
			log.Printf("Info: License lookup cancelled for %s: %v\n", name, ctx.Err())
			return Licenses[licenseUnknownKey]
		default:
		}

		if provider.Supports(name) {
			key, err := provider.GetLicense(ctx, name, version, branch, url)
			if err != nil {
				log.Printf("Warning: %s provider couldn't get license for %s: %v\n", provider.Name(), name, err)
				return Licenses[licenseUnknownKey]
			}

			// Check if the key exists in our license map
			if lic, ok := Licenses[key]; ok {
				return lic
			}

			// If license key from API doesn't match our map, log it and return unknown
			log.Printf("Warning: unknown license key '%s' for %s from %s provider\n", key, name, provider.Name())
			return Licenses[licenseUnknownKey]
		}
	}

	// No provider supports this import path
	log.Printf("Info: No license provider available for %s\n", name)
	return Licenses[licenseUnknownKey]
}
