package hegel

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ptr "k8s.io/utils/pointer"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateDeployment(ctx context.Context, client ctrlruntimeclient.Client, ns string) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "hegel",
			Namespace: ns,
			Labels: map[string]string{
				"app": "hegel",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: ptr.Int32(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":   "hegel",
					"stack": "tinkerbell",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":   "hegel",
						"stack": "tinkerbell",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "hegel",
							Image:           "quay.io/tinkerbell/hegel:v0.8.0",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Args:            []string{"--data-model", "kubernetes", "--kube-namespace", ns, "--http-port", "50061"},
							Env: []corev1.EnvVar{
								{
									Name: "HEGEL_TRUSTED_PROXIES",
									// TODO: pass the TRUSTED_PROXIES as a command line
									Value: "10.244.0.0/24,10.244.1.0/24,10.244.2.0/24",
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse("64Mi"),
									corev1.ResourceCPU:    resource.MustParse("10m"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse("128Mi"),
									corev1.ResourceCPU:    resource.MustParse("500m"),
								},
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: int32(50061),
									Name:          "hegel-http",
								},
							},
						},
					},
					ServiceAccountName: serviceAccountName,
				},
			},
		},
	}

	if err := client.Create(ctx, deployment); err != nil {
		if kerrors.IsAlreadyExists(err) {
			return nil
		}

		return err
	}

	return nil
}
