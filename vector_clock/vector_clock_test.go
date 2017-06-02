package vector_clock

import (
	"reflect"
	"testing"
)

func TestNewVc(t *testing.T) {
	newVc := NewVc()
	if newVc == nil || newVc.Store == nil {
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
					Store: map[string]uint64{
						"A": 3,
						"B": 4,
						"C": 1,
					},
				},
				&VC{
					Store: map[string]uint64{
						"A": 0,
						"B": 0,
						"C": 2,
					},
				},
			},
			want: &VC{
				Store: map[string]uint64{
					"A": 3,
					"B": 4,
					"C": 2,
				},
			},
		},
		{
			args: args{
				&VC{
					Store: map[string]uint64{
						"A": 3,
					},
				},
				&VC{
					Store: map[string]uint64{
						"B": 1,
						"C": 2,
					},
				},
			},
			want: &VC{
				Store: map[string]uint64{
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
					Store: map[string]uint64{
						"A": 3,
						"B": 4,
						"C": 1,
					},
				},
				&VC{
					Store: map[string]uint64{
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
					Store: map[string]uint64{
						"A": 1,
						"B": 2,
						"C": 1,
					},
				},
				&VC{
					Store: map[string]uint64{
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
					Store: map[string]uint64{
						"B": 1,
						"C": 1,
					},
				},
				&VC{
					Store: map[string]uint64{
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
					Store: map[string]uint64{
						"B": 3,
						"C": 3,
					},
				},
				&VC{
					Store: map[string]uint64{
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

func TestEqual(t *testing.T) {
	vc1 := NewVc()
	vc2 := NewVc()
	if Equal(vc1, vc1) != true {
		t.Error("Should be equal to itself")
	}
	if Equal(vc1, vc2) != true {
		t.Error("Empty vector clocks should be equal")
	}
	vc1.Incr("node1")
	if Equal(vc1, vc2) != false {
		t.Error("VCs of different length shold not be equal")
	}
	vc2.Incr("node2")
	if Equal(vc1, vc2) != false {
		t.Error("All keys should be equal")
	}
	vc1.Incr("node1")
	vc2.Incr("node1")
	if Equal(vc1, vc2) != false {
		t.Error("All values should be equal")
	}
	vc2.Incr("node1")
	vc1.Incr("node2")
	if Equal(vc1, vc2) != true {
		t.Error("VCs are expected to be equal")
	}

}

func TestGetStore(t *testing.T) {
	vc := NewVc()
	vc.Incr("Node")
	store := vc.GetStore()
	if &store == &(vc.Store) {
		t.Error("Stores are expected to be different")
	}
	if v, ok := store["Node"]; !ok || v != 1 {
		t.Error("Stores are expected to contain same values")
	}
}

func TestCorrolary(t *testing.T) {
	vc1 := &VC{
		Store: map[string]uint64{
			"a": 1,
			"b": 2,
		},
	}

	vc2 := &VC{
		Store: map[string]uint64{
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
