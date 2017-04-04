// Copyright © 2017 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause
package secret

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/vmware/virtual-security-module/config"
	"github.com/vmware/virtual-security-module/context"
	"github.com/vmware/virtual-security-module/model"
	"github.com/vmware/virtual-security-module/vds"
	"github.com/vmware/virtual-security-module/vks"
)

var sm *SecretManager

func TestMain(m *testing.M) {
	cfg := config.GenerateTestConfig()

	ds, err := vds.GetDataStoreFromConfig(cfg)
	if err != nil {
		fmt.Printf("Failed to get data store from config: %v\n", err)
		os.Exit(1)
	}

	ks, err := vks.GetKeyStoreFromConfig(cfg)
	if err != nil {
		fmt.Printf("Failed to get key store from config: %v\n", err)
		os.Exit(1)
	}

	sm = New()
	if err := sm.Init(context.NewModuleInitContext(cfg, ds, ks)); err != nil {
		fmt.Printf("Failed to initialize secret manager: %v\n", err)
		os.Exit(1)
	}
	defer sm.Close()

	apiTestSetup()
	defer apiTestCleanup()

	os.Exit(m.Run())
}

func TestCreateAndGetSecret(t *testing.T) {
	se := &model.SecretEntry{
		Id:             "id1",
		SecretData:     []byte("secret0"),
		Owner:          "user0",
		ExpirationTime: time.Now().Add(time.Hour),
	}

	id2, err := sm.CreateSecret(se)
	if err != nil {
		t.Fatalf("Failed to create secret: %v", err)
	}
	if len(id2) == 0 {
		t.Fatalf("Failed to create secret: returned id is empty")
	}

	se2, err := sm.GetSecret(id2)
	if err != nil {
		t.Fatalf("Failed to get secret for id %v: %v", id2, err)
	}

	if ok := model.SecretsEqual(se, se2); !ok {
		t.Fatalf("Created and retrieved secrets do not match: %v %v", se, se2)
	}

	if err := sm.DeleteSecret(id2); err != nil {
		t.Fatalf("Failed to delete secret for id %v: %v", id2, err)
	}
}
