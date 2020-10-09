package grawler

import "testing"

func TestSave(t *testing.T) {
	type args struct {
		p *Page
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Save(tt.args.p)
		})
	}
}
