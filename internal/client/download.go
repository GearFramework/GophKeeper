package client

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"strings"
	"time"
)

// DownloadEntity download binary file by GUID
func (c *GkClient) DownloadEntity() error {
	var err error
	dest := ""
	guid := ""
	fmt.Print("Enter GUID entity: ")
	_, err = fmt.Fscan(os.Stdin, &guid)
	if err != nil {
		return err
	}
	fmt.Print("Enter destination directory: ")
	_, err = fmt.Fscan(os.Stdin, &dest)
	if err != nil {
		return err
	}
	if dest == "" {
		dest = "./"
	}
	if _, err = os.Stat(dest); err != nil {
		return err
	}
	httpClient := &http.Client{
		Timeout: 3 * time.Second,
	}
	req, err := http.NewRequest(
		"GET", c.Conf.Addr+"/v1/entities/download/"+guid, nil,
	)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", c.token)
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return NewError(resp.StatusCode)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	cd := resp.Header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(cd)
	filename := params["filename"]
	if filename == "" {
		filename = "in.dat"
	}
	dest = strings.TrimRight(dest, "/") + "/" + filename
	fd, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		err := fd.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	_, err = io.Copy(fd, resp.Body)
	if err != nil {
		return err
	}
	fmt.Println("Download complete as " + dest)
	return nil
}
