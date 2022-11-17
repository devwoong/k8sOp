package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8sOp/channel"
	"log"
	"os"
	"strconv"
	"strings"
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

const targetDir string = "/home/server/apply"

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
	minute := now.Minute()

	c.createDeploymentsAll(hour, minute)
	c.deleteDeployments(hour, minute)

}

func (c *k8sService) deleteDeployments(hour int, minute int) {
	deploymentClient := c.clientset.AppsV1().Deployments("default")
	deployments, err := deploymentClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, deployment := range deployments.Items {
		if enable, exist := deployment.Annotations["scheduler.enable"]; !exist || enable != "true" {
			continue
		}
		if time, exist := deployment.Annotations["scheduler.shutdown"]; exist {
			times := strings.Split(time, ":")
			parseHour, err := strconv.ParseInt(times[0], 10, 64)
			if err != nil {
				continue
			}
			parseMinute, _ := strconv.ParseInt(times[1], 10, 64)
			if err != nil {
				continue
			}
			if parseHour == int64(hour) && parseMinute == int64(minute) {
				deletePolicy := metav1.DeletePropagationForeground
				if err := deploymentClient.Delete(context.TODO(), deployment.Name, metav1.DeleteOptions{
					PropagationPolicy: &deletePolicy,
				}); err != nil {
					continue
				}
				deployment.ResourceVersion = ""
				deployment.UID = ""
				c.deleteFileByDeploymentName(deployment.Name)
				c.createDeploymentFile(&deployment)
				fmt.Println("Deleted deployment: " + deployment.Name)
			}
		}
	}
}

func (c *k8sService) createDeploymentsAll(hour int, minute int) {
	deploymentClient := c.clientset.AppsV1().Deployments("default")
	files, err := ioutil.ReadDir(targetDir)
	if err != nil {
		panic(err)
	}
	for _, fileInfo := range files {
		file, err := os.Open(targetDir + "/" + fileInfo.Name())
		if err != nil {
			panic(err)
		}
		defer file.Close()

		b, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}

		deployemntJson, err := yaml.YAMLToJSON(b)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			continue
		}
		dep := &appsv1.Deployment{}
		json.Unmarshal(deployemntJson, &dep)
		if enable, exist := dep.Annotations["scheduler.enable"]; exist && enable == "true" {
			if time, exist := dep.Annotations["scheduler.startup"]; exist {
				existDeploy, _ := deploymentClient.Get(context.TODO(), dep.Name, metav1.GetOptions{})
				if existDeploy == nil {
					continue
				}
				times := strings.Split(time, ":")
				parseHour, err := strconv.ParseInt(times[0], 10, 64)
				if err != nil {
					continue
				}
				parseMinute, _ := strconv.ParseInt(times[1], 10, 64)
				if err != nil {
					continue
				}
				if parseHour == int64(hour) && parseMinute == int64(minute) {
					result, err := deploymentClient.Create(context.TODO(), dep, metav1.CreateOptions{})
					if err != nil {
						fmt.Printf("create %v", err)
					}
					fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
				}
			}
		}
	}

}

func (c *k8sService) createDeployments(name string) {
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
		dep := c.findDeploymentFromFileByDeploymentName(name)
		if dep == nil {
			return
		}
		result, err := deploymentClient.Create(context.TODO(), dep, metav1.CreateOptions{})
		if err != nil {
			fmt.Printf("create %v", err)
		}
		fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	}
}
func (c *k8sService) createDeploymentFile(deployment *appsv1.Deployment) {
	deployemntJson, err := json.Marshal(deployment)
	if err != nil {
		return
	}
	b, _ := yaml.JSONToYAML(deployemntJson)
	ioutil.WriteFile(targetDir+"/"+deployment.Name+".yaml", b, 0777)
}
func (c *k8sService) deleteFileByDeploymentName(appName string) {
	files, err := ioutil.ReadDir(targetDir)
	if err != nil {
		panic(err)
	}
	for _, fileInfo := range files {
		file, err := os.Open(targetDir + "/" + fileInfo.Name())
		if err != nil {
			panic(err)
		}
		defer file.Close()

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
		json.Unmarshal(deployemntJson, &dep)
		if dep.Name == appName {
			os.Remove(targetDir + "/" + fileInfo.Name())
			return
		}
	}
}

func (c *k8sService) findDeploymentFromFileByDeploymentName(appName string) *appsv1.Deployment {
	files, err := ioutil.ReadDir(targetDir)
	if err != nil {
		panic(err)
	}
	for _, fileInfo := range files {
		file, err := os.Open(targetDir + "/" + fileInfo.Name())
		if err != nil {
			panic(err)
		}
		defer file.Close()

		b, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}

		deployemntJson, err := yaml.YAMLToJSON(b)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return nil
		}
		dep := &appsv1.Deployment{}
		json.Unmarshal(deployemntJson, &dep)
		if dep.Name == appName {
			if enable, exist := dep.Annotations["scheduler.enable"]; !exist || enable != "true" {
				return nil
			}
			return dep
		}
	}
	return nil
}

func (c *k8sService) CheckAndBootDeployment() {
	for {
		serviceName := <-channel.CommonChannel.RequestChannel
		service, err := c.clientset.CoreV1().Services("default").Get(context.TODO(), serviceName, metav1.GetOptions{})
		if err != nil {
			panic(err)
		}
		appName := service.Spec.Selector["app"]
		c.createDeployments(appName)
	}

}
