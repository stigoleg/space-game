package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"

	"stellar-siege/game"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joho/godotenv"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile = flag.String("memprofile", "", "write memory profile to file")
	pprofAddr  = flag.String("pprof", "", "enable pprof server on address (e.g., :6060)")
)

func main() {
	flag.Parse()

	// Start CPU profiling if requested
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
		log.Printf("CPU profiling enabled, writing to %s", *cpuprofile)
	}

	// Start pprof HTTP server if requested
	if *pprofAddr != "" {
		go func() {
			log.Printf("Starting pprof server on %s", *pprofAddr)
			log.Println(http.ListenAndServe(*pprofAddr, nil))
		}()
	}

	// Detect if running from .app bundle and adjust working directory
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)

		// If running from .app/Contents/MacOS, change to Resources directory
		if strings.HasSuffix(exeDir, "Contents/MacOS") {
			resourcesDir := filepath.Join(exeDir, "../Resources")
			if err := os.Chdir(resourcesDir); err != nil {
				log.Printf("Warning: Could not change to Resources directory: %v", err)
			} else {
				log.Println("Running from .app bundle, using Resources directory")
			}
		}
	}

	// Load .env file for configuration (GitHub tokens, etc.)
	// Ignore error if .env doesn't exist - we'll fall back to config file
	_ = godotenv.Load()

	g := game.NewGame()

	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("STELLAR SIEGE - Defend the Frontier")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

	// Write memory profile if requested
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		log.Printf("Memory profile written to %s", *memprofile)
	}
}
