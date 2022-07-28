package main

import "testing"

func Test_extractRepoName(t *testing.T) {
	type args struct {
		repo string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			name: "normal case",
			args: args{
				repo: "git@github.com:threepipes/circleci-env.git",
			},
			want:  "circleci-env",
			want1: "threepipes",
		},
		{
			name: "trailing slash",
			args: args{
				repo: "git@github.com:threepipes/circleci-env.git/",
			},
			want:  "circleci-env",
			want1: "threepipes",
		},
		{
			name: "https url",
			args: args{
				repo: "https://github.com/user/repo.git",
			},
			want:  "repo",
			want1: "user",
		},
		{
			name: "https url (.git suffix omitted)",
			args: args{
				repo: "https://github.com/user/repo",
			},
			want:  "repo",
			want1: "user",
		},
		{
			name: "https url containing '.git' in repo name",
			args: args{
				repo: "https://github.com/user/repo.git.repo.git",
			},
			want:  "repo.git.repo",
			want1: "user",
		},
		{
			name: "https url containing '.git' in repo name (.git suffix omitted)",
			args: args{
				repo: "https://github.com/user/repo.git.repo",
			},
			want:  "repo.git.repo",
			want1: "user",
		},
		{
			name: "error pattern",
			args: args{
				repo: "repo.git",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := extractRepoName(tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractRepoName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractRepoName() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("extractRepoName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
