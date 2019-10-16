package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/google/go-containerregistry/pkg/v1/mutate"

	"github.com/google/go-containerregistry/pkg/crane"
)

type Foo interface {
}

func main() {
	pulledImage, err := crane.Pull(os.Args[1])
	if err != nil {
		log.Fatalf("could not pull image from url:")
	}

	c, _ := pulledImage.ConfigFile()
	c.Config.Labels = map[string]string{
		"foobar": "hello world",
	}

	//config := v1.Config{Labels: map[string]string{
	//	"foobar": "hello world",
	//}}

	newImage, err := mutate.Config(pulledImage, c.Config)
	if err != nil {
		log.Fatalf("could not pull image from url: %s", err)
	}

	dir, err := ioutil.TempDir("", "deplab-crane-")
	if err != nil {
		log.Fatalf("Could not create temp directory. %s", err)
	}
	c2, _ := newImage.ConfigFile()
	fmt.Println(c.Config)
	fmt.Println(c2.Config)

	fmt.Println(dir)

	if err := crane.Save(newImage, os.Args[1]+"-with-label", dir+"/image.tgz"); err != nil {
		panic(err)
	}
}
