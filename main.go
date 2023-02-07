package main

import (
	"context"
	"os"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	config "github.com/tehlers320/apiserver-load-tester/pkg/config"
	k8sattacks "github.com/tehlers320/apiserver-load-tester/pkg/k8sattacks"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func init() {
	switch os.Getenv("LOGGING_LEVEL") {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "trace":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}


func main() {
	config.InitConfig()
	k8s := config.K8s{Local: true}
	k8sClient, err := k8s.CreateClient()
	if err != nil {
		log.Error(err, "unable to start k8s client")
		os.Exit(1)
	}

	list, err := k8sClient.CoreV1().ConfigMaps("default").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Errorf("%s", err)
	}

	log.Infof("existing configmaps: %d", len(list.Items))

	makeConfigMaps := viper.GetStringSlice("configmaps_to_make")
	log.Infof("creating %s", makeConfigMaps)

	for _, cm := range makeConfigMaps {
		k8sattacks.CreateCM(k8sClient, "default", cm)
	}
	
	list, err = k8sClient.CoreV1().ConfigMaps("default").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Errorf("%s", err)
	}

	for _, cm := range makeConfigMaps {
		go func(cm string) {k8sattacks.ImAnIdiotAndKeepDDOSingMyself(k8sClient, "default", cm)}(cm)
	}

	log.Infof("existing configmaps: %d", len(list.Items))

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9090", nil)
}