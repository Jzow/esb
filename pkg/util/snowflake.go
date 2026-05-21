package util

import (
	"github.com/bwmarrin/snowflake"
	"time"
)

var node *snowflake.Node

func SetupSnowFlake(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	snowflake.NodeBits = 3 //support 8 machine
	snowflake.StepBits = 9 //support 512data per ms,41timestamp bit+3nodeId bit+9seq bit=53 javascript max int bit
	snowflake.Epoch = st.UnixNano() / 1000000
	node, err = snowflake.NewNode(machineID)
	return
}

func GetSnowFlakeId() int64 {
	return node.Generate().Int64()
}
