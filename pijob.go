package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var clientset *kubernetes.Clientset

const jobName = "test-job-clientgo"
const cronJobName = "test-cronjob-clientgo"
const hardCodedNamespace = "ollie"

func init() {

	kubeconfig := flag.String("kubeconfig", "/home/ollie/.kube/config", "location to your kubeconfig file")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		// handle error
		fmt.Printf("erorr %s building config from flags\n", err.Error())
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("error %s, getting inclusterconfig", err.Error())
		}
	}

	flag.Parse()

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
}

func createJob() {
	fmt.Println("creating job", jobName)

	// This is so ArgoCD can track jobs that are spawned from inside this go program:
	// https://argo-cd.readthedocs.io/en/stable/user-guide/resource_tracking/
	LabelsMap := map[string]string{
		"app.kubernetes.io/instance": os.Getenv("APPNAME"),
	}

	jobsClient := clientset.BatchV1().Jobs(hardCodedNamespace)
	job := &batchv1.Job{ObjectMeta: metav1.ObjectMeta{
		Name:        jobName,
		Annotations: LabelsMap,
	},
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:    "pi",
							Image:   "perl:5.34.0",
							Command: []string{"perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"},
						},
					},
					RestartPolicy: apiv1.RestartPolicyNever,
				},
			},
		},
	}
	_, err := jobsClient.Create(context.Background(), job, metav1.CreateOptions{})
	if err != nil {

		log.Fatal("failed to create job", err)

	}
	fmt.Println("created job successfully")

}

func main() {
	createJob()
	// Sleep for a while to pretend we are a running service
	time.Sleep(300 * time.Second)
}
