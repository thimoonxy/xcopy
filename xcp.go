package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/mitchellh/ioprogress"
	"github.com/wjkohnen/bwio"
)

var (
	src, dst *string
	err      error
	bw       *int
)
var verbose = new(bool)
var human = new(bool)
var progress = new(bool)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	start := time.Now().Unix()
	sep := sep_perOS()
	pwd, _ := os.Getwd()
	src = flag.String("s", pwd, "src path")
	dst = flag.String("d", pwd, "dst path")
	bw = flag.Int("l", 0, "Bandwitdh limts, e.g. 50 means 50MB/s")

	//	*verbose = false
	flag.BoolVar(verbose, "v", false, "Print Verbose.")
	flag.BoolVar(human, "h", false, "Print human readable output .")
	flag.BoolVar(progress, "p", false, "Print progress status .")

	flag.Parse()
	if *src == *dst {
		flag.Usage()
		fmt.Println("Error:\ndst not allowed same as src")
		os.Exit(1)
	}

	*bw = *bw << 20
	fmt.Printf("src: %s\ndst: %s\nCopying ...\n\n", *src, *dst)

	lsrc := flag.Lookup("s")
	if lsrc.Value.String() == lsrc.DefValue {
		flag.Usage()
		fmt.Println("\nError:\nsrc not found.")
		os.Exit(1)
	}
	ldst := flag.Lookup("d")
	if ldst.Value.String() == ldst.DefValue {
		flag.Usage()
		fmt.Println("\nError:\ndst not found.")
		os.Exit(1)
	}

	s_info, err := os.Stat(*src)
	//	perm := s_info.Mode()
	if err != nil {
		fmt.Println("Error:\n", err)
		os.Exit(1)
	}
	dst_folder, dst_is_folder := dst_folder_parse(*dst)
	_, err = os.Stat(dst_folder)
	//	fmt.Println("dst_folder in main:", dst_folder)
	if err != nil {
		//	dont use this	os.MkdirAll(*dst, perm), but raise err instead. User must mkdir dst folder before cp.
		fmt.Println("Error:\n", err)
		os.Exit(1)
	}

	if s_info.IsDir() {
		if dst_is_folder == false {
			fmt.Println("Error:\nsrc is dir, but dst is a file.")
			os.Exit(1)
		} else {
			walk(*src)
		}
	} else {
		if dst_is_folder {

			*dst = strings.TrimRight(*dst, sep) + sep + s_info.Name()
		}
		cp(*src, *dst)
	}
	//	time.Sleep(200 * time.Second)
	fmt.Printf("\n\nElapsed: %ds .", time.Now().Unix()-start)
}

func walk(path string) {
	i, _ := ioutil.ReadDir(path)
	sep := sep_perOS()
	for _, info := range i {
		perm := info.Mode()
		if info.IsDir() {
			next_src_path := path + sep + info.Name()
			next_dst_path := *dst + strings.Split(next_src_path, *src)[1]
			//			fmt.Println("mkdir ", next_dst_path)
			os.MkdirAll(next_dst_path, perm)
			walk(next_src_path)
		} else {
			src_fname := path + sep + info.Name()
			dst_fname := *dst + strings.Split(src_fname, *src)[1]
			dst_folder := *dst + strings.Split(path, *src)[1]
			os.MkdirAll(dst_folder, perm)
			//			fmt.Println("dst_folder ", dst_folder)
			cp(src_fname, dst_fname)
		}
	}

}

var byteUnits = []string{"B", "KB", "MB", "GB", "TB", "PB"}

func DrawTerminal(w io.Writer) ioprogress.DrawFunc {
	return ioprogress.DrawTerminalf(w, func(progress, total int64) string {
		return fmt.Sprintf("%s/%s", byteUnitStr(progress), byteUnitStr(total))
	})
}

func byteUnitStr(n int64) string {
	var unit string
	size := float64(n)
	for i := 1; i < len(byteUnits); i++ {
		if size < 1024 {
			unit = byteUnits[i-1]
			break
		}

		size = size / 1024
	}

	return fmt.Sprintf("%.2f %s", size, unit)
}

func cp(src_fname, dst_fname string) {

	var (
		draw  ioprogress.DrawFunc
		bwsrc *bwio.Reader
		bwdst *bwio.Writer
	)
	lbw := flag.Lookup("l")

	d, _ := os.Create(dst_fname)
	s, _ := os.Open(src_fname)
	src_stat, _ := s.Stat()
	if lbw.Value.String() != lbw.DefValue {

		bwsrc = bwio.NewReader(s, *bw)
		bwdst = bwio.NewWriter(d, *bw)

	}

	st := os.Stdout
	if *verbose == true {
		fmt.Fprintf(st, "%s\n", dst_fname)
	}

	if *human == true {
		draw = DrawTerminal(st)
	} else {
		draw = nil
	}

	progressR := &ioprogress.Reader{
		//		Reader:   src,
		Size:     src_stat.Size(),
		DrawFunc: draw,
	}

	if *progress == true && lbw.Value.String() == lbw.DefValue {
		progressR.Reader = s
		_, err = io.Copy(d, progressR)
	} else if *progress == true && lbw.Value.String() != lbw.DefValue {
		progressR.Reader = bwsrc
		//		fmt.Println(*bw)
		bwio.Copy(bwdst, progressR, *bw)
	} else if *progress == false && lbw.Value.String() != lbw.DefValue {
		//		fmt.Println(*bw)
		bwio.Copy(bwdst, bwsrc, *bw)

	} else {
		_, err = io.Copy(d, s)
	}

	if err != nil {
		fmt.Fprintf(st, "# Failed copying %s to %s .\n", src_fname, dst_fname)
	}
}
func sep_perOS() (sep string) {
	os, has := os.LookupEnv("OS")
	if has {
		os = strings.ToLower(os)
		windows := regexp.MustCompile(".*(windows).*")
		if len(windows.FindStringSubmatch(os)) >= 1 {
			sep = "\\"
		} else {
			sep = "/"
		}

	}
	return
}

func dst_folder_parse(d string) (folder string, dst_is_folder bool) {
	dst_is_folder = true
	sep := sep_perOS()
	// replace twice since '/' might be used in gitbash on windows
	d = strings.Replace(d, "\\", sep, len(d))
	d = strings.Replace(d, "/", sep, len(d))
	isfile := regexp.MustCompile(sep + "\\$")
	t := isfile.FindAllString(d, 1)

	if len(t) == 1 {
		folder = d
		dst_is_folder = true
	} else {
		fname_slice := strings.Split(d, sep)
		fname := fname_slice[len(fname_slice)-1]
		folder = strings.Split(d, fname)[0]
		dst_is_folder = false
		//		fmt.Println("folder in parse func:", folder)
	}
	return folder, dst_is_folder
}
