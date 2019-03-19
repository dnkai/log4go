// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)
/*
henly
增加json输出方式，json方式增加如下字段
1. 模块名A: 用户传入
2. 唯一标记B: 启动时间+进程id
*/

type (
	GlobalLogInfo struct {
		StartAt time.Time
		Pid     int
		B       string
	}
)

var (
	globalLogInfo GlobalLogInfo
)

func init() {
	globalLogInfo.StartAt = time.Now()
	globalLogInfo.Pid = os.Getpid()

	globalLogInfo.B = fmt.Sprintf("%s-%d",
		globalLogInfo.StartAt.Format("2006-01-02_15:04:05"), globalLogInfo.Pid)
}

// Known format codes:
// %T - Time (15:04:05 MST)
// %t - Time (15:04)
// %D - Date (2006/01/02)
// %d - Date (01/02/06)
// %L - Level (FNST, FINE, DEBG, TRAC, WARN, EROR, CRIT)
// %S - Source
// %M - Message
// Ignores unknown formats
// Recommended: "[%D %T] [%L] (%S) %M"
func FormatLogRecordJson(format string, modulename string, rec *LogRecord) string {
	if rec == nil {
		return "{nil}"
	}
	if len(format) == 0 {
		return "{}"
	}

	out := bytes.NewBuffer(make([]byte, 0, 64))
	secs := rec.Created.UnixNano() / 1e9

	cache := getFormatCache()
	if cache.LastUpdateSeconds != secs {
		month, day, year := rec.Created.Month(), rec.Created.Day(), rec.Created.Year()
		hour, minute, second := rec.Created.Hour(), rec.Created.Minute(), rec.Created.Second()
		updated := &formatCacheType{
			LastUpdateSeconds: secs,
			shortTime:         fmt.Sprintf("%02d:%02d", hour, minute),
			shortDate:         fmt.Sprintf("%02d/%02d/%02d", year%100, month, day),
			longTime:          fmt.Sprintf("%02d:%02d:%02d", hour, minute, second),
			longDate:          fmt.Sprintf("%04d/%02d/%02d", year, month, day),
		}
		cache = updated
		setFormatCache(updated)
	}

	// Split the string into pieces by % signs
	pieces := bytes.Split([]byte(format), []byte{'%'})

	tmp := struct {
		A  string `json:"A,omitempty"`
		B  string `json:"B,omitempty"`
		D  string `json:"D,omitempty"`
		Ds string `json:"d,omitempty"`
		T  string `json:"T,omitempty"`
		Ts string `json:"t,omitempty"`
		L  string `json:"L,omitempty"`
		S  string `json:"S,omitempty"`
		Ss string `json:"s,omitempty"`
		M  string `json:"M,omitempty"`
	}{}

	// Iterate over the pieces, replacing known formats
	for i, piece := range pieces {
		if i > 0 && len(piece) > 0 {
			switch piece[0] {
			case 'A':
				tmp.A = modulename
			case 'B':
				tmp.B = globalLogInfo.B
			case 'T':
				//out.WriteString(cache.longTime)
				tmp.T = cache.longTime
			case 't':
				//out.WriteString(cache.shortTime)
				tmp.Ts = cache.shortTime
			case 'D':
				//out.WriteString(cache.longDate)
				tmp.D = cache.longDate
			case 'd':
				//out.WriteString(cache.shortDate)
				tmp.Ds = cache.shortDate
			case 'L':
				//out.WriteString(levelStrings[rec.Level])
				tmp.L = levelStrings[rec.Level]
			case 'S':
				//out.WriteString(rec.Source)
				tmp.S = rec.Source
			case 's':
				slice := strings.Split(rec.Source, "/")
				//out.WriteString(slice[len(slice)-1])
				tmp.Ss = slice[len(slice)-1]
			case 'M':
				//out.WriteString(rec.Message)
				tmp.M = rec.Message
			}
			if len(piece) > 1 {
				//out.Write(piece[1:])
			}
		} else if len(piece) > 0 {
			//out.Write(piece)
		}
	}
	//out.WriteByte('\n')
	bb, _ := json.Marshal(tmp)
	out.Write(bb)
	out.WriteByte('\n')

	return out.String()
}
