// Package report implements the `lic report golang` (`lic r go`)command.
package report

import (
	"crypto/sha256"
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
	"github.com/tehcyx/lic/internal/licensereport"
	"github.com/tehcyx/lic/pkg/lic/core"
)

const (
	//GoModImportStart will cover the line in the go.mod file that starts indicating the libraries included: "required ("
	GoModImportStart = "^require \\(.*"
	//GoModImportEnd will cover the line in the go.mod file that ends indicating the libraries included: ")", as we go top to bottom in the file this is the first closing bracket after finding "required ("
	GoModImportEnd = "^\\).*"
	//GoModLineImport will cover imports on a line between "require (" and ")".
	GoModLineImport = "(\\S+|\\/|\\.)+"
	//GoModInlineImport will cover single line imports that are just "require github.com/user/repo".
	GoModInlineImport = "^require (\\S+|\\/|\\.)+ (\\S+|\\/|\\.)+"

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

	//StdLibraryList list of Standard Library imports as of go 1.11.5
	stdLibraryList = []string{
		"archive", "archive/tar", "archive/zip", "bufio", "builtin", "bytes", "compress", "compress/bzip2", "compress/flate",
		"compress/gzip", "compress/lzw", "compress/zlib", "container	", "container/heap", "container/list", "container/ring",
		"context", "crypto", "crypto/aes", "crypto/cipher", "crypto/des", "crypto/dsa", "crypto/ecdsa", "crypto/elliptic",
		"crypto/hmac", "crypto/md5", "crypto/rand", "crypto/rc4", "crypto/rsa", "crypto/sha1", "crypto/sha256", "crypto/sha512",
		"crypto/subtle", "crypto/tls", "crypto/x509", "crypto/x509/pkix", "database", "database/sql", "database/sql/driver",
		"debug", "debug/dwarf", "debug/elf", "debug/gosym", "debug/macho", "debug/pe", "debug/plan9obj", "encoding",
		"encoding/ascii85", "encoding/asn1", "encoding/base32", "encoding/base64", "encoding/binary", "encoding/csv",
		"encoding/gob", "encoding/hex", "encoding/json", "encoding/pem", "encoding/xml", "errors", "expvar", "flag", "fmt",
		"go", "go/ast", "go/build", "go/constant", "go/doc", "go/format", "go/importer", "go/parser", "go/printer",
		"go/scanner", "go/token", "go/types", "hash", "hash/adler32", "hash/crc32", "hash/crc64", "hash/fnv", "html",
		"html/template", "image", "image/color", "image/palette", "image/draw", "image/gif", "image/jpeg", "image/png", "index",
		"index/suffixarray", "io", "io/ioutil", "log", "log/syslog", "math", "math/big", "math/bits", "math/cmplx", "math/rand",
		"mime", "mime/multipart", "mime/quotedprintable", "net", "net/http", "net/http/cgi", "net/http/cookiejar", "net/http/fcgi",
		"net/http/httptest", "net/http/httptrace", "net/http/httputil", "net/http/pprof", "net/mail", "net/rpc", "net/rpc/jsonrpc",
		"net/smtp", "net/textproto", "net/url", "os", "os/exec", "os/signal", "os/user", "path", "path/filepath", "plugin",
		"reflect", "regexp", "regexp/syntax", "runtime", "runtime/cgo", "runtime/debug", "runtime/msan", "runtime/pprof",
		"runtime/race", "runtime/trace", "sort", "strconv", "strings", "sync", "sync/atomic", "syscall", "syscall/js", "testing",
		"testing/iotest", "testing/quick", "text", "text/scanner", "text/tabwriter", "text/template", "text/template/parse",
		"time", "unicode", "unicode/utf16", "unicode/utf8", "unsafe",
	}
)

//GolangReportOptions defines available options for the command
type GolangReportOptions struct {
	*Options
}

