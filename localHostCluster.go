package cluster

// xlCluster_go/localHostCluster.go

import (
	xi "github.com/jddixon/xlNodeID_go"
)

type LocalHostCluster struct {
	Cluster
}

func NewLocalHostCluster(name string, id *xi.NodeID, attrs uint64,
	maxSize, epCount uint32) (tc *LocalHostCluster, err error) {

	cl, err := NewCluster(name, id, attrs, maxSize, epCount)
	if err == nil {
		tc = &LocalHostCluster{
			Cluster: *cl,
		}
	}
	return
}

/**
 * Call OpenAcc() on all constituent Nodes, which will activate acceptors.
 */
func (lhc *LocalHostCluster) Start() (err error) {
	for i := 0; err == nil && i < len(lhc.ClMembers); i++ {
		err = lhc.ClMembers[i].OpenAcc()
	}
	return
}

/**
 * Call Close() on all constituent Nodes, which will close acceptors.
 */
func (lhc *LocalHostCluster) Stop() (err error) {

	for i := 0; err == nil && i < len(lhc.ClMembers); i++ {
		err = lhc.ClMembers[i].CloseAcc()
	}
	return
}
