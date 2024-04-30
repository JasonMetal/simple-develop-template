package time

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

var TypeArray = []string{
	"2006-01-02 15:04:05",
	"2006年01月02日 15时04分05秒",
	"2006/01/02 15:04:05",
	"2006/01-02",
	"2006/01/02",
	"2006年01月02日",
	"15:04:05",
	"15时04分05秒",
	"15:04",
	"15时04分",
	"04分05秒",
	"04:05",
	"2006年1月2日 3时4分5秒",
	"2006年1月2日 3时4分5秒",
	"2009-9-2 5:4:5",
}
var TimeArray = []string{
	"2022-04-28 09:19:29",
	"2022年04月28日 09时19分29秒",
	"2022/04/28 09:19:29",
	"2022-04-28",
	"2022/04/28",
	"2022年04月28日",
	"09:19:29",
	"09时19分29秒",
	"09:19",
	"09时19分",
	"19分29秒",
	"19:29",
	"2022年4月28日 9时19分29秒",
	"2022/04-28 9:19:29",
	"2022/04-28 9:19:29",
}

// 可忽略的测试
func TestNewDateTime(t *testing.T) {
	type args struct {
		dateTime   string
		timeLayout string
	}

	for k, layout := range TypeArray {
		times := TimeArray[k]
		//fmt.Println("格式：", layout, "时间：", times)
		eq, _ := time.ParseInLocation(times, layout, time.Local)
		if got := NewDateTime(layout, times); !reflect.DeepEqual(got, eq) {
			t.Errorf("CreateTime() = %v, want %v", got, eq)
		}
	}
}

func TestDate(t *testing.T) {

	var timestamp = int64(1651111878)

	for _, layout := range TypeArray {
		result := time.Unix(timestamp, 0).Format(layout)
		got := Date(timestamp, layout)
		if got != result {
			t.Errorf("Date() = %v, want %v", got, result)
		}
		//fmt.Printf("Date转换数据%v,真实数据%v  \n", got, result)

	}
}

func TestDiff(t *testing.T) {
	var (
		time1,
		time2,
		time3,
		time4,
		time5,
		time6,
		time7,
		time8 string
	)

	time1 = "2022-04-28 11:11:15"
	time2 = "2022-04-28 11:11:28"
	tmp1 := NewDateTime(time1, TypeArray[0])
	tmp2 := NewDateTime(time2, TypeArray[0])

	time3 = "2020-04-08 23:50:45"
	time4 = "2000-04-28 11:00:00"
	tmp3 := NewDateTime(time3, TypeArray[0])
	tmp4 := NewDateTime(time4, TypeArray[0])

	time5 = "2022-04-28"
	time6 = "2022-04-28"
	tmp5 := NewDateTime(time5, TypeArray[0])
	tmp6 := NewDateTime(time6, TypeArray[0])

	time7 = "13:35:00"
	time8 = "8:26:02"
	tmp7 := NewDateTime(time7, TypeArray[6])
	tmp8 := NewDateTime(time8, TypeArray[6])

	tests := []struct {
		name       string
		time1      time.Time
		time2      time.Time
		wantYear   int
		wantMonth  int
		wantDay    int
		wantHour   int
		wantMinute int
		wantSecond int
	}{
		// TODO: Add test cases.
		{
			"时间 差13秒",
			tmp1,
			tmp2,
			0, 0, 0, 0, 0, 13,
		},
		{
			"时间差19年11个月11天12小时50分钟45秒 结果给全0",
			tmp3,
			tmp4,
			0, 0, 0, 0, 0, 0,
		},
		{
			"全为0",
			tmp5,
			tmp6,
			0, 0, 0, 0, 0, 0,
		},
		{
			"全为0",
			tmp5,
			tmp6,
			0, 0, 0, 0, 0, 0,
		},
		{
			"间隔为5小时5分钟58秒",
			tmp7,
			tmp8,
			0, 0, 0, 5, 8, 58,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotYear, gotMonth, gotDay, gotHour, gotMinute, gotSecond := Diff(tt.time1, tt.time2)
			if gotYear != tt.wantYear {
				t.Errorf("Diff() gotYear = %v, want %v", gotYear, tt.wantYear)
			}
			if gotMonth != tt.wantMonth {
				t.Errorf("Diff() gotMonth = %v, want %v", gotMonth, tt.wantMonth)
			}
			if gotDay != tt.wantDay {
				t.Errorf("Diff() gotDay = %v, want %v", gotDay, tt.wantDay)
			}
			if gotHour != tt.wantHour {
				t.Errorf("Diff() gotHour = %v, want %v", gotHour, tt.wantHour)
			}
			if gotMinute != tt.wantMinute {
				t.Errorf("Diff() gotMinute = %v, want %v", gotMinute, tt.wantMinute)
			}
			if gotSecond != tt.wantSecond {
				t.Errorf("Diff() gotSecond = %v, want %v", gotSecond, tt.wantSecond)
			}
		})
	}
}

func TestStrToTime(t *testing.T) {
	//
	typeLength := len(TypeArray)

	for k, layout := range TypeArray {
		var wantR int64
		timeStr := TimeArray[k]
		want := NewDateTime(timeStr, layout)
		wantR = want.Unix()
		if typeLength == (k + 1) {
			wantR = int64(1651114766)
		}
		got := StrToTime(timeStr, layout)
		if got != wantR {
			t.Errorf("kye is %d StrToTime() = %v, want %v", k, got, want)
		}
	}
}

func TestTime(t *testing.T) {
	want := time.Now().Unix()
	tests := []struct {
		name string
		want int64
	}{
		// TODO: Add test cases.
		{
			"获取当前时间的时间戳",
			want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Time()
			if got != tt.want {
				t.Errorf("Time() = %v, want %v", got, tt.want)
			}
			fmt.Println("方法返回：", got)
		})
	}
	fmt.Println("得到的时间", want)
}
