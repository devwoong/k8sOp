package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8sOp/channel"
	"log"
	"os"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/yaml"
)

type k8sService struct {
	ctx       context.Context
	clientset *kubernetes.Clientset
}

var K8sService k8sService = k8sService{nil, nil}

func (c *k8sService) Init() {
	c.ctx = context.Background()

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	c.clientset = clientset

	// Demo to make sure it works
	pods, err := c.clientset.CoreV1().Pods("").List(c.ctx, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	log.Printf("There are %d pods in the cluster", len(pods.Items))

}

func (c *k8sService) Run() {
	var now = time.Now().UTC()
	now = now.In(time.FixedZone("KST", 9*60*60))
	hour := now.Hour()
	if hour >= 8 && hour < 21 {
		c.createDeployments("test-server", "./testServer.yaml")
	} else if hour >= 21 {
		deploymentClient := c.clientset.AppsV1().Deployments("default")
		deployments, err := deploymentClient.List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		for _, deployment := range deployments.Items {
			if deployment.Name == "scheduler" || deployment.Name == "proxy-server" {
				continue
			}
			deletePolicy := metav1.DeletePropagationForeground
			if err := deploymentClient.Delete(context.TODO(), deployment.Name, metav1.DeleteOptions{
				PropagationPolicy: &deletePolicy,
			}); err != nil {
				panic(err)
			}
			fmt.Println("Deleted deployment: " + deployment.Name)
		}
	}

}

func (c *k8sService) createDeployments(name string, fileName string) {
	deploymentClient := c.clientset.AppsV1().Deployments("default")

	deployments, err := deploymentClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	find := false
	for _, deployment := range deployments.Items {
		if deployment.Name == name {
			find = true
		}
	}

	if !find {
		file, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}

		b, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}

		deployemntJson, err := yaml.YAMLToJSON(b)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
		dep := &appsv1.Deployment{}
		// dep.Unmarshal(deployemntJson)
		json.Unmarshal(deployemntJson, &dep)
		fmt.Println(string(deployemntJson))

		result, err := deploymentClient.Create(context.TODO(), dep, metav1.CreateOptions{})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	}
}

func (c *k8sService) CheckAndBootDeployment() {
	for {
		serviceName := <-channel.CommonChannel.RequestChannel
		service, err := c.clientset.CoreV1().Services("default").Get(context.TODO(), serviceName, metav1.GetOptions{})
		if err != nil {
			panic(err)
		}
		appName := service.Spec.Selector["app"]
		c.createDeployments(appName, "./testServer.yaml")
	}

}
