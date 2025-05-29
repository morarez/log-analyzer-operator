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
	"bytes"
	"context"
	"io"
	"strings"

	ai "github.com/morarez/log-analyzer-operator/internal/ai"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	aiopsv1 "github.com/morarez/log-analyzer-operator/api/v1"
)

// LogAnalysisRequestReconciler reconciles a LogAnalysisRequest object
type LogAnalysisRequestReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	Clientset kubernetes.Interface
}

// +kubebuilder:rbac:groups=aiops.aiops.dev,resources=loganalysisrequests,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aiops.aiops.dev,resources=loganalysisrequests/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aiops.aiops.dev,resources=loganalysisrequests/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LogAnalysisRequest object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.4/pkg/reconcile
func (r *LogAnalysisRequestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Step 1: Fetch the LogAnalysisRequest
	var analysisReq aiopsv1.LogAnalysisRequest
	if err := r.Get(ctx, req.NamespacedName, &analysisReq); err != nil {
		log.Error(err, "unable to fetch LogAnalysisRequest")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Step 2: Default namespace if not provided
	targetNamespace := analysisReq.Spec.ObjectRef.Namespace
	if targetNamespace == "" {
		targetNamespace = analysisReq.Namespace
	}

	// Step 3: Retrieve pod logs
	podLogOpts := &corev1.PodLogOptions{
		TailLines: analysisReq.Spec.TailLines,
	}
	logReq := r.Clientset.CoreV1().Pods(targetNamespace).GetLogs(analysisReq.Spec.ObjectRef.Name, podLogOpts)
	logStream, err := logReq.Stream(ctx)
	if err != nil {
		log.Error(err, "unable to stream logs from pod", "pod", analysisReq.Spec.ObjectRef.Name)
		return ctrl.Result{}, err
	}
	defer logStream.Close()

	logBuf := new(bytes.Buffer)
	if _, err := io.Copy(logBuf, logStream); err != nil {
		log.Error(err, "error reading pod logs")
		return ctrl.Result{}, err
	}

	logs := logBuf.String()

	// Step 4: Run AI analysis
	diagnosis, err := ai.AnalyzeWithAI(logs, "gpt-4o")
	if err != nil {
		log.Error(err, "AI analysis failed")
		return ctrl.Result{}, err
	}

	// Step 5: Update status
	analysisReq.Status.Diagnosis = diagnosis
	analysisReq.Status.Timestamp = metav1.Now()
	analysisReq.Status.Resolved = strings.Contains(strings.ToLower(diagnosis), "resolved")

	if err := r.Status().Update(ctx, &analysisReq); err != nil {
		log.Error(err, "unable to update LogAnalysisRequest status")
		return ctrl.Result{}, err
	}

	log.Info("log analysis complete", "diagnosis", diagnosis)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LogAnalysisRequestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aiopsv1.LogAnalysisRequest{}).
		Named("loganalysisrequest").
		Complete(r)
}
