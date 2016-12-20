// Copyright © 2016 Mathias Monnerville <mathias.monnerville@keeneye.tech>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/textproto"
	"os"
	"path"
	"strings"

	"ecli/api"
	"ecli/config"
	"ecli/core"

	"github.com/spf13/cobra"

	"gopkg.in/mgo.v2/bson"
)

const (
	ChunkSize    = 512 // Chunk size is 512kB
	minChunkSize = 64  // kB
)

var (
	cfgDebug          bool
	cfgChunkSize      uint16
	cfgParentId       string
	cfgImageType      string
	cfgPixelSizeValue float64
	cfgPixelSizeUnit  string
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:     "upload IMAGE",
	Aliases: []string{"up"},
	Short:   "Upload an image to the platform",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			usageErrorExit(cmd, "Missing image file.")
		}
		endpoint, token, err := config.LoadSession()
		if err != nil {
			errorExit(err)
		}
		// Quick hack
		endpoint = strings.Replace(endpoint, "/api/v2", "/upload/files", 1)
		// Only Chunk size > minChunkSize allowed
		if cfgChunkSize < minChunkSize {
			errorExit(fmt.Errorf("Chunk size must be greater than %d", minChunkSize))
		}
		if cfgParentId == "" {
			// Get the root node to attach the image to (as a parent)
			wl, err := api.WorkList("")
			if err != nil {
				errorExit(err)
			}
			paths := wl["path"].([]interface{})
			top := paths[len(paths)-1]
			cfgParentId = top.(map[string]interface{})["id"].(string)
		}
		if err := upload(args[0], endpoint, token); err != nil {
			errorExit(err)
		}
	},
}

func upload(filename, endpoint, token string) error {
	// Get file size
	fi, err := os.Stat(filename)
	if err != nil {
		return err
	}
	size := fi.Size()

	buf := make([]byte, int(cfgChunkSize)*1024)
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Use a custom HTTP client to send back cookies in requests.
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	client := &http.Client{
		Jar: cookieJar,
	}

	if !cfgDebug {
		fmt.Printf("Uploading ")
	}

	var mimeType string // Computed on first chunk of data
	var offset int64
	k := 1
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if k == 1 {
			// Compute file's mimetype
			mimeType = http.DetectContentType(buf)
		}

		min, max := offset, offset+int64(n)-1
		if cfgDebug {
			fmt.Printf("chunk %d (%d-%d/%d) ...\n", k, min, max, size)
		}
		offset += int64(n)

		r, err := makeMultiPartChunkedRequest(path.Base(filename), endpoint, token,
			bson.ObjectIdHex(cfgParentId), min, max, size, mimeType, buf)
		if err != nil {
			return err
		}
		if !cfgDebug {
			resp, err := client.Do(r)
			if err != nil {
				return err
			}
			if err := core.DebugResponse(resp); err != nil {
				return err
			}
			fmt.Printf(".")
		}
		k++
	}
	fmt.Printf(" DONE\nImage %q successfully uploaded to the platform.\n", path.Base(filename))
	return nil
}

type slidePixelSize struct {
	// The value of the pixel size
	Value float64 `json:"value"`
	// The unit of the pixel size
	Unit string `json:"unit"`
}

// A label is used for slide annotation. Can tag a slide or 3d region
// in a slide.
type label struct {
	Name        string `json:"name"`
	Color       string `json:"color"` // Hex color
	Description string `json:"description"`
}

type uploadImageArgs struct {
	Token     string         `json:"token"`
	ParentId  bson.ObjectId  `json:"parentId"`
	SlideName string         `json:"slideName"`
	PixelSize slidePixelSize `json:"pixelSize"`
	Labels    []*label       `json:"labels"`
	ImageType string         `json:"imageType"`
}

type simpleReader struct {
	*bytes.Buffer
}

/*
	The generated HTTP request must look like:

	Recommanded request header:
		User-Agent: Engine CLI
		Content-Range: bytes 0-499999/973334
		Content-Disposition: attachment; filename="image_152.tiff"
		Content-Length: 500758
		Content-Type: multipart/form-data; boundary=---------------------------15160088341497182101124616286

	Body:
		-----------------------------15160088341497182101124616286
		Content-Disposition: form-data; name="extra"

		{"pixelSize":{"value":0.92,"unit":"um"},"slideName":"name","token":"atoken","parentId":"58518280e7798910f77ca485","labels":[{"name":"cell type 1","color":"#5cb85c","description":""}],"imageType":"generic-image"}
		-----------------------------15160088341497182101124616286
		Content-Disposition: form-data; name="files[]"; filename="annotations.json"
		Content-Type: image/tiff

		["a valid JSON annotations document"]
		-----------------------------15160088341497182101124616286

min-max specify the content range within the file.
*/
func makeMultiPartChunkedRequest(filename, endpoint, token string, parentId bson.ObjectId,
	min, max, size int64, contentType string, chunk []byte) (*http.Request, error) {

	buf := new(bytes.Buffer)
	mp := multipart.NewWriter(buf)
	defer mp.Close()

	// Creates part for extra data required by Coquelicot.
	partw, err := mp.CreateFormField("extra")
	if err != nil {
		return nil, err
	}
	args := new(uploadImageArgs)
	args.Token = token
	args.SlideName = filename
	args.ParentId = parentId
	args.ImageType = cfgImageType
	args.PixelSize = slidePixelSize{cfgPixelSizeValue, cfgPixelSizeUnit}
	data, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}
	if _, err := partw.Write(data); err != nil {
		return nil, err
	}

	// Creates part for JSON content. Uses CreatePart() instead of CreateFormFile() to set a custom
	// content type, not the default application/octet-stream.
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "files[]", filename))
	h.Set("Content-Type", contentType)
	partw, err = mp.CreatePart(h)
	if err != nil {
		return nil, err
	}
	if !cfgDebug {
		if _, err := partw.Write(chunk); err != nil {
			return nil, err
		}
	}
	r, _ := http.NewRequest("POST", endpoint, simpleReader{buf})

	r.Header.Set("User-Agent", fmt.Sprintf("ecli/%s", version))
	r.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	r.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", mp.Boundary()))
	// Coquelicot only reads (and stores) the part named "files[]". So the Content-Range must be the length
	// of that part's body, NOT the length of the body of all parts in the request.
	r.Header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", min, max, size))

	// Set to true to print the HTTP request. Beware that printing the request's body will make it unavailable
	// for later processing.
	if cfgDebug {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			return nil, err
		}
		fmt.Println(string(dump))
	}

	return r, nil
}

func init() {
	slideCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	uploadCmd.Flags().BoolVar(&cfgDebug, "debug", false, "Show request debugging info only")
	uploadCmd.Flags().Uint16Var(&cfgChunkSize, "chunk-size", ChunkSize, "Chunk size in kB")
	uploadCmd.Flags().StringVar(&cfgParentId, "parent-id", "", "Image's parent ID (worklist)")
	uploadCmd.Flags().StringVarP(&cfgImageType, "image-type", "t", "generic-image", "Image type")
	uploadCmd.Flags().Float64VarP(&cfgPixelSizeValue, "pixel-size", "p", 0, "Pixel size value")
	uploadCmd.Flags().StringVarP(&cfgPixelSizeUnit, "pixel-size-unit", "u", "um", "Pixel size unit")
}
