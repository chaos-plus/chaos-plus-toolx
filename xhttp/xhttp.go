package xhttp

import (
	"net/http"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

func GetHeader(r *http.Request, key string) []string {
	for k, v := range r.Header {
		if strings.EqualFold(k, key) {
			return v
		}
	}
	return nil
}

func GetQuery(r *http.Request, key string) []string {
	params := r.URL.Query()
	for k, v := range params {
		if strings.EqualFold(k, key) {
			return v
		}
	}
	return nil
}

func FindParams(r *http.Request, keys ...string) []string {
	set := mapset.NewSet[string]()
	for _, key := range keys {
		if v := GetHeader(r, key); len(v) > 0 {
			set.Append(v...)
		}
		if v := GetQuery(r, key); len(v) > 0 {
			set.Append(v...)
		}
	}
	set.Remove("")
	return set.ToSlice()
}

func GetTokens(r *http.Request) []string {
	return FindParams(r, "Authorization", "X-Token", "Token")
}

func GetTenantId(r *http.Request) string {
	ids := FindParams(r, "X-Tenant-Id", "Tenant-Id", "TenantId")
	if len(ids) > 0 {
		return ids[0]
	} else {
		return ""
	}
}

func GetUserAgent(r *http.Request) string {
	return r.Header.Get("User-Agent")
}

func GetClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
		if ip != "" {
			ips := strings.Split(ip, ",")
			ip = strings.TrimSpace(ips[0])
		}
	}
	if ip == "" {
		ip = r.Header.Get("Proxy-Client-IP")
	}
	if ip == "" {
		ip = r.Header.Get("WL-Proxy-Client-IP")
	}
	if ip == "" {
		ip = r.Header.Get("WL-Proxy-Client-IP")
	}
	if ip == "" {
		ip = r.RemoteAddr
		if idx := strings.Index(ip, ":"); idx != -1 {
			ip = ip[:idx]
		}
	}
	if ip == "" {
		return "unknown"
	}
	return ip
}
