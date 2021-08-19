// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fluentbit

import (
	"fmt"

	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *Reconciler) prometheusRules() (runtime.Object, reconciler.DesiredState, error) {
	if r.Logging.Spec.FluentbitSpec.Metrics != nil && r.Logging.Spec.FluentbitSpec.Metrics.PrometheusRules {
		objectMetadata := r.FluentbitObjectMeta(fluentbitServiceName)
		nsJobLabel := fmt.Sprintf(`job="%s", namespace="%s"`, objectMetadata.Name, objectMetadata.Namespace)

		return &v1.PrometheusRule{
			ObjectMeta: objectMetadata,
			Spec: v1.PrometheusRuleSpec{
				Groups: []v1.RuleGroup{{
					Name: "fluentbit",
					Rules: []v1.Rule{
						{
							Alert: "FluentbitTooManyErrors",
							Expr: intstr.IntOrString{
								Type:   intstr.String,
								StrVal: fmt.Sprintf("rate(fluentbit_output_retries_failed_total{%s}[10m]) > 0", nsJobLabel),
							},
							For: "10m",
							Labels: map[string]string{
								"service":  "fluentbit",
								"severity": "warning",
							},
							Annotations: map[string]string{
								"summary":     `Fluentbit too many errors.`,
								"description": `Fluentbit ({{ $labels.instance }}) is erroring.`,
							},
						},
					},
				},
				},
			},
		}, reconciler.StatePresent, nil
	}
	return &v1.PrometheusRule{
		ObjectMeta: r.FluentbitObjectMeta(fluentbitServiceName),
		Spec:       v1.PrometheusRuleSpec{},
	}, reconciler.StateAbsent, nil
}
