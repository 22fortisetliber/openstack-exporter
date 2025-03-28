package exporters

import (
	"strconv"

	"github.com/go-kit/log"
	"github.com/gophercloud/gophercloud/openstack/clustering/v1/clusters"
	"github.com/gophercloud/gophercloud/openstack/clustering/v1/nodes"
	"github.com/prometheus/client_golang/prometheus"
)

var clustering_status = []string{
	"INIT",
	"ACTIVE",
	"ERROR",
	"CRITICAL",
	"WARNING",
	"CREATING",
	"UPDATING",
	"DELETING",
	"RESIZING",
	"CHECKING",
	"RECOVERING",
}

var clustering_nodes_status = []string{
	"INIT",
	"ACTIVE",
	"ERROR",
	"WARNING",
	"CREATING",
	"UPDATING",
	"DELETING",
	"RECOVERING",
}

func mapClusteringStatus(current string) int {
	for idx, status := range clustering_status {
		if current == status {
			return idx
		}
	}
	return -1
}

func mapNodesStatus(current string) int {
	for idx, status := range clustering_nodes_status {
		if current == status {
			return idx
		}
	}
	return -1
}

type ClusteringExporter struct {
	BaseOpenStackExporter
}

var defaultClusteringMetrics = []Metric{
	{Name: "total_clusters", Fn: ListAllClustering},
	{Name: "cluster_status", Labels: []string{"id", "name", "profile_id", "status", "nodes_count", "desired_capacity", "min_size", "max_size", "project_id"}, Fn: nil},
	{Name: "total_nodes", Fn: ListAllClusteringNodes},
	{Name: "cluster_nodes_status", Labels: []string{"id", "name", "profile_id", "status", "cluster_id", "physical_id", "project_id"}, Fn: nil},
}

func NewClusteringExporter(config *ExporterConfig, logger log.Logger) (*ClusteringExporter, error) {
	exporter := ClusteringExporter{
		BaseOpenStackExporter{
			Name:           "clustering",
			ExporterConfig: *config,
			logger:         logger,
		},
	}
	for _, metric := range defaultClusteringMetrics {
		if exporter.isDeprecatedMetric(&metric) {
			continue
		}
		if !exporter.isSlowMetric(&metric) {
			exporter.AddMetric(metric.Name, metric.Fn, metric.Labels, metric.DeprecatedVersion, nil)
		}
	}
	return &exporter, nil
}

func ListAllClustering(exporter *BaseOpenStackExporter, ch chan<- prometheus.Metric) error {
	var allClusters []clusters.Cluster

	allPagesClusters, err := clusters.List(exporter.Client, clusters.ListOpts{}).AllPages()
	if err != nil {
		return err
	}
	allClusters, err = clusters.ExtractClusters(allPagesClusters)
	if err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(exporter.Metrics["total_clusters"].Metric,
		prometheus.GaugeValue, float64(len(allClusters)))
	//Clustering senlin status metrics
	for _, cluster := range allClusters {
		ch <- prometheus.MustNewConstMetric(exporter.Metrics["cluster_status"].Metric,
			prometheus.GaugeValue, float64(1), cluster.ID, cluster.Name,
			cluster.ProfileID, cluster.Status, strconv.Itoa(len(cluster.Nodes)), strconv.Itoa(cluster.DesiredCapacity), strconv.Itoa(cluster.MinSize), strconv.Itoa(cluster.MaxSize), cluster.Project)
	}
	return nil
}

func ListAllClusteringNodes(exporter *BaseOpenStackExporter, ch chan<- prometheus.Metric) error {
	var allNodes []nodes.Node
	allPagesNodes, err := nodes.List(exporter.Client, nodes.ListOpts{}).AllPages()
	if err != nil {
		return err
	}
	allNodes, err = nodes.ExtractNodes(allPagesNodes)
	if err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(exporter.Metrics["total_nodes"].Metric,
		prometheus.GaugeValue, float64(len(allNodes)))
	//Node senlin status metrics
	for _, node := range allNodes {
		ch <- prometheus.MustNewConstMetric(exporter.Metrics["cluster_nodes_status"].Metric,
			prometheus.GaugeValue, float64(1), node.ID, node.Name,
			node.ProfileID, node.Status, node.ClusterID, node.PhysicalID, node.Project)
	}
	return nil
}
