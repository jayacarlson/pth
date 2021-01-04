package pth

import (
	"testing"

	"github.com/jayacarlson/dbg"
)

// These can really only be checked by eye as everyones dir structure is different

func TestPaths(t *testing.T) {
	dbg.Note("-: " + AsRealPath(""))
	dbg.Note(".: " + AsRealPath("."))
	dbg.Note("./:" + AsRealPath("./"))
	dbg.Note("./foo: " + AsRealPath("./foo"))
	dbg.Note("foo: " + AsRealPath("foo"))
	dbg.Note("~: " + AsRealPath("~/foo/bar/../boo/bar"))
	dbg.Note("~: " + AsRealPath("~foo/bar/../boo/bar"))
	dbg.Note("^: " + AsRealPath("^"))
	dbg.Note("@: " + AsRealPath("@"))
	dbg.Note("$: " + AsRealPath("$"))
	dbg.Note("p: " + AsRealPath("/test/path/file.ext"))
	dbg.Note("x,y,z: " + AsRealPath("foo", "bar", "boo"))

	d, f, e := Split("/test/path/file.ext")
	dbg.Message("d: %s <- /test/path/", d)
	dbg.Message("f: %s <- file", f)
	dbg.Message("e: %s <- .ext", e)
	d, f, e = Split("/test/path/fileNoExt")
	dbg.Message("d: %s <- /test/path/", d)
	dbg.Message("f: %s <- fileNoExt", f)
	dbg.Message("e: %s <- --no extension--", e)
	d, f, e = Split("TestNoPathOrExt")
	dbg.Message("d: %s <- --no dirpath--", d)
	dbg.Message("f: %s <- TestNoPathOrExt", f)
	dbg.Message("e: %s <- --no extension--", e)
	d, f, e = Split("TestNoPath.ext")
	dbg.Message("d: %s <- --no dirpath--", d)
	dbg.Message("f: %s <- TestNoPath", f)
	dbg.Message("e: %s <- .ext", e)

	dbg.Warning("x: %s <- .ext", Ext("/test/path/file.ext"))
	dbg.Warning("x: %s <- --no extension--", Ext("/test/path/file"))

	SetNumberedPath(1, "$Foobar")
	SetNumberedPath(2, "~Foobar")
	dbg.Note("1: " + AsRealPath("#1/data"))
	dbg.Note("2: " + AsRealPath("#2/text"))
}
