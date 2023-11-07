package smee

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	serviceAccountName = "smee"
)

func CreateServiceAccount(ctx context.Context, client ctrlruntimeclient.Client, ns string) error {
	sa := &corev1.ServiceAccount{
		ObjectMeta: v1.ObjectMeta{
			Name:      serviceAccountName,
			Namespace: ns,
		},
	}

	if err := client.Create(ctx, sa); err != nil {
		if kerrors.IsAlreadyExists(err) {
			return nil
		}

		return err
	}

	return nil
}
