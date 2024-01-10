package measure

//Тест на получение метрик (проверяем по увеличению счетчика,
//так как для остального не знаем конкретных значений)
import (
	"reflect"
	"testing"

	"github.com/Azcarot/Metrics/internal/storage"
)

func TestGetMetrics(t *testing.T) {
	type args struct {
		m storage.MemStorage
	}
	tests := struct {
		name string
		args args
		want storage.Counter
	}{
		name: "testcounter",
		args: args{},
		want: 1,
	}

	t.Run(tests.name, func(t *testing.T) {
		if got := CollectMetrics(tests.args.m); !reflect.DeepEqual(got.Countermem["PollCount"], tests.want) {
			t.Errorf("GetMetrics() = %v, want %v", got.Countermem["PollCount"], tests.want)
		}
	})
}
