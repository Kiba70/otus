package netstat

import (
	"testing"
	"time"
)

const (
	netstatInput = `Active Internet connections (servers and established)
Proto Recv-Q Send-Q Local Address           Foreign Address         State       User       Inode      PID/Program name    
tcp        0      0 127.0.0.1:46749         0.0.0.0:*               LISTEN      hev        184402348  1793155/code-384ff7 
tcp        0      0 0.0.0.0:10050           0.0.0.0:*               LISTEN      zabbix     175588060  3664077/zabbix_agen 
tcp        0      0 0.0.0.0:111             0.0.0.0:*               LISTEN      root       20470      1/systemd           
tcp        0      0 127.0.0.1:53            0.0.0.0:*               LISTEN      root       26521      1041/dnsmasq        
tcp        0      0 0.0.0.0:22              0.0.0.0:*               LISTEN      root       26580      1044/sshd           
tcp        0    320 10.125.81.41:22         10.74.156.153:48236     ESTABLISHED root       188078594  2945687/sshd: hev  
tcp        0      0 10.125.81.41:22         10.74.156.153:51896     ESTABLISHED root       184402075  1793091/sshd: hev  
tcp        0      0 127.0.0.1:46749         127.0.0.1:49260         ESTABLISHED hev        184402378  1793155/code-384ff7 
tcp        0      0 127.0.0.1:49260         127.0.0.1:46749         ESTABLISHED hev        184402373  1793094/sshd: hev@n 
tcp        0      0 10.125.81.41:10050      10.60.145.12:49786      TIME_WAIT   root       0          -                   
udp        0      0 0.0.0.0:60050           0.0.0.0:*                           root       28955      1151/rsyslogd       
udp        0      0 127.0.0.1:53            0.0.0.0:*                           root       26520      1041/dnsmasq        
udp        0      0 0.0.0.0:111             0.0.0.0:*                           root       20471      1/systemd           
udp        0      0 0.0.0.0:54558           0.0.0.0:*                           root       28239      1151/rsyslogd       
udp        0      0 127.0.0.1:323           0.0.0.0:*                           root       26288      1038/chronyd        `
	netstatInput2 = `Active Internet connections (servers and established)
Proto Recv-Q Send-Q Local Address           Foreign Address         State       User       Inode      PID/Program name
tcp        0      0 10.255.255.254:53       0.0.0.0:*               LISTEN      root       20490      -
tcp        0      0 127.0.0.1:37227         0.0.0.0:*               LISTEN      hev        28833      413/node
tcp        0      0 127.0.0.54:53           0.0.0.0:*               LISTEN      systemd-resolve 25756      199/systemd-resolve
tcp        0      0 127.0.0.53:53           0.0.0.0:*               LISTEN      systemd-resolve 25754      199/systemd-resolve
tcp        0      0 0.0.0.0:23              0.0.0.0:*               LISTEN      root       21735      213/inetd
tcp        0      0 127.0.0.1:37227         127.0.0.1:36308         ESTABLISHED hev        6667166    413/node
tcp        0      0 127.0.0.1:36308         127.0.0.1:37227         ESTABLISHED hev        6679367    62414/node
tcp        0      0 127.0.0.1:37227         127.0.0.1:40940         ESTABLISHED hev        6666963    522/node
tcp        0      0 127.0.0.1:40940         127.0.0.1:37227         ESTABLISHED hev        6682094    62080/node
udp        0      0 127.0.0.54:53           0.0.0.0:*                           systemd-resolve 25755      199/systemd-resolve
udp        0      0 127.0.0.53:53           0.0.0.0:*                           systemd-resolve 25753      199/systemd-resolve
udp        0      0 10.255.255.254:53       0.0.0.0:*                           root       20489      -
udp        0      0 127.0.0.1:323           0.0.0.0:*                           root       19468      -`
)

func TestIntegrate(t *testing.T) {
	t.Run("Запускаем parser и готовим данные", func(t *testing.T) {
		go parser()
		chToParser <- []byte(netstatInput)
		time.Sleep(100 * time.Millisecond)
		close(chToParser)
	})
}
