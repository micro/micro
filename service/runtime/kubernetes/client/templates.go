// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Original source: github.com/micro/go-micro/v3/util/kubernetes/client/templates.go

package client

var templates = map[string]string{
	"deployment":      deploymentTmpl,
	"service":         serviceTmpl,
	"namespace":       namespaceTmpl,
	"secret":          secretTmpl,
	"serviceaccount":  serviceAccountTmpl,
	"networkpolicies": networkPolicyTmpl,
	"networkpolicy":   networkPolicyTmpl,
	"resourcequota":   resourceQuotaTmpl,
}

var deploymentTmpl = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: "{{ .Metadata.Name }}"
  namespace: "{{ .Metadata.Namespace }}"
  labels:
    {{- with .Metadata.Labels }}
    {{- range $key, $value := . }}
    {{ $key }}: "{{ $value }}"
    {{- end }}
    {{- end }}
  annotations:
    {{- with .Metadata.Annotations }}
    {{- range $key, $value := . }}
    {{ $key }}: "{{ $value }}"
    {{- end }}
    {{- end }}
spec:
  replicas: {{ .Spec.Replicas }}
  selector:
    matchLabels:
      {{- with .Spec.Selector.MatchLabels }}
      {{- range $key, $value := . }}
      {{ $key }}: "{{ $value }}"
      {{- end }}
      {{- end }}
  template:
    metadata:
      labels:
        {{- with .Spec.Template.Metadata.Labels }}
        {{- range $key, $value := . }}
        {{ $key }}: "{{ $value }}"
        {{- end }}
        {{- end }}
      annotations:
        {{- with .Spec.Template.Metadata.Annotations }}
        {{- range $key, $value := . }}
        {{ $key }}: "{{ $value }}"
        {{- end }}
        {{- end }}
    spec: 
      runtimeClassName: {{ .Spec.Template.PodSpec.RuntimeClassName }}
      serviceAccountName: {{ .Spec.Template.PodSpec.ServiceAccountName }}
      containers:
      {{- with .Spec.Template.PodSpec.Containers }}
      {{- range . }}
        - name: {{ .Name }}
          env:
          {{- with .Env }}
          {{- range . }}
          - name: "{{ .Name }}"
            value: "{{ .Value }}"
          {{- if .ValueFrom }}
          {{- with .ValueFrom }}
            valueFrom: 
              {{- if .SecretKeyRef }}
              {{- with .SecretKeyRef }}
              secretKeyRef:
                key: {{ .Key }}
                name: {{ .Name }}
                optional: {{ .Optional }}
              {{- end }}
              {{- end }}
          {{- end }}
          {{- end }}
          {{- end }}
          {{- end }}
          args:
          {{- range .Args }}
          - {{.}}
          {{- end }}
          command:
          {{- range .Command }}
          - {{.}}
          {{- end }}
          image: {{ .Image }}
          imagePullPolicy: IfNotPresent
          ports:
          {{- with .Ports }}
          {{- range . }}
          - containerPort: {{ .ContainerPort }}
            name: {{ .Name }}
          {{- end }}
          {{- end }}
          {{- if .ReadinessProbe }}
          {{- with .ReadinessProbe }}
          readinessProbe:
            {{- with .TCPSocket }}
            tcpSocket:
              {{- if .Host }}
              host: {{ .Host }}
              {{- end }}
              port: {{ .Port }}
            {{- end }}
            initialDelaySeconds: {{ .InitialDelaySeconds }}
            periodSeconds: {{ .PeriodSeconds }}
          {{- end }}
          {{- end }}
          {{- if .Resources }}
          {{- with .Resources }}
          resources:
            {{- if .Limits }}
            {{- with .Limits }}
            limits:
              {{- if .Memory }}
              memory: {{ .Memory }}
              {{- end }}
              {{- if .CPU }}
              cpu: {{ .CPU }}
              {{- end }}
              {{- if .EphemeralStorage }}
              ephemeral-storage: {{ .EphemeralStorage }}
              {{- end }}
            {{- end }}
            {{- end }}
            {{- if .Requests }}
            {{- with .Requests }}
            requests:
              {{- if .Memory }}
              memory: {{ .Memory }}
              {{- end }}
              {{- if .CPU }}
              cpu: {{ .CPU }}
              {{- end }}
              {{- if .EphemeralStorage }}
              ephemeral-storage: {{ .EphemeralStorage }}
              {{- end }}
            {{- end }}
            {{- end }}
          {{- end }}
          {{- end }}
          volumeMounts:
          {{- with .VolumeMounts }}
          {{- range . }}
            - name: {{ .Name }}
              mountPath: {{ .MountPath }}
          {{- end }}
          {{- end }}
      {{- end }}
      {{- end }} 
      volumes:
      {{- with .Spec.Template.PodSpec.Volumes }}
      {{- range . }}
        - name: {{ .Name }}
          persistentVolumeClaim:
            claimName: {{ .PersistentVolumeClaim.ClaimName }}
      {{- end }}
      {{- end }}
