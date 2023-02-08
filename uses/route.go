package uses

import (
	"jetiny/sgin/common"

	"github.com/gin-gonic/gin"
)

const gRoutesKey = "routesKey"

func withRoute(r *gin.Engine, routes []*common.Route) {
	handler := r.Use(func(c *gin.Context) {
		c.Set(gRoutesKey, routes)
		route := getRoute(c)
		if route != nil {
			if acceptRoute(c, route) {
				c.Next()
			}
		} else {
			c.Next()
		}
	})
	for _, route := range routes {
		handler.Handle(route.Method, route.Path, route.Handle)
	}
}

func acceptRoute(c *gin.Context, route *common.Route) bool {
	if getAppModel(c) == nil {
		c.AbortWithError(gErrAppCodeInvalid.Error().GinError())
		return false
	}
	if route.EnsureAuth {
		tokenValue := getUserToken(c)
		if tokenValue == nil {
			c.AbortWithError(gErrAuthTokenExpired.Error().GinError())
			return false
		}
		if tokenValue.IsTokenExpired() {
			c.AbortWithError(gErrAuthTokenExpired.Error().GinError())
			return false
		}
	}
	return true
}

const routeKey = "routeKey"

func getRoute(c *gin.Context) *common.Route {
	value, exists := c.Get(routeKey)
	if exists {
		return value.(*common.Route)
	}
	routes := c.MustGet(gRoutesKey).([]*common.Route)
	for _, route := range routes {
		if route.Path == c.Request.URL.Path {
			c.Set(routeKey, route)
			return route
		}
	}
	return nil
}
