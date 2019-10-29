package ovn

import (
	. "github.com/onsi/gomega"
	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/config"
	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/factory"
	ovntest "github.com/ovn-org/ovn-kubernetes/go-controller/pkg/testing"
	util "github.com/ovn-org/ovn-kubernetes/go-controller/pkg/util"
	"github.com/urfave/cli"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

const (
	k8sTCPLoadBalancerIP = "k8s_tcp_load_balancer"
	k8sUDPLoadBalancerIP = "k8s_udp_load_balancer"
)

type FakeOVN struct {
	fakeClient *fake.Clientset
	watcher    *factory.WatchFactory
	controller *Controller
}

func (o *FakeOVN) start(ctx *cli.Context, fexec *ovntest.FakeExec, objects ...runtime.Object) {
	err := util.SetExec(fexec)
	Expect(err).NotTo(HaveOccurred())

	_, err = config.InitConfig(ctx, fexec, nil)
	Expect(err).NotTo(HaveOccurred())

	o.fakeClient = fake.NewSimpleClientset(objects...)
	o.init()
}

func (o *FakeOVN) restart() {
	o.watcher.Shutdown()
	o.init()
}

func (o *FakeOVN) init() {
	var err error

	o.watcher, err = factory.NewWatchFactory(o.fakeClient, make(chan struct{}))
	Expect(err).NotTo(HaveOccurred())

	o.controller = NewOvnController(o.fakeClient, o.watcher, nil)
	o.controller.portGroupSupport = true
}
