package worker

import "fmt"

type SessionState int

func (s SessionState) String() string {
	if s < 0 || len(StateNames) <= int(s) {
		return fmt.Sprintf("SessionState(%d)", int(s))
	}
	return StateNames[int(s)]
}

func (s *SessionState) UnmarshalJSON(d []byte) error {
	str := string(d)
	for i, name := range StateNames {
		if name == str {
			*s = SessionState(i)
			return nil
		}
	}
	return fmt.Errorf("%s is invalid as SessionState", str)
}

func (s SessionState) MarshalJSON() ([]byte, error) {
	if s < 0 || len(StateNames) <= int(s) {
		return nil, fmt.Errorf("invalid state value: %d", int(s))
	}
	return []byte(StateNames[int(s)]), nil
}

const (
	StateInitialized SessionState = iota
	StateStarted
	StateFinished
	StateError
)

var (
	StateNames = []string{
		"initialized",
		"started",
		"error",
	}
)
