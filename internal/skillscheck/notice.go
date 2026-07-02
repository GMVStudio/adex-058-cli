package skillscheck

import (
	"fmt"
	"sync/atomic"
)

type StaleNotice struct {
	Current string `json:"current"`
	Target  string `json:"target"`
}

func (s *StaleNotice) Message() string {
	return fmt.Sprintf(
		"adex skills %s out of sync with binary %s, run: adex update",
		s.Current, s.Target,
	)
}

var pending atomic.Pointer[StaleNotice]

func SetPending(n *StaleNotice) { pending.Store(n) }

func GetPending() *StaleNotice { return pending.Load() }
