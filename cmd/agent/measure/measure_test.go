package measure

//Тест на получение метрик (проверяем по увеличению счетчика,
//так как для остального не знаем конкретных значений)
import (
	"reflect"
	"testing"

	"github.com/Azcarot/Metrics/cmd/types"
)

func TestGetMetrics(t *testing.T) {
	type args struct {
		m types.MemStorage
	}
	tests := []struct {
		name string
		args args
		want types.Counter
	}{{
		name: "testcounter",
		args: args{},
		want: 1,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CollectMetrics(tt.args.m); !reflect.DeepEqual(got.Countermem["PollCount"], tt.want) {
				t.Errorf("GetMetrics() = %v, want %v", got.Countermem["PollCount"], tt.want)
			}
		})
	}
}
