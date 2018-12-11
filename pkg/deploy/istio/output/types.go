package output

import (
	"github.com/rancher/norman/types"
	"github.com/rancher/rio/pkg/apply"
	"github.com/rancher/rio/pkg/settings"
	"github.com/rancher/rio/types/apis/networking.istio.io/v1alpha3"
	riov1 "github.com/rancher/rio/types/apis/rio.cattle.io/v1"
	"github.com/rancher/rio/types/client/rio/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Deployment struct {
	Enabled         bool
	UseLoadBalancer bool
	Ports           []string
	Stacks          map[string]*riov1.Stack
	Services        map[string]*v1.Service
	Gateways        map[string]*Gateway
	VirtualService  map[string]*VirtualService
}

func NewDeployment() *Deployment {
	return &Deployment{
		Stacks:   map[string]*riov1.Stack{},
		Services: map[string]*v1.Service{},
		Gateways: map[string]*Gateway{},
	}
}

func (d *Deployment) Deploy(ns, groupID string) error {
	ad := apply.Data{
		GroupID: groupID,
	}

	ad.Add(ns, v1.GroupName, client.StackType, d.Stacks)
	ad.Add(settings.IstioExternalLBNamespace, v1.GroupName, "Service", d.Services)
	ad.Add(ns, v1alpha3.GroupName, "Gateway", d.Gateways)

	return ad.Apply()
}

type Gateway struct {
	types.Namespaced

	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec interface{} `json:"spec,omitempty"`
}

func (g *Gateway) DeepCopyObject() runtime.Object {
	panic("not implemented")
}

type VirtualService struct {
	types.Namespaced

	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec interface{} `json:"spec,omitempty"`
}

func (v *VirtualService) DeepCopyObject() runtime.Object {
	panic("not implemented")
}

type Pod struct {
	types.Namespaced

	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec interface{} `json:"spec,omitempty"`
}

func (p *Pod) DeepCopyObject() runtime.Object {
	panic("not implemented")
}

type Service struct {
	types.Namespaced

	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec interface{} `json:"spec,omitempty"`
}

func (s *Service) DeepCopyObject() runtime.Object {
	panic("not implemented")
}

type DestinationRule struct {
	types.Namespaced

	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec interface{} `json:"spec,omitempty"`
}

func (d *DestinationRule) DeepCopyObject() runtime.Object {
	panic("not implemented")
}
