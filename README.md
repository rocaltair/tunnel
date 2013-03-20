tunnel
======

a cryptor based on RC4 algorithm for network transmitting.

if you got something to transfer using insecure protocol,
such as ftp or nc, you can use this for secure transmitting
with nothing changing.

--------------------------------------------------------
sample

client:

	./tunnel -f 5050 -t yourhost:8080 -p passwd &

	tar -czf - yourfiles | nc localhost 5050 

server:

	./tunnel -f 8080 -t 3030 -p passwd -s &

	nc -l 3030 | tar -xzf - 

datas between client and server:8080 are encrypted

--------------------------------------------------------
sample

ssh -DN 7070 username@vpshost &

it can't be a good enough way to establish a proxy, for GFW can sniff that out.

a better way:

server :

	ssh -DN localhost:7070 username@localhost &

	./tunnel -f localhost:7070 -t 22 -p passwd -s

client :

	./tunnel -f 6060 -t vpshost:7070 -p passwd

now you can get a proxy connetion based on socks 5 on port 6060 to avoid GFW's sniffing

better, you need to change some configurations for sshd,
change default Port 22 to Port 127.0.0.1:22(local connections only) or other(s)

---------------------------------

Usage of tunnel:

	-d=false: using debug mode

	-f="127.0.0.1:6060": endpoint where stream from,fmt:'host:port|:port|port'

	-p="passwd": your passwd

	-s=false: as a server

	-t="127.0.0.1:5432": endpoint where stream to,fmt:'host:port|:port|port'

enjoy it!

