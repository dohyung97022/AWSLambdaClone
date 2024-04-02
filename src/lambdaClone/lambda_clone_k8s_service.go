package lambdaClone

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	network "k8s.io/api/networking/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"src/k8s"
	"time"
)

func getName(lambda *Lambda) string {
	return "lambda-" + lambda.Id
}

func getNameById(id *primitive.ObjectID) string {
	return "lambda-" + id.Hex()
}

func getRuntime(lambda *Lambda) (*Runtime, error) {
	runtimes, err := getRuntimes()
	if err != nil {
		return nil, err
	}
	for _, runtime := range *runtimes {
		if lambda.Runtime == runtime.Runtime && lambda.Version == runtime.Version {
			return &runtime, nil
		}
	}
	return nil, errors.New("could not find matching image in lambda runtime collection")
}

func getLabels(lambda *Lambda) map[string]string {
	return map[string]string{"app": "lambda-clone-created", "key": lambda.Id}
}

func createDeployment(lambda *Lambda) error {
	runtime, err := getRuntime(lambda)
	if err != nil {
		return err
	}
	var replicas int32 = 1
	var revisionHistoryLimit int32 = 0
	const maxUnavailable int32 = 0
	const maxSurge int32 = 1
	namespace := core.NamespaceDefault
	labels := getLabels(lambda)
	name := getName(lambda)

	deployment := &apps.Deployment{
		ObjectMeta: meta.ObjectMeta{Name: name, Namespace: namespace, Labels: labels},
		Spec: apps.DeploymentSpec{
			Replicas: &replicas,
			Strategy: apps.DeploymentStrategy{
				Type: apps.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &apps.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{Type: intstr.Int, IntVal: maxUnavailable},
					MaxSurge:       &intstr.IntOrString{Type: intstr.Int, IntVal: maxSurge},
				},
			},
			RevisionHistoryLimit: &revisionHistoryLimit,
			Selector:             &meta.LabelSelector{MatchLabels: labels},
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{Labels: labels},
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name:            name,
							Image:           runtime.Image,
							ImagePullPolicy: core.PullAlways,
							Command:         []string{"sh", "-c", fmt.Sprintf("echo %v && "+runtime.RunCommand, time.Now(), lambda.Id)},
						},
					},
				},
			},
		},
	}
	_, err = k8s.Client.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, meta.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func updateDeployment(lambda *Lambda) error {
	runtime, err := getRuntime(lambda)
	if err != nil {
		return err
	}
	var replicas int32 = 2
	var revisionHistoryLimit int32 = 0
	const maxUnavailable int32 = 0
	const maxSurge int32 = 1
	namespace := core.NamespaceDefault
	labels := getLabels(lambda)
	name := getName(lambda)

	deployment := &apps.Deployment{
		ObjectMeta: meta.ObjectMeta{Name: name, Namespace: namespace, Labels: labels},
		Spec: apps.DeploymentSpec{
			Replicas: &replicas,
			Strategy: apps.DeploymentStrategy{
				Type: apps.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &apps.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{Type: intstr.Int, IntVal: maxUnavailable},
					MaxSurge:       &intstr.IntOrString{Type: intstr.Int, IntVal: maxSurge},
				},
			},
			RevisionHistoryLimit: &revisionHistoryLimit,
			Selector:             &meta.LabelSelector{MatchLabels: labels},
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{Labels: labels},
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name:            name,
							Image:           runtime.Image,
							ImagePullPolicy: core.PullAlways,
							Command:         []string{"sh", "-c", fmt.Sprintf("echo %v && "+runtime.RunCommand, time.Now(), lambda.Id)},
						},
					},
				},
			},
		},
	}
	_, err = k8s.Client.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, meta.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func deleteDeployment(id *primitive.ObjectID) error {
	name := getNameById(id)
	namespace := core.NamespaceDefault

	err := k8s.Client.AppsV1().Deployments(namespace).Delete(context.TODO(), name, meta.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func createService(lambda *Lambda) error {
	const port int32 = 443
	name := getName(lambda)
	labels := getLabels(lambda)
	namespace := core.NamespaceDefault

	service := &core.Service{
		ObjectMeta: meta.ObjectMeta{Name: name, Namespace: namespace, Labels: labels},
		Spec: core.ServiceSpec{
			Type:     core.ServiceTypeNodePort,
			Selector: labels,
			Ports: []core.ServicePort{{
				Port:       port,
				TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: port},
			}},
		},
	}
	_, err := k8s.Client.CoreV1().Services(namespace).Create(context.TODO(), service, meta.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func deleteService(id *primitive.ObjectID) error {
	name := getNameById(id)
	namespace := core.NamespaceDefault

	err := k8s.Client.CoreV1().Services(namespace).Delete(context.TODO(), name, meta.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func updateIngress(lambda *Lambda) error {
	ingressName := "lambda-clone-ingress"
	namespace := core.NamespaceDefault

	ingress, err := k8s.Client.NetworkingV1().Ingresses(namespace).Get(context.TODO(), ingressName, meta.GetOptions{})
	if err != nil {
		return err
	}

	host := "lambda-clone-endpoint.dev-doe.com"
	pathTypePrefix := network.PathTypePrefix
	const port int32 = 443
	var path = "/endpoint/" + lambda.Id
	serviceName := getName(lambda)

	ingress.Spec.Rules = append(ingress.Spec.Rules, network.IngressRule{
		Host: host,
		IngressRuleValue: network.IngressRuleValue{
			HTTP: &network.HTTPIngressRuleValue{
				Paths: []network.HTTPIngressPath{
					{
						Path:     path,
						PathType: &pathTypePrefix,
						Backend: network.IngressBackend{
							Service: &network.IngressServiceBackend{
								Name: serviceName,
								Port: network.ServiceBackendPort{Number: port},
							},
						},
					},
				},
			},
		},
	})
	_, err = k8s.Client.NetworkingV1().Ingresses(namespace).Update(context.TODO(), ingress, meta.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func deleteIngress(id *primitive.ObjectID) error {
	ingressName := "lambda-clone-ingress"
	namespace := core.NamespaceDefault

	ingress, err := k8s.Client.NetworkingV1().Ingresses(namespace).Get(context.TODO(), ingressName, meta.GetOptions{})
	if err != nil {
		return err
	}

	host := "lambda-clone-endpoint.dev-doe.com"
	serviceName := getNameById(id)

	var foundIndex int
	for i, rule := range ingress.Spec.Rules {
		if rule.Host != host {
			continue
		}
		if rule.IngressRuleValue.HTTP.Paths[0].Backend.Service.Name != serviceName {
			continue
		}
		foundIndex = i
		break
	}
	ingress.Spec.Rules = append(ingress.Spec.Rules[:foundIndex], ingress.Spec.Rules[foundIndex+1:]...)

	_, err = k8s.Client.NetworkingV1().Ingresses(namespace).Update(context.TODO(), ingress, meta.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}
