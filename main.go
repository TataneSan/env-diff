package main

import (
	"os"

	"github.com/TataneSan/env-diff/internal/domain"
	"github.com/TataneSan/env-diff/internal/infrastructure"
	"github.com/TataneSan/env-diff/internal/presentation"
)

func main() {
	printer := presentation.NewPrinter()
	collector := infrastructure.NewEnvCollector()
	differ := infrastructure.NewEnvDiffer()

	if len(os.Args) < 2 {
		printer.PrintUsage()
		os.Exit(1)
	}

	var ctxA, ctxB *domain.EnvContext
	var err error

	switch os.Args[1] {
	case "file":
		if len(os.Args) != 4 {
			printer.PrintUsage()
			os.Exit(1)
		}
		ctxA, err = collector.CollectFromFile(os.Args[2])
		if err != nil {
			printer.PrintError(err.Error())
			os.Exit(1)
		}
		ctxB, err = collector.CollectFromFile(os.Args[3])
		if err != nil {
			printer.PrintError(err.Error())
			os.Exit(1)
		}
	case "current":
		if len(os.Args) != 3 {
			printer.PrintUsage()
			os.Exit(1)
		}
		ctxA = collector.CollectCurrent()
		ctxB, err = collector.CollectFromFile(os.Args[2])
		if err != nil {
			printer.PrintError(err.Error())
			os.Exit(1)
		}
	case "cmd":
		if len(os.Args) != 4 {
			printer.PrintUsage()
			os.Exit(1)
		}
		ctxA, err = collector.CollectFromCommand(os.Args[2])
		if err != nil {
			printer.PrintError(err.Error())
			os.Exit(1)
		}
		ctxB, err = collector.CollectFromCommand(os.Args[3])
		if err != nil {
			printer.PrintError(err.Error())
			os.Exit(1)
		}
	default:
		printer.PrintUsage()
		os.Exit(1)
	}

	diffs := differ.Compare(ctxA, ctxB)
	printer.PrintDiff(diffs, ctxA.Name, ctxB.Name)
}
