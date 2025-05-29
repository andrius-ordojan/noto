package main

import "testing"

func Test_listFiles(t *testing.T) {
	tests := []struct {
		name string
		dir  string
		want []string
	}{
		{
			name: "Test with empty directory",
			dir:  "testdata/empty",
			want: []string{},
		},
		{
			name: "Test with single file",
			dir:  "testdata/single",
			want: []string{"testdata/single/file1.md"},
		},
		{
			name: "Test with multiple files",
			dir:  "testdata/multiple",
			want: []string{"testdata/multiple/file1.md", "testdata/multiple/file2.md"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := listFiles(tt.dir)

			if len(got) != len(tt.want) {
				t.Errorf("listFiles() = %v, want %v", got, tt.want)
			}

			for i, file := range got {
				if file != tt.want[i] {
					t.Errorf("listFiles() got[%d] = %v, want %v", i, file, tt.want[i])
				}
			}
		})
	}
}
