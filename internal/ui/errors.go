package ui

import (
	"errors"
	"fmt"

	"github.com/manifoldco/promptui"
)

// ErrCancelled is returned when user cancels an operation
var ErrCancelled = errors.New("cancelled")

// IsCancelError checks if the error is a cancellation error
func IsCancelError(err error) bool {
	if err == nil {
		return false
	}

	// Check for promptui cancellation errors
	if err == promptui.ErrInterrupt || err == promptui.ErrEOF {
		return true
	}

	// Check for our custom cancellation error
	if errors.Is(err, ErrCancelled) {
		return true
	}

	// Legacy check for string-based cancellation (for compatibility)
	if err.Error() == "cancelled" {
		return true
	}

	return false
}

// HandleCancelError processes cancellation errors consistently
func HandleCancelError(err error) error {
	if IsCancelError(err) {
		fmt.Println("\nCancelled.")
		return nil
	}
	return err
}

// WrapCancelError converts promptui cancellation errors to our standard error
func WrapCancelError(err error) error {
	if err == nil {
		return nil
	}

	if err == promptui.ErrInterrupt || err == promptui.ErrEOF {
		return ErrCancelled
	}

	return err
}
