package NCMB

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func SetUpFile() NCMB {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	applicationKey, clientKey := os.Getenv("APPLICATION_KEY"), os.Getenv("CLIENT_KEY")
	ncmb := Initialize(applicationKey, clientKey)
	return ncmb
}

func TearDownFile(ncmb NCMB) {
	fmt.Printf("TearDown %s\n", ncmb)
}

func TestFileUploadText(t *testing.T) {
	ncmb := SetUpFile()
	file := ncmb.File("test.txt")
	text := []byte("Hello, NCMB!")
	file.Bytes = text
	bol, err := file.Upload()
	if err != nil {
		t.Errorf("file.Save() = %T, %s", bol, err)
	}
	TearDownFile(ncmb)
}

func TestFileUploadImage(t *testing.T) {
	ncmb := SetUpFile()
	file := ncmb.File("ncmb.png")
	f, err := os.Open("ncmb.png")
	data, err := ioutil.ReadAll(f)
	if err != nil {
		t.Errorf("f.Read() = %s", err)
	}
	file.Bytes = data
	bol, err := file.Upload()
	if err != nil {
		t.Errorf("file.Save() = %T, %s", bol, err)
	}
	if d, err := file.GetDate("createDate"); err != nil {
		t.Errorf("file.GetDate() = %s, want not nil (%s)", d, err)
	}
	TearDownFile(ncmb)
}

func TestFileDownloadText(t *testing.T) {
	ncmb := SetUpFile()
	file := ncmb.File("test.txt")
	text := []byte("Hello, NCMB!")
	file.Bytes = text
	bol, err := file.Upload()
	if err != nil {
		t.Errorf("file.Save() = %T, %s", bol, err)
	}
	file.Fetch()
	data, err := file.Download()
	if err != nil {
		t.Errorf("file.Download() = %s", err)
	}
	if string(data) != string(text) {
		t.Errorf("file.Download() = %s, want %s", string(data), string(text))
	}
	TearDownFile(ncmb)
}

func TestFileDelete(t *testing.T) {
	ncmb := SetUpFile()
	file := ncmb.File("test.txt")
	text := []byte("Hello, NCMB!")
	file.Bytes = text
	bol, err := file.Upload()
	if err != nil {
		t.Errorf("file.Save() = %T, %s", bol, err)
	}
	bol, err = file.Delete()
	if err != nil {
		t.Errorf("file.Delete() = %T, %s", bol, err)
	}
	TearDownFile(ncmb)
}

func TestFileUploadWithAcl(t *testing.T) {
	ncmb := SetUpFile()
	file := ncmb.File("test2.txt")
	text := []byte("Hello, NCMB!")
	file.Bytes = text
	acl := ncmb.Acl()
	acl.SetPublicReadAccess(true)
	acl.SetPublicWriteAccess(false)
	file.SetAcl(acl)
	bol, err := file.Upload()
	if err != nil {
		t.Errorf("file.Save() = %T, %s", bol, err)
	}
	bol, err = file.Delete()
	if err == nil {
		t.Errorf("file.Delete() = %T, %s", bol, err)
	}
	TearDownFile(ncmb)
}

func TestFileUpdateText(t *testing.T) {
	ncmb := SetUpFile()
	file := ncmb.File("test.txt")
	text := []byte("Hello, NCMB!")
	file.Bytes = text
	bol, err := file.Upload()
	if err != nil {
		t.Errorf("file.Save() = %T, %s", bol, err)
	}
	acl, err := file.GetAcl()
	if err != nil {
		t.Errorf("file.GetAcl() = %T, %s", acl, err)
	}
	acl.SetRoleReadAccess("Admin", true)
	file.SetAcl(acl)
	bol, err = file.Update()
	if err != nil {
		t.Errorf("file.Update() = %T, %s", bol, err)
	}
	TearDownFile(ncmb)
}
