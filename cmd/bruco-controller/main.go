package main

import (
	"flag"
	"time"

	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	brucocontroller "github.com/ferama/bruco/pkg/kube/controller"
	clientset "github.com/ferama/bruco/pkg/kube/generated/clientset/versioned"
	informers "github.com/ferama/bruco/pkg/kube/generated/informers/externalversions"
	"github.com/ferama/bruco/pkg/kube/signals"
)

var (
	masterURL  string
	kubeconfig string
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	brucoClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building example clientset: %s", err.Error())
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	brucoInformerFactory := informers.NewSharedInformerFactory(brucoClient, time.Second*30)
	brucoProjectInformerFactory := informers.NewSharedInformerFactory(brucoClient, time.Second*30)

	controller := brucocontroller.NewBrucoController(kubeClient, brucoClient,
		kubeInformerFactory.Apps().V1().Deployments(),
		kubeInformerFactory.Core().V1().Services(),
		kubeInformerFactory.Core().V1().ConfigMaps(),
		brucoInformerFactory.Bruco().V1alpha1().Brucos())

	projectController := brucocontroller.NewBrucoProjectController(
		kubeClient,
		brucoClient,
		brucoInformerFactory.Bruco().V1alpha1().Brucos(),
		brucoProjectInformerFactory.Bruco().V1alpha1().BrucoProjects(),
	)
	// notice that there is no need to run Start methods in a separate goroutine. (i.e. go kubeInformerFactory.Start(stopCh)
	// Start method is non-blocking and runs all registered informers in a dedicated goroutine.
	kubeInformerFactory.Start(stopCh)
	brucoInformerFactory.Start(stopCh)
	brucoProjectInformerFactory.Start(stopCh)

	go func() {
		if err = projectController.Run(2, stopCh); err != nil {
			klog.Fatalf("Error running project controller: %s", err.Error())
		}
	}()

	if err = controller.Run(2, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
