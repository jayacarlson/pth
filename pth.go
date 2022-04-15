package pth

import (
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/jayacarlson/dbg"
	"github.com/jayacarlson/env"
)

// ----------------------------------------------------------------------------

/*
	A utility file to allow apps to use consistent internal paths for cross
	platform development.  All paths use '/' as the separator internally,
	being replaced with the platform specific separator after calling the
	'AsRealPath' function.

	Can accept any of the following leading meta path sequences
	   (nothing)  <pwd>                     	empty string supplies Working Directory
	   ./blah     <pwd>/blah					. supplies Working Directory
	   ~/blah     <home>/blah					~ supplies users HOME dir
	   ^/blah     <gopath>/blah					^ supplies the systen GOPATH
	   @/blah     <expath>/blah             	@ supplies the executables path
	   $/blah     <srcpath>/blah            	$ supplies the source file path
	   #N/blah    <numberedPath(0..9)>/blah 	#N supplies numbered path, see SetNumberedPath

	As some examples:		linux					windows
		"~/Documents"		<home>/Documents		c:\User\{user}\Documents
		"./subdir"			<pwd>/subdir			<pwd>\subdir   (e.g.  d:\data\subdir)
*/

var (
	envSet = getEnv() // doing this gets the environment vars before any init() function(s) are called

	numberedPaths [10]string
	paths         struct {
		// Names must be capitalized to be set by the ReadEnvVars
		Home   string // Home directory -- must be retrieved from environment
		GoPath string // Base directory to all go related code
	}
)

// getEnv -- run as variable assignment to be assured it is run before all 'init' methods; as some may call into here
func getEnv() bool {
	env.ReadEnvVars(&paths)

	// validate we have some needed values, supplying defaults as needed
	if paths.Home == "" {
		if env.IsWindows() {
			paths.Home = "c:/Users/" + env.User()
		}

		// -- other systems -- other ways to set 'Home'

		dbg.FatalIf(paths.Home == "", "Init: Home not configured")
	}
	if paths.GoPath == "" {
		// not fatal, goPath isn't vital, just a convenience
		dbg.Note("Init: GoPath not configured")
	} else {
		// fix GoPath incase there are multiple entries -- only use the 1st
		paths.GoPath = strings.Split(paths.GoPath, string(os.PathListSeparator))[0]
	}
	return true
}

// cleanup a dir string, removing needless chars
func cleanupDir(p string) string {
	p = strings.ReplaceAll(p, "/./", "/")
	p = strings.ReplaceAll(p, "//", "/")
	t := len(p)
	if t > 0 && '/' == p[t-1] {
		p = p[:t-1]
	}
	return p
}

// Exists tests if the given path exists.
func Exists(filePath string) bool {
	filePath = AsRealPath(filePath)
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// set the numbered meta path for use with 'AsRealPath'
func SetNumberedPath(n int, path string) {
	dbg.ChkTruX(n < 10 && n >= 0, "Illegal path index (only 0..9)")
	p := asReal(path)
	if '/' != os.PathSeparator {
		p = strings.ReplaceAll(p, "/", string(os.PathSeparator))
	}
	numberedPaths[n] = p
}

/*
	return a real path for the given internal path

	can pass a list of strings for joining; e.g. AsRealPath("fred", "barney", "wilma", "betty")
	and get "fred/barney/wilma/betty"
	-- with 'fred' allowed to have the meta path tokens above
*/
func AsRealPath(l ...string) string {
	p := asReal(strings.Join(l, "/"))
	p = cleanupDir(p)
	if '/' != os.PathSeparator {
		p = strings.Replace(p, "/", string(os.PathSeparator), -1)
	}
	return p
}

// make a directory path, with all parents as needed
func MakePath(path string) error {
	p := asReal(path)
	if '/' != os.PathSeparator {
		p = strings.ReplaceAll(p, "/", string(os.PathSeparator))
	}
	return os.MkdirAll(p, 0775)
}

// returns just the file extention of the path (if it contains one)
func Ext(srcPath string) string {
	return path.Ext(srcPath)
}

// returns dir, filename & ext for the given srcPath
//		e.g. /foo/bar/boo.ext -> "/foo/bar"  "boo"  & ".ext"
//	NOTE: if there is no .extention, the 'file' may be the last dir in a path
func Split(src string) (dir, file, ext string) {
	dir, file = path.Split(src)
	dir = cleanupDir(dir)
	ext = path.Ext(src)
	file = file[:len(file)-len(ext)]
	return
}

// returns a path consisting of the dir/file.ext
func Join(dir, file, ext string) string {
	if file == "" { // silly to pass dir w/o file, but just in case...
		ext = ""
	} else {
		if ext != "" && ext[0] != '.' {
			ext = "." + ext
		}
	}
	return path.Join(dir, file+ext)
}

// returns a path with home '~' if it can
func Homify(p string) string {
	p = AsRealPath(p)
	if len(p) >= len(paths.Home) {
		l := len(paths.Home)
		if p[:l] == paths.Home {
			return "~" + p[l:]
		}
	}
	return p
}

func PathToFilename(p string) string {
	p = strings.ReplaceAll(Homify(p), string(os.PathSeparator), "@")
	if '@' == p[len(p)-1] {
		p = p[:len(p)-1]
	}
	if 2 < len(p) && "~@" == p[:2] {
		p = "~" + p[1:]
	}
	//if '~' == p[0] {
	//	p = "-" + p[1:]
	//}
	return p
}

func asReal(p string) string {
	// internally we use '/' as the separator, later it will be replaced with true path separator
	if len(p) > 0 && p[0] == '/' {
		// shortcircuit true paths '/blah'
		return p
	}
	if p == "" || p == "." || (len(p) > 1 && p[:2] == "./") {
		wd, err := os.Getwd()
		dbg.FatalIfErr(err, "GetWD err: %v\n", err)
		if p == "" {
			return wd
		}
		return path.Join(wd, p[1:])
	} else if p[0] == '~' { // currently not going to deal with ~user/blah
		return path.Join(paths.Home, p[1:])
	} else if p[0] == '$' {
		_, file, _, ok := runtime.Caller(2)
		dbg.FatalIf(!ok, "runtime.Caller failed")
		return path.Join(path.Dir(file), p[1:])
	} else if p[0] == '@' {
		ex, err := os.Executable()
		dbg.FatalIfErr(err, "os.Executable failed")
		return path.Join(path.Dir(ex), p[1:])
	} else if len(p) > 1 && p[0] == '#' {
		n := p[1] - '0'
		dbg.ChkTruX(n < 10 && n >= 0, "Illegal path index (only 0..9)")
		dbg.ChkTruX(numberedPaths[n] != "", "Path index %d unset", n)
		return path.Join(asReal(numberedPaths[n]), p[2:])
	} else if p[0] == '^' {
		return path.Join(paths.GoPath, p[1:])
	}
	// to test on Windows, is this returning os normalized paths?  c:/foo/bar => c:\foo\bar
	return path.Join(p)
}
