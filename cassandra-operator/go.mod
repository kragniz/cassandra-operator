module github.com/sky-uk/cassandra-operator/cassandra-operator

go 1.12

require (
	github.com/PaesslerAG/gval v0.1.1 // indirect
	github.com/PaesslerAG/jsonpath v0.1.0
	github.com/evanphx/json-patch v4.2.0+incompatible // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/gofrs/flock v0.7.1 // indirect
	github.com/golang/groupcache v0.0.0-20181024230925-c65c006176ff // indirect
	github.com/googleapis/gnostic v0.2.0 // indirect
	github.com/imdario/mergo v0.3.6 // indirect
	github.com/kubernetes-sigs/controller-tools v0.1.10
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/prometheus/client_golang v0.9.2
	github.com/prometheus/common v0.1.0
	github.com/prometheus/procfs v0.0.0-20190104112138-b1a0a9a36d74 // indirect
	github.com/robfig/cron v1.1.0
	github.com/sirupsen/logrus v1.3.0
	github.com/sky-uk/licence-compliance-checker v1.1.0
	github.com/spf13/cobra v0.0.3
	github.com/theckman/go-flock v0.7.0
	golang.org/x/lint v0.0.0-20190409202823-959b441ac422
	golang.org/x/oauth2 v0.0.0-20181203162652-d668ce993890 // indirect
	golang.org/x/time v0.0.0-20181108054448-85acf8d2951c // indirect
	golang.org/x/tools v0.0.0-20190501045030-23463209683d
	google.golang.org/appengine v1.4.0 // indirect
	gopkg.in/src-d/go-license-detector.v2 v2.0.1 // indirect
	k8s.io/api v0.0.0-20190515023547-db5a9d1c40eb
	k8s.io/apiextensions-apiserver v0.0.0-20190515024537-2fd0e9006049 // indirect
	k8s.io/apimachinery v0.0.0-20190515023456-b74e4c97951f
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/code-generator v0.0.0-20190511023357-639c964206c2
	k8s.io/utils v0.0.0-20190506122338-8fab8cb257d5 // indirect
	sigs.k8s.io/controller-tools v0.2.0-alpha.2 // indirect
)

replace k8s.io/client-go => k8s.io/client-go v11.0.0+incompatible

replace k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190216013122-f05b8decd79c

replace sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.2.0-alpha.2
