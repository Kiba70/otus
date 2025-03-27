package netstat

import (
	"log/slog"
	"testing"

	"otus/internal/storage"

	"github.com/stretchr/testify/require"
)

//nolint:lll
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
udp        0      0 127.0.0.1:323           0.0.0.0:*                           root       26288      1038/chronyd        
`
)

func TestNetstat(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	dataMon = storage.New[Netstat]()

	t.Run("Готовим данные и прогоняем parser", func(t *testing.T) {
		chToParser = make(chan []byte, 10)

		chToParser <- []byte(netstatInput)
		close(chToParser)

		parser()

		stat, err := GetSum(1)
		require.Nil(t, err)
		require.Equal(t, 2, len(stat.Conn))
		require.Equal(t, 10, len(stat.Socket))
	})
}
