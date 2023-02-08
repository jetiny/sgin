package uses

import "github.com/gin-gonic/gin"

func withTask(tasks map[string][]gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		r := getRoute(c)
		if r != nil {
			if arr, ok := tasks[r.Name]; ok {
				for _, v := range arr {
					v(c)
				}
			}
		}
	}
}
