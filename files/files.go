package files

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

var print = fmt.Println

//GetFileContentType check file type
func GetFileContentType(out *os.File) (string, error) {

	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {

		return "", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

//OSReadDir List of Folders
func OSReadDir(root string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {

		f, _ := os.Open(root + file.Name())

		defer f.Close()

		_, err = GetFileContentType(f)

		if err != nil {
			files = append(files, file.Name())

		}

	}
	return files, nil
}

//OSReadFile List of Files
func OSReadFile(root string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}
	for _, file := range fileInfo {
		f, _ := os.Open(root + file.Name())
		defer f.Close()

		_, err := GetFileContentType(f)

		if err == nil {
			path := strings.Replace(root, "public", "file", 1)
			files = append(files, path+file.Name())

		}

	}
	return files, nil
}

//Upload files on server
func Upload(c *gin.Context) {

	type Filepaths struct {
		Filepath []string
	}
	var filepath Filepaths
	_ = filepath
	var paths []string
	_ = paths

	form, err := c.MultipartForm()

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	files := form.File["files"]

	folder := c.PostForm("folder")
	subfolder := c.PostForm("subfolder")

	print(subfolder, len(files))

	os.Mkdir(fmt.Sprintf("public/%s/", folder), os.ModePerm)

	switch {
	case len(subfolder) < 1:
		for _, file := range files {

			if err := c.SaveUploadedFile(file, fmt.Sprintf("public/%s/%s", folder, file.Filename)); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
				return
			}

			cmd := exec.Command("/home/kot/CovertIFCtoGLTF/IfcConvert " + "--use-element-guids" + fmt.Sprintf("public/%s/%s", folder, file.Filename) + " " + fmt.Sprintf("public/%s/%s", folder, strings.Replace(file.Filename, "ifc", "dae", 1)))
			print(cmd.Output())
			paths = append(paths, fmt.Sprintf("public/%s/%s", folder, strings.Replace(file.Filename, "ifc", "dae", 1)))

		}
		filepath = Filepaths{
			Filepath: paths,
		}
	case len(subfolder) > 0:
		os.Mkdir(fmt.Sprintf("public/%s/%s", folder, subfolder), os.ModePerm)
		for _, file := range files {

			if err := c.SaveUploadedFile(file, fmt.Sprintf("public/%s/%s/%s", folder, subfolder, file.Filename)); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
				return
			}

			cmd := exec.Command("/home/kot/CovertIFCtoGLTF/IfcConvert " + "--use-element-guids" + fmt.Sprintf("public/%s/%s/%s", folder, subfolder, file.Filename) + " " + fmt.Sprintf("public/%s/%s/%s", folder, subfolder, strings.Replace(file.Filename, "ifc", "dae", 1)))
			print(cmd.Output())

			paths = append(paths, fmt.Sprintf("/file/%s/%s/%s", folder, subfolder, strings.Replace(file.Filename, "ifc", "dae", 1)))

		}
		filepath = Filepaths{
			Filepath: paths,
		}
	}

	jsonData, err := json.Marshal(filepath)
	_ = jsonData
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	}

	c.JSON(http.StatusOK, gin.H{"filepath": string(jsonData)})

}

//Fileslist of files on server
func Fileslist(c *gin.Context) {

	folder := c.PostForm("folder")
	subfolder := c.PostForm("subfolder")

	switch {
	case len(string(subfolder)) < 1:
		root := fmt.Sprintf("public/%s/", folder)

		dir, err := OSReadDir(root)

		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		}

		files, err := OSReadFile(root)

		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		}
		c.JSON(http.StatusOK, gin.H{"subfolders": dir, "files": files})
	case len(string(subfolder)) > 0:

		root := fmt.Sprintf("public/%s/%s/", folder, subfolder)
		dir, err := OSReadDir(root)

		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		}

		files, err := OSReadFile(root)

		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		}
		c.JSON(http.StatusOK, gin.H{"subfolders": dir, "files": files})
	}

}

//Mkrmsubfolders Remove and make folder and files
func Mkrmsubfolders(c *gin.Context) {
	doit := c.PostForm("doit")
	folder := c.PostForm("folder")
	subfolders := c.PostFormArray("subfolders")

	switch {
	case doit == "rm":
		for _, subfolder := range subfolders {
			print(subfolder)
			err := os.RemoveAll(fmt.Sprintf("public/%s/%s", folder, subfolder))
			if err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("stat: %s", err.Error()))
			}
		}
	case doit == "mk":
		for _, subfolder := range subfolders {
			print(subfolder)
			err := os.Mkdir(fmt.Sprintf("public/%s/%s", folder, subfolder), os.ModePerm)
			if err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("stat: %s", err.Error()))
			}
		}

	}

}

//Rmfiles Remove and make folder and files
func Rmfiles(c *gin.Context) {

	files := c.PostFormArray("files")
	for _, file := range files {

		err := os.Remove(strings.Replace(file, "file", "public", 1))
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("stat: %s", err.Error()))
		}
	}

}
