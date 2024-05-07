/*
Copyright 2024.

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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	gatewayv1alpha1 "github.com/shinemost/test-operator/api/v1alpha1"
)

const (
	Require_Namespace = "test-ns"
	Requre_Replicas   = 2
)

// MyProxyReconciler reconciles a MyProxy object
type MyProxyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=gateway.shinemost.top,resources=myproxies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=gateway.shinemost.top,resources=myproxies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=gateway.shinemost.top,resources=myproxies/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MyProxy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *MyProxyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// 先查找CRD对象，此处是MyProxy 是否存在
	myProxy := &gatewayv1alpha1.MyProxy{}
	err := r.Get(ctx, req.NamespacedName, myProxy)

	if err != nil {
		// 只有在对象已经被删除的情况下错误才会被忽略
		if errors.IsNotFound(err) {
			logger.Info("Resource not found. Error ignored as the resource must have been deleted.")
			return ctrl.Result{}, nil
		}
		// 其他错误直接返回
		logger.Error(err, "获取 MyProxy 实例发生错误")
		return ctrl.Result{}, err
	}

	// 如果CRD对象存在，判断部署是否存在
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: myProxy.Spec.Name, Namespace: Require_Namespace}, found)

	// 如果没有则新建
	if err != nil && errors.IsNotFound(err) {
		dep := r.deploymentForExample(myProxy)
		logger.Info("新建一个部署", "命名空间", dep.Namespace, "部署名称", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			logger.Error(err, "新建部署报错", "命名空间", dep.Namespace, "部署名称", dep.Name)
			return ctrl.Result{}, err
		}
		// 通知控制器重新入队，等待调谐
		return ctrl.Result{Requeue: true}, nil

	} else if err != nil {
		logger.Error(err, "获取部署报错")
		return ctrl.Result{}, err
	}

	// 如果部署存在，比较目前部署的replicas与requre_replicas数量是否一致,如果不一致，调谐到requre_replicas
	if *found.Spec.Replicas != Requre_Replicas {
		logger.Info("目前部署实例数", found.Status.Replicas, "期望部署实例数", Requre_Replicas)
		var replicas int32 = Requre_Replicas
		found.Spec.Replicas = &replicas
		err = r.Update(ctx, found)
		if err != nil {
			logger.Error(err, "更新部署实例数报错")
		}
		return ctrl.Result{Requeue: true}, err
	}
	return ctrl.Result{}, nil
}

// 创建deployment的私有方法，所属命名空间为test_ns
func (r *MyProxyReconciler) deploymentForExample(myproxy *gatewayv1alpha1.MyProxy) *appsv1.Deployment {
	dep := &appsv1.Deployment{}
	dep.Namespace = Require_Namespace
	dep.Name = myproxy.Spec.Name
	var replicas int32 = Requre_Replicas
	labels := map[string]string{
		"test_label": myproxy.Spec.Name,
	}
	dep.Spec = appsv1.DeploymentSpec{
		Replicas: &replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: labels,
		},
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "nginx",
						Image: "nginx",
					},
				},
			},
		},
	}
	dep.Labels = labels
	dep.Spec.Template.Labels = labels
	return dep
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyProxyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gatewayv1alpha1.MyProxy{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
