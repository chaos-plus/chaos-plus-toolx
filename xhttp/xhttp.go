package xhttp

import (
	"net"
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
	// 1. 从标准代理头中获取
	headers := []string{
		"X-Real-IP",
		"X-Forwarded-For",
		"Proxy-Client-IP",
		"WL-Proxy-Client-IP",
	}

	for _, h := range headers {
		ipList := r.Header.Get(h)
		if ipList == "" {
			continue
		}
		// 取第一个非空 IP（X-Forwarded-For 可能是逗号分隔）
		for _, ip := range strings.Split(ipList, ",") {
			ip = strings.TrimSpace(ip)
			if ip != "" && ip != "unknown" {
				return ip
			}
		}
	}

	// 2. 从 RemoteAddr 获取
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr // fallback 原值
	}

	if ip == "" || ip == "unknown" {
		ip = r.RemoteAddr // fallback 原值
	}

	if ip == "::1" {
		ip = "127.0.0.1"
	}

	return ip
}
