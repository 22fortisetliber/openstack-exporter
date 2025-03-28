package exporters

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

type ClusteringTestSuite struct {
	BaseOpenStackTestSuite
}

var clusteringExpectedUp = `
# HELP openstack_clustering_cluster_status cluster_status
# TYPE openstack_clustering_cluster_status gauge
openstack_clustering_cluster_status{id="7d85f602-a948-4a30-afd4-e84f47471c15",name="cluster1",profile_id="edc63d0a-2ca4-48fa-9854-27926da76a4a",status="ACTIVE",nodes_count="2",desired_capacity="3",min_size="3",max_size="-1",project_id=a"6e18cc2bdbeb48a5b3cad2dc499f6804"} 1
# HELP openstack_clustering_total_clusters total_clusters
# TYPE openstack_clustering_total_clusters gauge
openstack_clustering_total_clusters 1
# HELP openstack_clustering_node_status cluster_nodes_status
# TYPE openstack_clustering_node_status gauge
openstack_clustering_node_status{id="82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1",name="node-e395be1e-002",profile_id="d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",status="ACTIVE",cluster_id="e395be1e-8d8e-43bb-bd6c-943eccf76a6d",physical_id="66a81d68-bf48-4af5-897b-a3bfef7279a8",project_id="eee0b7c083e84501bdd50fb269d2a10e"} 1
# HELP openstack_clustering_total_nodes total_nodes
# TYPE openstack_clustering_total_nodes gauge
openstack_clustering_total_nodes 1
# HELP openstack_clustering_up up
# TYPE openstack_clustering_up gauge
openstack_clustering_up 1
`

func (suite *ClusteringTestSuite) TestClusteringExporter() {
	err := testutil.CollectAndCompare(*suite.Exporter, strings.NewReader(clusteringExpectedUp))
	assert.NoError(suite.T(), err)
}