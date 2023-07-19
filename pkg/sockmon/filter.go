package sockmon

import (
	"net/netip"
	"net/url"
	"strconv"
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
