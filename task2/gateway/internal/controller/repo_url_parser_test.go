package controller

import "testing"

func TestParseGitHubRepoURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		raw       string
		wantOwner string
		wantRepo  string
		wantErr   bool
	}{
		{
			name:      "valid https",
			raw:       "https://github.com/suvorovrain/golang-course",
			wantOwner: "suvorovrain",
			wantRepo:  "golang-course",
		},
		{
			name:      "valid host without scheme",
			raw:       "github.com/suvorovrain/golang-course",
			wantOwner: "suvorovrain",
			wantRepo:  "golang-course",
		},
		{
			name:      "valid url with .git",
			raw:       "https://github.com/suvorovrain/golang-course.git",
			wantOwner: "suvorovrain",
			wantRepo:  "golang-course",
		},
		{
			name:    "invalid empty",
			raw:     "",
			wantErr: true,
		},
		{
			name:    "invalid non github host",
			raw:     "https://gitlab.com/suvorovrain/golang-course",
			wantErr: true,
		},
		{
			name:    "invalid scheme",
			raw:     "ftp://github.com/suvorovrain/golang-course",
			wantErr: true,
		},
		{
			name:    "invalid extra path segments",
			raw:     "https://github.com/suvorovrain/golang-course/tree/main",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotOwner, gotRepo, err := ParseGitHubRepoURL(tc.raw)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got owner=%q repo=%q", gotOwner, gotRepo)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotOwner != tc.wantOwner || gotRepo != tc.wantRepo {
				t.Fatalf("unexpected parse result: owner=%q repo=%q", gotOwner, gotRepo)
			}
		})
	}
}