`

var serviceTmpl = `
apiVersion: v1
kind: Service
metadata:
  name: "{{ .Metadata.Name }}"
  namespace: "{{ .Metadata.Namespace }}"
  labels:
    {{- with .Metadata.Labels }}
    {{- range $key, $value := . }}
    {{ $key }}: "{{ $value }}"
    {{- end }}
    {{- end }}
spec:
  selector:
    {{- with .Spec.Selector }}
    {{- range $key, $value := . }}
    {{ $key }}: "{{ $value }}"
    {{- end }}
    {{- end }}
  ports:
  {{- with .Spec.Ports }}
  {{- range . }}
  - name: "{{ .Name }}"
    port: {{ .Port }}
    protocol: {{ .Protocol }}
  {{- end }}
  {{- end }}
`

var namespaceTmpl = `
apiVersion: v1
kind: Namespace
metadata:
  name: "{{ .Metadata.Name }}"
  labels:
    {{- with .Metadata.Labels }}
    {{- range $key, $value := . }}
    {{ $key }}: "{{ $value }}"
    {{- end }}
    {{- end }}
`

var secretTmpl = `
apiVersion: v1
kind: Secret
type: "{{ .Type }}"
metadata:
  name: "{{ .Metadata.Name }}"
  namespace: "{{ .Metadata.Namespace }}"
  labels:
    {{- with .Metadata.Labels }}
    {{- range $key, $value := . }}
    {{ $key }}: "{{ $value }}"
    {{- end }}
    {{- end }}
data:
  {{- with .Data }}
  {{- range $key, $value := . }}
  {{ $key }}: "{{ $value }}"
  {{- end }}
  {{- end }}
`

var serviceAccountTmpl = `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: "{{ .Metadata.Name }}"
  labels:
    {{- with .Metadata.Labels }}
    {{- range $key, $value := . }}
    {{ $key }}: "{{ $value }}"
    {{- end }}
    {{- end }}
imagePullSecrets:
{{- with .ImagePullSecrets }}
{{- range . }}
- name: "{{ .Name }}"
{{- end }}
{{- end }}
`

var networkPolicyTmpl = `
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: "{{ .Metadata.Name }}"
  namespace: "{{ .Metadata.Namespace }}"
  labels:
    {{- with .Metadata.Labels }}
    {{- range $key, $value := . }}
    {{ $key }}: "{{ $value }}"
    {{- end }}
    {{- end }}
spec:
  podSelector:
    matchLabels:
  ingress:
  - from: # Allow pods in this namespace to talk to each other
    - podSelector: {}
  - from: # Allow pods in the namespaces bearing the specified labels to talk to pods in this namespace:
    - namespaceSelector:
        matchLabels:
          {{- with (index (index .Spec.Ingress 1).From 0).NamespaceSelector.MatchLabels }}
          {{- range $key, $value := . }}
          {{ $key }}: "{{ $value }}"
          {{- end }}
          {{- end }}
`

var resourceQuotaTmpl = `
apiVersion: v1
kind: ResourceQuota
metadata:
  name: "{{ .Metadata.Name }}"
  namespace: "{{ .Metadata.Namespace }}"
  labels:
    {{- with .Metadata.Labels }}
    {{- range $key, $value := . }}
    {{ $key }}: "{{ $value }}"
    {{- end }}
    {{- end }}
spec:
  hard:
  {{- if .Spec.Hard }}
  {{- with .Spec.Hard }}
    {{- if .LimitsMemory }}
    limits.memory: {{ .LimitsMemory }}
    {{- end }}
    {{- if .LimitsCPU }}
    limits.cpu: {{ .LimitsCPU }}
    {{- end }}
    {{- if .LimitsEphemeralStorage }}
    limits.ephemeral-storage: {{ .LimitsEphemeralStorage }}
    {{- end }}
    {{- if .RequestsMemory }}
    requests.memory: {{ .RequestsMemory }}
    {{- end }}
    {{- if .RequestsCPU }}
    requests.cpu: {{ .RequestsCPU }}
    {{- end }}
    {{- if .RequestsEphemeralStorage }}
    requests.ephemeral-storage: {{ .RequestsEphemeralStorage }}
    {{- end }}
  {{- end }}
  {{- end }}
`
