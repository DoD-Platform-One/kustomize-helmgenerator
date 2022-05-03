package main

import (
	"testing"
)

func Test_processHelmGenerator(t *testing.T) {
	type args struct {
		fn string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"HelmGenerator: Valid",
			args{"testdata/generator.yaml"},
			`---
# Source: mocha/templates/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: mocha
spec:
  type: ClusterIP
  ports:
  - port: 80
    targetPort: 80
    protocol: TCP
    name: http
---
# Source: mocha/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mocha
spec:
  replicas: 99
  template:
    spec:
      containers:
        - name: mocha
          image: "donkers:1.16.0"
          imagePullPolicy: Always
          ports:
            - name: http
              containerPort: 80
              protocol: TCP
`,
			false,
		},
		{
			"HelmGenerator: Chart Fail",
			args{"testdata/generator-chart-fail.yaml"},
			``,
			true,
		},
		{
			"HelmGenerator: Values File Fail",
			args{"testdata/generator-values-file-fail.yaml"},
			``,
			true,
		},
		{
			"HelmGenerator: Values Inline Fail",
			args{"testdata/generator-values-fail.yaml"},
			``,
			true,
		},
		{
			"HelmGenerator: Sops Fail",
			args{"testdata/generator-sops-fail.yaml"},
			``,
			true,
		},
		{
			"HelmGenerator: Template Fail",
			args{"testdata/generator-template-fail.yaml"},
			``,
			true,
		},
		{
			"HelmGenerator: File Fail",
			args{"testdata/no-such-file.yaml"},
			``,
			true,
		},
		{
			"HelmGenerator: YAML Fail",
			args{"testdata/generator-bad-yaml.yaml"},
			``,
			true,
		},
		{
			"HelmGenerator: Hooks",
			args{"testdata/generator-hooks.yaml"},
			`---
# Source: mocha/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mocha
spec:
  replicas: 99
  template:
    spec:
      containers:
        - name: mocha
          image: "donkers:1.16.0"
          imagePullPolicy: Always
          ports:
            - name: http
              containerPort: 80
              protocol: TCP

---
apiVersion: v1
kind: Service
metadata:
  name: mocha
  annotations:
    "helm.sh/hook": post-install
    "helm.sh/hook-weight": "-5"
spec:
  type: ClusterIP
  ports:
  - port: 80
    targetPort: 80
    protocol: TCP
    name: http`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processHelmGenerator(tt.args.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("processHelmGenerator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("processHelmGenerator() got = %v, want %v", got, tt.want)
			}
		})
	}
}
