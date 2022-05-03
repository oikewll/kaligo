package database

import (
    "errors"
    "time"

    klogs "github.com/owner888/kaligo/logs"
)

// logger 输出 database 日志
type logger struct {
    klogs.Logger
    SlowThreshold             time.Duration // 慢查询阈值
    IgnoreRecordNotFoundError bool
}

// logs 默认的日志输出
var logs = &logger{
    klogs.New("DB", klogs.LevelDefault, nil),
    200 * time.Millisecond,
    true,
}

// Trace 跟踪 SQL 执行时间
func (l *logger) Trace(db *DB, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
    elapsed := time.Since(begin)
    sql, rows := fc()
    if err != nil && (!errors.Is(err, ErrRecordNotFound) || !l.IgnoreRecordNotFoundError) {
        l.Errorf("%v %d %s %s", elapsed, rows, err, sql)
    } else if elapsed > l.SlowThreshold && l.SlowThreshold != 0 {
        l.Warnf("%v %d %s", elapsed, rows, sql)
    } else {
        l.Infof("%v %d %s", elapsed, rows, sql)
    }
}
