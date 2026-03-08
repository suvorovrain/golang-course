package main

import
(
	"fmt"
	"os"
	"flag"
)

func usage() {
    fmt.Fprintf(os.Stderr, "usage: %s [url_of_github_repo]\n", os.Args[0])
    flag.PrintDefaults()
    os.Exit(2)
}


func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
        fmt.Println("Input file is missing.");
        os.Exit(1);
    }
	args := flag.Args()
    fmt.Printf("opening %s\n", args[0]);

}
