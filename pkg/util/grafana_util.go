// Copyright (c) 2020 Red Hat, Inc.

package util

import (
	"encoding/hex"
	"hash/fnv"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"k8s.io/klog"
)

const (
	defaultAdmin = "WHAT_YOU_ARE_DOING_IS_VOIDING_SUPPORT_0000000000000000000000000000000000000000000000000000000000000000"
)

// GenerateUID generates UID for customized dashboard
func GenerateUID(namespace string, name string) string {
	uid := namespace + "-" + name
	if len(uid) > 40 {
		hasher := fnv.New128a()
		hasher.Write([]byte(uid))
		uid = hex.EncodeToString(hasher.Sum(nil))
	}
	return uid
}

// GetHTTPClient returns http client
func getHTTPClient() *http.Client {
	transport := &http.Transport{}
	client := &http.Client{Transport: transport}
	return client
}

// SetRequest ...
func SetRequest(method string, url string, body io.Reader) ([]byte, int) {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-User", defaultAdmin)

	resp, err := getHTTPClient().Do(req)
	for {
		if err == nil {
			break
		}
		klog.Error("failed to send HTTP request. Retry in 5 seconds", "error", err)
		time.Sleep(5)
		resp, err = getHTTPClient().Do(req)
	}
	/* 	if err != nil {
		klog.Error("failed to send HTTP request", "error", err)
		return nil, 500
	} */

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		klog.Info("failed to parse response body", "error", err)
	} //else {
	// 	klog.Info("Succeed to parse response body", "Response body", string(respBody))
	// }
	return respBody, resp.StatusCode
}
