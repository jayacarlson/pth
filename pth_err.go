package pth

import (
	"errors"
	"os"

	"github.com/jayacarlson/dbg"
)

var (
	Err_NotExist    = errors.New("File/dir doesn't exist")
	Err_Permission  = errors.New("Permission Denied")
	Err_ReadOnly    = errors.New("Read only file")
	Err_WriteOnly   = errors.New("Write only file")
	Err_RangeError  = errors.New("Range error")
	Err_IllegalFile = errors.New("Illegal file - bad data")
	Err_InvalidData = errors.New("Invalid data")
)

func fileErr(oerr, err error) error {
	switch t := err.(type) {
	case *os.PathError:
		if os.IsNotExist(err) {
			return Err_NotExist
		}
		if os.IsPermission(err) {
			return Err_Permission
		}
		dbg.Message("errs.fileErr - Path err: %v", oerr)
	default:
		dbg.Message("Err is: %v", t)
	}
	return oerr
}

func ChkFileErr(err error) error {
	if err != nil {
		return fileErr(err, err)
	}
	return err
}
