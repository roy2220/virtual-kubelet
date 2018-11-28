// +build e2e

package e2e

import (
	"fmt"
	"testing"

	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"
)

// TestGetStatsSummary creates a pod having two containers and queries the /stats/summary endpoint of the virtual-kubelet.
// It expects this endpoint to return stats for the current node, as well as for the aforementioned pod and each of its two containers.
func TestGetStatsSummary(t *testing.T) {
	// Create a pod with prefix "nginx-0-" having three containers.
	pod, err := f.CreatePod(f.CreateDummyPodObjectWithPrefix("nginx-0-", "foo", "bar", "baz"))
	if err != nil {
		t.Fatal(err)
	}
	// Delete the "nginx-0-X" pod after the test finishes.
	defer func() {
		if err := f.DeletePod(pod.Namespace, pod.Name); err != nil && !apierrors.IsNotFound(err) {
			t.Error(err)
		}
	}()

	// Wait for the "nginx-0-X" pod to be reported as running and ready.
	if err := f.WaitUntilPodReady(pod.Namespace, pod.Name); err != nil {
		t.Fatal(err)
	}

	// Grab the stats from the provider.
	stats, err := f.GetStatsSummary()
	if err != nil {
		t.Fatal(err)
	}

	// Make sure that we've got stats for the current node.
	if stats.Node.NodeName != f.NodeName {
		t.Fatalf("expected stats for node %s, got stats for node %s", f.NodeName, stats.Node.NodeName)
	}

	// Make sure the "nginx-0-X" pod exists in the slice of PodStats.
	idx, err := findPodInPodStats(stats, pod)
	if err != nil {
		t.Fatal(err)
	}

	// Make sure that we've got stats for all the containers in the "nginx-0-X" pod.
	desiredContainerStatsCount := len(pod.Spec.Containers)
	currentContainerStatsCount := len(stats.Pods[idx].Containers)
	if currentContainerStatsCount != desiredContainerStatsCount {
		t.Fatalf("expected stats for %d containers, got stats for %d containers", desiredContainerStatsCount, currentContainerStatsCount)
	}
}

// TestPodLifecycle creates two pods and verifies that the provider has been asked to create them.
// Then, it deletes one of the pods and verifies that the provider has been asked to delete it.
// These verifications are made using the /stats/summary endpoint of the virtual-kubelet, by checking for the presence or absence of the pods.
// Hence, the provider being tested must implement the PodMetricsProvider interface.
func TestPodLifecycle(t *testing.T) {
	// Create a pod with prefix "nginx-0-" having a single container.
	pod0, err := f.CreatePod(f.CreateDummyPodObjectWithPrefix("nginx-0-", "foo"))
	if err != nil {
		t.Fatal(err)
	}
	// Delete the "nginx-0-X" pod after the test finishes.
	defer func() {
		if err := f.DeletePod(pod0.Namespace, pod0.Name); err != nil && !apierrors.IsNotFound(err) {
			t.Error(err)
		}
	}()

	// Create a pod with prefix "nginx-1-" having a single container.
	pod1, err := f.CreatePod(f.CreateDummyPodObjectWithPrefix("nginx-1-", "bar"))
	if err != nil {
		t.Fatal(err)
	}
	// Delete the "nginx-1-Y" pod after the test finishes.
	defer func() {
		if err := f.DeletePod(pod1.Namespace, pod1.Name); err != nil && !apierrors.IsNotFound(err) {
			t.Error(err)
		}
	}()

	// Wait for the "nginx-0-X" pod to be reported as running and ready.
	if err := f.WaitUntilPodReady(pod0.Namespace, pod0.Name); err != nil {
		t.Fatal(err)
	}
	// Wait for the "nginx-1-Y" pod to be reported as running and ready.
	if err := f.WaitUntilPodReady(pod1.Namespace, pod1.Name); err != nil {
		t.Fatal(err)
	}

	// Grab the stats from the provider.
	stats, err := f.GetStatsSummary()
	if err != nil {
		t.Fatal(err)
	}

	// Make sure the "nginx-0-X" pod exists in the slice of PodStats.
	if _, err := findPodInPodStats(stats, pod0); err != nil {
		t.Fatal(err)
	}

	// Make sure the "nginx-1-Y" pod exists in the slice of PodStats.
	if _, err := findPodInPodStats(stats, pod1); err != nil {
		t.Fatal(err)
	}

	// Delete the "nginx-1" pod.
	if err := f.DeletePod(pod1.Namespace, pod1.Name); err != nil {
		t.Fatal(err)
	}

	// Wait for the "nginx-1-Y" pod to be reported as having been marked for deletion.
	if err := f.WaitUntilPodDeleted(pod1.Namespace, pod1.Name); err != nil {
		t.Fatal(err)
	}

	// Grab the stats from the provider.
	stats, err = f.GetStatsSummary()
	if err != nil {
		t.Fatal(err)
	}

	// Make sure the "nginx-1-Y" pod DOES NOT exist in the slice of PodStats anymore.
	if _, err := findPodInPodStats(stats, pod1); err == nil {
		t.Fatalf("expected to NOT find pod \"%s/%s\" in the slice of pod stats", pod1.Namespace, pod1.Name)
	}
}

// findPodInPodStats returns the index of the specified pod in the .pods field of the specified Summary object.
// It returns an error if the specified pod is not found.
func findPodInPodStats(summary *v1alpha1.Summary, pod *v1.Pod) (int, error) {
	for i, p := range summary.Pods {
		if p.PodRef.Namespace == pod.Namespace && p.PodRef.Name == pod.Name && string(p.PodRef.UID) == string(pod.UID) {
			return i, nil
		}
	}
	return -1, fmt.Errorf("failed to find pod \"%s/%s\" in the slice of pod stats", pod.Namespace, pod.Name)
}