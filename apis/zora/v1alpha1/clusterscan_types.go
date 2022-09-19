// Copyright 2022 Undistro Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/getupio-undistro/zora/pkg/apis"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// +kubebuilder:validation:Minimum=0
// +kubebuilder:validation:Maximum=6
type DayOfWeek int

type Schedule struct {
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=23
	HourlyRep int `json:"hourlyRep,omitempty"`
	// +kubebuilder:validation:MaxItems=7
	DaysOfWeek []DayOfWeek `json:"daysOfWeek,omitempty"`
	// +kubebuilder:validation:Pattern=`^(2[0-3]|[0-1]?[0-9]):[0-5]?[0-9]$`
	StartTime *string `json:"startTime,omitempty"`
}

// ClusterScanSpec defines the desired state of ClusterScan
type ClusterScanSpec struct {
	// ClusterRef is a reference to a Cluster in the same namespace
	ClusterRef corev1.LocalObjectReference `json:"clusterRef"`

	// This flag tells the controller to suspend subsequent executions, it does
	// not apply to already started executions.  Defaults to false.
	Suspend *bool `json:"suspend,omitempty"`

	Schedule *Schedule `json:"schedule"`

	// The list of Plugin references that are used to scan the referenced Cluster.  Defaults to 'popeye'
	Plugins []PluginReference `json:"plugins,omitempty"`
}

type PluginReference struct {
	// Name is unique within a namespace to reference a Plugin resource.
	Name string `json:"name"`

	// Namespace defines the space within which the Plugin name must be unique.
	Namespace string `json:"namespace,omitempty"`

	// This flag tells the controller to suspend subsequent executions, it does
	// not apply to already started executions.  Defaults to false.
	Suspend *bool `json:"suspend,omitempty"`

	// The schedule in Cron format for this Plugin, see https://en.wikipedia.org/wiki/Cron.
	Schedule *Schedule `json:"schedule,omitempty"`

	// List of environment variables to set in the Plugin container.
	Env []corev1.EnvVar `json:"env,omitempty"`
}

func (in *PluginReference) PluginKey(defaultNamespace string) types.NamespacedName {
	ns := in.Namespace
	if ns == "" {
		ns = defaultNamespace
	}
	return types.NamespacedName{Name: in.Name, Namespace: ns}
}

// ClusterScanStatus defines the observed state of ClusterScan
type ClusterScanStatus struct {
	apis.Status `json:",inline"`

	// Information of the last scans of plugins
	Plugins map[string]*PluginScanStatus `json:"plugins,omitempty"`

	// Comma separated list of plugins
	PluginNames string `json:"pluginNames,omitempty"`

	// Suspend field value from ClusterScan spec
	Suspend bool `json:"suspend"`

	// Information when was the last time the job was scheduled.
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty"`

	// Information when was the last time the job was finished.
	LastFinishedTime *metav1.Time `json:"lastFinishedTime,omitempty"`

	// Status of the last finished scan. Complete or Failed
	LastFinishedStatus string `json:"lastFinishedStatus,omitempty"`

	// Status of the last scan. Active, Complete or Failed
	LastStatus string `json:"lastStatus,omitempty"`

	// Information when was the last time the job successfully completed.
	LastSuccessfulTime *metav1.Time `json:"lastSuccessfulTime,omitempty"`

	// Time when the next job will schedule.
	NextScheduleTime *metav1.Time `json:"nextScheduleTime,omitempty"`

	// Total of ClusterIssues reported in the last successful scan
	TotalIssues *int `json:"totalIssues,omitempty"`
}

// GetPluginStatus returns a PluginScanStatus of a plugin
func (in *ClusterScanStatus) GetPluginStatus(name string) *PluginScanStatus {
	if in.Plugins == nil {
		in.Plugins = make(map[string]*PluginScanStatus)
	}
	if _, ok := in.Plugins[name]; !ok {
		in.Plugins[name] = &PluginScanStatus{}
	}
	return in.Plugins[name]
}

