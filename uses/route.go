package uses

import (
	"github.com/jetiny/sgin/common"

	"github.com/gin-gonic/gin"
)

const gRoutesKey = "routesKey"

func withRoute(r *gin.Engine, routes []*common.Route, tokenHandlers map[string]common.TokenHandler) {
	handler := r.Use(func(c *gin.Context) {
		c.Set(gRoutesKey, routes)
		route := getRoute(c)
		if route != nil {
			token := c.Request.Header.Get(gEnvHeaderAccessToken.String())
			if token != "" {
				c.Set(tokenKey, tokenHandlers[route.TokenKey](c, token))
			} else {
				c.Set(tokenKey, nil)
			}
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
	if !route.NoAppCode {
		if getAppModel(c) == nil {
			c.AbortWithError(gErrAppCodeInvalid.Error().GinError())
			return false
		}
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
		if tokenValue.TokenKey != route.TokenKey {
			c.AbortWithError(gErrAuthTokenKeyInvalid.Error().GinError())
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
		if route.Path == c.FullPath() && route.Method == c.Request.Method {
			c.Set(routeKey, route)
			return route
		}
	}
	return nil
}
