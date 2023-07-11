/*
Copyright 2023.

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

package instrumentation

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.5.0"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/newrelic-experimental/newrelic-agent-operator/api/v1alpha1"
	"github.com/newrelic-experimental/newrelic-agent-operator/pkg/constants"
)

const (
	volumeName        = "newrelic-instrumentation"
	initContainerName = "newrelic-instrumentation"
)

type sdkInjector struct {
	client client.Client
	logger logr.Logger
}

func (i *sdkInjector) inject(ctx context.Context, insts languageInstrumentations, ns corev1.Namespace, pod corev1.Pod, containerName string) corev1.Pod {
	if len(pod.Spec.Containers) < 1 {
		return pod
	}

	// We search for specific container to inject variables and if no one is found
	// We fallback to first container
	var index = 0
	for idx, ctnair := range pod.Spec.Containers {
		if ctnair.Name == containerName {
			index = idx
		}
	}

	if insts.Java != nil {
		newrelic := *insts.Java
		var err error
		i.logger.V(1).Info("injecting Java instrumentation into pod", "newrelic-namespace", newrelic.Namespace, "newrelic-name", newrelic.Name)
		pod, err = injectJavaagent(newrelic.Spec.Java, pod, index)
		if err != nil {
			i.logger.Info("Skipping javaagent injection", "reason", err.Error(), "container", pod.Spec.Containers[index].Name)
		} else {
			pod = i.injectNewrelicConfig(ctx, newrelic, ns, pod, index)
		}
	}
	if insts.NodeJS != nil {
		newrelic := *insts.NodeJS
		var err error
		i.logger.V(1).Info("injecting NodeJS instrumentation into pod", "newrelic-namespace", newrelic.Namespace, "newrelic-name", newrelic.Name)
		pod, err = injectNodeJSSDK(newrelic.Spec.NodeJS, pod, index)
		if err != nil {
			i.logger.Info("Skipping NodeJS SDK injection", "reason", err.Error(), "container", pod.Spec.Containers[index].Name)
		} else {
			pod = i.injectNewrelicConfig(ctx, newrelic, ns, pod, index)
		}
	}
	if insts.Python != nil {
		newrelic := *insts.Python
		var err error
		i.logger.V(1).Info("injecting Python instrumentation into pod", "newrelic-namespace", newrelic.Namespace, "newrelic-name", newrelic.Name)
		pod, err = injectPythonSDK(newrelic.Spec.Python, pod, index)
		if err != nil {
			i.logger.Info("Skipping Python SDK injection", "reason", err.Error(), "container", pod.Spec.Containers[index].Name)
		} else {
			pod = i.injectNewrelicConfig(ctx, newrelic, ns, pod, index)
		}
	}
	if insts.DotNet != nil {
		newrelic := *insts.DotNet
		var err error
		i.logger.V(1).Info("injecting DotNet instrumentation into pod", "newrelic-namespace", newrelic.Namespace, "newrelic-name", newrelic.Name)
		pod, err = injectDotNetSDK(newrelic.Spec.DotNet, pod, index)
		if err != nil {
			i.logger.Info("Skipping DotNet SDK injection", "reason", err.Error(), "container", pod.Spec.Containers[index].Name)
		} else {
			pod = i.injectNewrelicConfig(ctx, newrelic, ns, pod, index)
		}
	}
	return pod
}

func (i *sdkInjector) injectNewrelicConfig(ctx context.Context, newrelic v1alpha1.Instrumentation, ns corev1.Namespace, pod corev1.Pod, index int) corev1.Pod {
	container := &pod.Spec.Containers[index]
	resourceMap := i.createResourceMap(ctx, newrelic, ns, pod, index)
	idx := getIndexOfEnv(container.Env, constants.EnvNewRelicAppName)
	if idx == -1 {
		container.Env = append(container.Env, corev1.EnvVar{
			Name:  constants.EnvNewRelicAppName,
			Value: chooseServiceName(pod, resourceMap, index),
		})
	}
	idx = getIndexOfEnv(container.Env, constants.EnvNewRelicLicenseKey)
	if idx == -1 {
		optional := true
		container.Env = append(container.Env, corev1.EnvVar{
			Name: constants.EnvNewRelicLicenseKey,
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: "newrelic-key-secret"},
					Key:                  "new_relic_license_key",
					Optional:             &optional,
				},
			},
		})
	}
	return pod
}

func chooseServiceName(pod corev1.Pod, resources map[string]string, index int) string {
	if name := resources[string(semconv.K8SDeploymentNameKey)]; name != "" {
		return name
	}
	if name := resources[string(semconv.K8SStatefulSetNameKey)]; name != "" {
		return name
	}
	if name := resources[string(semconv.K8SJobNameKey)]; name != "" {
		return name
	}
	if name := resources[string(semconv.K8SCronJobNameKey)]; name != "" {
		return name
	}
	if name := resources[string(semconv.K8SPodNameKey)]; name != "" {
		return name
	}
	return pod.Spec.Containers[index].Name
}

// createResourceMap creates resource attribute map.
// User defined attributes (in explicitly set env var) have higher precedence.
func (i *sdkInjector) createResourceMap(ctx context.Context, newrelic v1alpha1.Instrumentation, ns corev1.Namespace, pod corev1.Pod, index int) map[string]string {
	// get existing resources env var and parse it into a map
	existingRes := map[string]bool{}

	res := map[string]string{}

	k8sResources := map[attribute.Key]string{}
	k8sResources[semconv.K8SNamespaceNameKey] = ns.Name
	k8sResources[semconv.K8SContainerNameKey] = pod.Spec.Containers[index].Name
	// Some fields might be empty - node name, pod name
	// The pod name might be empty if the pod is created form deployment template
	k8sResources[semconv.K8SPodNameKey] = pod.Name
	k8sResources[semconv.K8SPodUIDKey] = string(pod.UID)
	k8sResources[semconv.K8SNodeNameKey] = pod.Spec.NodeName
	for k, v := range k8sResources {
		if !existingRes[string(k)] && v != "" {
			res[string(k)] = v
		}
	}
	return res
}

func getIndexOfEnv(envs []corev1.EnvVar, name string) int {
	for i := range envs {
		if envs[i].Name == name {
			return i
		}
	}
	return -1
}

func validateContainerEnv(envs []corev1.EnvVar, envsToBeValidated ...string) error {
	for _, envToBeValidated := range envsToBeValidated {
		for _, containerEnv := range envs {
			if containerEnv.Name == envToBeValidated {
				if containerEnv.ValueFrom != nil {
					return fmt.Errorf("the container defines env var value via ValueFrom, envVar: %s", containerEnv.Name)
				}
				break
			}
		}
	}
	return nil
}
