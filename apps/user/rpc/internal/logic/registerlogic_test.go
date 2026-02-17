package logic

import (
	"context"
	"easy-chat/apps/user/rpc/user"
	"testing"
)

func TestRegisterLogic_Register(t *testing.T) {
	type args struct {
		in *user.RegisterReq
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		args args
		// Named output parameters for target function.
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Register success",
			args: args{
				in: &user.RegisterReq{
					Phone:    "13800000000",
					Nickname: "test",
					Password: "123456",
					Avatar:   "png.jpg",
					Sex:      1,
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewRegisterLogic(context.Background(), svcCtx)
			got, gotErr := l.Register(tt.args.in)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Register() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Register() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Log(tt.name, got)
			}
		})
	}
}
