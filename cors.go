package webapp

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type (
	// CORSConfig defines the config for CORS middleware.
	CORSConfig struct {

		// AllowOrigin defines a list of origins that may access the resource.
		// Optional. Default value []string{"*"}.
		AllowOrigins []string `yaml:"allow_origins"`

		// AllowOriginFunc is a custom function to validate the origin. It takes the
		// origin as an argument and returns true if allowed or false otherwise. If
		// an error is returned, it will panic the handler. If this option is
		// set, AllowOrigins is ignored.
		// Optional.
		AllowOriginFunc func(origin string) (bool, error) `yaml:"allow_origin_func"`

		// AllowMethods defines a list methods allowed when accessing the resource.
		// This is used in response to a preflight request.
		// Optional. Default value defaultCORSConfig.AllowMethods.
		AllowMethods []string `yaml:"allow_methods"`

		// AllowHeaders defines a list of request headers that can be used when
		// making the actual request. This is in response to a preflight request.
		// Optional. Default value []string{}.
		AllowHeaders []string `yaml:"allow_headers"`

		// AllowCredentials indicates whether or not the response to the request
		// can be exposed when the credentials flag is true. When used as part of
		// a response to a preflight request, this indicates whether or not the
		// actual request can be made using credentials.
		// Optional. Default value false.
		// Security: avoid using `AllowCredentials = true` with `AllowOrigins = *`.
		// See http://blog.portswigger.net/2016/10/exploiting-cors-misconfigurations-for.html
		AllowCredentials bool `yaml:"allow_credentials"`

		// ExposeHeaders defines a whitelist headers that clients are allowed to
		// access.
		// Optional. Default value []string{}.
		ExposeHeaders []string `yaml:"expose_headers"`

		// MaxAge indicates how long (in seconds) the results of a preflight request
		// can be cached.
		// Optional. Default value 0.
		MaxAge int `yaml:"max_age"`
	}
)

var (
	// defaultCORSConfig is the default CORS middleware config.
	defaultCORSConfig = CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}
)

// CORS returns a Cross-Origin Resource Sharing (CORS) middleware.
// See: https://developer.mozilla.org/en/docs/Web/HTTP/Access_control_CORS
func cors() http.Handler {
	return corsWithConfig(defaultCORSConfig)
}

// CORSWithConfig returns a CORS middleware with config.
// See: `CORS()`.
func corsWithConfig(config CORSConfig) http.Handler {
	if len(config.AllowOrigins) == 0 {
		config.AllowOrigins = defaultCORSConfig.AllowOrigins
	}

	if len(config.AllowMethods) == 0 {
		config.AllowMethods = defaultCORSConfig.AllowMethods
	}

	allowOriginPatterns := []string{}
	for _, origin := range config.AllowOrigins {
		pattern := regexp.QuoteMeta(origin)
		pattern = strings.Replace(pattern, "\\*", ".*", -1)
		pattern = strings.Replace(pattern, "\\?", ".", -1)
		pattern = "^" + pattern + "$"
		allowOriginPatterns = append(allowOriginPatterns, pattern)
	}

	allowMethods := strings.Join(config.AllowMethods, ",")
	allowHeaders := strings.Join(config.AllowHeaders, ",")
	exposeHeaders := strings.Join(config.ExposeHeaders, ",")
	maxAge := strconv.Itoa(config.MaxAge)

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		origin := req.Header.Get("Origin")
		allowOrigin := ""

		res.Header().Add("Vary", "Origin")

		// No Origin provided. This is (probably) not request from actual browser - stop returning no content
		if origin == "" {
			res.WriteHeader(http.StatusNoContent)
			return
		}

		if config.AllowOriginFunc != nil {
			allowed, err := config.AllowOriginFunc(origin)
			if err != nil {
				panic(err)
			}
			if allowed {
				allowOrigin = origin
			}
		} else {
			// Check allowed origins
			for _, o := range config.AllowOrigins {
				if o == "*" && config.AllowCredentials {
					allowOrigin = origin
					break
				}
				if o == "*" || o == origin {
					allowOrigin = o
					break
				}
				if matchSubdomain(origin, o) {
					allowOrigin = origin
					break
				}
			}

			checkPatterns := false
			if allowOrigin == "" {
				// to avoid regex cost by invalid (long) domains (253 is domain name max limit)
				if len(origin) <= (253+3+5) && strings.Contains(origin, "://") {
					checkPatterns = true
				}
			}
			if checkPatterns {
				for _, re := range allowOriginPatterns {
					if match, _ := regexp.MatchString(re, origin); match {
						allowOrigin = origin
						break
					}
				}
			}
		}

		// Origin not allowed
		if allowOrigin == "" {
			res.WriteHeader(http.StatusNoContent)
			return
		}

		res.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		if config.AllowCredentials {
			res.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Preflight request
		res.Header().Add("Vary", "Access-Control-Request-Method")
		res.Header().Add("Vary", "Access-Control-Request-Headers")

		if allowMethods != "" {
			res.Header().Set("Access-Control-Allow-Methods", allowMethods)
		}

		if exposeHeaders != "" {
			res.Header().Set("Access-Control-Expose-Headers", exposeHeaders)
		}

		if allowHeaders != "" {
			res.Header().Set("Access-Control-Allow-Headers", allowHeaders)
		} else {
			h := req.Header.Get("Access-Control-Request-Headers")
			if h != "" {
				res.Header().Set("Access-Control-Allow-Headers", h)
			}
		}
		if config.MaxAge > 0 {
			res.Header().Set("Access-Control-Max-Age", maxAge)
		}

		res.WriteHeader(http.StatusNoContent)
	})
}

func matchScheme(domain, pattern string) bool {
	didx := strings.Index(domain, ":")
	pidx := strings.Index(pattern, ":")
	return didx != -1 && pidx != -1 && domain[:didx] == pattern[:pidx]
}

// matchSubdomain compares authority with wildcard
func matchSubdomain(domain, pattern string) bool {
	if !matchScheme(domain, pattern) {
		return false
	}
	didx := strings.Index(domain, "://")
	pidx := strings.Index(pattern, "://")
	if didx == -1 || pidx == -1 {
		return false
	}
	domAuth := domain[didx+3:]
	// to avoid long loop by invalid long domain
	if len(domAuth) > 253 {
		return false
	}
	patAuth := pattern[pidx+3:]

	domComp := strings.Split(domAuth, ".")
	patComp := strings.Split(patAuth, ".")
	for i := len(domComp)/2 - 1; i >= 0; i-- {
		opp := len(domComp) - 1 - i
		domComp[i], domComp[opp] = domComp[opp], domComp[i]
	}
	for i := len(patComp)/2 - 1; i >= 0; i-- {
		opp := len(patComp) - 1 - i
		patComp[i], patComp[opp] = patComp[opp], patComp[i]
	}

	for i, v := range domComp {
		if len(patComp) <= i {
			return false
		}
		p := patComp[i]
		if p == "*" {
			return true
		}
		if p != v {
			return false
		}
	}
	return false
}
