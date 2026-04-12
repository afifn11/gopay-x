package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func ReverseProxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetURL, err := url.Parse(target)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid service URL"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "service temporarily unavailable",
				"service": targetURL.Host,
			})
		}

		// Update request URL
		c.Request.URL.Host = targetURL.Host
		c.Request.URL.Scheme = targetURL.Scheme
		c.Request.Host = targetURL.Host

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}