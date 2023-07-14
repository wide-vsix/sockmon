package sockmon

import (
	"testing"
)

func TestParse(t *testing.T) {
	s := "UNCONN  1        0            [2001:1::1]:52556     [2001:2::10]:9090   " +
		"ino:0 sk:1dfcc ts sack wscale:7,7 rto:204 " +
		"rtt:0.029/0.016 ato:40 mss:1364 pmtu:1436 rcvmss:536 " +
		"advmss:1364 cwnd:10 bytes_sent:81 bytes_acked:83 bytes_received:191 " +
		"segs_out:7 segs_in:5 data_segs_out:1 data_segs_in:2 send 3762758621bps " +
		"lastsnd:4 pacing_rate 7525517240bps delivery_rate 682000000bps " +
		"delivered:3 rcv_space:13640 rcv_ssthresh:64172 minrtt:0.008"

	sock, err := ParseSsOutput(s)
	if err != nil {
		panic(err)
	}

	if !isValidOutput(s, sock) {
		t.Errorf("test-fail value=%+v", sock)
	}
}