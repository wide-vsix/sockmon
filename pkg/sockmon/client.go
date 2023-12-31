package sockmon

import (
	"encoding/json"
	"io"
	"net/http"
	"net/netip"
)

func GetLocalCache() (map[string]Socket, error) {
	resp, err := http.Get("http://localhost:8931")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	m := map[string]Socket{}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func GetLocalCacheByTuple(fiveTupleStr string) (*Socket, error) {
	cache, err := GetLocalCache()
	if err != nil {
		return nil, err
	}

	val, ok := cache[fiveTupleStr]
	if !ok {
		return nil, nil
	}
	return &val, nil
}

func GetLocalCacheByDst(dst string) ([]Socket, error) {
	socks := []Socket{}
	if cache, err := GetLocalCache(); err == nil {
		for _, val := range cache {
			dstaddr, _ := netip.ParseAddr(dst)
			if val.Dst == dstaddr {
				socks = append(socks, val)
			}
		}
	}
	return socks, nil
}

func GetLocalCacheByPrefix(prefix netip.Prefix) ([]Socket, error) {
	socks := []Socket{}
	if cache, err := GetLocalCache(); err == nil {
		for _, val := range cache {
			if prefix.Contains(val.Dst) {
				socks = append(socks, val)
			}
		}
	}
	return socks, nil
}
