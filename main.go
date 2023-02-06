package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"celsorodrigues.co.uk/gitclones/k8s"
	"celsorodrigues.co.uk/gitclones/ppk"
	"github.com/go-git/go-billy/v5/memfs"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {

	//create key test
	myPrivatekey, myPublicKey, err := ppk.CreateKey()
	if err != nil {
		log.Println("unable to create the key:", err)
	}
	secret := ppk.CreateSecret(myPrivatekey, myPublicKey)

	ctx := context.Background()
	clientSet := k8s.GetClient()
	_, err = clientSet.CoreV1().Secrets("dev-data").Create(ctx, secret, metav1.CreateOptions{})
	if err != nil {
		if k8serror.IsAlreadyExists(err) {
			log.Println("namespace already exists")
		} else {
			log.Println("error creating secret: ", err.Error())
		}
	}

	envKey, ok := os.LookupEnv("DEPLOY_KEY")
	if !ok {
		log.Println("ENV not found")
	}
	key := strings.Replace(envKey, "\\n", "\n", -1)

	publicKey, err := ppk.GetKey(key)
	if err != nil {
		log.Fatal(err)
	}
	store := memory.NewStorage()
	fs := memfs.New()

	_, err = gogit.Clone(store, fs, &gogit.CloneOptions{
		Auth:          publicKey,
		URL:           "git@github.com:celsoRodrigues/argotest.git",
		Progress:      os.Stdout,
		ReferenceName: plumbing.ReferenceName("refs/heads/main"),
		SingleBranch:  true,
		Depth:         1,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("clone success")
}
