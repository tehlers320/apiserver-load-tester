package k8sattacks

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"

	log "github.com/sirupsen/logrus"
	metrics "github.com/tehlers320/apiserver-load-tester/pkg/metrics"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)


var (
	//go:embed testdata/basic-configmap.yaml
	basicConfigMap []byte
	// Yeah yeah, whatever no error
	cm, _ = configMapFromYaml(basicConfigMap) 
)


func ImAnIdiotAndKeepDDOSingMyself(k8sClient *kubernetes.Clientset, namespace string, name string) {
	for {
		imAnIdiotAndKeepDDOSingMyself(k8sClient, namespace, name)
		time.Sleep(viper.GetDuration("delay_between_attacks"))
	}

}

func imAnIdiotAndKeepDDOSingMyself(k8sClient *kubernetes.Clientset, namespace string, name string) {
	var pre time.Time
	var post float64

	pre = time.Now()
	metrics.ResourceCount("configmap", name, "update")
	updatetime := fmt.Sprintf("%d", time.Now().UnixNano())
	cm.SetAnnotations(map[string]string{"lastupdated": updatetime })
	_, err := k8sClient.CoreV1().ConfigMaps(namespace).Update(context.TODO(), cm, v1.UpdateOptions{})
	if err != nil {
		log.Errorf("%s", err)
	}
	post = postTime(pre)
	log.Infof("cm: latency for update on %s %f seconds", name, post)
}


func CreateCM(k8sclient kubernetes.Interface, ns string, name string) {
	var pre time.Time
	var post float64

	pre = time.Now()
	metrics.ResourceCount("configmap", name, "create")
	_, err := createConfigMapFromBytes(k8sclient, ns, name)
	// We just want to log but NOT return
	if err != nil {
		metrics.ResourceErrors("configmap", name, "create")
		log.Errorf("cm: %s", err)
	}

	post = postTime(pre)
	log.Infof("cm: latency for create on %s %f seconds", name, post)

}


func createConfigMapFromBytes(client kubernetes.Interface, namespace string, name string) (*corev1.ConfigMap, error) {
	cm.Name = name
	ccm, err := client.CoreV1().ConfigMaps(namespace).Create(context.TODO(), cm, v1.CreateOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "create configmap from file, error creating configmap in namespace '%s'", namespace)
	}

	return ccm, nil
}


func postTime(pre time.Time) float64 {
	return time.Since(pre).Seconds()
}

func configMapFromYaml(y []byte) (*corev1.ConfigMap, error) {
	j, err := yaml.ToJSON(y)
	if err != nil {
		return nil, errors.Wrap(err, "configMapFromYaml- error converting yaml to json")
	}
	cm := corev1.ConfigMap{}
	err = json.Unmarshal(j, &cm)
	if err != nil {
		return nil, errors.Wrap(err, "configMapFromYaml- unmarshaling error")
	}

	return &cm, nil
}