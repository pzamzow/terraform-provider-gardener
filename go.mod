module github.com/kyma-incubator/terraform-provider-gardener

go 1.13

require (
	contrib.go.opencensus.io/exporter/ocagent v0.4.12 // indirect
	github.com/Azure/azure-pipeline-go v0.2.2 // indirect
	github.com/Azure/go-autorest v11.5.0+incompatible // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/NYTimes/gziphandler v1.1.1 // indirect
	github.com/containerd/continuity v0.0.0-20200107194136-26c1120b8d41 // indirect
	github.com/coreos/bbolt v1.3.3 // indirect
	github.com/docker/spdystream v0.0.0-20181023171402-6480d4af844c // indirect
	github.com/elazarl/goproxy v0.0.0-20191011121108-aa519ddbe484 // indirect
	github.com/gardener/gardener v1.3.1-0.20200423115721-07c09fe6cf3e
	github.com/gardener/gardener-extension-provider-aws v1.6.2
	github.com/gardener/gardener-extension-provider-azure v1.6.1
	github.com/gardener/gardener-extension-provider-gcp v1.5.2
	github.com/gardener/gardener-extensions v1.2.1 // indirect
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8 // indirect
	github.com/go-ini/ini v1.46.0 // indirect
	github.com/gobuffalo/genny v0.0.0-20190403191548-3ca520ef0d9e // indirect
	github.com/golang/gddo v0.0.0-20190904175337-72a348e765d2 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/google/go-cmp v0.4.0
	github.com/gorilla/mux v1.7.0 // indirect
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/gotestyourself/gotestyourself v2.2.0+incompatible // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.11.3 // indirect
	github.com/hashicorp/go-hclog v0.9.2 // indirect
	github.com/hashicorp/go-plugin v1.0.1 // indirect
	github.com/hashicorp/hil v0.0.0-20190212132231-97b3a9cdfa93 // indirect
	github.com/hashicorp/terraform v0.12.13
	github.com/hashicorp/yamux v0.0.0-20190923154419-df201c70410d // indirect
	github.com/mattn/go-ieproxy v0.0.0-20190805055040-f9202b1cfdeb // indirect
	github.com/ory/go-acc v0.1.0 // indirect
	github.com/ory/x v0.0.94 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.6.2 // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	gopkg.in/ini.v1 v1.52.0 // indirect
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	sigs.k8s.io/structured-merge-diff v1.0.2 // indirect
)

replace (
	github.com/census-instrumentation/opencensus-proto v0.1.0-0.20181214143942-ba49f56771b8 => github.com/census-instrumentation/opencensus-proto v0.0.3-0.20181214143942-ba49f56771b8
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655 // kubernetes-1.16.0
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90 // kubernetes-1.16.0

)
