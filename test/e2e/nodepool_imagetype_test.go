//go:build e2e

package e2e

import (
	"testing"

	. "github.com/onsi/gomega"
	hyperv1 "github.com/openshift/hypershift/api/hypershift/v1beta1"
	e2eutil "github.com/openshift/hypershift/test/e2e/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NodePoolImageTypeTest struct {
	DummyInfraSetup
	hostedCluster *hyperv1.HostedCluster
}

func NewNodePoolImageTypeTest(hostedCluster *hyperv1.HostedCluster) *NodePoolImageTypeTest {
	return &NodePoolImageTypeTest{
		hostedCluster: hostedCluster,
	}
}

func (it *NodePoolImageTypeTest) Setup(t *testing.T) {
	// Skip test for non-AWS platforms since ImageType is currently AWS-specific
	if globalOpts.Platform != hyperv1.AWSPlatform {
		t.Skip("test is only supported for AWS platform")
	}
	if e2eutil.IsLessThan(e2eutil.Version419) {
		t.Skip("test only supported from version 4.19")
	}
	t.Log("Starting test NodePoolImageTypeTest")
}

func (it *NodePoolImageTypeTest) BuildNodePoolManifest(defaultNodepool hyperv1.NodePool) (*hyperv1.NodePool, error) {
	nodePool := &hyperv1.NodePool{
		ObjectMeta: metav1.ObjectMeta{
			Name:      it.hostedCluster.Name + "-test-imagetype",
			Namespace: it.hostedCluster.Namespace,
		},
	}
	defaultNodepool.Spec.DeepCopyInto(&nodePool.Spec)

	// Set 1 replica and Windows ImageType
	nodePool.Spec.Replicas = &oneReplicas
	nodePool.Spec.Platform.AWS.InstanceType = "m5.metal"
	nodePool.Spec.Platform.AWS.ImageType = hyperv1.ImageTypeWindows

	return nodePool, nil
}

func (it *NodePoolImageTypeTest) Run(t *testing.T, nodePool hyperv1.NodePool, nodes []corev1.Node) {
	g := NewWithT(t)

	t.Logf("NodePool created with ImageType=%s, Replicas=%d",
		nodePool.Spec.Platform.AWS.ImageType, *nodePool.Spec.Replicas)

	// Verify that ImageType is correctly set in the NodePool
	g.Expect(nodePool.Spec.Platform.AWS.ImageType).To(Equal(hyperv1.ImageTypeWindows),
		"NodePool should have ImageType=Windows")

	// Verify that nodes are running Windows OS
	// The framework already waited for nodes to be Ready
	g.Expect(nodes).NotTo(BeEmpty(), "Expected at least one node to be provisioned")

	for _, node := range nodes {
		g.Expect(node.Status.NodeInfo.OperatingSystem).To(Equal("windows"),
			"Node %s should be running Windows OS", node.Name)
		t.Logf("✓ Node %s is running Windows OS", node.Name)
	}

	t.Log("NodePool ImageType test passed successfully")
}
