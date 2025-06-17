package ui

import (
	"errors"
	"testing"

	"github.com/manifoldco/promptui"
)

func TestIsCancelError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "promptui ErrInterrupt",
			err:  promptui.ErrInterrupt,
			want: true,
		},
		{
			name: "promptui ErrEOF",
			err:  promptui.ErrEOF,
			want: true,
		},
		{
			name: "ErrCancelled",
			err:  ErrCancelled,
			want: true,
		},
		{
			name: "wrapped ErrCancelled",
			err:  errors.New("cancelled"),
			want: true,
		},
		{
			name: "other error",
			err:  errors.New("some other error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCancelError(tt.err); got != tt.want {
				t.Errorf("IsCancelError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWrapCancelError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr error
	}{
		{
			name:    "nil error",
			err:     nil,
			wantErr: nil,
		},
		{
			name:    "promptui ErrInterrupt",
			err:     promptui.ErrInterrupt,
			wantErr: ErrCancelled,
		},
		{
			name:    "promptui ErrEOF", 
			err:     promptui.ErrEOF,
			wantErr: ErrCancelled,
		},
		{
			name:    "other error",
			err:     errors.New("some error"),
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WrapCancelError(tt.err)
			if tt.wantErr == nil {
				if got != nil {
					t.Errorf("WrapCancelError() = %v, want nil", got)
				}
			} else if got == nil {
				t.Errorf("WrapCancelError() = nil, want %v", tt.wantErr)
			} else if got.Error() != tt.wantErr.Error() {
				t.Errorf("WrapCancelError() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestHandleCancelError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr bool
	}{
		{
			name:    "cancel error returns nil",
			err:     ErrCancelled,
			wantErr: false,
		},
		{
			name:    "promptui interrupt returns nil",
			err:     promptui.ErrInterrupt,
			wantErr: false,
		},
		{
			name:    "other error is returned",
			err:     errors.New("some error"),
			wantErr: true,
		},
		{
			name:    "nil error returns nil",
			err:     nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := HandleCancelError(tt.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleCancelError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}