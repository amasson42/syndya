package main

import (
	"flag"
	"fmt"
	"syndya/internal/MatchFinder"
	"syndya/pkg/Models"
)

type Params struct {
	filepath string
	mode     string
}

func parseParams() *Params {

	filepath := flag.String("file", "assets/matchup.lua", "path to the lua script")
	mode := flag.String("mode", "matchfinder", "algorithm to test [ matchfinder | matchdeployer ]")

	flag.Parse()

	return &Params{
		filepath: *filepath,
		mode:     *mode,
	}
}

func main() {

	params := parseParams()

	if params.mode == "matchfinder" {
		testMatchfinder(params.filepath)
	}
	if params.mode == "matchdeployer" {
		testMatchdeployer(params.filepath)
	}

}

func testMatchfinder(filepath string) {
	players := Models.NewSearchingPlayersList()

	id := players.CreateSearchingPlayer()
	players.UpdateSearchingPlayerMetadata(id, "rating", "1000")

	mf, err := MatchFinder.NewMatchFinder(players, false, filepath, false)

	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	mf.RunOnce()
}

func testMatchdeployer(filepath string) {

}