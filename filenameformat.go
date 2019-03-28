package log4go

import (
	"fmt"
	"os"
	"time"
	"bytes"
	"strings"
)

/*
henly:
fname增加一种自定义方式，如果判断带%就是新的方式，
%P : 替换成pid
%T : 替换成当前时间
 */

func (w *FileLogWriter) isFileNameFormat() bool {
	if index := strings.Index(w.filename_format, "%"); index == -1 {
		return false
	}
	return true
}

func (w *FileLogWriter) buildFileNameByFileNameFormat() string {
	pieces := bytes.Split([]byte(w.filename_format), []byte{'%'})

	out := bytes.NewBuffer(make([]byte, 0, 64))
	for i, piece := range pieces {
		if i > 0 && len(piece) > 0 {
			switch piece[0] {
			case 'P':
				// process_id
				out.WriteString(fmt.Sprintf("%d", os.Getpid()))
			case 'T':
				// current time
				out.WriteString(time.Now().Format("2006-01-02_15-04-05"))
			}

			if len(piece) > 1 {
				out.Write(piece[1:])
			}
		} else if len(piece) > 0 {
			out.Write(piece)
		}
	}

	return out.String()
}

// If this is called in a threaded context, it MUST be synchronized
func (w *FileLogWriter) intRotateByFileNameFormat() error {
	// Close any log file that may be open
	if w.file != nil {
		fmt.Fprint(w.file, FormatLogRecord(w.trailer, &LogRecord{Created: time.Now()}))
		w.file.Close()
	}

	// If we are keeping log files, move it to the next available number
	if w.rotate {
		//_, err := os.Lstat(w.filename)
		//if err == nil { // file exists
		//	// Find the next available number
		//	num := 1
		//	fname := ""
		//	if w.daily && time.Now().Day() != w.daily_opendate {
		//		yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		//
		//		for ; err == nil && num <= w.maxbackup; num++ {
		//			fname = w.filename + fmt.Sprintf(".%s.%03d", yesterday, num)
		//			_, err = os.Lstat(fname)
		//		}
		//		// return error if the last file checked still existed
		//		if err == nil {
		//			return fmt.Errorf("Rotate: Cannot find free log number to rename %s\n", w.filename)
		//		}
		//	} else {
		//		num = w.maxbackup - 1
		//		for ; num >= 1; num-- {
		//			fname = w.filename + fmt.Sprintf(".%d", num)
		//			nfname := w.filename + fmt.Sprintf(".%d", num+1)
		//			_, err = os.Lstat(fname)
		//			if err == nil {
		//				os.Rename(fname, nfname)
		//			}
		//		}
		//	}
		//
		//	w.file.Close()
		//	// Rename the file to its newfound home
		//	err = os.Rename(w.filename, fname)
		//	if err != nil {
		//		return fmt.Errorf("Rotate: %s\n", err)
		//	}
		//}

		w.filename = w.buildFileNameByFileNameFormat()
	}

	// Open the log file
	fd, err := os.OpenFile(w.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	w.file = fd

	now := time.Now()
	fmt.Fprint(w.file, FormatLogRecord(w.header, &LogRecord{Created: now}))

	// Set the daily open date to the current date
	w.daily_opendate = now.Day()

	// initialize rotation values
	w.maxlines_curlines = 0
	w.maxsize_cursize = 0

	return nil
}