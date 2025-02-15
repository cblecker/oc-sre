module github.com/cblecker/oc-sre

go 1.12

require (
	github.com/aws/aws-sdk-go v1.24.4
	github.com/pkg/browser v0.0.0-20180916011732-0a3d74bf9ce4
	github.com/sethvargo/go-password v0.1.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	k8s.io/apimachinery v0.0.0-20190717022731-0bb8574e0887
	k8s.io/cli-runtime v0.0.0-20190717024643-59adbd30f884
	k8s.io/client-go v0.0.0-20190717023132-0c47f9da0001
	k8s.io/kubectl v0.0.0-20190717025603-972620de1df0
)

replace (
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20181025213731-e84da0312774
	golang.org/x/net => golang.org/x/net v0.0.0-20190206173232-65e2d4e15006
	golang.org/x/sync => golang.org/x/sync v0.0.0-20181108010431-42b317875d0f
	golang.org/x/sys => golang.org/x/sys v0.0.0-20190209173611-3b5209105503
	golang.org/x/text => golang.org/x/text v0.3.1-0.20181227161524-e6919f6577db
	golang.org/x/tools => golang.org/x/tools v0.0.0-20190313210603-aa82965741a9
	k8s.io/api => k8s.io/api v0.0.0-20190717022910-653c86b0609b
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190717022731-0bb8574e0887
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20190717024643-59adbd30f884
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190717023132-0c47f9da0001
	k8s.io/kubectl => k8s.io/kubectl v0.0.0-20190717025603-972620de1df0
)
