package tests

import (
	"testing"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
)

func configLogger(log *logger.Logger) *logger.Logger {
	// Use the group filters to show specific algorithm logs while debugging.
	//log = log.Show(`analyze`)
	//log = log.Show(`converter`)
	log = log.Show(`files`)
	//log = log.Show(`inheritance`)
	//log = log.Show(`instantiator`)
	//log = log.Show(`generateInterfaces`)
	log = log.Show(`packages`)
	//log = log.Show(`usages`)
	return log
}

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
func Test_T0012(t *testing.T) { newTest(t, `test0012`).abstract().full() }

func Test_T0013(t *testing.T) { newTest(t, `test0013`).abstract().full() }
func Test_T0014(t *testing.T) { newTest(t, `test0014`).abstract().dump().full() }
func Test_T0015(t *testing.T) { newTest(t, `test0015`).abstract().full() }

func Test_T0016(t *testing.T) { newTest(t, `test0016`).abstract().full() }
func Test_T0017(t *testing.T) { newTest(t, `test0017`).abstract().full() }
func Test_T0018(t *testing.T) { newTest(t, `test0018`).abstract().full() }
