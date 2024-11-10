package client

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// UploadEntity binary file
func (c *GkClient) UploadEntity() error {
	var err error
	fp := ""
	guid := ""
	fmt.Print("Enter filepath: ")
	_, err = fmt.Fscan(os.Stdin, &fp)
	if err != nil {
		return err
	}
	if _, err = os.Stat(fp); err != nil {
		return err
	}
	fmt.Print("Enter GUID entity: ")
	_, err = fmt.Fscan(os.Stdin, &guid)
	if err != nil {
		return err
	}
	return c.sendBinaryData(fp, guid)
}

func (c *GkClient) sendBinaryData(fp string, guid string) error {
	fmt.Printf("Uploading file %s as %s\n", fp, guid)
	data, err := os.Open(fp)
	if err != nil {
		return err
	}
	var fileBody bytes.Buffer
	writer := multipart.NewWriter(&fileBody)
	filePart, err := writer.CreateFormFile("file", filepath.Base(fp))
	if err != nil {
		fmt.Println("CreateFormFile error : ", err)
		return err
	}
	_, err = io.Copy(filePart, data)

	if err != nil {
		fmt.Println("io.Copy error : ", err)
		return err
	}
	err = writer.Close()
	if err != nil {
		fmt.Println("Writer close error : ", err)
		return err
	}
	defer func() {
		if err := data.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	err = writer.Close()
	if err != nil {
		fmt.Println("Writer close error : ", err)
		return err
	}
	httpClient := &http.Client{
		Timeout: 3 * time.Second,
	}
	req, err := http.NewRequest("PUT", c.Conf.Addr+"/v1/entities/"+guid, &fileBody)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", c.token)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return NewError(resp.StatusCode)
	}
	return nil
}
