// Package report implements the `lic report golang` (`lic r go`)command.
package report

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/tehcyx/lic/internal/golang/godep"
	"github.com/tehcyx/lic/internal/golang/gomod"
	"github.com/tehcyx/lic/internal/golang/gopath"
	"github.com/tehcyx/lic/internal/report"

	"github.com/spf13/cobra"

	"github.com/tehcyx/lic/internal/fileop"
	"github.com/tehcyx/lic/pkg/lic/core"
)

var (
	//DefaultWhitelistResources default list of acceptable imports that will get auto-parsed and checked for licenses
	DefaultWhitelistResources = []string{"github.com", "gopkg.in", "golang.org"}

	//StdLibraryList list of Standard Library imports as of go 1.11.5
	stdLibraryList = map[string]string{
		"archive": "", "archive/tar": "", "archive/zip": "", "bufio": "", "builtin": "", "bytes": "", "compress": "", "compress/bzip2": "", "compress/flate": "",
		"compress/gzip": "", "compress/lzw": "", "compress/zlib": "", "container	": "", "container/heap": "", "container/list": "", "container/ring": "",
		"context": "", "crypto": "", "crypto/aes": "", "crypto/cipher": "", "crypto/des": "", "crypto/dsa": "", "crypto/ecdsa": "", "crypto/elliptic": "",
		"crypto/hmac": "", "crypto/md5": "", "crypto/rand": "", "crypto/rc4": "", "crypto/rsa": "", "crypto/sha1": "", "crypto/sha256": "", "crypto/sha512": "",
		"crypto/subtle": "", "crypto/tls": "", "crypto/x509": "", "crypto/x509/pkix": "", "database": "", "database/sql": "", "database/sql/driver": "",
		"debug": "", "debug/dwarf": "", "debug/elf": "", "debug/gosym": "", "debug/macho": "", "debug/pe": "", "debug/plan9obj": "", "encoding": "",
		"encoding/ascii85": "", "encoding/asn1": "", "encoding/base32": "", "encoding/base64": "", "encoding/binary": "", "encoding/csv": "",
		"encoding/gob": "", "encoding/hex": "", "encoding/json": "", "encoding/pem": "", "encoding/xml": "", "errors": "", "expvar": "", "flag": "", "fmt": "",
		"go": "", "go/ast": "", "go/build": "", "go/constant": "", "go/doc": "", "go/format": "", "go/importer": "", "go/parser": "", "go/printer": "",
		"go/scanner": "", "go/token": "", "go/types": "", "hash": "", "hash/adler32": "", "hash/crc32": "", "hash/crc64": "", "hash/fnv": "", "html": "",
		"html/template": "", "image": "", "image/color": "", "image/palette": "", "image/draw": "", "image/gif": "", "image/jpeg": "", "image/png": "", "index": "",
		"index/suffixarray": "", "io": "", "io/ioutil": "", "log": "", "log/syslog": "", "math": "", "math/big": "", "math/bits": "", "math/cmplx": "", "math/rand": "",
		"mime": "", "mime/multipart": "", "mime/quotedprintable": "", "net": "", "net/http": "", "net/http/cgi": "", "net/http/cookiejar": "", "net/http/fcgi": "",
		"net/http/httptest": "", "net/http/httptrace": "", "net/http/httputil": "", "net/http/pprof": "", "net/mail": "", "net/rpc": "", "net/rpc/jsonrpc": "",
		"net/smtp": "", "net/textproto": "", "net/url": "", "os": "", "os/exec": "", "os/signal": "", "os/user": "", "path": "", "path/filepath": "", "plugin": "",
		"reflect": "", "regexp": "", "regexp/syntax": "", "runtime": "", "runtime/cgo": "", "runtime/debug": "", "runtime/msan": "", "runtime/pprof": "",
		"runtime/race": "", "runtime/trace": "", "sort": "", "strconv": "", "strings": "", "sync": "", "sync/atomic": "", "syscall": "", "syscall/js": "", "testing": "",
		"testing/iotest": "", "testing/quick": "", "text": "", "text/scanner": "", "text/tabwriter": "", "text/template": "", "text/template/parse": "",
		"time": "", "unicode": "", "unicode/utf16": "", "unicode/utf8": "", "unsafe": "",
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
	cmd.Flags().StringVarP(&o.ProjectName, "project-name", "", "", "Name of scan target")

	cmd.Flags().BoolVarP(&o.StdLib, "stdlib", "s", true, "Should go dependencies be part of the output (default: true)")

	return cmd
}

//Run runs the command
// Scan has paths this could go:
//		1) If there's a go.mod file, check go.mod file for dependencies and versions of these
//		2) If there's at least one Gopkg.toml/Gopkg.lock file, check Gopkg.lock(s) for all dependencies and versions
//		3) If there's no go.mod file, check $GOPATH and make assumption based on that
func (o *GolangReportOptions) Run() error {
	if o.SrcPath != "" {
		if err := fileop.Exists(o.SrcPath); err != nil {
			err := fmt.Errorf("path '%s' does not exist or you don't have the proper access rights", o.SrcPath)
			log.Printf("%s\n", err.Error())
			os.Exit(1)
		}
	} else {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			err := fmt.Errorf("couldn't get application path")
			log.Printf("%s\n", err.Error())
			os.Exit(1)
		}
		o.SrcPath = dir
	}
	if o.ProjectVersion == "n/a" {
		cmd := exec.Command("git", "describe", "--tags", "--always")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		o.ProjectVersion = out.String()
	}

	proj := report.NewProjectReport()

	// 1) go.mod file EXISTS
	errGoMod := gomod.Collect(proj, o.SrcPath)
	if errGoMod != nil {
		log.Printf("Info: %s. Fallback to godep.", errGoMod.Error())
	}
	// 2) go.mod DOES NOT EXIST but Gopkg.lock exist
	if len(proj.Imports) == 0 {
		errGoDep := godep.Collect(proj, o.SrcPath)
		if errGoDep != nil {
			log.Printf("Info: %s. Fallback to file parsing.", errGoDep.Error())
		}
	}
	// 3) go.mod file or Gopkg.lock DOES NOT EXIST
	if len(proj.Imports) == 0 {
		errGoPath := gopath.Collect(proj, o.SrcPath)
		if errGoPath != nil {
			log.Printf("Info: %s.", errGoPath.Error())
		}
	}

	if len(proj.Imports) == 0 {
		err := fmt.Errorf("can't run on source folder: '%s'", o.SrcPath)
		log.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	h := sha256.New()
	h.Write([]byte(proj.Name + proj.Version))
	proj.Hash = fmt.Sprintf("%x", (h.Sum(nil)))
	proj.Version = o.ProjectVersion

	for _, imp := range proj.Imports {
		imp.License = report.Licenses["na"]

		if _, ok := stdLibraryList[imp.Name]; ok { // only execute if library is not stdlib
			// reference standard library in the report
			imp.Version = "Standard Library"
			proj.ValidatedLicenses[imp.Name] = imp
		} else {
			var whitelistViolation bool
			whitelistViolation = true
			for _, whitelist := range DefaultWhitelistResources {
				if strings.Contains(imp.Name, whitelist) {
					parsedURL, err := url.Parse("https://" + imp.Name)
					if err != nil {
						log.Printf("not a url: %s", imp.Name)
						continue
					}
					imp.ParsedURL = parsedURL.String() //TODO: call url and actually validate License
					proj.ValidatedLicenses[imp.Name] = imp
					whitelistViolation = false //TODO: collect all illegal imports
				}
			}
			if whitelistViolation {
				proj.Violations[imp.Name] = imp
			}
		}

		h := sha256.New()
		h.Write([]byte(imp.Name + imp.Version))
		imp.Hash = fmt.Sprintf("%x", (h.Sum(nil)))
	}

	proj.PrintReport()
	if len(proj.Violations) > 0 {
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
