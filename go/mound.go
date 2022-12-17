package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"syscall"

	"encoding/json"

	"github.com/delaemon/go-uuidv4"
)

var MOUND_DIR string = "/tmp/mound_data"

type Mound struct {
	Did     string        `json:"did"`
	Program string        `json:"program"`
	Version string        `json:"version"`
	Status  int           `json:"status"`
	Blobs   []interface{} `json:"blobs"` // Holds int or string name for blob
	Links   []string      `json:"links"`
}

type Blob struct {
	Mound  *Mound
	BlobNo int
}

func setup(moundDir string) {
	MOUND_DIR = moundDir
}

func makeSemver(major int, minor int, patch int) string {
	return fmt.Sprintf("semver|%d.%d.%d", major, minor, patch)
}

func makeMound(program string, version string) (Mound, error) {
	did, err := uuidv4.Generate()
	if err != nil {
		return Mound{}, err
	}
	mound := Mound{
		Did:     did,
		Program: program,
		Version: version,
		Status:  -1,
		Blobs:   []interface{}{},
		Links:   []string{},
	}
	if err1 := mound._writeDoc(); err1 != nil {
		return Mound{}, err1
	}
	return mound, nil
}

func (mound *Mound) _writeDoc() error {
	dir := path.Join(MOUND_DIR, mound.Did[0:2], mound.Did[2:4], mound.Did[4:6], mound.Did[6:8])
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	name := path.Join(dir, "doc")
	data, err := json.Marshal(mound)
	if err != nil {
		return err
	}
	data = append(data, 0x0A)
	if err := os.WriteFile(name, data, os.ModePerm&0o666); err != nil {
		return err
	}
	return nil
}

func (mound *Mound) close(status int) error {
	mound.Status = status
	if err := mound._writeDoc(); err != nil {
		return err
	}
	return nil
}

func (mound *Mound) link(sourceDID string) error {
	// Mound only produces uuidv4 UUIDs and this Validate() ensures the version is 4.
	if !uuidv4.Validate(sourceDID) {
		return fmt.Errorf("sourceDID is not a UUID")
	}
	duplicate := false
	for i := 0; i < len(mound.Links); i += 1 {
		if mound.Links[i] == sourceDID {
			duplicate = true
			break
		}
	}
	if !duplicate {
		mound.Links = append(mound.Links, sourceDID)
	}
	return mound._writeDoc()
}

func (mound *Mound) blob(blobName ...string) (Blob, error) {
	if len(blobName) > 1 {
		return Blob{}, fmt.Errorf("One or zero arguments to blob()")
	}
	blobNo := len(mound.Blobs)
	var name interface{} = blobNo
	if len(blobName) > 0 {
		name = blobName[0]
	}
	blobfname := path.Join(MOUND_DIR, mound.Did[0:2], mound.Did[2:4], mound.Did[4:6], mound.Did[6:8], fmt.Sprintf("%d", blobNo))
	if err := os.WriteFile(blobfname, []byte{}, os.ModePerm&0o666); err != nil {
		return Blob{}, err
	}
	mound.Blobs = append(mound.Blobs, name)
	if err := mound._writeDoc(); err != nil {
		return Blob{}, err
	}
	return Blob{Mound: mound, BlobNo: blobNo}, nil
}

func (blob *Blob) Print(argv ...any) error {
	blobfname := path.Join(MOUND_DIR, blob.Mound.Did[0:2], blob.Mound.Did[2:4], blob.Mound.Did[4:6], blob.Mound.Did[6:8], fmt.Sprintf("%d", blob.BlobNo))
	text := fmt.Sprint(argv...)
	fout, err := os.OpenFile(blobfname, syscall.O_WRONLY|syscall.O_APPEND|syscall.O_CREAT, os.ModePerm&0o666)
	if err != nil {
		return err
	}
	_, err1 := fout.Write([]byte(text))
	if err2 := fout.Close(); err2 != nil && err1 == nil {
		err1 = err2
	}
	return err1
}

func (blob *Blob) Println(argv ...any) error {
	blobfname := path.Join(MOUND_DIR, blob.Mound.Did[0:2], blob.Mound.Did[2:4], blob.Mound.Did[4:6], blob.Mound.Did[6:8], fmt.Sprintf("%d", blob.BlobNo))
	text := fmt.Sprintln(argv...)
	fout, err := os.OpenFile(blobfname, syscall.O_WRONLY|syscall.O_APPEND|syscall.O_CREAT, os.ModePerm&0o666)
	if err != nil {
		return err
	}
	_, err1 := fout.Write([]byte(text))
	if err2 := fout.Close(); err2 != nil && err1 == nil {
		err1 = err2
	}
	return err1
}

func (blob *Blob) Printf(format string, argv ...any) error {
	blobfname := path.Join(MOUND_DIR, blob.Mound.Did[0:2], blob.Mound.Did[2:4], blob.Mound.Did[4:6], blob.Mound.Did[6:8], fmt.Sprintf("%d", blob.BlobNo))
	text := fmt.Sprintf(format, argv...)
	fout, err := os.OpenFile(blobfname, syscall.O_WRONLY|syscall.O_APPEND|syscall.O_CREAT, os.ModePerm&0o666)
	if err != nil {
		return err
	}
	_, err1 := fout.Write([]byte(text))
	if err2 := fout.Close(); err2 != nil && err1 == nil {
		err1 = err2
	}
	return err1
}

func main() {
	mound, err := makeMound("mound-go", makeSemver(1, 0, 0))
	if err != nil {
		log.Fatal(err)
	}

	b0, err := mound.blob()
	if err != nil {
		log.Fatal(err)
	}
	if err := b0.Println("Hello, Go!"); err != nil {
		log.Fatal(err)
	}
	if err := b0.Println("Hello, Go!"); err != nil {
		log.Fatal(err)
	}

	b1, err := mound.blob("test")
	if err != nil {
		log.Fatal(err)
	}
	if err := b1.Println("This is a test"); err != nil {
		log.Fatal(err)
	}

	if err := mound.close(0); err != nil {
		log.Fatal(err)
	}

	out, err := json.Marshal(mound)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
}
