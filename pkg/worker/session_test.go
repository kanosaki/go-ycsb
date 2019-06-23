package worker

import (
	"testing"

	"github.com/magiconair/properties"
)

func TestCreateSession(t *testing.T) {
	sess, err := CreateJob("a", properties.LoadMap(map[string]string{}), "workload_x", "foobar")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

}
