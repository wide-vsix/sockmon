# Sockmon
Sockmon is CLI tool to collect TCP socket statistics and make them available via REST API and DB,etc. 
This tool just parses the result of the iproute2 ss command,  but it is useful.

The initial implementation was developed by Hiroki Shirokura (a.k.a @slankdev) as a sub tool for another project.

## getting started

- build
```
$ make
```

- start sockmon daemon
```
$ bin/sockmon 
```

- get local cache by another processes via REST API
```
$ curl localhost:8931
yas-nyan@mptcp:~$ curl localhost:8931
{
  "2001:db8::1/64:ff9b::8efa:c60e/6/53066/80": {
    "Timestamp": "2023-07-19T13:29:39.47647492+09:00",
    "Src": "2001:db8::1",
    "Dst": "64:ff9b::8efa:c60e",
    "Protocol": 6,
    "Sport": 53066,
    "Dport": 80,
    "ExtId": 0,
    "Ext": {
      "ID": 0,
      "CreatedAt": "0001-01-01T00:00:00Z",
      "UpdatedAt": "0001-01-01T00:00:00Z",
      "DeletedAt": null,
      "Ino": 0,
      "Sk": "6068",
      "Ts": true,
      "Sack": true,
      "Ecn": false,
      "WscaleSnd": 8,
      "WscaleRcv": 7,
      "Rto": 208,
      "Rtt": 6.968,
      "RttVar": 2.283,
      "Reordering": -1,
      "Ato": 40,
      "Mss": 1400,
      "Pmtu": 9000,
      "Rcvmss": 1400,
      "Advmss": 8928,
      "Cwnd": 10,
      "Ssthresh": -1,
      "BytesSent": 79,
      "BytesRetrans": -1,
      "BytesAcked": 81,
      "BytesReceived": 20358,
      "SegsOut": 20,
      "SegsIn": 19,
      "DataSegsOut": 1,
      "DataSegsIn": 16,
      "Send": 16073479,
      "Lastsnd": 540,
      "Lastrcv": 8,
      "Lastack": -1,
      "PacingRate": 32142920,
      "DeliveryRate": 1701344,
      "Delivered": 3,
      "AppLimited": false,
      "Busy": 12,
      "RwndLimited": "-1",
      "ReordSeen": -1,
      "Retrans": -1,
      "RetransTotal": -1,
      "DsackDups": -1,
      "Rcvrtt": 30.187,
      "RcvSpace": 56608,
      "RcvSsthresh": 56608,
      "Minrtt": 6.247
    }
  }
}

```

- you can filter record with query params.

```
- 
curl localhost:8931/?dport=5201
curl localhost:8931/?src=2001:db8::1
curl localhost:8931/?src=2001:db8::1&dst=64:ff9b::192.0.2.1
```

## config and args

- args
```
$ sockmon -h
Usage:
  sockmon [flags]

Flags:
  -b, --bind-address string   Use: sockmon --bind-address <Address:Port> or sockmon -b <Address:Port>  (default ":8931")
  -c, --config string         Use: sockmon --config <CONFIG_PATH> or sockmon -c <CONFIG_PATH>  Various file formats such as YML, TOML and JSON are available.
  -D, --debug                 Use: sockmon --debug or sockmon -D to enable debug mode
  -d, --dump-file string      Use: sockmon --dump-file <FILENAME> or sockmon -d <FILENAME> (by default, it does not dump to file.) 
  -e, --error-file string     Use: sockmon --error-file <FILENAME> or sockmon -e <FILENAME> (by default, it does not dump to file.) 
  -f, --filter string         Use: sockmon --filter '<FILTER>' or sockmon -f '<FILTER>' ss filter.  Please take a look at the iproute2 official documentation. e.g. dport = :80 
  -h, --help                  help for sockmon
  -p, --postgres string       Use: sockmon --postgres 'postgres://user:password@localhost:5432/dbname' or sockmon -p 'postgres://user:password@localhost:5432/dbname' 
  -v, --version               version for sockmon
```

- config format
```yml
dump-file: /var/log/sockmon.log
error-file: /var/log/sockmon.error.log
bind-address: '[::]:8931'
postgres: 'postgres://DB_USER:DB_PASSWORD@localhost:5432/DB_NAME'
filter: 'dport = :443'
debug: true
```