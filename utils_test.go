package pegasus

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
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

func TestConf_GetConf(t *testing.T) {
	tests := []struct {
		name string
		c    *Conf
		want *Conf
	}{
		// TODO: Add test cases.
		{
			"解析配置测试",
			&Conf{},
			&Conf{
				Name: "client1",
				Port: 8845,
				Server: &Server{
					IP:   "127.0.0.1",
					Port: 9986,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetConf("./client/conf.yml"); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Conf.getConf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConf_Validate(t *testing.T) {
	tests := []struct {
		name string
		c    *Conf
		want bool
	}{
		// TODO: Add test cases.
		{
			"配置校验测试",
			&Conf{
				Name: "client1",
				Port: 8845,
				Server: &Server{
					IP:   "",
					Port: 0,
				},
			},
			false,
		},
		{
			"配置校验测试",
			&Conf{
				Name: "client1",
				Port: 8845,
				Server: &Server{
					IP:   "127.0.0.1",
					Port: 9986,
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Validate(); got != tt.want {
				t.Errorf("Conf.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAvailableIPAddress(t *testing.T) {
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"获取本机IP测试",
			[]string{"192.168.66.45",
				"172.17.181.81"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAvailableIPAddress()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAvailableIPAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAvailableIPAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
