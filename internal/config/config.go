package config

// Config holds all configurable data for the lic tool
type Config struct {
	Golang  GolangConfig
	License LicenseConfig
}

// GolangConfig holds Golang-specific configuration
type GolangConfig struct {
	// WhitelistDomains is the list of acceptable import domains that will get auto-parsed and checked for licenses
	WhitelistDomains []string
	// StdLibPackages is the map of standard library packages
	StdLibPackages map[string]struct{}
}

// LicenseConfig holds license-related configuration
type LicenseConfig struct {
	// Licenses maps SPDX license identifiers to license information
	// This is currently managed by the license package but could be moved here
	// for centralized configuration
}

// Default returns the default configuration
func Default() *Config {
	return &Config{
		Golang: GolangConfig{
			WhitelistDomains: DefaultWhitelistDomains(),
			StdLibPackages:   DefaultStdLibPackages(),
		},
	}
}

// DefaultWhitelistDomains returns the default list of whitelisted domains
func DefaultWhitelistDomains() []string {
	return []string{
		"github.com",
		"gopkg.in",
		"golang.org",
	}
}

// DefaultStdLibPackages returns the default standard library packages as of Go 1.24
func DefaultStdLibPackages() map[string]struct{} {
	packages := []string{
		// Archive
		"archive", "archive/tar", "archive/zip",
		// Bufio
		"bufio", "builtin", "bytes",
		// Compress
		"compress", "compress/bzip2", "compress/flate", "compress/gzip", "compress/lzw", "compress/zlib",
		// Container
		"container", "container/heap", "container/list", "container/ring",
		// Context
		"context",
		// Crypto
		"crypto", "crypto/aes", "crypto/cipher", "crypto/des", "crypto/dsa", "crypto/ecdsa", "crypto/elliptic",
		"crypto/hmac", "crypto/md5", "crypto/rand", "crypto/rc4", "crypto/rsa", "crypto/sha1", "crypto/sha256",
		"crypto/sha512", "crypto/subtle", "crypto/tls", "crypto/x509", "crypto/x509/pkix",
		// Database
		"database", "database/sql", "database/sql/driver",
		// Debug
		"debug", "debug/dwarf", "debug/elf", "debug/gosym", "debug/macho", "debug/pe", "debug/plan9obj",
		// Embed (Go 1.16+)
		"embed",
		// Encoding
		"encoding", "encoding/ascii85", "encoding/asn1", "encoding/base32", "encoding/base64", "encoding/binary",
		"encoding/csv", "encoding/gob", "encoding/hex", "encoding/json", "encoding/pem", "encoding/xml",
		// Errors
		"errors",
		// Expvar
		"expvar",
		// Flag
		"flag",
		// Fmt
		"fmt",
		// Go
		"go", "go/ast", "go/build", "go/constant", "go/doc", "go/format", "go/importer", "go/parser",
		"go/printer", "go/scanner", "go/token", "go/types",
		// Hash
		"hash", "hash/adler32", "hash/crc32", "hash/crc64", "hash/fnv", "hash/maphash",
		// Html
		"html", "html/template",
		// Image
		"image", "image/color", "image/color/palette", "image/draw", "image/gif", "image/jpeg", "image/png",
		// Index
		"index", "index/suffixarray",
		// Io
		"io", "io/fs", "io/ioutil",
		// Log
		"log", "log/slog", "log/syslog",
		// Math
		"math", "math/big", "math/bits", "math/cmplx", "math/rand",
		// Mime
		"mime", "mime/multipart", "mime/quotedprintable",
		// Net
		"net", "net/http", "net/http/cgi", "net/http/cookiejar", "net/http/fcgi", "net/http/httptest",
		"net/http/httptrace", "net/http/httputil", "net/http/pprof", "net/mail", "net/netip", "net/rpc",
		"net/rpc/jsonrpc", "net/smtp", "net/textproto", "net/url",
		// Os
		"os", "os/exec", "os/signal", "os/user",
		// Path
		"path", "path/filepath",
		// Plugin
		"plugin",
		// Reflect
		"reflect",
		// Regexp
		"regexp", "regexp/syntax",
		// Runtime
		"runtime", "runtime/cgo", "runtime/coverage", "runtime/debug", "runtime/metrics", "runtime/pprof",
		"runtime/trace",
		// Sort
		"sort",
		// Strconv
		"strconv",
		// Strings
		"strings",
		// Sync
		"sync", "sync/atomic",
		// Syscall
		"syscall", "syscall/js",
		// Testing
		"testing", "testing/fstest", "testing/iotest", "testing/quick",
		// Text
		"text", "text/scanner", "text/tabwriter", "text/template", "text/template/parse",
		// Time
		"time",
		// Unicode
		"unicode", "unicode/utf16", "unicode/utf8",
		// Unsafe
		"unsafe",
	}

	// Convert to map for O(1) lookups
	stdLib := make(map[string]struct{}, len(packages))
	for _, pkg := range packages {
		stdLib[pkg] = struct{}{}
	}
	return stdLib
}

// IsStdLib checks if a package is part of the standard library
func (c *GolangConfig) IsStdLib(pkg string) bool {
	_, ok := c.StdLibPackages[pkg]
	return ok
}

// IsWhitelisted checks if an import domain is whitelisted
func (c *GolangConfig) IsWhitelisted(domain string) bool {
	for _, whitelistDomain := range c.WhitelistDomains {
		if domain == whitelistDomain {
			return true
		}
	}
	return false
}
