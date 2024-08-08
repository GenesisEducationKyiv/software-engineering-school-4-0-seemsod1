package architecture_test

import (
	"strings"
	"testing"

	"golang.org/x/tools/go/packages"
)

// loadPackages function to load package information
func loadPackages(patterns ...string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedImports,
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, err
	}
	return pkgs, nil
}

// checkDependencies function to check if a package depends on restricted packages
func checkDependencies(t *testing.T, pkg *packages.Package, restricted []string) {
	for _, imp := range pkg.Imports {
		for _, r := range restricted {
			if strings.HasPrefix(imp.PkgPath, r) {
				t.Errorf("Package %s should not depend on %s", pkg.PkgPath, r)
			}
		}
	}
}

func TestRateAPIDependencies(t *testing.T) {
	pkgs, err := loadPackages("github.com/seemsod1/api-project/internal/rateapi")
	if err != nil {
		t.Fatalf("Failed to load packages: %v", err)
	}

	for _, pkg := range pkgs {
		restricted := []string{
			"github.com/seemsod1/api-project/internal/handler",
			"github.com/seemsod1/api-project/internal/notifier",
			"github.com/seemsod1/api-project/internal/scheduler",
		}
		checkDependencies(t, pkg, restricted)
	}
}

func TestNotifierDependencies(t *testing.T) {
	pkgs, err := loadPackages("github.com/seemsod1/api-project/internal/notifier")
	if err != nil {
		t.Fatalf("Failed to load packages: %v", err)
	}

	for _, pkg := range pkgs {
		restricted := []string{
			"github.com/seemsod1/api-project/internal/handler",
			"github.com/seemsod1/api-project/internal/rateapi",
			"github.com/seemsod1/api-project/internal/scheduler",
			"github.com/seemsod1/api-project/internal/storage/dbrepo",
		}
		checkDependencies(t, pkg, restricted)
	}
}

func TestSchedulerDependencies(t *testing.T) {
	pkgs, err := loadPackages("github.com/seemsod1/api-project/internal/scheduler")
	if err != nil {
		t.Fatalf("Failed to load packages: %v", err)
	}

	for _, pkg := range pkgs {
		restricted := []string{
			"github.com/seemsod1/api-project/internal/handler",
			"github.com/seemsod1/api-project/internal/rateapi",
			"github.com/seemsod1/api-project/internal/notifier",
		}
		checkDependencies(t, pkg, restricted)
	}
}

func TestDBRepoDependencies(t *testing.T) {
	pkgs, err := loadPackages("github.com/seemsod1/api-project/internal/storage/dbrepo")
	if err != nil {
		t.Fatalf("Failed to load packages: %v", err)
	}

	for _, pkg := range pkgs {
		restricted := []string{
			"github.com/seemsod1/api-project/internal/handler",
			"github.com/seemsod1/api-project/internal/rateapi",
			"github.com/seemsod1/api-project/internal/notifier",
			"github.com/seemsod1/api-project/internal/scheduler",
		}
		checkDependencies(t, pkg, restricted)
	}
}

func TestHandlerDependencies(t *testing.T) {
	pkgs, err := loadPackages("github.com/seemsod1/api-project/internal/handler")
	if err != nil {
		t.Fatalf("Failed to load packages: %v", err)
	}

	for _, pkg := range pkgs {
		restricted := []string{
			"github.com/seemsod1/api-project/internal/repository",
			"github.com/seemsod1/api-project/internal/logger",
			"github.com/seemsod1/api-project/internal/rateapi",
			"github.com/seemsod1/api-project/internal/notifier",
			"github.com/seemsod1/api-project/internal/scheduler",
			"github.com/seemsod1/api-project/internal/storage/dbrepo",
		}
		checkDependencies(t, pkg, restricted)
	}
}
