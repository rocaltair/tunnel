tunnel
======

a cryptor based on RC4 algorithm for network transmission.

if you got something to transfer using unsafe protocol,
such as ftp or nc, you can use this for safe transfering 
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

you can also do some interesting things like below:
--------------------------------------------------------
sample

ssh -DN 7070 username@host &

it can't be a good enough way to proxy, for GFW can sniffer that.

a better way

server :
ssh -DN 7070 username@localhost &
./tunnel -f 127.0.0.1:7070 -t 22 -p passwd -s

client :
./tunnel -f 6060 -t vps_host:7070 -p passwd

now you can get a proxy connetion based on socks 5 on port 6060 avoid GFW's sniffering

for better secury, you need to change some configurations for sshd
change Port 22 to Port 127.0.0.1:22 or other(s)

---------------------------------

Usage of tunnel:
  -d=false: using debug mode
  -f="127.0.0.1:6060": endpoint where stream from,fmt:'host:port|:port|port'
  -p="passwd": your passwd
  -s=false: as a server
  -t="127.0.0.1:5432": endpoint where stream to,fmt:'host:port|:port|port'

enjoy it!

