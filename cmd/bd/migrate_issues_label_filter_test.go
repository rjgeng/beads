package main

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestParseMigrateIssuesLabelsRejectsEmptyFilter(t *testing.T) {
	for _, value := range []string{"", "   ", ",,"} {
		t.Run(value, func(t *testing.T) {
			cmd := &cobra.Command{Use: "issues"}
			cmd.Flags().StringSlice("label", nil, "")
			if err := cmd.Flags().Set("label", value); err != nil {
				t.Fatal(err)
			}

			_, err := parseMigrateIssuesLabels(cmd)
			if err == nil {
				t.Fatalf("parseMigrateIssuesLabels(%q) error = %v, want empty-label refusal", value, err)
			}
		})
	}
}

func TestParseMigrateIssuesLabelsNormalizesUsableFilter(t *testing.T) {
	cmd := &cobra.Command{Use: "issues"}
	cmd.Flags().StringSlice("label", nil, "")
	if err := cmd.Flags().Set("label", "a,, b "); err != nil {
		t.Fatal(err)
	}

	got, err := parseMigrateIssuesLabels(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Fatalf("labels = %q, want [a b]", got)
	}
}

func TestParseMigrateIssuesLabelsAllowsOmittedFilter(t *testing.T) {
	cmd := &cobra.Command{Use: "issues"}
	cmd.Flags().StringSlice("label", nil, "")

	got, err := parseMigrateIssuesLabels(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("labels = %q, want empty", got)
	}
}
