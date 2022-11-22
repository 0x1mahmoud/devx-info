// Code generated by cue get go. DO NOT EDIT.

//cue:generate cue get go k8s.io/pod-security-admission/api

package api

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Attributes exposes the admission request parameters consumed by the PodSecurity admission controller.
#Attributes: _

// AttributesRecord is a simple struct implementing the Attributes interface.
#AttributesRecord: {
	Name:        string
	Namespace:   string
	Kind:        schema.#GroupVersionKind
	Resource:    schema.#GroupVersionResource
	Subresource: string
	Operation:   admissionv1.#Operation
	Object:      runtime.#Object
	OldObject:   runtime.#Object
	Username:    string
}
