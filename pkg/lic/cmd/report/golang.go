// Package report implements the ``lic report golang'' command.
package report

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/tehcyx/lic/internal/fileop"
	"github.com/tehcyx/lic/pkg/lic/core"
)

const (
	//GoModImportStart will cover the line in the go.mod file that starts indicating the libraries included: "required ("
	GoModImportStart = "^require \\(.*"
	//GoModImportEnd will cover the line in the go.mod file that ends indicating the libraries included: ")", as we go top to bottom in the file this is the first closing bracket after finding "required ("
	GoModImportEnd = "^\\).*"
	//GoModLineImport will cover imports on a line between "require (" and ")". Imports will be of the format words, forward slash or "."
	GoModLineImport = "(\\S+|\\/|\\.)+"

	//GoFileImportStart indicates a multiline or single line import, either "import (" or "import \""
	GoFileImportStart = "^import (\\(|\").*"
	//GoFileImportEnd indicates the end of imports, either by a closing bracket ")", a variable definition, a function definition or struct definition
	GoFileImportEnd = "^(\\)|var|func|type).*"
	//GoFileLineImport indicates the import found in a multiline import between the double quotes
	GoFileLineImport = "\"(\\S+|\\/|\\.)+\""

	//GoFileExtension holds the pattern for the file extensions that should be included for import scans
	GoFileExtension = ".*\\.go$"
)

var (
	//DefaultWhitelistResources default list of acceptable imports that will get auto-parsed and checked for licenses
	DefaultWhitelistResources = []string{"github.com", "gopkg.in", "golang.org"}
)

//GolangReportOptions defines available options for the command
type GolangReportOptions struct {
	*ReportOptions
	WhitelistSourcesRegex []string
	BlacklistSourcesRegex []string
}

//NewGolangReportOptions creates options with default values
func NewGolangReportOptions(o *core.Options) *GolangReportOptions {
	return &GolangReportOptions{ReportOptions: NewReportOptions(o)}
}

//NewGolangReportCmd creates a new report command
func NewGolangReportCmd(o *GolangReportOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "golang",
		Short:   "Generates a report of current working directory",
		Long:    `Taking in consideration the source on the current path and checking for all licenses, generating a report output in the shell.`,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
		Aliases: []string{"go"},
	}

	cmd.Flags().BoolVarP(&o.Upload, "upload", "u", false, "Upload report to specified report endpoint to capture continuously")
	cmd.Flags().StringVarP(&o.UploadEndpoint, "upload-endpoint", "", "", "URL of the endpoint to report results of the scans")

	cmd.Flags().StringVarP(&o.SrcPath, "src-path", "", "", "Local path of sources to scan")
	cmd.Flags().BoolVarP(&o.HTMLOutput, "html-output", "o", false, "Specifies if results should be published as .html-file stored in current path")

	return cmd
}

//Run runs the command
func (o *GolangReportOptions) Run() error {
	//TODO (IF I need my own source codes actual package name [I assume I do to filter out self-imports])
	//
	// Two ways:
	//		1) If there's a go.mod file, check for "module" line and read the packages path
	//		2) If there's no go.mod file, check $GOPATH and make assumption based on that
	//
	if o.SrcPath != "" {
		if err := fileop.Exists(o.SrcPath); err != nil {
			log.Printf("Path '%s' does not exist or you don't have the proper access rights.\n", o.SrcPath)
			os.Exit(1)
		}
	} else {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Println("Couldn't get application path, exiting")
			os.Exit(1)
		}
		o.SrcPath = dir
	}
	var imports map[string]int
	imports = make(map[string]int)

	goModPath := path.Join(o.SrcPath, "go.mod")
	if err := fileop.Exists(goModPath); err == nil {
		// 1) PATH go.mod file EXISTS
		imports, err = fileop.ReadImports(goModPath, GoModImportStart, GoModImportEnd, GoModLineImport)
		if err != nil {
			log.Println("Error reading imports from go.mod file. Reading file tree now.")
		}
	} else if goPath := os.Getenv("GOPATH"); goPath != "" {
		// 2) PATH go.mod file DOES NOT EXIST
		var packageName string
		if match, _ := regexp.MatchString(goPath+"/src/.*", o.SrcPath); match {
			re := regexp.MustCompile(goPath + "/src/(.*)")
			res := re.FindStringSubmatch(o.SrcPath)
			if len(res) > 1 {
				packageName = res[1]
			} else {
				fmt.Println("Can't make assumption on package name, continuing without.")
			}
		} else {
			fmt.Println("Can't run on source folder. Project is not on $GOPATH.")
			os.Exit(1)
		}
		files, err := fileop.FilesInPath(o.SrcPath, ".*\\.go")
		if err != nil {
			fmt.Println("Couldn't read files in $GOPATH", packageName)
			os.Exit(1)
		}
		for _, f := range files {
			res, err := parseImportsGo(f)
			if err != nil {
				fmt.Println("Couldn't parse go imports from files")
				os.Exit(1)
			}
			for p := range res {
				if !strings.Contains(p, packageName) {
					p = strings.Replace(p, "\"", "", -1)
					imports[p] = 1
				}
			}
			// imports[f] = 1
		}
	} else {
		fmt.Println("Can't run on source folder. Project doesn't have a go.mod file or $GOPATH is not specified.")
		os.Exit(1)
	}

	var urlMap map[string]int
	urlMap = make(map[string]int)

	for u := range imports {
		for _, whitelist := range DefaultWhitelistResources {
			if strings.Contains(u, whitelist) {
				parsedURL, err := url.Parse("https://" + u)
				if err != nil {
					fmt.Printf("not a url: %s", u)
					continue
				}
				urlMap[parsedURL.String()] = 1
			}
		}
	}

	fmt.Printf("%+v\n", urlMap)

	// for u := range urlMap {
	// 	visitURL(u)
	// }

	//TODO: Check go.mod/go.sum, if they don't exist, open each file and check the imports statement

	return nil
}

func checkGitHubDependency(url string) (string, error) {
	re := regexp.MustCompile("^github.com(\\/\\w+\\/\\w+).*")
	result := re.FindStringSubmatch(url)
	if len(result) > 1 {
		return "http://api.github.com/repos" + result[1] + "/license", nil
	}
	return "", fmt.Errorf("Not a GitHub url: %s", url)
}

func visitURL(url string) {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	fmt.Println("query url", url)
	resp, err := netClient.Get(url + "/master/LICENSE")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func parseImportsGo(file string) (map[string]int, error) {
	fset := token.NewFileSet() // positions are relative to fset

	// Parse src but stop after processing the imports.
	f, err := parser.ParseFile(fset, file, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}
	var imports map[string]int
	imports = make(map[string]int)

	// Print the imports from the file's AST.
	for _, s := range f.Imports {
		imports[s.Path.Value] = 1
	}
	return imports, nil
}
