package uses

import (
	"sgin/common"

	"github.com/gin-gonic/gin"
)

const logKey = "dbLogKey"

type LogTargetTable int

func (s LogTargetTable) Get(c *gin.Context) *dbLog {
	return getDbLog(c).Table(s)
}

const (
	LogTableNone LogTargetTable = 0 // 没有表
)

type dbLog struct {
	c       *gin.Context
	tableId LogTargetTable
}

var gLogHandle common.LogHandler

func (l *dbLog) Table(tableId LogTargetTable) *dbLog {
	return &dbLog{
		c:       l.c,
		tableId: tableId,
	}
}

func (l *dbLog) LogInsert(primaryId int64, data any) *dbLog {
	gLogHandle(l.c, int(l.tableId), common.LogOpInsert, primaryId, data)
	return l
}

func (l *dbLog) LogUpdate(primaryId int64, data any) *dbLog {
	gLogHandle(l.c, int(l.tableId), common.LogOpUpdate, primaryId, data)
	return l
}

func (l *dbLog) Log(primaryId int64, data any) *dbLog {
	gLogHandle(l.c, int(l.tableId), common.LogOpNormal, primaryId, data)
	return l
}

func (l *dbLog) LogDelete(primaryId int64, data any) *dbLog {
	gLogHandle(l.c, int(l.tableId), common.LogOpDelete, primaryId, data)
	return l
}

func (l *dbLog) LogDrop(primaryId int64, data any) *dbLog {
	gLogHandle(l.c, int(l.tableId), common.LogOpDrop, primaryId, data)
	return l
}

func withDbLog(handler common.LogHandler) gin.HandlerFunc {
	gLogHandle = handler
	return func(c *gin.Context) {
		c.Set(logKey, &dbLog{c: c, tableId: LogTableNone})
		c.Next()
	}
}

func getDbLog(c *gin.Context) *dbLog {
	value := c.MustGet(logKey)
	return value.(*dbLog)
}
