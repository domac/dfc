package app

import (
	"fmt"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	filepath := "../conf/base.conf"
	cfg, err := LoadConfig(filepath)

	if err != nil {
		t.Fatal(err)
	}

	if cfg == nil {
		t.Fatal("load error")
	}

	if cfg.Name != "DFC_NODE" {
		t.Fatal("load error")
	}

	fmt.Printf("%v", cfg)

	pp, _ := cfg.GetParentPeerNodes()
	fmt.Printf("parents count : %d\n", len(pp))
}
