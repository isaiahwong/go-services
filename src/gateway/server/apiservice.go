package server

import (
	"encoding/json"
	"fmt"
	"strings"
)

type (
	// APIAuthentication defines a set of rules
	// to which gateway should adhere
	APIAuthentication struct {
		required bool
		// exclude excludes routes from authentication if authentication is enabled
		exclude []string
	}

	// APIService encapsulates k8s admissionregistration.k8s.io/v1
	// into a simpler format
	APIService struct {
		Path           string            // Path path to route to service '/api/v1/payment
		DNSPath        string            // DNSPath DNS to service
		Ports          map[string]string // Ports exposed ports of a service
		Authentication APIAuthentication // Authentication rules
		ServiceName    string
		Namespace      string
		APIversion     string
		ResourceType   string
	}

	labels struct {
		ResourceType string `json:"resourceType"`
	}

	annotations struct {
		Config string `json:"config"`
	}

	metadata struct {
		Namespace   string      `json:"namespace"`
		Labels      labels      `json:"labels"`
		Name        string      `json:"name"`
		Annotations annotations `json:"annotations"`
	}

	ports struct {
	}

	spec struct {
		Ports ports `json:"ports"`
	}

	k8sObject struct {
		Metadata metadata `json:"metadata"`
		// Spec     spec     `json:"spec"`
	}

	kind struct {
		Kind string `json:"kind"`
	}

	admissionRequest struct {
		Operation string    `json:"operation"`
		Kind      kind      `json:"kind"`
		Name      string    `json:"name"`
		Namespace string    `json:"Namespace"`
		Object    k8sObject `json:"object"`
	}

	admission struct {
		Request admissionRequest `json:"request"`
	}
)

type objectKind string

//  object kinds
const (
	TypeService    objectKind = "service"
	TypeDeployment            = "deployment"
)

type operation string

// K8S operation type
const (
	Create operation = "CREATE"
	Delete           = "DELETE"
)

func constructAPI(data []byte, s *APIService) error {
	var a *admission = &admission{}

	err := json.Unmarshal(data, a)
	if err != nil {
		return err
	}

	// end function if K8S object is not type service
	if strings.ToLower(a.Request.Kind.Kind) != string(TypeService) {
		return nil
	}

	// TODO filter api-service labels

	switch a.Request.Operation {
	case string(Create):
		createAPI(a)
	}

	return nil
}

func createAPI(a *admission) {
	fmt.Println(a)
}