// SyncStatus fills ClusterScan status and time fields based on PluginStatus
func (in *ClusterScanStatus) SyncStatus() {
	var names []string
	var failed, active, complete int
	in.NextScheduleTime = nil
	for n, p := range in.Plugins {
		names = append(names, n)
		if in.LastScheduleTime == nil || in.LastScheduleTime.Before(p.LastScheduleTime) {
			in.LastScheduleTime = p.LastScheduleTime
		}
		if in.LastFinishedTime == nil || in.LastFinishedTime.Before(p.LastFinishedTime) {
			in.LastFinishedTime = p.LastFinishedTime
		}
		if in.LastSuccessfulTime == nil || in.LastSuccessfulTime.Before(p.LastSuccessfulTime) {
			in.LastSuccessfulTime = p.LastSuccessfulTime
		}
		if in.NextScheduleTime == nil || p.NextScheduleTime.Before(in.NextScheduleTime) {
			in.NextScheduleTime = p.NextScheduleTime
		}
		if p.LastStatus == "Active" {
			active++
		}
		switch p.LastFinishedStatus {
		case string(batchv1.JobFailed):
			failed++
		case string(batchv1.JobComplete):
			complete++
		}
	}

	if failed > 0 {
		in.LastFinishedStatus = string(batchv1.JobFailed)
		in.LastStatus = string(batchv1.JobFailed)
	}
	if failed == 0 && complete > 0 {
		in.LastFinishedStatus = string(batchv1.JobComplete)
		in.LastStatus = string(batchv1.JobComplete)
	}
	if active > 0 {
		in.LastStatus = "Active"
	}

	sort.Strings(names)
	in.PluginNames = strings.Join(names, ",")
}

// LastScanIDs returns a list of all the last scan IDs
func (in *ClusterScanStatus) LastScanIDs(successful bool) []string {
	lastScans := make([]string, 0, len(in.Plugins))
	for _, ps := range in.Plugins {
		sid := ps.LastScanID
		if successful {
			sid = ps.LastSuccessfulScanID
		}
		if sid != "" {
			lastScans = append(lastScans, sid)
		}
	}
	return lastScans
}

