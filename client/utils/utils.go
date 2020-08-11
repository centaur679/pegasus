package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Conf client配置
type Conf struct {
	Name   string  `yaml:"name"`
	Port   int     `yaml:"port"`
	Server *Server `yaml:"server"`
}

// Server 配置
type Server struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

// GetConf 读取配置
func (c *Conf) GetConf(filePath string) *Conf {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return c
}

// Validate  校验配置
func (c *Conf) Validate() bool {
	switch {
	case c.Port <= 0, c.Server == nil, c.Server.IP == "", c.Server.Port <= 0:
		return false
	default:
		return true
	}
}

// Compress 压缩
func Compress(origin, prefix, dest string) error {
	file, err := os.Open(origin)
	if err != nil {
		return err
	}
	defer file.Close()

	fw, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer fw.Close()

	gw := gzip.NewWriter(fw)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	compress(file, prefix, tw)
	return nil
}

func compress(file *os.File, prefix string, tw *tar.Writer) error {

	info, err := file.Stat()
	if err != nil {
		return err
	}

	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, tw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = prefix + "/" + header.Name
		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(tw, file)
		file.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// DeCompress 解压
func DeCompress(tarFile, dest string) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		filename := dest + hdr.Name
		file, err := createFile(filename)
		if err != nil {
			return err
		}
		io.Copy(file, tr)
	}
	return nil
}

func createFile(name string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}

// GetAvailableIPAddress 获取本地可用IP地址
func GetAvailableIPAddress() ([]string, error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return []string{}, err
	}

	var res []string
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						fmt.Println(ipnet.IP.String())
						res = append(res, ipnet.IP.String())
					}
				}
			}
		}
	}

	return res, nil
}
