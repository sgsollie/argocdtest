package main

import (

    "context"

    "flag"

    "fmt"

    "log"

    "time"

    "path/filepath"
    batchv1 "k8s.io/api/batch/v1"

    apiv1 "k8s.io/api/core/v1"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

    "k8s.io/client-go/kubernetes"

    "k8s.io/client-go/tools/clientcmd"

    "k8s.io/client-go/util/homedir"

)


var clientset *kubernetes.Clientset
const jobName = "test-job-clientgo"
const cronJobName = "test-cronjob-clientgo"

func init() {
    var kubeconfig *string
    if home := homedir.HomeDir(); home != "" {

        kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")

    } else {

        kubeconfig = flag.String("kubeconfig","  ", "absolute path to the kubeconfig file")

    }

    flag.Parse()

    config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
    if err != nil {

        log.Fatal(err)

    }

    clientset, err = kubernetes.NewForConfig(config)
    if err != nil {

        log.Fatal(err)

    }
}


func createJob() {
    fmt.Println("creating job", jobName)
    jobsClient := clientset.BatchV1().Jobs(apiv1.NamespaceDefault)
    job := &batchv1.Job{ObjectMeta: metav1.ObjectMeta{
        Name: jobName,
    },
        Spec: batchv1.JobSpec{
            Template: apiv1.PodTemplateSpec{
                Spec: apiv1.PodSpec{
                    Containers: []apiv1.Container{
                        {
                            Name: "pi",
                            Image: "perl:5.34.0",
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
