package internal

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

//MainLoop for simpFTP reading and dispatching commands
func MainLoop(srv io.Closer) {

	inputReader := bufio.NewReader(os.Stdin)

	var input []string
	var command string
	for command != "quit\r\n" {
		ip, err := inputReader.ReadString('\n')
		input = strings.Split(ip, " ")
		command = input[0]
		if err != nil {
			log.Println(err)
		}
		switch command {
		case "see\r\n":
			see()
		case "sd":
			seedir(input[1][:len(input[1])-2])
		case "maked":
			maked(input[1][:len(input[1])-2])
		case "makef":
			makef(input[1][:len(input[1])-2])
		case "del":
			del(input[1][:len(input[1])-2])
		}

	}

	defer func() {
		err := srv.Close()

		if err != nil {
			log.Println(err)
		}
	}()

}

func see() {
	var err error

	info, err := ioutil.ReadDir(".")
	data := make([][]string, 1)
	for i, f := range info {
		if len(data) == cap(data) {
			n := make([][]string, cap(data)+1)
			copy(n, data)
			data = n
		}
		year, month, day := f.ModTime().Date()
		date := fmt.Sprintf("%d/%d/%d", day, month, year)
		perm := fmt.Sprintf("%#o", f.Mode().Perm())
		size := fmt.Sprintf("%d", f.Size())
		var ftype string
		if f.IsDir() {
			ftype = "dir"
		} else {
			ftype = "file"
		}
		data[i] = []string{perm, f.Name(), ftype, size, date}
	}

	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"perm", "name", "type", "size", "mod"})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()

	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()
}

func seedir(dir string) {
	var err error

	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	err = os.Chdir(dir)
}

func maked(name string) {
	var err error
	defer func() {
		log.Println(err)
	}()
	err = os.Mkdir(name, os.ModeDir)
}

func del(target string) {

	var err error
	defer func() {
		log.Println(err)
	}()
	err = os.Remove(target)
}

func makef(name string) {

	var err error
	defer func() {
		log.Println(err)
	}()
	file, err := os.Create(name)

	err = file.Close()

}
