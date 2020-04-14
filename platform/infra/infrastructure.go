// Package infra provides functions for orchestrating a Micro platform
package infra

import (
	"strings"
)

// Task describes an individual task
type Task interface {
	Validate() error
	Plan() error
	Apply() error
	Finalise() error
	Destroy() error
}

// Step is a list of parallisable tasks
type Step []Task

// ExecutePlan carries out a plan on steps
func ExecutePlan(steps []Step) error {
	for _, step := range steps {
		for _, t := range step {
			defer t.Finalise()
			if err := t.Validate(); err != nil {
				return err
			}
		}
	}
	return nil
}

// ExecuteApply carries out an apply on steps
func ExecuteApply(steps []Step) error {
	for _, step := range steps {
		for _, t := range step {
			defer t.Finalise()
			if err := t.Validate(); err != nil {
				return err
			}
			if err := t.Apply(); err != nil {
				return err
			}
		}
	}
	return nil
}

// ExecuteDestroy destroys steps
func ExecuteDestroy(steps []Step) error {
	// Find any kubeconfig steps; we need them to destroy the resources
	for _, s := range steps {
		for _, task := range s {
			switch t := task.(type) {
			case *TerraformModule:
				if strings.Contains(t.Source, "kubeconfig") {
					defer t.Finalise()
					if err := t.Validate(); err != nil {
						return err
					}
					if err := t.Apply(); err != nil {
						return err
					}
					t.Variables["kubernetes"] = "none"
					defer t.Destroy()
				}
			}
		}
	}
	for i := len(steps) - 1; i >= 0; i-- {
		for _, task := range steps[i] {
			switch t := task.(type) {
			case *TerraformModule:
				// Skip any kubeconfig steps
				if !strings.Contains(t.Source, "kubeconfig") {
					defer t.Finalise()
					if err := t.Validate(); err != nil {
						return err
					}
					if err := t.Destroy(); err != nil {
						return err
					}
				}
			default:
				defer t.Finalise()
				if err := t.Validate(); err != nil {
					return err
				}
				if err := t.Destroy(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
