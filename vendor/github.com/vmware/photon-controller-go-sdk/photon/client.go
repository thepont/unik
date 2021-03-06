// Copyright (c) 2016 VMware, Inc. All Rights Reserved.
//
// This product is licensed to you under the Apache License, Version 2.0 (the "License").
// You may not use this product except in compliance with the License.
//
// This product may include a number of subcomponents with separate copyright notices and
// license terms. Your use of these subcomponents is subject to the terms and conditions
// of the subcomponent's license, as noted in the LICENSE file.

package photon

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"strings"
	"time"
)

// Represents stateless context needed to call photon APIs.
type Client struct {
	options           ClientOptions
	httpClient        *http.Client
	Endpoint          string
	AuthEndpoint      string
	Status            *StatusAPI
	Tenants           *TenantsAPI
	Tasks             *TasksAPI
	Projects          *ProjectsAPI
	Flavors           *FlavorsAPI
	Images            *ImagesAPI
	Disks             *DisksAPI
	VMs               *VmAPI
	Hosts             *HostsAPI
	Deployments       *DeploymentsAPI
	ResourceTickets   *ResourceTicketsAPI
	Networks          *NetworksAPI
	Clusters          *ClustersAPI
	Auth              *AuthAPI
	AvailabilityZones *AvailabilityZonesAPI
}

// Represents Tokens
type TokenOptions struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IdToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
}

// Options for Client
type ClientOptions struct {
	// When using the Tasks.Wait APIs, defines the duration of how long
	// the SDK should continue to poll the server. Default is 30 minutes.
	// TasksAPI.WaitTimeout() can be used to specify timeout on
	// individual calls.
	TaskPollTimeout time.Duration

	// Whether or not to ignore any TLS errors when talking to photon,
	// false by default.
	IgnoreCertificate bool

	// List of root CA's to use for server validation
	// nil by default.
	RootCAs *x509.CertPool

	// For tasks APIs, defines the delay between each polling attempt.
	// Default is 100 milliseconds.
	TaskPollDelay time.Duration

	// For tasks APIs, defines the number of retries to make in the event
	// of an error. Default is 3.
	TaskRetryCount int

	// Tokens for user authentication. Default is empty.
	TokenOptions *TokenOptions
}

// Creates a new photon client with specified options. If options
// is nil, default options will be used.
func NewClient(endpoint string, authEndpoint string, options *ClientOptions) (c *Client) {
	defaultOptions := &ClientOptions{
		TaskPollTimeout:   30 * time.Minute,
		TaskPollDelay:     100 * time.Millisecond,
		TaskRetryCount:    3,
		TokenOptions:      &TokenOptions{},
		IgnoreCertificate: false,
		RootCAs:           nil,
	}

	if options != nil {
		if options.TaskPollTimeout != 0 {
			defaultOptions.TaskPollTimeout = options.TaskPollTimeout
		}
		if options.TaskPollDelay != 0 {
			defaultOptions.TaskPollDelay = options.TaskPollDelay
		}
		if options.TaskRetryCount != 0 {
			defaultOptions.TaskRetryCount = options.TaskRetryCount
		}
		if options.TokenOptions != nil {
			defaultOptions.TokenOptions = options.TokenOptions
		}
		if options.RootCAs != nil {
			defaultOptions.RootCAs = options.RootCAs
		}
		defaultOptions.IgnoreCertificate = options.IgnoreCertificate
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: defaultOptions.IgnoreCertificate,
			RootCAs:            defaultOptions.RootCAs},
	}

	endpoint = strings.TrimRight(endpoint, "/")
	authEndpoint = strings.TrimRight(authEndpoint, "/")

	c = &Client{Endpoint: endpoint, AuthEndpoint: authEndpoint, httpClient: &http.Client{Transport: tr}}
	// Ensure a copy of options is made, rather than using a pointer
	// which may change out from underneath if misused by the caller.
	c.options = *defaultOptions
	c.Status = &StatusAPI{c}
	c.Tenants = &TenantsAPI{c}
	c.Tasks = &TasksAPI{c}
	c.Projects = &ProjectsAPI{c}
	c.Flavors = &FlavorsAPI{c}
	c.Images = &ImagesAPI{c}
	c.Disks = &DisksAPI{c}
	c.VMs = &VmAPI{c}
	c.Hosts = &HostsAPI{c}
	c.Deployments = &DeploymentsAPI{c}
	c.ResourceTickets = &ResourceTicketsAPI{c}
	c.Networks = &NetworksAPI{c}
	c.Clusters = &ClustersAPI{c}
	c.Auth = &AuthAPI{c}
	c.AvailabilityZones = &AvailabilityZonesAPI{c}
	return
}

// Creates a new photon client with specified options and http.Client.
// Useful for functional testing where http calls must be mocked out.
// If options is nil, default options will be used.
func NewTestClient(endpoint string, authEndpoint string, options *ClientOptions, httpClient *http.Client) (c *Client) {
	c = NewClient(endpoint, authEndpoint, options)
	c.httpClient = httpClient
	return
}
