// Copyright © 2017 The Keen Eye Developers
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
	"net/url"
	"os"
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/keeneyetech/ecli/api"
	"github.com/keeneyetech/ecli/config"
	"github.com/keeneyetech/ecli/core"

	"github.com/spf13/cobra"

	"gopkg.in/mgo.v2/bson"
)

const (
	ChunkSize    = 512 // Chunk size is 512kB
	minChunkSize = 64  // kB
)

var (
	cfgChunkSize        uint16
	cfgDebug            bool
	cfgImageFormat      string
	cfgParentId         string
	cfgPixelSizeUnit    string
	cfgPixelSizeValue   float64
	cfgSlideDescription string
	cfgSlideName        string
	cfgSlideLabels      []string
)

// Label info to send with upload parameters.
var customLabels []*api.Label

// Jquery File Upload uses the builtin JS encodeURI() to format the filename before sending it.
// Example: cochléa (crop).tif would be sent as cochl%C3%A9a%20(crop).tif.
// See jquery.fileupload.js:456
// The engine needs the same escape mecanism to compute the same temporary chunk filename.
func encodeURI(s string) string {
	// Javascript's encodeURI() does not encode the following chars, but Go does so we need
	// to decode them afterwards.
	nE := map[byte]string{
		';': "%3B", ',': "%2C", '/': "%2F", '?': "%3F", ':': "%3A",
		'@': "%40", '&': "%26", '=': "%3D", '+': "%2B", '$': "%24",
		'!': "%21", '*': "%2A", '\'': "%27", '(': "%28", ')': "%29",
	}
	out := strings.Replace(url.QueryEscape(s), "+", "%20", -1)
	for a, code := range nE {
		out = strings.Replace(out, code, string(a), -1)
	}
	return out
}

func checkLabels(labels []string) ([]*api.Label, error) {
	if len(labels) == 0 {
		return nil, nil
	}
	customLabels = make([]*api.Label, len(labels))
	for k, name := range labels {
		l, err := api.OneLabel(name)
		if err != nil {
			return nil, fmt.Errorf("Label %q not found, please create it first. See `ecli label create`.", name)
		}
		customLabels[k] = &api.Label{
			l["name"].(string),
			l["color"].(string),
			l["description"].(string),
		}
	}
	return customLabels, nil
}

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:     "upload IMAGE",
	Aliases: []string{"up"},
	Short:   "Upload an image to the platform",
	Long: `For example, uploading a TIFF image with a 1um pixel size and apply 2 "retina"
and "core" labels on it can be performed with

  slide upload myfile.tiff -f tiff -p 1 -l "retina" -l "core"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			usageErrorExit(cmd, "Missing image file.")
		}
		var err error
		customLabels, err = checkLabels(cfgSlideLabels)
		if err != nil {
			log.Fatal(err)
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
		if cfgParentId == "root" {
			// Get the root node to attach the image to (as a parent)
			wl, err := api.WorkList("")
			if err != nil {
				errorExit(err)
			}
			paths := wl["path"].([]interface{})
			top := paths[len(paths)-1]
			cfgParentId = top.(map[string]interface{})["id"].(string)
		}
		filename := args[0]
		if cfgSlideName == "" {
			cfgSlideName = path.Base(filename)
		}
		if cfgSlideDescription == "" {
			cfgSlideDescription = "Uploaded with ecli " + core.Version
		}
		fmt.Println("Uploading, please wait ...")
		if err := upload(filename, endpoint, token); err != nil {
			errorExit(err)
		}
		fmt.Printf("\nImage %q successfully uploaded to the platform.\n", path.Base(filename))
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

	var mimeType string // Computed on first chunk of data
	var offset int64

	k := 1
	numChunks := (size / int64(int(cfgChunkSize)*1024)) + 1
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
			fmt.Printf("  Uploading %.1f%%\r", float64(k)/float64(numChunks)*100)
		}
		k++
	}
	return nil
}

type slidePixelSize struct {
	// The value of the pixel size
	Value float64 `json:"value"`
	// The unit of the pixel size
	Unit string `json:"unit"`
}

type uploadImageArgs struct {
	Token       string         `json:"token"`
	ParentId    bson.ObjectId  `json:"parentId"`
	PixelSize   slidePixelSize `json:"pixelSize"`
	Labels      []*api.Label   `json:"labels"`
	ImageFormat string         `json:"imageFormat"`
	SlideName   string         `json:"slideName"`
	Description string         `json:"description"`
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

		{"pixelSize":{"value":0.92,"unit":"um"},"slideName":"name","token":"atoken","parentId":"58518280e7798910f77ca485","labels":[{"name":"cell type 1","color":"#5cb85c","description":""}],"imageFormat":"tiff"}
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
	args.SlideName = cfgSlideName
	args.Description = cfgSlideDescription
	args.ParentId = parentId
	args.ImageFormat = cfgImageFormat
	args.PixelSize = slidePixelSize{cfgPixelSizeValue, cfgPixelSizeUnit}
	args.Labels = customLabels

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

	r.Header.Set("User-Agent", fmt.Sprintf("ecli/%s", core.Version))
	r.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, encodeURI(filename)))
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
	uploadCmd.Flags().StringVarP(&cfgParentId, "group-id", "g", "root", "Image's group ID (parent)")
	uploadCmd.Flags().StringVarP(&cfgImageFormat, "image-format", "f", "tiff", "Image Format")
	uploadCmd.Flags().Float64VarP(&cfgPixelSizeValue, "pixel-size", "p", 0, "Pixel size value")
	uploadCmd.Flags().StringVarP(&cfgPixelSizeUnit, "pixel-size-unit", "u", "um", "Pixel size unit")
	uploadCmd.Flags().StringVarP(&cfgSlideName, "name", "n", "", "Image name (default filename)")
	uploadCmd.Flags().StringVarP(&cfgSlideDescription, "description", "d", "", "Image description (default ecli version)")
	uploadCmd.Flags().StringArrayVarP(&cfgSlideLabels, "label", "l", []string{}, "Label to apply on image")
}
