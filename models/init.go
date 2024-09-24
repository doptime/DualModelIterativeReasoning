package models

import cmap "github.com/orcaman/concurrent-map/v2"

var Models = cmap.New[*Model]()

func init() {
	Models.Set(SLM1.ModelName, SLM1)
	Models.Set(SLM2.ModelName, SLM2)
}
