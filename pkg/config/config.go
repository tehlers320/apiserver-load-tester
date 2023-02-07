package config

import (
	"bytes"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes"
)



var (
	defaultConfig = []byte(`
logging.level: "info"
delay_between_attacks: "20ms"
configmaps_to_make:
- uno
- dos
- tres
- cuatro
- cinco
- seis
- siete
- ocho
- nueve
- diez
`)

)

func InitConfig() error {
	var err error
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("ALT")
	var cfg []byte = nil
	cfg = defaultConfig
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	err = viper.ReadConfig(bytes.NewBuffer(cfg))
	if err != nil {
		return err
	}
	return err
}

type K8s struct {
	Local bool
}

func (k8s *K8s) CreateClient() (*kubernetes.Clientset, error) {
	config, err := k8s.buildConfig()

	if err != nil {
		return nil, errors.Wrapf(err, "error setting up cluster config")
	}

	return kubernetes.NewForConfig(config)
}

func (k8s *K8s) buildConfig() (*rest.Config, error) {
	if k8s.Local {
		log.Debug("Using local kubeconfig.")
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	log.Debug("Using in cluster kubeconfig.")
	return rest.InClusterConfig()
}