//NewGolangReportOptions creates options with default values
func NewGolangReportOptions(o *core.Options) *GolangReportOptions {
	return &GolangReportOptions{Options: NewReportOptions(o)}
}

//NewGolangReportCmd creates a new report command
func NewGolangReportCmd(o *GolangReportOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "golang",
		Short:   "Generates a report of current working directory or specified path",
		Long:    `Taking in consideration the source on the current path and checking for all licenses, generating a report output in the shell.`,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
		Aliases: []string{"go"},
	}

	cmd.Flags().BoolVarP(&o.Upload, "upload", "u", false, "Upload report to specified report endpoint to capture continuously")
	cmd.Flags().StringVarP(&o.UploadEndpoint, "upload-endpoint", "", "", "URL of the endpoint to report results of the scans")

	cmd.Flags().StringVarP(&o.SrcPath, "src", "", "", "Local path of sources to scan")
	cmd.Flags().BoolVarP(&o.HTMLOutput, "html-output", "o", false, "Specifies if results should be published as .html-file stored in current path")

	cmd.Flags().StringVarP(&o.ProjectVersion, "project-version", "", "n/a", "Version of scan target")

	return cmd
}

//Run runs the command
// Scan has to exclusive paths this could go:
//		1) If there's a go.mod file, check for "module" line and read the packages path
//		2) If there's no go.mod file, check $GOPATH and make assumption based on that
func (o *GolangReportOptions) Run() error {
	if o.SrcPath != "" {
		fmt.Println(o.SrcPath)
		if err := fileop.Exists(o.SrcPath); err != nil {
			err := fmt.Errorf("path '%s' does not exist or you don't have the proper access rights", o.SrcPath)
			log.Printf("%s\n", err.Error())
			return err
		}
	} else {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			err := fmt.Errorf("couldn't get application path, exiting")
			log.Printf("%s\n", err.Error())
			return err
		}
		o.SrcPath = dir
	}

	var imports map[string]string
	var packageName string

	imports = make(map[string]string)

	goModPath := path.Join(o.SrcPath, "go.mod")
	if err := fileop.Exists(goModPath); err == nil {
		// 1) PATH go.mod file EXISTS
		imports, err = fileop.ReadImports(goModPath, GoModImportStart, GoModImportEnd, GoModLineImport, GoModInlineImport)
		if err != nil {
			log.Println("Error reading imports from go.mod file. Reading file tree now.")
		}
		packageName, err = fileop.ReadModPackageName(goModPath)
		if err != nil {
			log.Println("Couldn't read package name from go.mod file. Ignoring.")
		}
	}

	if goPath := os.Getenv("GOPATH"); goPath != "" && len(imports) == 0 {
		// 2) PATH go.mod file DOES NOT EXIST
		if match, _ := regexp.MatchString(goPath+"/src/.*", o.SrcPath); match {
			re := regexp.MustCompile(goPath + "/src/(.*)")
			res := re.FindStringSubmatch(o.SrcPath)
			if len(res) > 1 {
				packageName = res[1]
			} else {
				fmt.Println("Can't make assumption on package name, continuing without.")
			}
		} else {
			err := fmt.Errorf("can't run on source folder: '%s'. project is not on $GOPATH", o.SrcPath)
			log.Printf("%s\n", err.Error())
			return err
		}
		files, err := fileop.FilesInPath(o.SrcPath, ".*\\.go")
		if err != nil {
			err := fmt.Errorf("couldn't read files in $GOPATH")
			log.Printf("%s\n", err.Error())
			return err
		}
		for _, f := range files {
			res, err := parseImportsGo(f)
			if err != nil {
				err := fmt.Errorf("couldn't parse go imports from files")
				log.Printf("%s\n", err.Error())
				return err
			}
			for p := range res {
				if !strings.Contains(p, packageName) {
					p = strings.Replace(p, "\"", "", -1)
					imports[p] = ""
				}
			}
		}
	}
	if len(imports) == 0 {
		err := fmt.Errorf("can't run on source folder: '%s'. project doesn't have a go.mod file or $GOPATH is not specified", o.SrcPath)
		log.Printf("%s\n", err.Error())
		return err
	}

	var resultReport licensereport.LicenseReport
	resultReport.ProjectID = packageName
	resultReport.ProjectVersion = o.ProjectVersion

	h := sha256.New()
	h.Write([]byte(resultReport.ProjectID + resultReport.ProjectVersion))
	resultReport.ProjectHash = fmt.Sprintf("%x", (h.Sum(nil)))

	// filter out Standard Library from imports, the rest should be URLs
	for _, k := range stdLibraryList {
		if _, ok := imports[k]; ok {
			// reference standard library in the report
			var res licensereport.LicenseResult
			res.License = licensereport.Licenses["na"]
			res.ProjectID = k
			res.ProjectVersion = imports[k]
			h := sha256.New()
			h.Write([]byte(res.ProjectID + res.ProjectVersion))
			res.ProjectHash = fmt.Sprintf("%x", (h.Sum(nil)))
			resultReport.ValidatedLicenses = append(resultReport.ValidatedLicenses, res)

			// delete it from scan worthy import list
			delete(imports, k)
		}
	}

	var urlMap map[string]string
	urlMap = make(map[string]string)

	for u := range imports {
		var whitelistViolation bool
		whitelistViolation = true
		for _, whitelist := range DefaultWhitelistResources {
			if strings.Contains(u, whitelist) {
				parsedURL, err := url.Parse("https://" + u)
				if err != nil {
					fmt.Printf("not a url: %s", u)
					continue
				}
				urlMap[u] = parsedURL.String()
				whitelistViolation = false // TODO collect all illegal imports
			}
		}
		if whitelistViolation {
			var res licensereport.LicenseResult
			res.License = licensereport.Licenses["na"]
			res.ProjectID = u
			res.ProjectVersion = ""
			h := sha256.New()
			h.Write([]byte(res.ProjectID + res.ProjectVersion))
			res.ProjectHash = fmt.Sprintf("%x", (h.Sum(nil)))
			resultReport.Violations = append(resultReport.Violations, res)
		}
	}

	for u := range urlMap {
		// visitURL(urlMap[u])
		var res licensereport.LicenseResult
		res.License = licensereport.Licenses["na"]
		res.ProjectID = u
		res.ProjectVersion = imports[u]
		h := sha256.New()
		h.Write([]byte(res.ProjectID + res.ProjectVersion))
		res.ProjectHash = fmt.Sprintf("%x", (h.Sum(nil)))
		resultReport.ValidatedLicenses = append(resultReport.ValidatedLicenses, res)
	}

	printReport(resultReport)
	if len(resultReport.Violations) > 0 {
		os.Exit(1)
	}

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

func printReport(rep licensereport.LicenseReport) {
	fmt.Printf("Report for %s %s\n", rep.ProjectID, rep.ProjectVersion)
	fmt.Printf("Generated project hash: %s\n", rep.ProjectHash)
	fmt.Println("")
	if len(rep.ValidatedLicenses) == 1 {
		fmt.Printf("During the scan there was %d external dependency found:\n", len(rep.ValidatedLicenses))
	} else {
		fmt.Printf("During the scan there were %d external dependencies found:\n", len(rep.ValidatedLicenses))
	}
	for _, licen := range rep.ValidatedLicenses {
		fmt.Printf("\tImport: %s, Version: %s\n", licen.ProjectID, licen.ProjectVersion)
	}
	violLen := len(rep.Violations)
	if violLen > 0 {
		if violLen == 1 {
			fmt.Printf("Additionally %d blacklisted import was found:\n", violLen)
		} else {
			fmt.Printf("Additionally %d blacklisted imports were found:\n", violLen)
		}
	}
	for _, viol := range rep.Violations {
		fmt.Printf("\tImport: %s, Version: %s\n", viol.ProjectID, viol.ProjectVersion)
	}
}
