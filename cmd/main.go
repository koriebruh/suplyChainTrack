package main

import (
	_ "github.com/koriebruh/suplyChainTrack/docs"
	"github.com/koriebruh/suplyChainTrack/route"
)

// @title Supply Chain API
// @version 1.0
// @description API supply chain tracking system
// @host localhost:8080
// @BasePath /api/v1
func main() {
	route.RunApplicationContext()
}
