package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	flag "github.com/ogier/pflag"

	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
)

func main() {
	var fpath, bucketname, objname string

	flag.Usage = func() {
		fmt.Println("usage: s3put <file>\n\nUpload a file to S3.\n\nYou need to have AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY defined in your ENV.\n\nOPTIONS\n")
		flag.PrintDefaults()
	}

	flag.StringVarP(&bucketname, "bucket", "b", "mm-tests", "The bucket to upload to")
	flag.StringVarP(&objname, "name", "n", "", "The name of the uploaded object (defaults to value passed to -f/--file)")

	flag.Parse()

	if flag.NArg() > 0 {
		fpath = flag.Args()[0]
	}

	if fpath == "" {
		fmt.Println("abort: you must specify a file to upload")
		os.Exit(1)
	}

	var auth aws.Auth
	var err error

	if auth, err = aws.EnvAuth(); err != nil {
		log.Fatal(err)
	}

	c := s3.New(auth, aws.USEast)
	b := c.Bucket(bucketname)

	f, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}

	st, err := os.Stat(fpath)
	if err != nil {
		log.Fatal(err)
	}

	if objname == "" {
		objname = filepath.Base(fpath)
	}

	if err = b.PutReader(objname, f, st.Size(), "binary/octet-stream", s3.PublicRead); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s -> %s/%s\n", fpath, bucketname, objname)
	fmt.Println(b.URL(objname))
	// err := b.PutReader(flag.Args()[0], "binary/octet-stream", perm)
}
