package gourd

import (
	"testing"
)

func Test_dividePath(t *testing.T) {
	testPath1 := "/user/:username/info"
	testSubPaths1 := []string{"/user", "/:username", "/info"}
	subPaths1 := dividePath(testPath1)
	if len(testSubPaths1) != len(subPaths1) {
		t.Errorf("Size Mismatched!\n")
	}
	for key, subPath := range subPaths1 {
		if subPath != testSubPaths1[key] {
			t.Errorf("Expect %s but get %s\n", testSubPaths1[key], subPath)
		}
	}
}

func Test_Router(t *testing.T) {
	ifFind := false
	rt := newRouterManager()
	_,ifFind,_ = rt.findRouter("/")
	if ifFind != false {
		t.Errorf("Expect find none but get sth\n")
	}
	rt.addRouter("/", nil)
	_,ifFind,_ = rt.findRouter("/")
	if ifFind != true {
		t.Errorf("Expect find / but get nothing\n")
	}
	_,ifFind,_ = rt.findRouter("/admin/curled")
	if ifFind != false {
		t.Errorf("Expect find none but get sth\n")
	}
	// 动态路由测试
	rt.addRouter("/admin/:user", nil)
	_,ifFind,_ = rt.findRouter("/admin/curled")
	if ifFind != true {
		t.Errorf("Expect find /admin.curled but get nothing\n")
	}
	// 重复添加路由测试
	err := rt.addRouter("/admin/:sth", nil)
	if err == nil {
		t.Errorf("Expect find /admin.curled but get nothing\n")
	} else {
		t.Log(err)
	}
}
