package pth

import (
	"testing"

	"github.com/jayacarlson/dbg"
)

// These can really only be checked by eye as everyones dir structure is different

func TestPaths(t *testing.T) {
	dbg.Info("-: " + AsRealPath(""))
	dbg.Info(".: " + AsRealPath("."))
	dbg.Info("./:" + AsRealPath("./"))
	dbg.Info("./foo: " + AsRealPath("./foo"))
	dbg.Info("foo: " + AsRealPath("foo"))
	dbg.Info("~: " + AsRealPath("~/foo/bar/../boo/bar"))
	dbg.Info("~: " + AsRealPath("~foo/bar/../boo/bar"))
	dbg.Info("//: " + AsRealPath("~/foo/bar/", "/boo/", "/bar"))
	dbg.Info("^: " + AsRealPath("^"))
	dbg.Info("@: " + AsRealPath("@"))
	dbg.Info("$: " + AsRealPath("$"))
	dbg.Info("p: " + AsRealPath("/test/path/file.ext"))
	dbg.Info("foo/: " + AsRealPath("/test/path/foo/"))
	dbg.Info("x,y,z: " + AsRealPath("foo", "bar", "boo"))

	d, f, e := Split("/test/path/file.ext")
	dbg.Info("Split: /test/path/file.ext")
	dbg.Message("d: %s <- /test/path/", d)
	dbg.Message("f: %s <- file", f)
	dbg.Message("e: %s <- .ext", e)
	d, f, e = Split("/test/path/fileNoExt")
	dbg.Info("Split: /test/path/fileNoExt")
	dbg.Message("d: %s <- /test/path/", d)
	dbg.Message("f: %s <- fileNoExt", f)
	dbg.Message("e: %s <- --no extension--", e)
	d, f, e = Split("TestNoPathOrExt")
	dbg.Info("Split: TestNoPathOrExt")
	dbg.Message("d: %s <- --no dirpath--", d)
	dbg.Message("f: %s <- TestNoPathOrExt", f)
	dbg.Message("e: %s <- --no extension--", e)
	d, f, e = Split("/root/")
	dbg.Info("Split: /root/")
	dbg.Message("d: %s <- /root", d)
	dbg.Message("f: %s <- --no file--", f)
	dbg.Message("e: %s <- --no extension--", e)
	d, f, e = Split("TestJustFile.ext")
	dbg.Info("Split: TestJustFile.ext")
	dbg.Message("d: %s <- --no dirpath--", d)
	dbg.Message("f: %s <- TestJustFile", f)
	dbg.Message("e: %s <- .ext", e)

	dbg.Warning("x: %s <- .ext", Ext("/test/path/file.ext"))
	dbg.Warning("x: %s <- --no extension--", Ext("/test/path/file"))

	SetNumberedPath(1, "$Foobar")
	SetNumberedPath(2, "~Foobar")
	dbg.Info("1: " + AsRealPath("#1/data"))
	dbg.Info("2: " + AsRealPath("#2/text"))

	dbg.Info("Homified ~: %s", Homify("~/testdir"))
	dbg.Info("Homified /: %s", Homify("/root/testdir"))
	dbg.Info("Homified _: %s", Homify("testdir"))

	dbg.Info("Safe ~/...: %s", PathToFilename("~/foo/bar/boo"))
	dbg.Info("Safe ./...: %s", PathToFilename("./foo/bar/boo"))
	dbg.Info("Safe  /...: %s", PathToFilename("/foo/bar/boo"))

}
