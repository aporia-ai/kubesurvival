module github.com/aporia-ai/kubesurvival/v2

go 1.15

require (
	github.com/containerd/containerd v1.2.5 // indirect
	github.com/coreos/bbolt v1.3.6 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/cpuguy83/strongerrors v0.2.1 // indirect
	github.com/cristim/ec2-instances-info v0.0.0-20210201160642-80270dab05f8
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/evanphx/json-patch v0.5.2 // indirect
	github.com/golang/groupcache v0.0.0-20190129154638-5b532d6fd5ef // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/google/gofuzz v0.0.0-20170612174753-24818f796faf // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/googleapis/gnostic v0.2.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.14.6 // indirect
	github.com/hashicorp/golang-lru v0.5.1 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/json-iterator/go v1.1.6 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	github.com/onsi/ginkgo v1.16.4 // indirect
	github.com/onsi/gomega v1.13.0 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/pfnet-research/k8s-cluster-simulator v0.0.0-20190415111150-55e4108275b4
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v0.9.2 // indirect
	github.com/prometheus/common v0.2.0 // indirect
	github.com/prometheus/procfs v0.0.0-20190319124303-40f3c57fb198 // indirect
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/spf13/afero v1.2.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.3.2 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/tmc/grpc-websocket-proxy v0.0.0-20201229170055-e5319fda7802 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/bbolt v1.3.6 // indirect
	go.uber.org/zap v1.17.0 // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	google.golang.org/appengine v1.5.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/grpc v1.38.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/square/go-jose.v2 v2.3.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools v2.2.0+incompatible // indirect
	k8s.io/api v0.0.0-20190313115550-3c12c96769cc
	k8s.io/apiextensions-apiserver v0.0.0-20190320070711-2af94a2a482f // indirect
	k8s.io/apimachinery v0.0.0-20190320104356-82cbdc1b6ac2
	k8s.io/apiserver v0.0.0-20190321025803-be70ee97012b // indirect
	k8s.io/client-go v11.0.0+incompatible // indirect
	k8s.io/cloud-provider v0.0.0-20190313124351-c76aa0a348b5 // indirect
	k8s.io/csi-translation-lib v0.0.0-20190313124639-7f5cabc6aac8 // indirect
	k8s.io/klog v0.2.0 // indirect
	k8s.io/kube-openapi v0.0.0-20190320154901-5e45bb682580 // indirect
	k8s.io/kubernetes v1.14.0-rc.1
	k8s.io/utils v0.0.0-20190308190857-21c4ce38f2a7 // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect

)

replace (
	github.com/coreos/bbolt v1.3.6 => go.etcd.io/bbolt v1.3.6
	github.com/pfnet-research/k8s-cluster-simulator => github.com/zorro786/k8s-cluster-simulator v0.0.0-20190415111150-55e4108275b4
	go.etcd.io/bbolt v1.3.6 => github.com/coreos/bbolt v1.3.6
	google.golang.org/grpc v1.38.0 => google.golang.org/grpc v1.26.0
)
