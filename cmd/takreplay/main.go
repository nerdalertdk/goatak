package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/kdudkov/goatak/pkg/cot"
	"github.com/kdudkov/goatak/pkg/cotproto"
	"google.golang.org/protobuf/proto"
)

func main() {
	var format = flag.String("format", "", "dump format (text|json|gpx)")
	var uid = flag.String("uid", "", "uid to show")
	var typ = flag.String("type", "", "type to show")

	flag.Parse()

	file := flag.Arg(0)

	if file == "" {
		fmt.Println("usege: takreplay <filename>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	f, err := os.Open(file)

	if err != nil {
		panic(err)
	}

	var dmp Dumper

	switch *format {
	case "", "text":
		dmp = new(TextDumper)
	case "json":
		dmp = new(JsonDumper)
	case "json2":
		dmp = new(Json2Dumper)
	case "gpx":
		if *uid == "" {
			fmt.Println("need uid to make gpx")
			os.Exit(1)
		}
		dmp = &GpxDumper{name: *uid}
	default:
		fmt.Printf("invalid format %s\n", *format)
		os.Exit(1)
	}

	if err := readFile(f, *uid, *typ, dmp); err != io.EOF {
		fmt.Println(err)
	}
}

func readFile(f *os.File, uid, typ string, dmp Dumper) error {
	dmp.Start()
	defer dmp.Stop()

	lenBuf := make([]byte, 2)
	for {
		if _, err := io.ReadFull(f, lenBuf); err != nil {
			return err
		}

		n := uint32(lenBuf[0]) + uint32(lenBuf[1])*256
		buf := make([]byte, n)

		if _, err := io.ReadFull(f, buf); err != nil {
			return err
		}

		m := new(cotproto.TakMessage)
		if err := proto.Unmarshal(buf, m); err != nil {
			return err
		}

		if uid != "" && m.GetCotEvent().GetUid() != uid {
			continue
		}

		if typ != "" && !cot.MatchPattern(m.GetCotEvent().GetType(), typ) {
			continue
		}

		d, err := cot.DetailsFromString(m.GetCotEvent().GetDetail().GetXmlDetail())
		msg := &cot.CotMessage{TakMessage: m, Detail: d}
		if err = dmp.Process(msg); err != nil {
			return err
		}
	}
}
