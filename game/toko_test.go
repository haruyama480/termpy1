package game

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haruyama480/termpy1/pu2"
)

func TestNewToko(t *testing.T) {
	type args struct {
		width  int
		height int
		seed   int64
	}
	tests := []struct {
		name          string
		args          args
		want          *Toko
		wantYamaHead5 string
	}{
		{
			name: "normal",
			args: args{
				seed: 2,
			},
			want: &Toko{
				rec: pu2.NewSoloRecord(pu2.NewYama(pu2.TsumoLoop, 2)),
				tax: 2,
				tay: pu2.Y12,
			},
			wantYamaHead5: "OX\nX0\nX0\n0X\nIX",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewToko(tt.args.width, tt.args.height, tt.args.seed)

			opts := cmp.AllowUnexported(Toko{}, pu2.Yama{}, pu2.SoloRecord{})
			if diff := cmp.Diff(got, tt.want, opts); diff != "" {
				t.Errorf("NewToko() mismatch (-got +want):\n%s", diff)
			}

			gotString := got.rec.Yama.String()
			if !strings.HasPrefix(gotString, tt.wantYamaHead5) {
				t.Errorf("NewToko() = %v, want %v", gotString[:len(tt.wantYamaHead5)], tt.wantYamaHead5)
			}
		})
	}
}
