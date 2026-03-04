/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"time"

	helmv2 "github.com/fluxcd/helm-controller/api/v2"
	"github.com/fluxcd/pkg/apis/meta"
	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	libutils "github.com/openmcp-project/openmcp-operator/lib/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openmcp-project/controller-utils/pkg/clusters"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	apiv1alpha1 "github.com/openmcp-project/service-provider-aso/api/v1alpha1"
	spruntime "github.com/openmcp-project/service-provider-aso/pkg/runtime"
)

const (
	// HelmReleaseName is the name of the Flux' HelmRelease object created by the controller.
	HelmReleaseName = "azure-service-operator"

	// HelmRepoistoryName is the name of the Flux' HelmRepository object created by the controller.
	HelmRepositoryName = "azure-service-operator"

	// ASOSystemNamespace is the default namespace on the target MCP cluster to install the Azure Servie Operator.
	ASOSystemNamespace = "azureserviceoperator-system"

	// HelmChartName is the name of the Helm chart
	HelmChartName = "azure-service-operator"
)

// AzureServiceOperatorReconciler reconciles a AzureServiceOperator object
type AzureServiceOperatorReconciler struct {
	// OnboardingCluster is the cluster where this controller watches AzureServiceOperator resources and reacts to their changes.
	OnboardingCluster *clusters.Cluster
	// PlatformCluster is the cluster where this controller is deployed and configured.
	PlatformCluster *clusters.Cluster
	// PodNamespace is the namespace where this controller is deployed in.
	PodNamespace string
}

// CreateOrUpdate is called on every add or update event
func (r *AzureServiceOperatorReconciler) CreateOrUpdate(ctx context.Context, svcobj *apiv1alpha1.AzureServiceOperator, providerConfig *apiv1alpha1.ProviderConfig, clusters spruntime.ClusterContext) (ctrl.Result, error) {
	l := logf.FromContext(ctx)
	spruntime.StatusProgressing(svcobj, "Reconciling", "Reconcile in progress")

	tenantNamespace, err := libutils.StableMCPNamespace(svcobj.Name, svcobj.Namespace)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to determine stable namespace for OCM instance: %w", err)
	}

	l.Info("checking tenantNamespace", "namespace", tenantNamespace)

	// 1. Create Flux OCIRepository resource
	if err := r.createOrUpdateHelmRepository(ctx, svcobj, providerConfig, tenantNamespace); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create HelmRepository resource on Platform cluster", err)
	}

	// 2. Create Flux HelmRelease resource
	if err = r.createOrUpdateHelmRelease(ctx, svcobj, providerConfig, clusters.MCPAccessSecretKey, tenantNamespace); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create HelmRelease resource on Platform cluster", err)
	}

	l.Info("Done reconciling AzureServiceOperator resource", "name", svcobj.Name)

	spruntime.StatusReady(svcobj)
	return ctrl.Result{}, nil
}

// Delete is called on every delete event
func (r *AzureServiceOperatorReconciler) Delete(ctx context.Context, obj *apiv1alpha1.AzureServiceOperator, _ *apiv1alpha1.ProviderConfig, clusters spruntime.ClusterContext) (ctrl.Result, error) {
	l := logf.FromContext(ctx)
	spruntime.StatusTerminating(obj)

	// 1. Delete HelmRelease object from Platform cluster

	// 2. Delete HelmRepository object from Platform cluster

	// object still exists
	return ctrl.Result{
		RequeueAfter: time.Second * 10,
	}, nil
}

func (r *AzureServiceOperatorReconciler) createOrUpdateHelmRepository(ctx context.Context, svcobj *apiv1alpha1.AzureServiceOperator, providerConfig *apiv1alpha1.ProviderConfig, namespace string) error {
	helmRepository := createHelmRepository(providerConfig, svcobj.Spec.Version, namespace)
	managedObj := &sourcev1.HelmRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      helmRepository.Name,
			Namespace: helmRepository.Namespace,
		},
	}
	l := logf.FromContext(ctx)
	l.Info("creating HelmRepository", "object", helmRepository)
	if _, err := ctrl.CreateOrUpdate(ctx, r.PlatformCluster.Client(), managedObj, func() error {
		managedObj.Spec = helmRepository.Spec
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func createHelmRepository(providerConfig *apiv1alpha1.ProviderConfig, version, namespace string) *sourcev1.HelmRepository {
	return &sourcev1.HelmRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      HelmRepositoryName,
			Namespace: namespace,
		},
		Spec: sourcev1.HelmRepositorySpec{
			Interval: metav1.Duration{Duration: time.Minute},
			URL:      providerConfig.Spec.HelmChartLocation,
		},
	}
}

func (r *AzureServiceOperatorReconciler) createOrUpdateHelmRelease(ctx context.Context, svcobj *apiv1alpha1.AzureServiceOperator, providerConfig *apiv1alpha1.ProviderConfig, mcpAccessSecret client.ObjectKey, namespace string) error {
	helmRelease := createHelmRelease(providerConfig, mcpAccessSecret, svcobj.Spec.Version, namespace)
	managedObj := &helmv2.HelmRelease{
		ObjectMeta: metav1.ObjectMeta{
			Name:      helmRelease.Name,
			Namespace: helmRelease.Namespace,
		},
	}
	l := logf.FromContext(ctx)
	l.Info("creating HelmRelease", "object", helmRelease)
	if _, err := ctrl.CreateOrUpdate(ctx, r.PlatformCluster.Client(), managedObj, func() error {
		managedObj.Spec = helmRelease.Spec
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func createHelmRelease(providerConfig *apiv1alpha1.ProviderConfig, mcpAccessSecret client.ObjectKey, version, namespace string) *helmv2.HelmRelease {
	return &helmv2.HelmRelease{
		ObjectMeta: metav1.ObjectMeta{
			Name:      HelmReleaseName,
			Namespace: namespace,
		},
		Spec: helmv2.HelmReleaseSpec{
			Interval: metav1.Duration{Duration: time.Minute},
			Chart: &helmv2.HelmChartTemplate{
				Spec: helmv2.HelmChartTemplateSpec{
					Chart:   HelmChartName,
					Version: version,
					SourceRef: helmv2.CrossNamespaceObjectReference{
						Kind: "HelmRepository",
						Name: HelmRepositoryName, // HelmRepository.ObjectMeta.Name from referenced Object
					},
				},
			},
			ReleaseName:      HelmReleaseName,
			TargetNamespace:  ASOSystemNamespace,
			StorageNamespace: ASOSystemNamespace,
			KubeConfig: &meta.KubeConfigReference{
				SecretRef: &meta.SecretKeyReference{
					Name: mcpAccessSecret.Name,
					Key:  "kubeconfig",
				},
			},
		},
	}
}
