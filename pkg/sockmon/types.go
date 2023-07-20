package sockmon

import (
	"fmt"
	"net/netip"
	"time"

	"gorm.io/gorm"
)

type Socket struct {
	Timestamp time.Time
	Src       netip.Addr `gorm:"type:inet"`
	Dst       netip.Addr `gorm:"type:inet"`
	Protocol  int
	Sport     int
	Dport     int
	Ext       SocketExtendedInformation
}

func (sock Socket) Key() string {
	return fmt.Sprintf("%s/%s/%d/%d/%d",
		sock.Src, sock.Dst, sock.Protocol, sock.Sport, sock.Dport)
}

func (sock Socket) String() string {
	return fmt.Sprintf("%s > %s proto %d tcp %d > %d RTT:%f",
		sock.Src, sock.Dst, sock.Protocol, sock.Sport, sock.Dport,
		sock.Ext.Rtt)
}

type SocketExtendedInformation struct {
	Ino           int     // "ino:0",
	Sk            string  // "sk:1e3ed",
	Ts            bool    // "ts",
	Sack          bool    // "sack",
	Ecn           bool    // "ecn",
	Ecnseen       bool    // "ecnseen",
	WscaleSnd     int     // "wscale:7,7",
	WscaleRcv     int     // "wscale:7,7",
	Rto           int     // "rto:204",
	Rtt           float32 // "rtt:0.03/0.016", rtt:<rtt>/<rttvar>
	RttVar        float32 // "rtt:0.03/0.016", rtt:<rtt>/<rttvar>
	Reordering    int     // "reordering:3",
	Ato           int     // "ato:40",
	Mss           int     // "mss:1364",
	Pmtu          int     // "pmtu:1436",
	Rcvmss        int     // "rcvmss:536",
	Advmss        int     // "advmss:1364",
	Cwnd          int     // "cwnd:10",
	Ssthresh      int     // "ssthresh:23",
	BytesSent     int     // "bytes_sent:81",
	BytesRetrans  int     // "bytes_retrans:29",
	BytesAcked    int     // "bytes_acked:83",
	BytesReceived int     // "bytes_received:191",
	SegsOut       int     // "segs_out:7",
	SegsIn        int     // "segs_in:5",
	DataSegsOut   int     // "data_segs_out:1",
	DataSegsIn    int     // "data_segs_in:2",
	Send          int64   // "send", "3637333333",
	Lastsnd       int     // "lastsnd:4",
	Lastrcv       int     // "lastrcv:12",
	Lastack       int     // "lastack:20",
	PacingRate    int64   // "pacing_rate", "7244481320",
	DeliveryRate  int64   // "delivery_rate", "574315784",
	Delivered     int     // "delivered:3",
	AppLimited    bool    // "app_limited",
	Busy          int     // "busy:100",
	RwndLimited   string  // "rwnd_limited:12ms(1.1%)",
	ReordSeen     int     // "reord_seen:2",
	Retrans       int     // "retrans:0/1, retrans:<retrans>/<retrans_total>"
	RetransTotal  int     // "retrans:0/1, retrans:<retrans>/<retrans_total>"
	DsackDups     int     // "dsack_dups:1",
	Rcvrtt        float32 // "rcvrtt:19.786",
	RcvSpace      int     // "rcv_space:13640",
	RcvSsthresh   int     // "rcv_ssthresh:64172",
	Minrtt        float32 // "minrtt:0.008",
	Notsent       int     // "notsent:391678",
	Sacked        int     // "sacked:113",
	Lost          int     //lost:1
}

// for gorm
type SockmonStat struct {
	gorm.Model
	Timestamp time.Time
	Src       netip.Addr `gorm:"type:inet"`
	Dst       netip.Addr `gorm:"type:inet"`
	Protocol  int
	Sport     int
	Dport     int
	SocketExtendedInformation
}
