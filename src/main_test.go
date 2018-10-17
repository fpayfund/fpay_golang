package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	//zlog.SetTagLevel(zlog.TRACE, "fpay/(*FPAY)")
	os.Exit(m.Run())
}
