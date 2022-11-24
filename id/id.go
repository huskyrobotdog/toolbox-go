package id

import (
	"github.com/bwmarrin/snowflake"
	"github.com/huskyrobotdog/toolbox-go/inner"
	"github.com/huskyrobotdog/toolbox-go/log"
)

var Instance *snowflake.Node

func Initialization(node int64) {
	sf, err := snowflake.NewNode(node)
	if err != nil {
		log.Instance.Fatal(err.Error())
	}
	Instance = sf
	inner.Debug("id initialization complete")
}
