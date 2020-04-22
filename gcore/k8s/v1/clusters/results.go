package clusters

import (
	"fmt"
	"net"
	"time"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/k8s/v1/pools"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/task/v1/tasks"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a cluster resource.
func (r commonResult) Extract() (*ClusterWithPool, error) {
	var s ClusterWithPool
	err := r.ExtractInto(&s)
	return &s, err
}

// Extract is a function that accepts a result and extracts a cluster config.
func (r ConfigResult) Extract() (*Config, error) {
	var c Config
	err := r.ExtractInto(&c)
	return &c, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Cluster.
type GetResult struct {
	commonResult
}

// ConfigResult represents the result of kubernetes config
type ConfigResult struct {
	gcorecloud.Result
}

// Cluster represents a cluster structure.
type Cluster struct {
	StatusReason       string            `json:"status_reason"`
	APIAddress         *gcorecloud.URL   `json:"api_address"`
	CoeVersion         string            `json:"coe_version"`
	ContainerVersion   string            `json:"container_version"`
	DiscoveryURL       *gcorecloud.URL   `json:"discovery_url"`
	HealthStatusReason map[string]string `json:"health_status_reason"`
	ProjectID          string            `json:"project_id"`
	UserID             string            `json:"user_id"`
	NodeAddresses      []net.IP          `json:"node_addresses"`
	MasterAddresses    []net.IP          `json:"master_addresses"`
	FixedNetwork       string            `json:"fixed_network"`
	FixedSubnet        string            `json:"fixed_subnet"`
	FloatingIPEnabled  bool              `json:"floating_ip_enabled"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          *time.Time        `json:"updated_at"`
	Faults             map[string]string `json:"faults"`
	*ClusterList
}

type ClusterWithPool struct {
	*Cluster
	Pools []pools.ClusterPool `json:"pools"`
}

// Config represents a k8s config structure.
type Config struct {
	Config string `json:"config,omitempty"`
}

// Cluster represents a cluster structure in list response.
type ClusterList struct {
	UUID              string            `json:"uuid"`
	Name              string            `json:"name"`
	ClusterTemplateID string            `json:"cluster_template_id"`
	KeyPair           string            `json:"keypair"`
	NodeCount         int               `json:"node_count"`
	MasterCount       int               `json:"master_count"`
	DockerVolumeSize  int               `json:"docker_volume_size"`
	Labels            map[string]string `json:"labels"`
	MasterFlavorID    string            `json:"master_flavor_id"`
	FlavorID          string            `json:"flavor_id"`
	CreateTimeout     int               `json:"create_timeout"`
	StackID           string            `json:"stack_id"`
	Status            string            `json:"status"`
	HealthStatus      string            `json:"health_status"`
	Version           string            `json:"version"`
}

type ClusterListWithPool struct {
	Pools []pools.ClusterListPool `json:"pools"`
	*ClusterList
}

// VersionPage is the page returned by a pager when traversing over a collection of clusters.
type ClusterPage struct {
	pagination.LinkedPageBase
}

// VersionPage is the page returned by a pager when traversing over a collection of cluster versions.
type VersionPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of clusters has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r ClusterPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a ClusterPage struct is empty.
func (r ClusterPage) IsEmpty() (bool, error) {
	is, err := ExtractClusters(r)
	return len(is) == 0, err
}

// IsEmpty checks whether a VersionPage struct is empty.
func (r VersionPage) IsEmpty() (bool, error) {
	is, err := ExtractVersions(r)
	return len(is) == 0, err
}

// ExtractCluster accepts a Page struct, specifically a ClusterPage struct,
// and extracts the elements into a slice of ClusterListWithPool structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractClusters(r pagination.Page) ([]ClusterListWithPool, error) {
	var s []ClusterListWithPool
	err := ExtractClustersInto(r, &s)
	return s, err
}

// ExtractVersions accepts a Page struct, specifically a VersionPage struct,
// and extracts the elements into a slice of strings.
func ExtractVersions(r pagination.Page) ([]string, error) {
	var s []string
	err := ExtractVersionInto(r, &s)
	return s, err
}

func ExtractClustersInto(r pagination.Page, v interface{}) error {
	return r.(ClusterPage).Result.ExtractIntoSlicePtr(v, "results")
}

func ExtractVersionInto(r pagination.Page, v interface{}) error {
	return r.(VersionPage).Result.ExtractIntoSlicePtr(v, "results")
}

type ClusterTaskResult struct {
	Clusters []string `json:"clusters"`
}

func ExtractClusterIDFromTask(task *tasks.Task) (string, error) {
	var result ClusterTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode cluster information in task structure: %w", err)
	}
	if len(result.Clusters) == 0 {
		return "", fmt.Errorf("cannot decode cluster information in task structure: %w", err)
	}
	return result.Clusters[0], nil
}