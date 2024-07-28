package tests

import "testing"

func Test_T0001(t *testing.T) { newTest(t, `test0001`).abstract().full() }
func Test_T0002(t *testing.T) { newTest(t, `test0002`).abstract().full() }
func Test_T0003(t *testing.T) { newTest(t, `test0003`).abstract().full() }

func Test_T0004(t *testing.T) { newTest(t, `test0004`).abstract().full() }
func Test_T0005(t *testing.T) { newTest(t, `test0005`).abstract(`cats.go`).full() }
func Test_T0006(t *testing.T) { newTest(t, `test0006`).abstract(`cats.go`).save().partial() }

func Test_T0007(t *testing.T) { newTest(t, `test0007`).abstract().full() }
func Test_T0008(t *testing.T) { newTest(t, `test0008`).abstract().full() }
func Test_T0009(t *testing.T) { newTest(t, `test0009`).abstract().full() }

func Test_T0010(t *testing.T) { newTest(t, `test0010`).abstract().full() }
func Test_T0011(t *testing.T) { newTest(t, `test0011`).abstract().full() }
