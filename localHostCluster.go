package cluster

// xlCluster_go/localHostCluster.go

import (
	xi "github.com/jddixon/xlNodeID_go"
)

type LocalHostCluster struct {
	Cluster
}

func NewLocalHostCluster(name string, id *xi.NodeID, attrs uint64,
	maxSize, epCount uint32) (tc *Cluster, err error) {

	return NewCluster(name, id, attrs, maxSize, epCount)
}

/**
 * Call Run() on all constituent Nodes, which will activate acceptors.
 */
func (lhc *LocalHostCluster) Run() (err error) {
	for i := 0; err == nil && i < len(lhc.ClMembers); i++ {
		err = lhc.ClMembers[i].Run()
	}
	return
}

/**
 * Call Close() on all constituent Nodes, which will close acceptors.
 */
func (lhc *LocalHostCluster) Close() (err error) {

	for i := 0; err == nil && i < len(lhc.ClMembers); i++ {
		err = lhc.ClMembers[i].Close()
	}
	return
}
