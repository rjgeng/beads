package main

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/steveyegge/beads/internal/utils"
)

// rejectEmptyLabelFilter reports whether flagName was supplied but every
// element in raw is blank, in which case the caller meant "match nothing"
// (however it phrased that), not "no filter" -- and "no filter" is exactly
// what an empty NormalizeLabels result means everywhere downstream. Using
// cmd.Flags().Changed keeps that distinction, which the normalized slice
// itself has already destroyed by the time any caller sees it.
//
// Shared by every command that takes a label filter -- ready, list, count,
// orphans, blocked, stale and search: each gathers this flag's raw value off
// its own command and calls this before doing anything else with it, so the
// same supplied-but-empty filter is refused everywhere it can be spelled
// rather than in just the one command it was first noticed on. Three of those
// gatherers serve two routes each (blockedFilterFromFlags for blocked's
// direct and proxied paths, parseStaleLabelFilter for stale's, and
// parseSearchLabelFilter for search's), so the refusal covers the
// proxied-server frontend without being written twice.
//
// Pattern and regular-expression filters use the string-specific helper below
// because their values are matched rather than normalized into a label list.
func rejectEmptyLabelFilter(cmd *cobra.Command, flagName string, raw []string) error {
	if !cmd.Flags().Changed(flagName) {
		return nil
	}
	if len(utils.NormalizeLabels(raw)) == 0 {
		return HandleErrorRespectJSON("--%s was supplied but contains no usable label", flagName)
	}
	return nil
}

// rejectEmptyLabelMatchFilter rejects a pattern or regular-expression flag
// that was explicitly supplied without a usable value. Unlike label-list
// filters, these values are intentionally left otherwise unchanged for their
// downstream matcher.
func rejectEmptyLabelMatchFilter(cmd *cobra.Command, flagName, raw string) error {
	if cmd.Flags().Changed(flagName) && strings.TrimSpace(raw) == "" {
		return HandleErrorRespectJSON("--%s was supplied but is empty", flagName)
	}
	return nil
}
