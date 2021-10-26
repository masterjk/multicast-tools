
# multicast-tools

Tools to join a source-specific multicast, and send/receive payload.

## Usage

### Via Docker

Sender:
```
$ docker run --net=host \
	josephkiok/multicast-sender:latest /sender \
	-ifname eth0 \
	-group-address 232.5.5.5 \
	-port 12000 \
	-source-ip 10.10.10.10 \
	-message 'Hello!'
```

Receiver:
```
$ docker run --net=host \
	josephkiok/multicast-receiver:latest /receiver \
	-ifname eth0 \
	-group-address 232.5.5.5 \
	-port 12000 \
	-source-ip 10.10.10.10
```

### Developing Via Go

```
$ go run cmd/sender/sender.go
Starting multicast-sender...

Usage of sender:
  -group-address string
    	multicast group address (range: 232.0.0.0/8)
  -ifname string
    	interface name (ex: eth0)
  -message string
    	message string to multicast
  -port int
    	multicast port
  -source-ip string
    	multicast source IP
```

```
$ go run cmd/receiver/receiver.go
Starting multicast-receiver...

Usage of receiver:
  -group-address string
    	multicast group address (range: 232.0.0.0/8)
  -ifname string
    	interface name (ex: eth0)
  -port int
    	multicast port
  -source-ip string
    	multicast source IP
```