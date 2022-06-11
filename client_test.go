package cli

import (
	"reflect"
	"testing"

	"github.com/grezar/go-circleci"
)

func Test_getIntersection(t *testing.T) {
	type args struct {
		vars  []string
		items []*circleci.ProjectVariable
	}
	tests := []struct {
		name string
		args args
		want []*circleci.ProjectVariable
	}{
		// TODO: Add test cases.
		{
			name: "simple array",
			args: args{
				vars: []string{"ENV_0", "ENV_1", "ENV_2"},
				items: []*circleci.ProjectVariable{
					&circleci.ProjectVariable{Name: "ENV_1", Value: "xxxxenv"},
				},
			},
			want: []*circleci.ProjectVariable{
				&circleci.ProjectVariable{Name: "ENV_1", Value: "xxxxenv"},
			},
		},
		{
			name: "two intersected items",
			args: args{
				vars: []string{"ENV_0", "ENV_1", "ENV_2"},
				items: []*circleci.ProjectVariable{
					&circleci.ProjectVariable{Name: "ENV_1", Value: "xxxxenv"},
					&circleci.ProjectVariable{Name: "ENV_2", Value: "xxxxenv"},
					&circleci.ProjectVariable{Name: "ENV_3", Value: "xxxxenv"},
				},
			},
			want: []*circleci.ProjectVariable{
				&circleci.ProjectVariable{Name: "ENV_1", Value: "xxxxenv"},
				&circleci.ProjectVariable{Name: "ENV_2", Value: "xxxxenv"},
			},
		},
		{
			name: "no intersection",
			args: args{
				vars: []string{"ENV_0", "ENV_1", "ENV_2"},
				items: []*circleci.ProjectVariable{
					&circleci.ProjectVariable{Name: "ENV_3", Value: "xxxxenv"},
				},
			},
			want: []*circleci.ProjectVariable{},
		},
		{
			name: "no envs",
			args: args{
				vars:  []string{"ENV_0", "ENV_1", "ENV_2"},
				items: []*circleci.ProjectVariable{},
			},
			want: []*circleci.ProjectVariable{},
		},
		{
			name: "no vars",
			args: args{
				vars: []string{},
				items: []*circleci.ProjectVariable{
					&circleci.ProjectVariable{Name: "ENV_0", Value: "xxxxenv"},
				},
			},
			want: []*circleci.ProjectVariable{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getIntersection(tt.args.vars, tt.args.items); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getIntersection() = %v, want %v", got, tt.want)
			}
		})
	}
}
