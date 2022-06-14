package cli

import (
	"reflect"
	"testing"

	"github.com/grezar/go-circleci"
)

func Test_getFoundAndNotFoundVariables(t *testing.T) {
	type args struct {
		vars  []string
		items []*circleci.ProjectVariable
	}
	tests := []struct {
		name  string
		args  args
		want  []*circleci.ProjectVariable
		want1 []string
	}{
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
			want1: []string{
				"ENV_0", "ENV_2",
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
			want1: []string{
				"ENV_0",
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
			want1: []string{
				"ENV_0", "ENV_1", "ENV_2",
			},
		},
		{
			name: "no envs",
			args: args{
				vars:  []string{"ENV_0", "ENV_1", "ENV_2"},
				items: []*circleci.ProjectVariable{},
			},
			want:  []*circleci.ProjectVariable{},
			want1: []string{"ENV_0", "ENV_1", "ENV_2"},
		},
		{
			name: "no vars",
			args: args{
				vars: []string{},
				items: []*circleci.ProjectVariable{
					&circleci.ProjectVariable{Name: "ENV_0", Value: "xxxxenv"},
				},
			},
			want:  []*circleci.ProjectVariable{},
			want1: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getFoundAndNotFoundVariables(tt.args.vars, tt.args.items)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDeletedAndNotDeletedVariables() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getDeletedAndNotDeletedVariables() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
