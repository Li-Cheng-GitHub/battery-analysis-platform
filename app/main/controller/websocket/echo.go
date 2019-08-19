// 用于测试 websocket

package websocket

import (
	"github.com/gin-gonic/gin"
)

func Echo(c *gin.Context) {
	conn, err := upgradeHttpConn(c.Writer, c.Request)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		if err = conn.WriteMessage(t, msg); err != nil {
			c.AbortWithStatus(500)
			return
		}
	}
}
