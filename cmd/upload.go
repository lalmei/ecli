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
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/textproto"
	"os"

	"github.com/spf13/cobra"

	"gopkg.in/mgo.v2/bson"
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
	},
}

/*
	The generated HTTP request must look like:

	Minimal request header:
		Content-Range: bytes 0-499999/973334
		Content-Disposition: attachment; filename="image_152.tiff"
		Content-Length: 500758
		Content-Type: multipart/form-data; boundary=---------------------------15160088341497182101124616286

	Body:
		-----------------------------15160088341497182101124616286
		Content-Disposition: form-data; name="extra"

		{"token":"atoken", "slideId": 123}
		-----------------------------15160088341497182101124616286
		Content-Disposition: form-data; name="files[]"; filename="annotations.json"
		Content-Type: application/json

		["a valid JSON annotations document"]
		-----------------------------15160088341497182101124616286
*/
func upload(file string) error {
	// Chunk size is 512k
	buf := make([]byte, 512*1024)
	f, err := os.Open("dd")
	if err != nil {
		return err
	}
	defer f.Close()
	for {
		_, err := f.Read(buf)
		if err != nil {
			return err
		}
		/*
			_, err = makeMultiPartChunkedRequest(file, "TOKEN HERE", buf)
			if err != nil {
				return err
			}
		*/
	}
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

type uploadSlideArgs struct {
	Token     string        `json:"token"`
	ParentId  bson.ObjectId `json:"parentId"`
	SlideName string        `json:"slideName"`
	// Path to the raw data. Can point to either a regular file
	// or directory.
	Path      string         `json:"path"`
	PixelSize slidePixelSize `json:"pixelSize"`
	// Date as a string so we can define custom format
	// requirement anytime.
	StartDate string   `json:"startDate"`
	filename  string   // Current filename being uploaded
	Labels    []*label `json:"labels"`
	Kind      string   `json:"imageType"` // Image type, if any
}

type simpleReader struct {
	*bytes.Buffer
}

func makeMultiPartChunkedRequest(filename string, token string, content []byte) (*http.Request, error) {
	buf := new(bytes.Buffer)
	mp := multipart.NewWriter(buf)
	defer mp.Close()
	mtype := "application/json"

	// Creates part for extra data required by Coquelicot.
	partw, err := mp.CreateFormField("extra")
	if err != nil {
		return nil, err
	}
	args := new(uploadSlideArgs)
	args.Token = token
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
	h.Set("Content-Type", mtype)
	partw, err = mp.CreatePart(h)
	if err != nil {
		return nil, err
	}
	if _, err := partw.Write(content); err != nil {
		return nil, err
	}

	r, _ := http.NewRequest("POST", "/endpoint", simpleReader{buf}) // FIXME

	r.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	r.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", mp.Boundary()))
	// Coquelicot only reads (and stores) the part named "files[]". So the Content-Range must be the length
	// of that part's body, NOT the length of the body of all parts in the request.
	r.Header.Set("Content-Range", fmt.Sprintf("bytes 0-%d/%d", len(content)-1, len(content)))

	// Set to true to print the HTTP request. Beware that printing the request's body will make it unavailable
	// for later processing.
	debug := false
	if debug {
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
}
