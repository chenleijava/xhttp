package registry

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/netx"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/rest"
	"os"
	"strings"
)

const (
	allEths  = "0.0.0.0"
	envPodIP = "POD_IP"
)

// RegisterRest
//
//	@Description:
//	@param etcd
//	@param svrConf
//	@return error
func RegisterRest(etcd discov.EtcdConf, svrConf rest.RestConf) error {
	err := etcd.Validate()
	logx.Must(err)

	listenOn := fmt.Sprintf("%s:%d", svrConf.Host, svrConf.Port)
	pubListenOn := figureOutListenOn(listenOn)
	var pubOpts []discov.PubOption
	if etcd.HasAccount() {
		pubOpts = append(pubOpts, discov.WithPubEtcdAccount(etcd.User, etcd.Pass))
	}
	if etcd.HasTLS() {
		pubOpts = append(pubOpts, discov.WithPubEtcdTLS(etcd.CertFile, etcd.CertKeyFile,
			etcd.CACertFile, etcd.InsecureSkipVerify))
	}

	key := fmt.Sprintf("/services/%s", svrConf.Name)
	pubClient := discov.NewPublisher(etcd.Hosts, key, fmt.Sprintf("http://%s", pubListenOn), pubOpts...)
	proc.AddShutdownListener(func() {
		pubClient.Stop()
	})

	return pubClient.KeepAlive()
}

func figureOutListenOn(listenOn string) string {
	fields := strings.Split(listenOn, ":")
	if len(fields) == 0 {
		return listenOn
	}

	host := fields[0]
	if len(host) > 0 && host != allEths {
		return listenOn
	}

	ip := os.Getenv(envPodIP)
	if len(ip) == 0 {
		ip = netx.InternalIp()
	}
	if len(ip) == 0 {
		return listenOn
	}

	return strings.Join(append([]string{ip}, fields[1:]...), ":")
}
