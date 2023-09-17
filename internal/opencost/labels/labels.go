package labels

import (
	"github.com/opencost/opencost/pkg/kubecost"
	commonskuberneteslabels "github.com/plantoncloud-inc/go-commons/kubernetes/labels"
	"github.com/plantoncloud-inc/go-commons/pulumi/resourcemetadata/keys"
	log "github.com/sirupsen/logrus"
)

type WorkloadLabels struct {
	Company      string
	Product      string
	Environment  string
	ResourceType string
	ResourceId   string
}

func GetLabels(allocation *kubecost.AllocationJSON) *WorkloadLabels {
	labels := &WorkloadLabels{}
	companyPrometheusLabel := keys.WithPrometheusFormat(commonskuberneteslabels.Company)
	log.Debugf("checking if %s has company label %s", allocation.Name, companyPrometheusLabel)
	if _, ok := allocation.Properties.Labels[companyPrometheusLabel]; ok {
		log.Debugf("####: %s has %s for company label %s", allocation.Name, allocation.Properties.Labels[companyPrometheusLabel], companyPrometheusLabel)
		labels.Company = allocation.Properties.Labels[companyPrometheusLabel]
	}
	productPrometheusLabel := keys.WithPrometheusFormat(commonskuberneteslabels.Product)
	log.Debugf("checking if %s has product label %s", allocation.Name, productPrometheusLabel)
	if _, ok := allocation.Properties.Labels[productPrometheusLabel]; ok {
		log.Debugf("####: %s has %s for product label %s", allocation.Name, allocation.Properties.Labels[productPrometheusLabel], productPrometheusLabel)
		labels.Product = allocation.Properties.Labels[productPrometheusLabel]
	}
	environmentPrometheusLabel := keys.WithPrometheusFormat(commonskuberneteslabels.Environment)
	log.Debugf("checking if %s has product-env label %s", allocation.Name, environmentPrometheusLabel)
	if _, ok := allocation.Properties.Labels[environmentPrometheusLabel]; ok {
		log.Debugf("####: %s has %s for product-env label %s", allocation.Name, allocation.Properties.Labels[environmentPrometheusLabel], environmentPrometheusLabel)
		labels.Environment = allocation.Properties.Labels[environmentPrometheusLabel]
	}

	resourceTypePrometheusLabel := keys.WithPrometheusFormat(commonskuberneteslabels.ResourceType)
	log.Debugf("checking if %s has resource-type label %s", allocation.Name, resourceTypePrometheusLabel)
	if _, ok := allocation.Properties.Labels[resourceTypePrometheusLabel]; ok {
		log.Debugf("####: %s has %s for resource-type label %s", allocation.Name, allocation.Properties.Labels[resourceTypePrometheusLabel], resourceTypePrometheusLabel)
		labels.ResourceType = allocation.Properties.Labels[resourceTypePrometheusLabel]
	}

	resourceIdPrometheusLabel := keys.WithPrometheusFormat(commonskuberneteslabels.ResourceId)
	log.Debugf("checking if %s has resource-id label %s", allocation.Name, resourceIdPrometheusLabel)
	if _, ok := allocation.Properties.Labels[resourceIdPrometheusLabel]; ok {
		log.Debugf("####: %s has %s for resource-id label %s", allocation.Name, allocation.Properties.Labels[resourceIdPrometheusLabel], resourceIdPrometheusLabel)
		labels.ResourceId = allocation.Properties.Labels[resourceIdPrometheusLabel]
	}

	return labels
}
