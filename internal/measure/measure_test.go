package measure

//Тест на получение метрик (проверяем по увеличению счетчика,
//так как для остального не знаем конкретных значений)
import (
	"reflect"
	"testing"

	"github.com/Azcarot/Metrics/internal/storage"
)

func TestGetMetrics(t *testing.T) {

	tests := struct {
		name string
		want storage.Counter
	}{
		name: "testcounter",
		want: 1,
	}
	result := make(chan storage.MemStorage)
	t.Run(tests.name, func(t *testing.T) {
		go CollectMetrics(result)
		if got := <-result; !reflect.DeepEqual(got.Countermem["PollCount"], tests.want) {
			t.Errorf("GetMetrics() = %v, want %v", got.Countermem["PollCount"], tests.want)
		}
	})
}
