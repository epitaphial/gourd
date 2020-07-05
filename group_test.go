package gourd

import (
	"testing"
)

func Test_Group(t *testing.T) {
	engine := Gourd()
	if engine.prefix != "" {
		t.Errorf("Expect %s but get %s\n", "", engine.prefix)
	}
	grp, _ := engine.Group("/v1")
	grp2, _ := grp.Group("/v3")
	if grp2.prefix != "/v1/v3" {
		t.Errorf("Expect prefix %s but get %s\n", "/v1/v3", grp2.prefix)
	}
}
