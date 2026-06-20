package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ValentinTorassa/autoconfine/internal/config"
	"github.com/ValentinTorassa/autoconfine/internal/drift"
	"github.com/ValentinTorassa/autoconfine/internal/enforce"
	"github.com/ValentinTorassa/autoconfine/internal/generate"
	"github.com/ValentinTorassa/autoconfine/internal/learn"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "learn":
		os.Exit(runLearn(os.Args[2:]))
	case "generate":
		os.Exit(runGenerate(os.Args[2:]))
	case "enforce":
		os.Exit(runEnforce(os.Args[2:]))
	case "version":
		fmt.Println("autoconfine", config.Version)
		os.Exit(0)
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "Uso: autoconfine <learn|generate|enforce|version> [opciones]")
}

func runLearn(args []string) int {
	fs := flag.NewFlagSet("learn", flag.ExitOnError)
	image := fs.String("image", "", "imagen OCI a observar")
	duration := fs.Duration("duration", 30*time.Second, "duración de la fase de aprendizaje")
	out := fs.String("out", "autoconfine.trace.jsonl", "archivo de traza de salida")
	fs.Parse(args)

	cfg := learn.Config{
		Image:    *image,
		Duration: *duration,
		Output:   *out,
	}

	tracer := learn.NewTracer(cfg)
	if err := tracer.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "learn: %v\n", err)
		return 1
	}
	fmt.Printf("Traza guardada en %s\n", *out)
	return 0
}

func runGenerate(args []string) int {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	out := fs.String("out", "autoconfine.seccomp.json", "perfil seccomp de salida")
	fs.Parse(args)

	if fs.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "generate: se requiere ruta de traza")
		return 1
	}

	gen := generate.NewGenerator(generate.Config{})
	if err := gen.Generate(fs.Arg(0), *out); err != nil {
		fmt.Fprintf(os.Stderr, "generate: %v\n", err)
		return 1
	}
	fmt.Printf("Perfil seccomp guardado en %s\n", *out)
	return 0
}

func runEnforce(args []string) int {
	fs := flag.NewFlagSet("enforce", flag.ExitOnError)
	profile := fs.String("profile", "", "ruta al perfil seccomp generado")
	audit := fs.Bool("audit", false, "modo audit: reporta drift sin bloquear")
	fs.Parse(args)

	if *profile == "" {
		fmt.Fprintln(os.Stderr, "enforce: se requiere --profile")
		return 1
	}

	cfg := enforce.Config{
		ProfilePath: *profile,
		Audit:       *audit,
		DriftReporter: drift.NewJSONReporter(os.Stdout),
	}

	runner := enforce.NewRunner(cfg)
	if err := runner.Run(fs.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "enforce: %v\n", err)
		return 1
	}
	return 0
}
