/*
Copyright 2018 BlackRock, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package webhook

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/argoproj/argo-events/gateways"
	gwcommon "github.com/argoproj/argo-events/gateways/common"
)

// ValidateEventSource validates webhook event source
func (ese *WebhookEventSourceExecutor) ValidateEventSource(ctx context.Context, es *gateways.EventSource) (*gateways.ValidEventSource, error) {
	w, err := parseEventSource(es.Data)
	if err != nil {
		return &gateways.ValidEventSource{
			Reason:  fmt.Sprintf("%s. err: %s", gateways.ErrEventSourceParseFailed, err.Error()),
			IsValid: false,
		}, nil
	}

	if err = validateWebhook(w); err != nil {
		return &gateways.ValidEventSource{
			Reason:  err.Error(),
			IsValid: false,
		}, nil
	}
	return &gateways.ValidEventSource{
		IsValid: true,
		Reason:  "valid",
	}, nil
}

func validateWebhook(w *gwcommon.Webhook) error {
	if w == nil {
		return fmt.Errorf("%+v, configuration must be non empty", gateways.ErrInvalidEventSource)
	}

	switch w.Method {
	case http.MethodHead, http.MethodPut, http.MethodConnect, http.MethodDelete, http.MethodGet, http.MethodOptions, http.MethodPatch, http.MethodPost, http.MethodTrace:
	default:
		return fmt.Errorf("unknown HTTP method %s", w.Method)
	}

	if w.Endpoint == "" {
		return fmt.Errorf("endpoint can't be empty")
	}
	if w.Port == "" {
		return fmt.Errorf("port can't be empty")
	}

	if !strings.HasPrefix(w.Endpoint, "/") {
		return fmt.Errorf("endpoint must start with '/'")
	}

	if w.Port != "" {
		_, err := strconv.Atoi(w.Port)
		if err != nil {
			return fmt.Errorf("failed to parse server port %s. err: %+v", w.Port, err)
		}
	}
	return nil
}
