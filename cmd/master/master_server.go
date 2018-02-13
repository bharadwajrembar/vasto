package master

import (
	"net"
	"sync"

	"github.com/chrislusf/vasto/pb"
	"google.golang.org/grpc"
	"github.com/chrislusf/glog"
)

type MasterOption struct {
	Address *string
}

type masterServer struct {
	option               *MasterOption
	clientChans          *clientChannels
	clientsStat          *clientsStat
	topo                 *masterTopology
	keyspaceMutexMap     map[string]*sync.Mutex
	keyspaceMutexMapLock sync.Mutex
}

func RunMaster(option *MasterOption) {
	var ms = &masterServer{
		option:           option,
		clientChans:      newClientChannels(),
		clientsStat:      newClientsStat(),
		topo:             newMasterTopology(),
		keyspaceMutexMap: make(map[string]*sync.Mutex),
	}

	listener, err := net.Listen("tcp", *option.Address)
	if err != nil {
		glog.Fatal(err)
	}
	glog.V(0).Infof("Vasto master starts on %s\n", *option.Address)

	// m := cmux.New(listener)
	// grpcListener := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))

	ms.serveGrpc(listener)

	//if err := m.Serve(); !strings.Contains(err.Error(), "use of closed network connection") {
	//	panic(err)
	//}

}

func (ms *masterServer) serveGrpc(listener net.Listener) {
	grpcServer := grpc.NewServer()
	pb.RegisterVastoMasterServer(grpcServer, ms)
	grpcServer.Serve(listener)
}

func (ms *masterServer) lock(keyspace string) {
	ms.keyspaceMutexMapLock.Lock()
	mu, found := ms.keyspaceMutexMap[keyspace]
	if !found {
		mu = &sync.Mutex{}
		ms.keyspaceMutexMap[keyspace] = mu
	}
	ms.keyspaceMutexMapLock.Unlock()

	mu.Lock()
}

func (ms *masterServer) unlock(keyspace string) {
	ms.keyspaceMutexMapLock.Lock()
	mu, found := ms.keyspaceMutexMap[keyspace]
	ms.keyspaceMutexMapLock.Unlock()
	if !found {
		return
	}

	mu.Unlock()
}
