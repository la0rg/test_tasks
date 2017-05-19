package vector_clock

import (
	"reflect"
	"testing"
)

func TestNewVc(t *testing.T) {
	newVc := NewVc()
	if newVc == nil || newVc.store == nil {
		t.Error("VectorClock itself as well as its fields should not be nil.")
	}
}

func TestVC_Incr(t *testing.T) {
	vc := NewVc()
	vc.Incr("test")
	vc.Incr("test")
	if vc.Get("test") != 2 {
		t.Error("Value for node test expected to be 2")
	}
}

func TestMerge(t *testing.T) {
	type args struct {
		vc1 *VC
		vc2 *VC
	}
	tests := []struct {
		name string
		args args
		want *VC
	}{
		{
			args: args{
				NewVc(),
				NewVc(),
			},
			want: NewVc(),
		},
		{
			args: args{
				&VC{
					store: map[string]uint64{
						"A": 3,
						"B": 4,
						"C": 1,
					},
				},
				&VC{
					store: map[string]uint64{
						"A": 0,
						"B": 0,
						"C": 2,
					},
				},
			},
			want: &VC{
				store: map[string]uint64{
					"A": 3,
					"B": 4,
					"C": 2,
				},
			},
		},
		{
			args: args{
				&VC{
					store: map[string]uint64{
						"A": 3,
					},
				},
				&VC{
					store: map[string]uint64{
						"B": 1,
						"C": 2,
					},
				},
			},
			want: &VC{
				store: map[string]uint64{
					"A": 3,
					"B": 1,
					"C": 2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Merge(tt.args.vc1, tt.args.vc2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Merge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	type args struct {
		vc1 *VC
		vc2 *VC
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{
				&VC{
					store: map[string]uint64{
						"A": 3,
						"B": 4,
						"C": 1,
					},
				},
				&VC{
					store: map[string]uint64{
						"A": 0,
						"B": 0,
						"C": 2,
					},
				},
			},
			want: 0,
		},
		{
			args: args{
				&VC{
					store: map[string]uint64{
						"A": 1,
						"B": 2,
						"C": 1,
					},
				},
				&VC{
					store: map[string]uint64{
						"A": 0,
						"B": 1,
						"C": 1,
					},
				},
			},
			want: 1,
		},
		{
			args: args{
				&VC{
					store: map[string]uint64{
						"B": 1,
						"C": 1,
					},
				},
				&VC{
					store: map[string]uint64{
						"A": 2,
						"B": 4,
						"C": 1,
					},
				},
			},
			want: -1,
		},
		{
			args: args{
				&VC{
					store: map[string]uint64{
						"B": 3,
						"C": 3,
					},
				},
				&VC{
					store: map[string]uint64{
						"A": 2,
						"B": 4,
						"C": 1,
					},
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Compare(tt.args.vc1, tt.args.vc2); got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCorrolary(t *testing.T) {
	vc1 := &VC{
		store: map[string]uint64{
			"a": 1,
			"b": 2,
		},
	}

	vc2 := &VC{
		store: map[string]uint64{
			"c": 4,
			"b": 5,
		},
	}

	union := Merge(vc1, vc2)
	union.Incr("new node")

	if Compare(union, vc1) != 1 || Compare(union, vc2) != 1 {
		t.Error("Merge of VC and it's increment (on any node) should happens AFTER any of the sources VCs")
	}
}
