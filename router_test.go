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
	_, ifFind, _ = rt.findRouter("/")
	if ifFind != false {
		t.Errorf("Expect find none but get sth\n")
	}
	rt.addRouter("/", nil)
	_, ifFind, _ = rt.findRouter("/")
	if ifFind != true {
		t.Errorf("Expect find / but get nothing\n")
	}
	_, ifFind, _ = rt.findRouter("/admin/curled")
	if ifFind != false {
		t.Errorf("Expect find none but get sth\n")
	}
	// 动态路由测试
	rt.addRouter("/admin/:user", nil)
	_, ifFind, params := rt.findRouter("/admin/curled")
	if ifFind != true || params["user"] != "curled" {
		t.Errorf("Expect find /admin/curled but get %s\n", params["user"])
	} else {
		t.Log(params["user"])
	}
	// 动态路由顺序验证
	rt.addRouter("/admin/curled/:lover", nil)
	_, ifFind, params = rt.findRouter("/admin/curled/curled1")
	if ifFind != true || params["lover"] != "curled1" {
		t.Errorf("Expect find curled1 but get %s\n", params["lover"])
	} else {
		t.Log(params["lover"])
	}
	// 通配路由测试
	rt.addRouter("/static/*relapath", nil)
	_, ifFind, params = rt.findRouter("/static/js/bootstrap.js")
	if ifFind != true || params["relapath"] != "js/bootstrap.js" {
		t.Errorf("Expect find /js/bootstrap.js but get %s\n", params["relapath"])
	} else {
		t.Log(params["relapath"])
	}
	// 通配路由顺序验证
	rt.addRouter("/static/js/*relapath", nil)
	_, ifFind, params = rt.findRouter("/static/js/bootstrap.js")
	if ifFind != true || params["relapath"] != "bootstrap.js" {
		t.Errorf("Expect find bootstrap.js but get %s\n", params["relapath"])
	} else {
		t.Log(params["relapath"])
	}
	// 重复添加路由测试
	err := rt.addRouter("/admin/:sth", nil)
	if err == nil {
		t.Errorf("Expect get duplicate router but get nothing\n")
	} else {
		t.Log(err)
	}
	// 重复添加路由测试
	err = rt.addRouter("/static/*file", nil)
	if err == nil {
		t.Errorf("Expect get duplicate router but get nothing\n")
	} else {
		t.Log(err)
	}
	// 路由格式错误测试
	err = rt.addRouter("/admin/curl*ed", nil)
	if err == nil {
		t.Errorf("Expect get format error but get nothing\n")
	} else {
		t.Log(err)
	}
}
