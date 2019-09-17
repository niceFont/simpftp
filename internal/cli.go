package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/olekukonko/tablewriter"
)

//MainLoop for simpFTP reading and dispatching commands
func MainLoop(srv io.Closer) {

	var err error

	defer func() {
		checkErrors(err)
	}()
	inputReader := bufio.NewReader(os.Stdin)

	var input []string
	var args []string
	var command string
	for command != "quit" {
		var ip string
		ip, err = inputReader.ReadString('\n')
		input = strings.Split(ip, " ")
		command, args, err = strip(input)
		switch command {
		case "see":
			see()
		case "sd":
			seedir(args[0])
		case "maked":
			maked(args[0])
		case "makef":
			makef(args[0])
		case "del":
			del(args[0])
		case "move":
			move(args[0], args[1])
		}
	}

	err = srv.Close()

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
		checkErrors(err)
	}()
}

func seedir(dir string) {
	var err error
	defer func() {
		checkErrors(err)
	}()
	directories := strings.Split(dir, "/")
	if len(directories) > 0 {
		for _, directory := range directories {
			if directory != "." && directory != "" {
				err = os.Chdir(directory)
			}
		}
	}
}

func maked(name string) {
	var err error
	defer func() {
		checkErrors(err)
	}()
	err = os.Mkdir(name, os.ModeDir)
}

func del(target string) {

	var err error
	defer func() {
		checkErrors(err)
	}()
	err = os.Remove(target)
}

func makef(name string) {

	var err error
	defer func() {
		checkErrors(err)
	}()
	file, err := os.Create(name)

	err = file.Close()

}

func move(target, destination string) {
	var err error

	defer func() {
		checkErrors(err)
	}()

	var targetContent []byte

	t := strings.Split(target, "/")
	d := strings.Split(destination, "/")
	targetPathEnd := t[len(t)-1]
	destPathEnd := d[len(d)-1]
	nt := strings.Join(t[:len(t)-1], "/")
	nd := strings.Join(d[:len(d)-1], "/")

	seedir(nt)

	files, err := ioutil.ReadDir(".")

	//Navigating to target Directory.
	//Reading File and navigating back to root.
	for _, f := range files {
		if f.Name() == targetPathEnd {
			if !f.IsDir() {
				targetContent, err = ioutil.ReadFile(targetPathEnd)
				del(f.Name())

				//Navigating Back based on the amount of Directories traveled.
				for i := 0; i < len(t)-2; i++ {
					seedir("..")
				}
			} else {
				err = errors.New("No such file found in " + target)
				return
			}
		}
	}

	seedir(nd)

	//Navigating to destination Directory.
	//Checking if the End of the give Path is a File or a Dir
	//If its a File that already exists, we delete the old one and create a file with the
	//Content of the target File
	files, err = ioutil.ReadDir(".")
	for _, f := range files {
		if f.Name() == destPathEnd {
			if !f.IsDir() {
				del(destPathEnd)
				file, _ := os.Create(destPathEnd)
				_, err = file.Write(targetContent)
				err = file.Sync()
				err = file.Close()

				for i := 0; i < len(d)-2; i++ {
					seedir("..")
				}
				return
			}
			seedir(destPathEnd)
			file, _ := os.Create(targetPathEnd)
			_, err = file.Write(targetContent)
			err = file.Sync()
			err = file.Close()
			//Navigating Back based on the amount of Directories traveled.
			for i := 0; i < len(d)-1; i++ {
				seedir("..")
			}
			return

		}
	}
	file, _ := os.Create(destPathEnd)
	_, err = file.Write(targetContent)
	err = file.Sync()
	err = file.Close()

}

//maybe func rename()

func strip(input []string) (string, []string, error) {
	if len(input) == 1 {

		if runtime.GOOS == "linux" {
			return strings.Trim(input[0], "\n"), nil, nil
		}
		if runtime.GOOS == "windows" {

			return strings.Trim(input[0], "\r\n"), nil, nil
		}
		err := errors.New("internal error: couldn't recognize system")
		return "", nil, err
	}

	a := make([]string, 1)
	if runtime.GOOS == "linux" {
		for i, arg := range input[1:] {
			if len(input[1:])-1 == i {
				a = append(a, arg[:len(arg)-1])
			} else {
				a = append(a, arg)
			}
		}
		return input[0], a[1:], nil
	}
	if runtime.GOOS == "windows" {

		for i, arg := range input[1:] {
			if len(input[1:])-1 == i {
				a = append(a, arg[:len(arg)-2])
			} else {
				a = append(a, arg)
			}
		}
		return input[0], a, nil
	}
	err := errors.New("internal error: couldn't recognize system")
	return "", nil, err
}

func checkErrors(err error) {
	if err != nil {
		log.Println(err)
	}
}
