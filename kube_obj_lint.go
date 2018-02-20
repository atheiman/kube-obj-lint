package main

import (
    "github.com/ghodss/yaml"
    "io/ioutil"
    "log"
    "os"
)

type Pod struct {
    APIVersion string `json:"apiVersion"`
    Kind       string `json:"kind"`
    Metadata   struct {
        Name string `json:"name"`
    } `json:"metadata"`
    Spec struct {
        Volumes []struct {
            Name string `json:"name"`
            Nfs  struct {
                Server string `json:"server"`
                Path   string `json:"path"`
            } `json:"nfs"`
        } `json:"volumes"`
        Containers []struct {
            Image string `json:"image"`
            Name  string `json:"name"`
            Ports []struct {
                ContainerPort int    `json:"containerPort"`
                Name          string `json:"name"`
                Protocol      string `json:"protocol"`
            } `json:"ports"`
            Resources struct {
                Requests struct {
                    CPU    string `json:"cpu"`
                    Memory string `json:"memory"`
                } `json:"requests"`
                Limits struct {
                    CPU    string `json:"cpu"`
                    Memory string `json:"memory"`
                } `json:"limits"`
            } `json:"resources"`
            VolumeMounts []struct {
                MountPath string `json:"mountPath"`
                Name      string `json:"name"`
            } `json:"volumeMounts"`
            LivenessProbe struct {
                HTTPGet struct {
                    Path string `json:"path"`
                    Port int    `json:"port"`
                } `json:"httpGet"`
                InitialDelaySeconds int `json:"initialDelaySeconds"`
                TimeoutSeconds      int `json:"timeoutSeconds"`
                PeriodSeconds       int `json:"periodSeconds"`
                FailureThreshold    int `json:"failureThreshold"`
            } `json:"livenessProbe"`
            ReadinessProbe struct {
                HTTPGet struct {
                    Path string `json:"path"`
                    Port int    `json:"port"`
                } `json:"httpGet"`
                InitialDelaySeconds int `json:"initialDelaySeconds"`
                TimeoutSeconds      int `json:"timeoutSeconds"`
                PeriodSeconds       int `json:"periodSeconds"`
                FailureThreshold    int `json:"failureThreshold"`
            } `json:"readinessProbe"`
        } `json:"containers"`
    } `json:"spec"`
}

func verifyVolNames(pod Pod) {
    // get all the declared volumes
    var volNames []string 
    for _, vol := range pod.Spec.Volumes {
        volNames = append(volNames, vol.Name)
    }

    // get all the container volume mount name references
    var mountNameRefs []string
    for _, c := range pod.Spec.Containers {
        for _, vm := range c.VolumeMounts {
            mountNameRefs = append(mountNameRefs, vm.Name)
            // while we're already iterating, confirm the references are valid
            if !stringInSlice(vm.Name, volNames) {
                log.Fatalf("attempting to mount undeclared volume '%s' in container '%s'\n",
                           vm.Name,
                           c.Name)
            }
        }
    }

    // verify all declared volumes are used
    for _, volName := range volNames {
        if !stringInSlice(volName, mountNameRefs) {
            log.Printf("volume '%s' declared but not mounted in any container\n",
                       volName)
        }
    }
}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func main() {
    if len(os.Args) < 2 {
        log.Println("Missing parameter, provide file name!")
        log.Fatalf("Usage: %s OBJECT_YAML_FILE\n", os.Args[0])
    }
    // read pod definition from cli Args
    content, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }

    var pod Pod
    err = yaml.Unmarshal(content, &pod)
    if err != nil {
        log.Printf("err: %v\n", err)
        return
    }
    log.Printf("Checking pod definition '%s'\n", pod.Metadata.Name)

    verifyVolNames(pod)

    log.Printf("OK!\n")
}
