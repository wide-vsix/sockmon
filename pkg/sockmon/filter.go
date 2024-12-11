package sockmon

import (
	"net/netip"
	"net/url"
	"strconv"
	"time"
)

func FilterByParams(params url.Values, socks map[string]Socket) map[string]Socket {
	c := map[string]Socket{}
	for k, v := range socks {
		c[k] = v
	}
	for k, v := range params {
		switch k {
		case "src":
			src, err := netip.ParseAddr(v[0])
			if err != nil {
				continue
			}
			c = FilterLocalCacheBySrc(src, c)
		case "dst":
			dst, err := netip.ParseAddr(v[0])
			if err != nil {
				continue
			}
			c = FilterLocalCacheByDst(dst, c)
		case "sport":
			sport, err := strconv.Atoi(v[0])
			if err != nil {
				continue
			}
			c = FilterLocalCacheBySport(sport, c)
		case "dport":
			dport, err := strconv.Atoi(v[0])
			if err != nil {
				continue
			}
			c = FilterLocalCacheByDport(dport, c)
		case "recent_seconds":
			seconds, err := strconv.Atoi(v[0])
			if err != nil {
				continue
			}
			if seconds < 0 {
				continue
			}
			c = FilterLocalCacheByRecentSeconds(seconds, c)
		case "inbound":
			c = FilterLocalCacheyTrafficDirection(true, c)
		case "outbound":
			c = FilterLocalCacheyTrafficDirection(false, c)
		}
	}
	return c
}

func FilterLocalCacheBySrc(src netip.Addr, sockCache map[string]Socket) map[string]Socket {
	socks := map[string]Socket{}
	for key, val := range sockCache {
		if val.Src == src {
			socks[key] = val
		}
	}
	return socks
}
func FilterLocalCacheByDst(dst netip.Addr, sockCache map[string]Socket) map[string]Socket {
	socks := map[string]Socket{}
	for key, val := range sockCache {
		if val.Dst == dst {
			socks[key] = val
		}
	}
	return socks
}
func FilterLocalCacheBySport(sport int, sockCache map[string]Socket) map[string]Socket {
	socks := map[string]Socket{}
	for key, val := range sockCache {
		if val.Sport == sport {
			socks[key] = val
		}
	}
	return socks
}
func FilterLocalCacheByDport(dport int, sockCache map[string]Socket) map[string]Socket {
	socks := map[string]Socket{}
	for key, val := range sockCache {
		if val.Dport == dport {
			socks[key] = val
		}
	}
	return socks
}

func FilterLocalCacheByRecentSeconds(seconds int, sockCache map[string]Socket) map[string]Socket {
	socks := map[string]Socket{}
	duration := time.Duration(seconds) * time.Second
	cutoff := time.Now().Add(-duration)
	for key, val := range sockCache {
		if val.Timestamp.After(cutoff) {
			socks[key] = val
		}
	}
	return socks
}

func FilterLocalCacheyTrafficDirection(inbound bool, sockCache map[string]Socket) map[string]Socket {
	socks := make(map[string]Socket)
	for key, val := range sockCache {
		if inbound {
			// inbound: DataSegsOut > DataSegsIn
			if val.Ext.DataSegsOut > val.Ext.DataSegsIn {
				socks[key] = val
			}
		} else {
			// outbound: DataSegsIn > DataSegsOut
			if val.Ext.DataSegsIn > val.Ext.DataSegsOut {
				socks[key] = val
			}
		}
		// DataSegsIn == DataSegsOut の場合は何もしない (分類しない)
	}
	return socks
}
