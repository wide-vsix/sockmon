package sockmon

import (
	"encoding/json"
	"fmt"
	"net/netip"
	"os"
	"strconv"
	"strings"
	"time"
)

func ParseSsOutput(in string) (Socket, error) {
	sock := initializeSocket()
	items := strings.Fields(in)
	if len(items) < 5 {
		return sock, fmt.Errorf("invalid format type1")
	}

	// misc local functions
	pPort := func(in string) (int, error) {
		tmp := strings.Split(in, ":")
		portstr := tmp[len(tmp)-1]
		port, _ := strconv.Atoi(portstr)
		if port < 1 || port > 65535 {
			return 0, fmt.Errorf("Invalid IP input detected : %s", in)
		}
		return port, nil
	}
	pAddr := func(in string) (netip.Addr, error) {
		switch in {
		case "Local":
			return netip.Addr{}, fmt.Errorf("Invalid IP input detected : %s", in)
		default:
			idx := strings.LastIndex(in, ":")
			if idx < 0 {
				log.Fatal("Invalid ss output")
			}
			addr := in[:idx]
			addr = strings.Replace(addr, "[", "", -1)
			addr = strings.Replace(addr, "]", "", -1)
			ipa, err := netip.ParseAddr(addr)
			if err != nil {
				return netip.Addr{}, fmt.Errorf("Invalid IP input detected : %s", in)
			}
			return ipa, nil
		}
	}

	sock.Protocol = 6 // TODO(slankdev): this cli only works for tcp

	src, err := pAddr(items[3])
	if err != nil {
		return sock, err
	}
	sock.Src = src

	dst, err := pAddr(items[4])
	if err != nil {
		return sock, err
	}
	sock.Dst = dst

	sport, err := pPort(items[3])
	if err != nil {
		return sock, err
	}
	sock.Sport = sport

	dport, err := pPort(items[4])
	if err != nil {
		return sock, err
	}
	sock.Dport = dport

	items = items[5:]

	pStr1 := func(item string) string { return strings.Split(item, ":")[1] }
	pInt1 := func(item string) int {
		val, err := strconv.Atoi(strings.Split(item, ":")[1])
		if err != nil {
			log.Fatalf("Invalid ss input. err: %s", err)
		}
		return val
	}
	pFloat1 := func(item string) float32 {
		val, err := strconv.ParseFloat(strings.Split(item, ":")[1], 32)
		if err != nil {
			log.Fatalf("Invalid ss input. err: %s", err)
		}
		return float32(val)
	}
	pInt2 := func(item string) (int, int) {
		item2 := strings.Split(item, ":")[1]
		item3 := strings.Split(item2, ",")
		val0, err0 := strconv.Atoi(item3[0])
		val1, err1 := strconv.Atoi(item3[1])
		if err0 != nil || err1 != nil {
			log.Fatalf("Invalid ss input. err: %s, err: %s", err0, err1)
		}
		return val0, val1
	}
	pFloat2 := func(item string) (float32, float32) {
		item1 := strings.Split(item, ":")
		if len(item1) < 2 {
			log.Errorf("pFloat2 invalid parse input type1 %s", item)
			return 0.0, 0.0
		}
		item2 := item1[1]
		item3 := strings.Split(item2, "/")
		if len(item3) < 2 {
			log.Errorf("pFloat2 invalid parse input type2 %s\n", item)
			return 0.0, 0.0
		}
		val0, err0 := strconv.ParseFloat(item3[0], 32)
		val1, err1 := strconv.ParseFloat(item3[1], 32)
		if err0 != nil || err1 != nil {
			log.Fatal("Invalid input type2")
		}
		return float32(val0), float32(val1)
	}
	pInt3 := func(item string) (int, int) {
		item2 := strings.Split(item, ":")[1]
		item3 := strings.Split(item2, "/")
		val0, err0 := strconv.Atoi(item3[0])
		val1, err1 := strconv.Atoi(item3[1])
		if err0 != nil || err1 != nil {
			log.Fatal("Invalid input type2")
		}
		return val0, val1
	}
	pInt64bps := func(item string) int64 {
		item = strings.Replace(item, "bps", "", -1)
		val, err := strconv.ParseInt(item, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		return val
	}
	pIntMs := func(item string) int {
		item2 := strings.Split(item, ":")[1]
		item3 := strings.Replace(item2, "ms", "", -1)
		val, err := strconv.Atoi(item3)
		if err != nil {
			log.Errorf("pIntMs invalid parse input type2 %s\n", item)
			return 0
		}
		return val
	}
	parseSendBufLimited := func(item string) (int, float64, error) {
		split := strings.Split(item, "(")
		if len(split) != 2 {
			return 0, 0, fmt.Errorf("invalid format")
		}

		durationStr := strings.TrimSuffix(split[0], "ms")
		duration, err := strconv.Atoi(strings.Split(durationStr, ":")[1])
		if err != nil {
			return 0, 0, err
		}

		percentageStr := strings.TrimSuffix(split[1], "%)")
		percentage, err := strconv.ParseFloat(percentageStr, 64)
		if err != nil {
			return 0, 0, err
		}

		return duration, percentage / 100, nil // convert parcemtage to ratio
	}

	for idx := 0; idx < len(items); idx++ {
		item := items[idx]
		switch {
		case strings.Contains(item, ":"):

			switch {
			case strings.Contains(item, "rcvmss"):
				sock.Ext.Rcvmss = pInt1(item)
			case strings.Contains(item, "advmss"):
				sock.Ext.Advmss = pInt1(item)
			case strings.Contains(item, "cwnd"):
				sock.Ext.Cwnd = pInt1(item)
			case strings.Contains(item, "bytes_sent"):
				sock.Ext.BytesSent = pInt1(item)
			case strings.Contains(item, "bytes_retrans"):
				sock.Ext.BytesRetrans = pInt1(item)
			case strings.Contains(item, "bytes_acked"):
				sock.Ext.BytesAcked = pInt1(item)
			case strings.Contains(item, "bytes_received"):
				sock.Ext.BytesReceived = pInt1(item)
			case strings.Contains(item, "data_segs_out"):
				sock.Ext.DataSegsOut = pInt1(item)
			case strings.Contains(item, "data_segs_in"):
				sock.Ext.DataSegsIn = pInt1(item)
			case strings.Contains(item, "segs_out"):
				sock.Ext.SegsOut = pInt1(item)
			case strings.Contains(item, "segs_in"):
				sock.Ext.SegsIn = pInt1(item)
			case strings.Contains(item, "lastsnd"):
				sock.Ext.Lastsnd = pInt1(item)
			case strings.Contains(item, "lastrcv"):
				sock.Ext.Lastrcv = pInt1(item)
			case strings.Contains(item, "lastack"):
				sock.Ext.Lastack = pInt1(item)
			case strings.Contains(item, "delivered"):
				sock.Ext.Delivered = pInt1(item)
			case strings.Contains(item, "busy"):
				sock.Ext.Busy = pIntMs(item)
			case strings.Contains(item, "rwnd_limited"):
				sock.Ext.RwndLimited = pStr1(item)
			case strings.Contains(item, "reord_seen"):
				sock.Ext.ReordSeen = pInt1(item)
			case strings.Contains(item, "retrans"):
				sock.Ext.Retrans, sock.Ext.RetransTotal = pInt3(item)
			case strings.Contains(item, "reordering"):
				sock.Ext.Reordering = pInt1(item)
			case strings.Contains(item, "dsack_dups"):
				sock.Ext.DsackDups = pInt1(item)
			case strings.Contains(item, "rcv_rtt"):
				sock.Ext.Rcvrtt = pFloat1(item)
			case strings.Contains(item, "rcv_space"):
				sock.Ext.RcvSpace = pInt1(item)
			case strings.Contains(item, "rcv_ssthresh"):
				sock.Ext.RcvSsthresh = pInt1(item)
			case strings.Contains(item, "ssthresh"):
				sock.Ext.Ssthresh = pInt1(item)
			case strings.Contains(item, "minrtt"):
				sock.Ext.Minrtt = pFloat1(item)
			case strings.Contains(item, "ino"):
				sock.Ext.Ino = pInt1(item)
			case strings.Contains(item, "sk"):
				sock.Ext.Sk = pStr1(item)
			case strings.Contains(item, "wscale"):
				sock.Ext.WscaleSnd, sock.Ext.WscaleRcv = pInt2(item)
			case strings.Contains(item, "rto"):
				sock.Ext.Rto = pInt1(item)
			case strings.Contains(item, "rtt"):
				sock.Ext.Rtt, sock.Ext.RttVar = pFloat2(item)
			case strings.Contains(item, "ato"):
				sock.Ext.Ato = pInt1(item)
			case strings.Contains(item, "mss"):
				sock.Ext.Mss = pInt1(item)
			case strings.Contains(item, "pmtu"):
				sock.Ext.Pmtu = pInt1(item)
			case strings.Contains(item, "sacked"):
				sock.Ext.Sacked = pInt1(item)
			case strings.Contains(item, "notsent"):
				sock.Ext.Notsent = pInt1(item)
			case strings.Contains(item, "lost"):
				sock.Ext.Lost = pInt1(item)
			case strings.Contains(item, "sndbuf_limited"):
				duration, ratio, err := parseSendBufLimited(item)
				if err != nil {
					log.Warnf("can not parse: %s\n", item)
					break
				}
				sock.Ext.SendbufLimitedDuration = duration
				sock.Ext.SendbufLimitedRatio = ratio

			default:
				log.Warnf("unknown key-value type %s\n", item)
			}

		case strings.Contains(item, "app_limited"):
			sock.Ext.AppLimited = true
		case item == "ts":
			sock.Ext.Ts = true
		case item == "sack":
			sock.Ext.Sack = true
		case item == "ecn":
			sock.Ext.Ecn = true
		case item == "ecn":
			sock.Ext.Ecn = true
		case item == "ecnseen":
			sock.Ext.Ecnseen = true
		case item == "send":
			sock.Ext.Send = pInt64bps(items[idx+1])
			idx++
		case item == "pacing_rate":
			sock.Ext.PacingRate = pInt64bps(items[idx+1])
			idx++
		case item == "delivery_rate":
			sock.Ext.DeliveryRate = pInt64bps(items[idx+1])
			idx++
		default:
			log.Warnf("unsupport key type  %s\n", item)
		}
	}
	sock.Timestamp = time.Now()
	return sock, nil
}

func ssDumpFile(sock Socket) error {
	s, err := json.Marshal(sock)
	if err != nil {

		log.Errorf("Marshal invalid ss type %s\n", sock)
		return err
	}
	file, err := os.OpenFile(dumpFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Errorf("Cannot open file %s", dumpFilename)
		return err
	}
	defer file.Close()
	fmt.Fprintln(file, string(s))
	return nil
}

// The following cases are supported
// (i) len(cache) < CACHE_SIZE
// (ii) len(cache) == CACHE_SIZE && sock.Key() is not new
// (iii) len(cache) == CACHE_SIZE && sock.Key() is new
func cacheStore(sock Socket) {
	if _, ok := cache[sock.Key()]; !ok && len(cache) == CACHE_SIZE {
		var lruFiveTuple string
		var lruTimestamp time.Time

		// TODO: (kanaya) This block should be refactored
		// Initialize lruTimestamp with current time
		for k, v := range cache {
			lruFiveTuple = k
			lruTimestamp = v.Timestamp
			break
		}

		for k, v := range cache {
			if v.Timestamp.Before(lruTimestamp) {
				lruFiveTuple = k
				lruTimestamp = v.Timestamp
			}
		}

		delete(cache, lruFiveTuple)
	}
	cache[sock.Key()] = sock
}

func initializeSocket() Socket {
	sock := Socket{}

	sock.Timestamp = time.Time{}
	sock.Src = netip.Addr{}
	sock.Dst = netip.Addr{}
	sock.Protocol = -1
	sock.Sport = -1
	sock.Dport = -1
	sock.Ext.Ino = -1
	sock.Ext.Sk = "-1"
	sock.Ext.Ts = false
	sock.Ext.Sack = false
	sock.Ext.Ecn = false
	sock.Ext.WscaleSnd = -1
	sock.Ext.WscaleRcv = -1
	sock.Ext.Rto = -1
	sock.Ext.Rtt = -1
	sock.Ext.RttVar = -1
	sock.Ext.Reordering = -1
	sock.Ext.Ato = -1
	sock.Ext.Mss = -1
	sock.Ext.Pmtu = -1
	sock.Ext.Rcvmss = -1
	sock.Ext.Advmss = -1
	sock.Ext.Cwnd = -1
	sock.Ext.Ssthresh = -1
	sock.Ext.BytesSent = -1
	sock.Ext.BytesRetrans = -1
	sock.Ext.BytesAcked = -1
	sock.Ext.BytesReceived = -1
	sock.Ext.SegsOut = -1
	sock.Ext.SegsIn = -1
	sock.Ext.DataSegsOut = -1
	sock.Ext.DataSegsIn = -1
	sock.Ext.Send = -1
	sock.Ext.Lastsnd = -1
	sock.Ext.Lastrcv = -1
	sock.Ext.Lastack = -1
	sock.Ext.PacingRate = -1
	sock.Ext.DeliveryRate = -1
	sock.Ext.Delivered = -1
	sock.Ext.AppLimited = false
	sock.Ext.Busy = -1
	sock.Ext.RwndLimited = "-1"
	sock.Ext.ReordSeen = -1
	sock.Ext.Retrans = -1
	sock.Ext.RetransTotal = -1
	sock.Ext.DsackDups = -1
	sock.Ext.Rcvrtt = -1
	sock.Ext.RcvSpace = -1
	sock.Ext.RcvSsthresh = -1
	sock.Ext.Minrtt = -1

	return sock
}

func isValidOutput(in string, sock Socket) bool {
	items := strings.Fields(in)
	items = items[5:]

	errString := ""
	for idx := 0; idx < len(items); idx++ {
		item := items[idx]
		switch {
		case strings.Contains(item, ":"):
			switch {
			case strings.Contains(item, "rcv_mss"):
				if sock.Ext.Rcvmss == -1 {
					errString += "rcv_mss "
				}
			case strings.Contains(item, "advmss"):
				if sock.Ext.Advmss == -1 {
					errString += "advmss "
				}
			case strings.Contains(item, "cwnd"):
				if sock.Ext.Cwnd == -1 {
					errString += "cwnd "
				}
			case strings.Contains(item, "bytes_sent"):
				if sock.Ext.BytesSent == -1 {
					errString += "bytes_sent "
				}
			case strings.Contains(item, "bytes_retrans"):
				if sock.Ext.BytesRetrans == -1 {
					errString += "bytes_retrans "
				}
			case strings.Contains(item, "bytes_acked"):
				if sock.Ext.BytesAcked == -1 {
					errString += "bytes_acked "
				}
			case strings.Contains(item, "bytes_received"):
				if sock.Ext.BytesReceived == -1 {
					errString += "bytes_received "
				}
			case strings.Contains(item, "data_segs_out"):
				if sock.Ext.DataSegsOut == -1 {
					errString += "data_segs_out "
				}
			case strings.Contains(item, "data_segs_in"):
				if sock.Ext.DataSegsIn == -1 {
					errString += "data_segs_in "
				}
			case strings.Contains(item, "segs_out"):
				if sock.Ext.SegsOut == -1 {
					errString += "segs_out "
				}
			case strings.Contains(item, "segs_in"):
				if sock.Ext.SegsIn == -1 {
					errString += "segs_in "
				}
			case strings.Contains(item, "lastsnd"):
				if sock.Ext.Lastsnd == -1 {
					errString += "lastsnd "
				}
			case strings.Contains(item, "lastrcv"):
				if sock.Ext.Lastrcv == -1 {
					errString += "lastrcv "
				}
			case strings.Contains(item, "lastack"):
				if sock.Ext.Lastack == -1 {
					errString += "lastack "
				}
			case strings.Contains(item, "delivered"):
				if sock.Ext.Delivered == -1 {
					errString += "delivered "
				}
			case strings.Contains(item, "busy"):
				if sock.Ext.Busy == -1 {
					errString += "busy "
				}
			case strings.Contains(item, "rwnd_limited"):
				if sock.Ext.RwndLimited == "-1" {
					errString += "rwnd_limited "
				}
			case strings.Contains(item, "reord_seen"):
				if sock.Ext.ReordSeen == -1 {
					errString += "reord_seen "
				}
			case strings.Contains(item, "retrans"):
				if sock.Ext.Retrans == -1 || sock.Ext.RetransTotal == -1 {
					errString += "retrans retrans_total "
				}
			case strings.Contains(item, "reordering"):
				if sock.Ext.Reordering == -1 {
					errString += "reordering "
				}
			case strings.Contains(item, "dsack_dups"):
				if sock.Ext.DsackDups == -1 {
					errString += "dsack_dups "
				}
			case strings.Contains(item, "rcv_rtt"):
				if sock.Ext.Rcvrtt == -1 {
					errString += "rcv_rtt "
				}
			case strings.Contains(item, "rcv_space"):
				if sock.Ext.RcvSpace == -1 {
					errString += "rcv_space "
				}
			case strings.Contains(item, "rcv_ssthresh"):
				if sock.Ext.RcvSsthresh == -1 {
					errString += "rcv_ssthresh "
				}
			case strings.Contains(item, "ssthresh"):
				if sock.Ext.Ssthresh == -1 {
					errString += "ssthresh "
				}
			case strings.Contains(item, "minrtt"):
				if sock.Ext.Minrtt == -1 {
					errString += "minrtt "
				}
			case strings.Contains(item, "ino"):
				if sock.Ext.Ino == -1 {
					errString += "ino "
				}
			case strings.Contains(item, "sk"):
				if sock.Ext.Sk == "-1" {
					errString += "sk "
				}
			case strings.Contains(item, "wscale"):
				if sock.Ext.WscaleSnd == -1 || sock.Ext.WscaleRcv == -1 {
					errString += "wcale_snd wscale_rcv "
				}
			case strings.Contains(item, "rto"):
				if sock.Ext.Rto == -1 {
					errString += "rto "
				}
			case strings.Contains(item, "rtt"):
				if sock.Ext.Rtt == -1 || sock.Ext.RttVar == -1 {
					errString += "rtt rttvar "
				}
			case strings.Contains(item, "ato"):
				if sock.Ext.Ato == -1 {
					errString += "ato "
				}
			case strings.Contains(item, "mss"):
				if sock.Ext.Mss == -1 {
					errString += "mss "
				}
			case strings.Contains(item, "pmtu"):
				if sock.Ext.Pmtu == -1 {
					errString += "pmtu "
				}
			}

		case strings.Contains(item, "app_limited"):
			if !sock.Ext.AppLimited {
				errString += "app_limited "
			}
		case item == "ts":
			if !sock.Ext.Ts {
				errString += "ts "
			}
		case item == "sack":
			if !sock.Ext.Sack {
				errString += "sack "
			}
		case item == "ecn":
			if !sock.Ext.Ecn {
				errString += "ecn "
			}
		case item == "send":
			if sock.Ext.Send == -1 {
				errString += "send "
			}
			idx++
		case item == "pacing_rate":
			if sock.Ext.PacingRate == -1 {
				errString += "pacing_rate "
			}
			idx++
		case item == "delivery_rate":
			if sock.Ext.DeliveryRate == -1 {
				errString += "delivery_rate "
			}
			idx++
		}
	}

	if errString != "" {
		// Parse error occured
		if errFilename != "" {
			// dump error file
			file, err := os.OpenFile(errFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				log.Errorf("Cannot open file %s", errFilename)
			}
			defer func() {
				if err := file.Close(); err != nil {
					log.Errorf("Cannot close file %s", errFilename)
				}
			}()
			if errString != "" {
				fmt.Fprintln(file, "--------------------------------------------")
				fmt.Fprintln(file, "------------Original SS Output------------")
				fmt.Fprintln(file, in)
				fmt.Fprintln(file, "------------Output sock------------")
				sockByte, err := json.Marshal(sock)
				if err != nil {
					log.Error("sock json marshal error.")
				}
				fmt.Fprintln(file, string(sockByte))
				fmt.Fprintln(file, "------------Error items------------")
				fmt.Fprintln(file, errString)
				fmt.Fprintln(file, "--------------------------------------------")
			}
		}
	}
	return errString == ""
}