// +k8s:deepcopy-gen=true
type PluginScanStatus struct {
	// Information when was the last time the job was scheduled.
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty"`

	// Information when was the last time the job was finished.
	LastFinishedTime *metav1.Time `json:"lastFinishedTime,omitempty"`

	// Information when was the last time the job successfully completed.
	LastSuccessfulTime *metav1.Time `json:"lastSuccessfulTime,omitempty"`

	// Time when the next job will schedule.
	NextScheduleTime *metav1.Time `json:"nextScheduleTime,omitempty"`

	// ID of the last plugin scan
	LastScanID string `json:"lastScanID,omitempty"`

	// ID of the last successful plugin scan
	LastSuccessfulScanID string `json:"lastSuccessfulScanID,omitempty"`

	// Status of the last plugin scan. Active, Complete or Failed
	LastStatus string `json:"lastStatus,omitempty"`

	// Status of the last finished plugin scan. Complete or Failed
	LastFinishedStatus string `json:"lastFinishedStatus,omitempty"`

	// LastErrorMsg contains a plugin error message from the last failed scan.
	LastErrorMsg string `json:"lastErrorMsg,omitempty"`

	// IssueCount holds the sum of ClusterIssues found in the last successful
	// scan.
	IssueCount *int `json:"issueCount,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName="cscan"
//+kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".spec.clusterRef.name",priority=0
//+kubebuilder:printcolumn:name="Suspend",type="boolean",JSONPath=".status.suspend",priority=0
//+kubebuilder:printcolumn:name="Plugins",type="string",JSONPath=".status.pluginNames",priority=0
//+kubebuilder:printcolumn:name="Last Status",type="string",JSONPath=".status.lastStatus",priority=0
//+kubebuilder:printcolumn:name="Last Schedule",type="date",JSONPath=".status.lastScheduleTime",priority=0
//+kubebuilder:printcolumn:name="Last Successful",type="date",JSONPath=".status.lastSuccessfulTime",priority=0
//+kubebuilder:printcolumn:name="Issues",type="integer",JSONPath=".status.totalIssues",priority=0
//+kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].status",priority=0
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",priority=0
//+kubebuilder:printcolumn:name="Hourly Repetition",type="string",JSONPath=".spec.schedule.hourlyRep",priority=1
//+kubebuilder:printcolumn:name="Days Of Week",type="string",JSONPath=".spec.schedule.daysOfWeek",priority=1
//+kubebuilder:printcolumn:name="Start Time",type="string",JSONPath=".spec.schedule.startTime",priority=1
//+kubebuilder:printcolumn:name="Next Schedule",type="string",JSONPath=".status.nextScheduleTime",priority=1

// ClusterScan is the Schema for the clusterscans API
//+genclient
//+genclient:onlyVerbs=list,get
//+genclient:noStatus
type ClusterScan struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterScanSpec   `json:"spec,omitempty"`
	Status ClusterScanStatus `json:"status,omitempty"`
}

func (in *ClusterScan) SetReadyStatus(status bool, reason, msg string) {
	s := metav1.ConditionFalse
	if status {
		s = metav1.ConditionTrue
	}
	in.Status.SetCondition(metav1.Condition{
		Type:               "Ready",
		Status:             s,
		ObservedGeneration: in.Generation,
		Reason:             reason,
		Message:            msg,
	})
}

func (in *ClusterScan) ClusterKey() types.NamespacedName {
	return types.NamespacedName{Name: in.Spec.ClusterRef.Name, Namespace: in.Namespace}
}

//+kubebuilder:object:root=true

// ClusterScanList contains a list of ClusterScan
type ClusterScanList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterScan `json:"items"`
}

func (r *Schedule) SplitStartTime() (int, int) {
	if r.StartTime == nil {
		now := time.Now().UTC()
		return now.Hour(), now.Minute()
	}
	tslice := strings.Split(*r.StartTime, ":")
	res := [2]int{0, 0}
	for c := 0; c < len(tslice); c++ {
		i, err := strconv.Atoi(tslice[c])
		if err != nil {
			i = 0
		}
		res[c] = i
	}
	return res[0], res[1]
}

func (r *Schedule) HourlyRepetitions() string {
	if r == nil {
		return ""
	}
	if r.HourlyRep == 1 {
		return "*"
	}
	h, _ := r.SplitStartTime()
	if r.HourlyRep == 0 {
		return strconv.Itoa(h)
	}

	reps := ""
	for c := h % r.HourlyRep; c <= 23; c += r.HourlyRep {
		reps += fmt.Sprintf("%d", c)
		if c+r.HourlyRep <= 23 {
			reps = fmt.Sprintf("%s,", reps)
		}
	}
	return reps
}

func (r *Schedule) DayOfWeekSeries() string {
	if len(r.DaysOfWeek) == 0 {
		return "*"
	}

	uniq := map[int]struct{}{}
	for c := 0; c < len(r.DaysOfWeek); c++ {
		uniq[int(r.DaysOfWeek[c])] = struct{}{}
	}
	dow := []string{}
	for d, _ := range uniq {
		dow = append(dow, strconv.Itoa(d))
	}
	return strings.Join(dow, ",")
}

func (r *Schedule) CronExpr() string {
	if r == nil {
		return ""
	}
	_, m := r.SplitStartTime()
	return fmt.Sprintf("%d %s * * %s", m, r.HourlyRepetitions(), r.DayOfWeekSeries())
}

func init() {
	SchemeBuilder.Register(&ClusterScan{}, &ClusterScanList{})
}
