package pegasus

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
)

func TestCompress(t *testing.T) {
	type args struct {
		origin string
		prefix string
		dest   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"压缩测试",
			args{
				"E:\\workspace\\pegasus\\file\\wxbot",
				"",
				"E:\\workspace\\pegasus\\file\\wxbot2.gz",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Compress(tt.args.origin, tt.args.prefix, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("Compress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeCompress(t *testing.T) {
	type args struct {
		tarFile string
		dest    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"解压测试",
			args{
				"E:\\workspace\\pegasus\\file\\wxbot2.gz",
				"E:\\workspace\\pegasus\\file\\dwxbot2",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeCompress(tt.args.tarFile, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("DeCompress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileCopy(t *testing.T) {
	os.MkdirAll("./abc", 0666)
	bytes, _ := ioutil.ReadFile("LICENSE")
	ioutil.WriteFile("./abc/filename", bytes, 0666)
}

func TestFileMode(t *testing.T) {
	um, _ := strconv.ParseInt(strconv.Itoa(777), 8, 0)
	fmt.Println(os.FileMode(um))
	fmt.Println(os.FileMode(0777))
	fmt.Println(0777)
}
