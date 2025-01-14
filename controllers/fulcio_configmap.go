package controllers

import (
	"context"

	client "sigs.k8s.io/controller-runtime/pkg/client"

	rhtasv1alpha1 "github.com/securesign/operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *SecuresignReconciler) ensureOIDCConfigMap(ctx context.Context, m *rhtasv1alpha1.Securesign, namespace string, configMapName string, component string) (*corev1.ConfigMap,
	error) {
	log := ctrllog.FromContext(ctx)
	log.Info("ensuring configmap")
	// Define a new ConfigMap object
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":     component,
				"app.kubernetes.io/instance": "trusted-artifact-signer",
			},
		},
		Data: map[string]string{
			"config.json": `{
				  "OIDCIssuers": {
					"http://keycloak-internal.keycloak-system.svc/auth/realms/sigstore": {
					"ClientID": "sigstore",
					"IssuerURL": "http://keycloak-internal.keycloak-system.svc/auth/realms/sigstore",
					"Type": "email"
					}
				}
			}`,
		},
	}

	// Check if this ConfigMap already exists else create it in the namespace
	err := r.Get(ctx, client.ObjectKey{Name: configMap.Name, Namespace: namespace}, configMap)
	// If the ConfigMap doesn't exist, create it but if it does, do nothing
	if err != nil {
		log.Info("Creating a new ConfigMap")
		err = r.Create(ctx, configMap)
		if err != nil {
			log.Error(err, "Failed to create new ConfigMap")
			return nil, err
		}
	}
	return configMap, nil
}
