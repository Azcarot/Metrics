package measure

//Тест на получение метрик (проверяем по увеличению счетчика,
//так как для остального не знаем конкретных значений)
import (
	"reflect"
	"testing"
)

func TestGetMetrics(t *testing.T) {
	type args struct {
		m MemStorage
	}
	tests := []struct {
		name string
		args args
		want Counter
	}{{
		name: "testcounter",
		args: args{},
		want: 1,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMetrics(tt.args.m); !reflect.DeepEqual(got.Countermem["PollCount"], tt.want) {
				t.Errorf("GetMetrics() = %v, want %v", got.Countermem["PollCount"], tt.want)
			}
		})
	}
}
